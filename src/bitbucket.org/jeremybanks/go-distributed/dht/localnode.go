package dht

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	weakrand "math/rand"
	"net"
	"time"
)

type LocalNode struct {
	Id                 NodeId
	Port               int
	Connection         *net.UDPConn
	Nodes              []*RemoteNode // TODO: proper spec-compliant routing-table
	OutstandingQueries map[string]*Query
}

func NewLocalNode() (local *LocalNode) {
	local = new(LocalNode)
	local.Id = GenerateNodeId()
	local.Nodes = []*RemoteNode{}
	local.Port = 1024 + weakrand.Intn(8192)
	local.OutstandingQueries = make(map[string]*Query)
	return local
}

func (local *LocalNode) String() string {
	return fmt.Sprintf("<LocalNode %v on :%v>", local.Id, local.Port)
}

// RPC Requests

type Query struct {
	TransactionId string
	Remote        *RemoteNode
	Result        chan *bencoding.Dict
	Err           chan error
}

func (local *LocalNode) query(remote *RemoteNode, message bencoding.Dict) (query *Query) {
	query = new(Query)
	query.Result = make(chan *bencoding.Dict)
	query.Err = make(chan error)
	query.Remote = remote

	go local.runQuery(remote, message, query)

	return query
}

func (local *LocalNode) runQuery(remote *RemoteNode, message bencoding.Dict, query *Query) {
	if message == nil {
		query.Err <- errors.New("message is nil")
	}

	if message["y"] == nil {
		query.Err <- errors.New("message missing type")
		return
	}

	transactionId := new([4]byte)
	_, err := rand.Read(transactionId[:])

	if err != nil {
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
	local.Connection.WriteTo(encodedMessage, &remote.Address)
}

type PingResult struct {
	Result chan *bencoding.Dict
	Err    chan error
}

func (local *LocalNode) Ping(remote *RemoteNode) (result *PingResult) {
	result = new(PingResult)
	result.Result = make(chan *bencoding.Dict)
	result.Err = make(chan error)

	query := local.query(remote, bencoding.Dict{
		"q": bencoding.String("ping"),
		"a": bencoding.Dict{
			"id": bencoding.String(local.Id),
		},
		"y": bencoding.String("q"),
	})

	go func() {
		select {
		case value := <-query.Result:
			remote.Id = NodeId((*value)["r"].(bencoding.Dict)["id"].(bencoding.String))

			result.Result <- value
		case err := <-query.Err:
			result.Err <- err
		}
	}()

	return result
}

// Bencoding

func LocalNodeFromBencodingDict(dict bencoding.Dict) (local *LocalNode) {
	local = new(LocalNode)

	panic("unmarshaling of LocalNode not implemented")
}

func (local *LocalNode) MarshalBencodingDict() (dict bencoding.Dict) {
	dict = bencoding.Dict{}

	panic("marshaling of LocalNode not implemented")
}

func (local *LocalNode) WriteBencodedTo(writer io.Writer) error {
	return local.MarshalBencodingDict().WriteBencodedTo(writer)
}

// Running

func (local *LocalNode) Run(terminated chan<- error) {
	// Main loop for LocalPeer's activity.
	// (Listening to replies and requests.)

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: local.Port,
	})
	defer conn.Close()

	if err != nil {
		terminated <- err
		return
	}

	local.Connection = conn

	response := new([1024]byte)

	for {
		n, remoteAddr, err := local.Connection.ReadFromUDP(response[:])

		if err != nil {
			fmt.Printf("Ignoring UDP read err: %v\n", err)
			continue
		}

		fmt.Printf("Got response?! %v from %v\n", string(response[:n]), remoteAddr)

		result, err := bencoding.Decode(response[:n])

		if err != nil {
			fmt.Printf("Ignoring un-bedecodable message: %v\n", err)
			continue
		}

		resultD, ok := result.(bencoding.Dict)

		if !ok {
			fmt.Printf("Ignoring bedecoded non-dict message: %v\n", err)
			return
		}

		transactionId := string(resultD["t"].(bencoding.String))

		query := local.OutstandingQueries[transactionId]

		query.Remote.LastResponseFrom = time.Now()

		query.Result <- &resultD
	}

	terminated <- errors.New("Not implemented")
	return
}
