package dht

import (
	"errors"
	"fmt"
	"github.com/jeremybanks/go-distributed/bencoding"
	"io"
	weakrand "math/rand"
	"net"
	"time"
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

	return dict.WriteBencodedTo(writer)
}

func (local *LocalNode) ToJsonable() (interface{}, error) {
	return bencoding.Dict(local.MarshalBencodingDict()).ToJsonable()
}

// Running

func (local *LocalNode) Run(terminate <-chan bool, terminated chan<- error) {
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

	rpcTerminated := make(chan error)
	rpcTerminate := make(chan bool)
	go local.runRpcListen(rpcTerminate, rpcTerminated)

	connectionTerminated := make(chan error)
	connectionTerminate := make(chan bool)
	go local.runConnection(connectionTerminate, connectionTerminated)

	select {
	case _ = <-terminate:
		// terminate sub-goroutines
		rpcTerminate <- true

		// notify caller of non-error termination
		terminated <- nil

	case err = <-rpcTerminated:
		logger.Printf("Fatal error from RPC goroutine: %v.\n", err)
		terminated <- err

	case err = <-connectionTerminated:
		logger.Printf("Fatal error from connection goroutine: %v.\n", err)
		terminated <- err
	}

}

func (local *LocalNode) runConnection(terminate <-chan bool, terminated chan<- error) {
	for {
		local.pingRandomNode()
		local.requestMoreNodes()

		logger.Printf("%v known nodes.\n", len(local.Nodes))

		time.Sleep(15 * time.Second)

		select {
		case _ = <-terminate:
			terminated <- nil
			break
		default:
		}
	}
}

/*
DHT: LocalNode running with 16 good remote nodes (123 unknown and 0 bad).
*/

func (local *LocalNode) pingRandomNode() {
	var randNode *RemoteNode
	randNodeOffset := weakrand.Intn(len(local.Nodes))
	i := 0

	for _, node := range local.Nodes {
		if i == randNodeOffset {
			randNode = node
			break
		}
		i++
	}

	logger.Printf("Pinging a random node: %v.\n", randNode)

	resultChan, errChan := local.Ping(randNode)

	timeoutChan := make(chan error)
	go func() {
		time.Sleep(10 * time.Second)
		timeoutChan <- errors.New("ping timed out")
	}()

	select {
	case _ = <-resultChan:
		logger.Printf("Successfully pinged %v.\n", randNode)

	case err := <-errChan:
		logger.Printf("Failed to ping %v: %v.\n", randNode, err)

	case err := <-timeoutChan:
		logger.Printf("Failed to ping %v: %v.\n", randNode, err)
	}
}

func (local *LocalNode) requestMoreNodes() {
	var randNode *RemoteNode
	randNodeOffset := weakrand.Intn(len(local.Nodes))
	i := 0

	for _, node := range local.Nodes {
		if i == randNodeOffset {
			randNode = node
			break
		}
		i++
	}

	target := GenerateNodeId()

	logger.Printf("Requesting new nodes around %v from %v.\n", target, randNode)

	resultChan, errChan := local.FindNode(randNode, target)

	timeoutChan := make(chan error)
	go func() {
		time.Sleep(10 * time.Second)
		timeoutChan <- errors.New("find nodes timed out")
	}()

	select {
	case _ = <-resultChan:
		logger.Printf("Successfully find nodes from %v.\n", randNode)

	case err := <-errChan:
		logger.Printf("Failed to find nodes from %v: %v.\n", randNode, err)

	case err := <-timeoutChan:
		logger.Printf("Failed to find nodes from %v: %v.\n", randNode, err)
	}
}
