package dht

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"fmt"
	"io"
	weakrand "math/rand"
)

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

func (local *LocalNode) String() string {
	return fmt.Sprintf("<LocalNode %v on :%v>", local.Id, local.Port)
}

// RPC Requests

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
