package dht

import (
	"errors"
	"fmt"
	"github.com/jeremybanks/go-distributed/bencoding"
	"github.com/jeremybanks/go-distributed/torrent"
	"io"
	weakrand "math/rand"
	"net"
	"time"
)

/*
localNodes is a DHT node implementation. Currently, it only supports the
client components of a node -- it does not maintain a proper routing table
and cannot respond to queries.
*/
type localNode struct {
	Id                 torrent.BTID
	Port               int
	Connection         *net.UDPConn
	Nodes              map[string]*RemoteNode
	OutstandingQueries map[string]*outstandingQuery
}

func newLocalNode() (local *localNode) {
	id, err := torrent.SecureRandomBTID()
	if err != nil {
		// You used up all the entropy!
		panic(err)
	}

	local = new(localNode)
	local.Id = id
	local.Port = 1024 + weakrand.Intn(8192)
	local.OutstandingQueries = make(map[string]*outstandingQuery)
	local.Nodes = map[string]*RemoteNode{}

	for _, node := range bootstrapNodes {
		local.AddOrGetRemoteNode(&node)
	}

	return local
}

func (local *localNode) AddOrGetRemoteNode(remote *RemoteNode) *RemoteNode {
	// If a node with the same address is already in .Nodes, returns that node.
	// Otherwise, add remote to .Nodes and return it.

	key := RemoteNodeKey(remote.Address)

	if existingRemote, ok := local.Nodes[key]; ok {
		remote = existingRemote
	} else {
		local.Nodes[key] = remote
	}

	return remote

}

func (local *localNode) String() string {
	return fmt.Sprintf("<localNode %v on :%v>", local.Id, local.Port)
}

func RemoteNodeKey(addr net.UDPAddr) string {
	return fmt.Sprintf("%v:%v", addr.IP, addr.Port)
}

// Bencoding

func localNodeFromBencodingDict(dict bencoding.Dict) (local *localNode) {
	local = new(localNode)

	local.Id = torrent.BTID(dict["Id"].(bencoding.String))
	local.Port = int(dict["Port"].(bencoding.Int))
	local.OutstandingQueries = make(map[string]*outstandingQuery)
	local.Nodes = map[string]*RemoteNode{}

	for _, nodeDict := range dict["Nodes"].(bencoding.List) {
		remote := RemoteNodeFromBencodingDict(nodeDict.(bencoding.Dict))
		local.AddOrGetRemoteNode(remote)
	}

	return local
}

func (local *localNode) MarshalBencodingDict() (dict bencoding.Dict) {
	dict = bencoding.Dict{}

	if local.Id != "" {
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

func (local *localNode) WriteBencodedTo(writer io.Writer) error {
	dict := local.MarshalBencodingDict()

	return dict.WriteBencodedTo(writer)
}

func (local *localNode) ToJsonable() (interface{}, error) {
	return bencoding.Dict(local.MarshalBencodingDict()).ToJsonable()
}

// Running

func (local *localNode) Run(terminate <-chan bool) (err error) {
	// Main loop for LocalPeer's activity.
	// (Listening to replies and requests.)

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: local.Port,
	})
	defer conn.Close()
	if err != nil {
		return
	}

	local.Connection = conn

	rpcTerminate := make(chan bool)
	go func() { local.runRpcListen(rpcTerminate) }()

	connectionTerminate := make(chan bool)
	go func() { local.runConnection(connectionTerminate) }()

	select {
	case _ = <-terminate:
		// terminate sub-goroutines
		rpcTerminate <- true
		close(rpcTerminate)
	}
	return
}

func (local *localNode) runConnection(terminate <-chan bool) {
	for {
		local.pingRandomNode()
		local.requestMoreNodes()

		info := (&localNodeClient{openNode: local}).ConnectionInfo()

		logger.Printf("localNode running with %v good nodes (%v unknown and %v bad).\n",
			info.GoodNodes, info.UnknownNodes, info.BadNodes)

		time.Sleep(15 * time.Second)

		select {
		case _ = <-terminate:
			break
		default:
		}
	}
}

/*
DHT: localNode running with 16 good remote nodes (123 unknown and 0 bad).
*/

func (local *localNode) pingRandomNode() {
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
		close(timeoutChan)
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

// Make this a client method, and add a saving loop to it.
func (local *localNode) requestMoreNodes() {
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

	target := torrent.WeakRandomBTID()

	logger.Printf("Requesting new nodes around %v from %v.\n", target, randNode)

	resultChan, errChan := local.FindNode(randNode, target)

	timeoutChan := make(chan error)
	go func() {
		time.Sleep(10 * time.Second)
		timeoutChan <- errors.New("find nodes timed out")
		close(timeoutChan)
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
