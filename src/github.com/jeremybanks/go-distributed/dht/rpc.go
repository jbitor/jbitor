package dht

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/jeremybanks/go-distributed/bencoding"
	"github.com/jeremybanks/go-distributed/torrent"
	"net"
	"time"
)

type Query struct {
	TransactionId string
	Remote        *RemoteNode
	Result        chan *bencoding.Dict
	Err           chan error
}

func decodeNodeAddress(encoded bencoding.String) (addr net.UDPAddr) {
	return net.UDPAddr{
		IP:   net.IPv4(encoded[0], encoded[1], encoded[2], encoded[3]),
		Port: int(encoded[4])<<8 + int(encoded[5]),
	}
}

func encodeNodeAddress(addr net.UDPAddr) (encoded bencoding.String) {
	if addr.Port >= (1<<32) || addr.Port < 0 {
		panic("Port out of bounds?")
	}

	ip4 := addr.IP.To4()

	return bencoding.String([]byte{
		ip4[0],
		ip4[1],
		ip4[2],
		ip4[3],
		byte((addr.Port >> 8) & 0xFF),
		byte(addr.Port & 0xFF),
	})
}

func (local *LocalNode) sendQuery(remote *RemoteNode, queryType string, arguments bencoding.Dict) (query *Query) {
	query = new(Query)
	query.Result = make(chan *bencoding.Dict)
	query.Err = make(chan error)
	query.Remote = remote

	if arguments == nil {
		arguments = bencoding.Dict{}
	}

	arguments["id"] = bencoding.String(local.Id)

	// XXX: assert that these keys are not already present?
	message := bencoding.Dict{
		"y": bencoding.String("q"),
		"q": bencoding.String(queryType),
		"a": arguments,
	}

	transactionId := new([4]byte)
	if _, err := rand.Read(transactionId[:]); err != nil {
		query.Err <- err
		return
	}

	query.TransactionId = string(transactionId[:])

	local.OutstandingQueries[query.TransactionId] = query

	message["t"] = bencoding.String(query.TransactionId)

	encodedMessage, err := bencoding.Encode(message)

	if err != nil {
		query.Err <- err
		return
	}

	remote.LastRequestTo = time.Now()

	go func() {
		// XXX:
		// Does this wait longer than necessary to send the packet?
		local.Connection.WriteTo(encodedMessage, &remote.Address)
	}()

	return query
}

func (local *LocalNode) Ping(remote *RemoteNode) (<-chan *bencoding.Dict, <-chan error) {
	pingResult := make(chan *bencoding.Dict)
	pingErr := make(chan error)

	query := local.sendQuery(remote, "ping", bencoding.Dict{})

	go func() {
		select {
		case value := <-query.Result:
			remote.Id = NodeId((*value)["id"].(bencoding.String))

			remote.ConsecutiveFailedQueries = 0

			pingResult <- value

		case err := <-query.Err:
			remote.ConsecutiveFailedQueries++
			pingErr <- err
		}
	}()

	return pingResult, pingErr
}

const PEER_CONTACT_INFO_LEN = 6
const NODE_CONTACT_INFO_ID_LEN = 20
const NODE_CONTACT_INFO_LEN = 26

func (local *LocalNode) decodeNodesString(nodesData bencoding.String, source *RemoteNode) ([]*RemoteNode, error) {
	result := make([]*RemoteNode, 0)

	for offset := 0; offset < len(nodesData); offset += NODE_CONTACT_INFO_LEN {
		nodeId := nodesData[offset : offset+NODE_CONTACT_INFO_ID_LEN]
		nodeAddress := decodeNodeAddress(nodesData[offset+NODE_CONTACT_INFO_ID_LEN : offset+NODE_CONTACT_INFO_LEN])

		resultRemote := RemoteNodeFromAddress(nodeAddress)
		resultRemote.Source = source
		resultRemote.Id = NodeId(nodeId)

		resultRemote = local.AddOrGetRemoteNode(resultRemote)

		result = append(result, resultRemote)
	}

	return result, nil
}

