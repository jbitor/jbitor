package dht

import (
	"crypto/rand"
	"github.com/jeremybanks/go-distributed/bencoding"
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

func (local *LocalNode) SendPing(remote *RemoteNode) (<-chan *bencoding.Dict, <-chan error) {
	pingResult := make(chan *bencoding.Dict)
	pingErr := make(chan error)

	query := local.sendQuery(remote, "ping", bencoding.Dict{})

	go func() {
		select {
		case value := <-query.Result:
			remote.Id = NodeId((*value)["id"].(bencoding.String))
			pingResult <- value
		case err := <-query.Err:
			pingErr <- err
		}
	}()

	return pingResult, pingErr
}

const PEER_CONTACT_INFO_LEN = 6
const NODE_CONTACT_INFO_ID_LEN = 20
const NODE_CONTACT_INFO_LEN = 26

func (local *LocalNode) FindNode(remote *RemoteNode, id NodeId) (<-chan []*RemoteNode, <-chan error) {
	findResult := make(chan []*RemoteNode)
	findErr := make(chan error)

	query := local.sendQuery(remote, "find_node", bencoding.Dict{
		"target": bencoding.String(id),
	})

	go func() {
		select {
		case value := <-query.Result:
			result := []*RemoteNode{}

			nodesData := (*value)["nodes"].(bencoding.String)

			for offset := 0; offset < len(nodesData); offset += NODE_CONTACT_INFO_LEN {
				nodeId := nodesData[offset : offset+NODE_CONTACT_INFO_ID_LEN]
				nodeAddress := decodeNodeAddress(nodesData[offset+NODE_CONTACT_INFO_ID_LEN : offset+NODE_CONTACT_INFO_LEN])

				remote := RemoteNodeFromAddress(nodeAddress)
				remote.Id = NodeId(nodeId)

				remote = local.AddOrGetRemoteNode(remote)

				result = append(result, remote)
			}

			findResult <- result
		case err := <-query.Err:
			findErr <- err
		}
	}()

	return findResult, findErr
}

func (local *LocalNode) SendGetPeers(remote *RemoteNode, id NodeId) (result <-chan *bencoding.Dict, err <-chan error) {
	logger.Fatalf("GetPeers() not implemented\n")
	return
}

func (local *LocalNode) SendAnnouncePeer(remote *RemoteNode, id NodeId) (result <-chan *bencoding.Dict, err <-chan error) {
	logger.Fatalf("AnnouncePeer() not implemented\n")
	return
}

func (local *LocalNode) RunRpcListen(rpcError chan<- error) {
	response := new([1024]byte)

	for {
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
