package dht

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
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

	query.TransactionId = "he"

	local.OutstandingQueries[query.TransactionId] = query

	message["t"] = bencoding.String(query.TransactionId)

	encodedMessage, err := bencoding.Encode(message)

	if err != nil {
		query.Err <- err
		return
	}

	if err != nil {
		query.Err <- err
		return
	}

	fmt.Printf("Sending %v...\n", string(encodedMessage))

	local.Connection.WriteTo(encodedMessage, &remote.Address)

	fmt.Printf("Sent. Listening...\n")

	response := new([1024]byte) // XXX: This can't be right.
	n, remoteAddr, err := local.Connection.ReadFromUDP(response[:])

	if err != nil {
		query.Err <- err
		return
	}

	fmt.Printf("Got response?! %v from %v\n", string(response[:n]), remoteAddr)

	result, err := bencoding.Decode(response[:n])

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
		"q": bencoding.String("ping"),
		"a": bencoding.Dict{
			"id": bencoding.String(local.Id),
		},
		"y": bencoding.String("q"),
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

	time.Sleep(5000)

	terminated <- errors.New("Not implemented")
	return
}
