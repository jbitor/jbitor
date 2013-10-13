package dht

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"fmt"
	"io"
	weakrand "math/rand"
	"net"
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

	rpcError := make(chan error)
	go local.RunRpcListen(rpcError)

	select {
	case err := <-rpcError:
		terminated <- err
	}
}
