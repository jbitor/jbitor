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
	Nodes              map[string]*RemoteNode // is not a spec-compliant routing-table
	OutstandingQueries map[string]*Query
}

// how do you de-duplicate/index nodes in a standard way?

func NewLocalNode() (local *LocalNode) {
	local = new(LocalNode)
	local.Id = GenerateNodeId()
	local.Port = 1024 + weakrand.Intn(8192)
	local.OutstandingQueries = make(map[string]*Query)
	local.Nodes = map[string]*RemoteNode{}
	return local
}

func (local *LocalNode) AddOrGetRemoteNode(remote *RemoteNode) *RemoteNode {
	// If a node with the same address is already in .Nodes, returns that node.
	// Otherwise, add remote to .Nodes and return it.

	key := remoteNodeKey(remote.Address)

	if existingRemote, ok := local.Nodes[key]; ok {
		remote = existingRemote
	} else {
		local.Nodes[key] = remote
	}

	return remote

}

func (local *LocalNode) String() string {
	return fmt.Sprintf("<LocalNode %v on :%v>", local.Id, local.Port)
}

func remoteNodeKey(addr net.UDPAddr) string {
	return fmt.Sprintf("%v:%v", addr.IP, addr.Port)
}

// Bencoding

func LocalNodeFromBencodingDict(dict bencoding.Dict) (local *LocalNode) {
	local = new(LocalNode)

	local.Id = NodeId(dict["Id"].(bencoding.String))
	local.Port = int(dict["Port"].(bencoding.Int))
	local.OutstandingQueries = make(map[string]*Query)
	local.Nodes = map[string]*RemoteNode{}

	for _, nodeDict := range dict["Nodes"].(bencoding.List) {
		remote := RemoteNodeFromBencodingDict(nodeDict.(bencoding.Dict))
		local.AddOrGetRemoteNode(remote)
	}

	return local
}

func (local *LocalNode) MarshalBencodingDict() (dict bencoding.Dict) {
	dict = bencoding.Dict{}

	if local.Id != UnknownNodeId {
		dict["Id"] = bencoding.String(local.Id)
	}

	dict["Port"] = bencoding.Int(local.Port)

	nodes := make(bencoding.List, len(local.Nodes))

	i := 0
	for _, node := range local.Nodes {
		nodes[i] = node.MarshalBencodingDict()
		i++
	}

	dict["Nodes"] = nodes

	return dict
}

func (local *LocalNode) WriteBencodedTo(writer io.Writer) error {
	dict := local.MarshalBencodingDict()

	logger.Printf("Prepared for serialization: %v\n", dict)

	return dict.WriteBencodedTo(writer)
}

func (local *LocalNode) ToJsonable() (interface{}, error) {
	return bencoding.Dict(local.MarshalBencodingDict()).ToJsonable()
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
