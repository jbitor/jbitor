package dht

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"errors"
	"fmt"
	"io"
	weakrand "math/rand"
	"net"
)

type LocalNode struct {
	Id    NodeId
	Port  int
	Nodes []*RemoteNode // TODO: proper spec-compliant routing-table
}

func NewLocalNode() (local *LocalNode) {
	local = new(LocalNode)
	local.Id = GenerateNodeId()
	local.Nodes = []*RemoteNode{}
	local.Port = 1024 + weakrand.Intn(8192)
	return local
}

func (local *LocalNode) String() string {
	return fmt.Sprintf("<LocalNode %v on :%v>", local.Id, local.Port)
}

// RPC Requests

type Query struct {
	TransactionId string
	Result        chan *bencoding.Dict
	Err           chan error
}

func (local *LocalNode) query(remote *RemoteNode, message bencoding.Dict) (query *Query) {
	query = new(Query)
	query.Result = make(chan *bencoding.Dict)
	query.Err = make(chan error)

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

	query.TransactionId = "hello"
	message["t"] = bencoding.String(query.TransactionId)

	encodedMessage, err := bencoding.Encode(message)

	if err != nil {
		query.Err <- err
		return
	}

	conn, err := net.DialUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: local.Port,
	}, &remote.Address)

	if err != nil {
		query.Err <- err
		return
	}

	conn.Write(encodedMessage)

	response := new([1024]byte) // XXX: This can't be right.
	_, err = conn.Read(response[:])

	if err != nil {
		query.Err <- err
		return
	}

	fmt.Printf("Got response?! %v\n", response[:])

	conn.Close() // hmmm...

	result, err := bencoding.Decode(response[:])

	if err != nil {
		query.Err <- err
		return
	}

	resultD, ok := result.(bencoding.Dict)

	if !ok {
		query.Err <- errors.New("bencoded result was not a dict")
		return
	}

	query.Result <- &resultD

}

func (local *LocalNode) Ping(remote *RemoteNode) (query *Query) {
	return local.query(remote, bencoding.Dict{
		"y":  bencoding.String("ping"),
		"id": bencoding.String(local.Id),
	})
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
	// Main loop for LocalPeer's activity

	terminated <- errors.New("Not implemented")
}
