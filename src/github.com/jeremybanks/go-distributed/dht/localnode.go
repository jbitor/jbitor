package dht

import (
	"fmt"
	"github.com/jeremybanks/go-distributed/bencoding"
	"github.com/jeremybanks/go-distributed/torrent"
	"io"
	weakrand "math/rand"
	"net"
	"sort"
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

type nodeOrderingByNearness struct {
	target    torrent.BTID
	nodes     []*RemoteNode
	distances [][5]uint32
}

func (local *localNode) nodeOrderingByNearnessFromTarget(target torrent.BTID) (ordering nodeOrderingByNearness) {
	ordering.target = target
	ordering.nodes = make([]*RemoteNode, len(local.Nodes))
	ordering.distances = make([][5]uint32, len(local.Nodes))

	i := 0
	for _, remoteNode := range local.Nodes {
		ordering.nodes[i] = remoteNode
		ordering.distances[i] = local.Id.XoredUint32s(remoteNode.Id)
		i++
	}

	return ordering
}

func (ordering nodeOrderingByNearness) Len() int {
	return len(ordering.nodes)
}

func (ordering nodeOrderingByNearness) Swap(i, j int) {
	tmpNode, tmpDistance := ordering.nodes[i], ordering.distances[i]
	ordering.nodes[i] = ordering.nodes[j]
	ordering.distances[i] = ordering.distances[j]
	ordering.nodes[j] = tmpNode
	ordering.distances[j] = tmpDistance
}

func (ordering nodeOrderingByNearness) Less(i, j int) bool {
	iDist, jDist := ordering.distances[i], ordering.distances[j]
	return (iDist[0] < jDist[0] || iDist[0] == jDist[0] &&
		(iDist[1] < jDist[i] || iDist[1] == jDist[1] &&
			(iDist[2] < jDist[2] || iDist[2] == jDist[2] &&
				(iDist[3] < jDist[3]))))

}

func (local *localNode) NodesByNearness(target torrent.BTID) []*RemoteNode {
	ordering := local.nodeOrderingByNearnessFromTarget(target)
	sort.Sort(ordering)
	return ordering.nodes[:]
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
	var conn *net.UDPConn
	// Main loop for LocalPeer's activity.
	// (Listening to replies and requests.)

	conn, err = net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: local.Port,
	})
	if err != nil {
		return
	}

	local.Connection = conn

	rpcTerminate := make(chan bool)
	go func() {
		local.rpcListenLoop(rpcTerminate)
	}()

	go func() {
		<-terminate

		close(rpcTerminate)
		conn.Close()

	}()

	return
}
