package dht

// Implements a BitTorrent DHT Node (peer), as described in
// http://www.bittorrent.org/beps/bep_0005.html.

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	weakrand "math/rand"
	"net"
)

// NodeId

type NodeId string // of length 20

var UnknownNodeId NodeId = NodeId("") // empty string is unkown node ID

func GenerateNodeId() NodeId {
	// Generates a new random NodeID.
	// Panics if unable to generate random number.

	bytes := new([20]byte)
	n, err := rand.Read(bytes[:])

	if n < 20 {
		panic("too few bytes generated for some reason?")
	}

	if err != nil {
		panic(err)
	}

	return NodeId(bytes[:])
}

func (id NodeId) String() string {
	if len(id) > 0 {
		return hex.EncodeToString([]byte(id))
	} else {
		return "[unknown ID]"
	}
}

// RemoteNode

type RemoteNode struct {
	Id      NodeId
	Address net.UDPAddr
}

func RemoteNodeFromAddress(address net.UDPAddr) (remote *RemoteNode) {
	// Creates a RemoteNode with a known address but an unknown ID.
	// You may want to .Ping() this node so that it learns its ID!
	remote = new(RemoteNode)
	remote.Id = UnknownNodeId
	remote.Address = address
	return remote
}

func GenerateFakeRemoteNode() (remote *RemoteNode) {
	remote = new(RemoteNode)
	remote.Id = GenerateNodeId()
	remote.Address = net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}
	return remote
}

func (remote *RemoteNode) String() string {
	return fmt.Sprintf("<RemoteNode %v at %v>", remote.Id, remote.Address)
}

// LocalNode

type LocalNode struct {
	Id    NodeId
	Port  int
	Nodes []*RemoteNode
}

func NewLocalNode() (local *LocalNode) {
	local = new(LocalNode)
	local.Id = GenerateNodeId()
	local.Nodes = []*RemoteNode{}
	local.Port = 1024 + weakrand.Intn(8192)
	return local
}

func UnmarshalBencodingDict(dict bencoding.Dict) (local *LocalNode) {
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

func (local *LocalNode) String() string {
	return fmt.Sprintf("<LocalNode %v on :%v>", local.Id, local.Port)
}

type QueryResult chan<- *bencoding.Dict // what is the convention?

func (local *LocalNode) query(remote *RemoteNode) QueryResult {
	panic("Not Implemented")
}

func (local *LocalNode) Ping(remote *RemoteNode) QueryResult {
	return local.query(remote)
}

func (local *LocalNode) FindClosest(id NodeId) *RemoteNode {
	return nil
}