func (local *LocalNode) FindNode(remote *RemoteNode, id NodeId) (<-chan []*RemoteNode, <-chan error) {
	findResult := make(chan []*RemoteNode)
	findErr := make(chan error)

	query := local.sendQuery(remote, "find_node", bencoding.Dict{
		"target": bencoding.String(id),
	})

	go func() {
		select {
		case value := <-query.Result:
			nodesData, ok := (*value)["nodes"].(bencoding.String)
			if !ok {
				remote.ConsecutiveFailedQueries++
				findErr <- errors.New(".nodes string does not exist")
				return
			}

			result, err := local.decodeNodesString(nodesData, remote)
			if err != nil {
				findErr <- err
				return
			}

			remote.ConsecutiveFailedQueries = 0
			findResult <- result

		case err := <-query.Err:
			remote.ConsecutiveFailedQueries++
			findErr <- err
		}
	}()

	return findResult, findErr
}

func (local *LocalNode) GetPeers(remote *RemoteNode, infoHash string) (<-chan []*torrent.RemotePeer, <-chan []*RemoteNode, <-chan error) {
	peersResult := make(chan []*torrent.RemotePeer)
	nodesResult := make(chan []*RemoteNode)
	getPeersErr := make(chan error)

	query := local.sendQuery(remote, "get_peers", bencoding.Dict{
		"info_hash": bencoding.String(infoHash),
	})

	go func() {
		select {
		case value := <-query.Result:
			peerData, peersOk := (*value)["values"].(bencoding.List)
			nodesData, nodesOk := (*value)["nodes"].(bencoding.String)

			if peersOk {
				result := make([]*torrent.RemotePeer, len(peerData))

				for i, data := range peerData {
					dataStr, ok := data.(bencoding.String)
					if !ok {
						getPeersErr <- errors.New(".values contained non-string")
						remote.ConsecutiveFailedQueries++
						return
					}

					addr := torrent.DecodePeerAddress(dataStr)
					result[i] = &torrent.RemotePeer{Address: addr}
				}

				remote.ConsecutiveFailedQueries = 0

				peersResult <- result
			} else if nodesOk {
				result, err := local.decodeNodesString(nodesData, remote)
				if err != nil {
					remote.ConsecutiveFailedQueries++
					getPeersErr <- err
					return
				}

				remote.ConsecutiveFailedQueries = 0
				nodesResult <- result
			} else {
				getPeersErr <- errors.New(fmt.Sprintf("response did not include peer or node list - %v", *value))
				remote.ConsecutiveFailedQueries++
			}

		case err := <-query.Err:
			remote.ConsecutiveFailedQueries++
			getPeersErr <- err
		}
	}()

	return peersResult, nodesResult, getPeersErr
}

func (local *LocalNode) AnnouncePeer(remote *RemoteNode, id NodeId) (result <-chan *bencoding.Dict, err <-chan error) {
	logger.Fatalf("AnnouncePeer() not implemented\n")
	return
}

func (local *LocalNode) runRpcListen(terminate <-chan bool, terminated chan<- error) {
	response := new([1024]byte)

	for {
		logger.Printf("Waiting for next incoming UDP message.\n")

		n, remoteAddr, err := local.Connection.ReadFromUDP(response[:])

		_ = remoteAddr

		if err != nil {
			logger.Printf("Ignoring UDP read err: %v\n", err)
			continue
		}

		result, err := bencoding.Decode(response[:n])

		if err != nil {
			logger.Printf("Ignoring un-bedecodable message: %v\n", err)
			continue
		}

		resultD, ok := result.(bencoding.Dict)

		if !ok {
			logger.Printf("Ignoring bedecoded non-dict message: %v\n", err)
			continue
		}

		transactionId := string(resultD["t"].(bencoding.String))

		query := local.OutstandingQueries[transactionId]

		query.Remote.LastResponseFrom = time.Now()

		resultBody, ok := resultD["r"].(bencoding.Dict)

		if !ok {
			logger.Printf("Ignoring response with non-dict contents.\n")
			continue
		}

		query.Result <- &resultBody

		delete(local.OutstandingQueries, transactionId)
	}
}
