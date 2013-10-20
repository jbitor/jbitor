package dht

import (
	"errors"
	"github.com/jeremybanks/go-distributed/torrent"
	"time"
)

// Any DHT queries sent to another node will time out after this long.
const QUERY_TIMEOUT = 10 * time.Second

// Information about how a client/node is connected to the DHT network.
type ConnectionInfo struct {
	GoodNodes    int
	UnknownNodes int
	BadNodes     int
}

/*
Client provides a high-level interface for interacting with the DHT.
Once a client has been opened it will continue to run asynchronously until it is closed.
*/
type Client interface {
	// Close release the client's port and unlocks its data file, leaving it unusable.
	Close() (err error)

	// Save saves the current state of the DHT to its data file.
	Save() (err error)

	// GetPeers attempts to find remote torrent peers downloading a given torrent.
	GetPeers(infoHash torrent.BTID) (peers []*torrent.RemotePeer, err error)

	// AnnouncePeer announces to the DHT that a local torrent peer is downloading a given torrent.
	AnnouncePeer(local *torrent.LocalPeer, infoHash torrent.BTID) (err error)

	ConnectionInfo() ConnectionInfo
}

/*
localNodeClient implements the Client interface for this package.
*/
type localNodeClient struct {
	// The localNode instance being used by this client.
	// This will be nil if the client is not open.
	openNode *localNode

	// Sending any value into this chanel will terminate the client's activity.
	terminate chan<- bool

	// A single value will be sent into this chanel when the client is
	// terminated. If the value is not nil, it will be a fatal error that
	// caused the client to be terminated.
	terminated <-chan error
}

/*
OpenClient instantiates a client whose state will be persisted at the specified path.

Existing state will be loaded if it exists, otherwise a new client will
be generated using a node a randomly-selected ID and port.

A filesystem lock will be used to ensure that only one Client may be open with
a given path at a time.
*/
func OpenClient(path string) (c Client, err error) {
	// TODO: use an exclusive open to lock path + ".lock"
	local := new(localNodeClient)

	terminate := make(chan bool)

	local.terminate = terminate

	go func() {

	}()

	c = Client(local)
	return
}

func (c *localNodeClient) Close() (err error) {
	if c.openNode == nil {
		return errors.New("dht.Client is not open.")
	}

	close(c.terminate)
	c.openNode = nil
	return
}

func (c *localNodeClient) Save() (err error) {
	panic("Save not implemented")
}

func (c *localNodeClient) GetPeers(infoHash torrent.BTID) (peers []*torrent.RemotePeer, err error) {
	panic("GetPeers not implemented")
}

func (c *localNodeClient) AnnouncePeer(local *torrent.LocalPeer, infoHash torrent.BTID) (err error) {
	panic("AnnouncePeer not implemented")
}

func (c *localNodeClient) ConnectionInfo() ConnectionInfo {
	info := ConnectionInfo{GoodNodes: 0, UnknownNodes: 0, BadNodes: 0}

	for _, node := range c.openNode.Nodes {
		switch node.Status() {
		case STATUS_UNKNOWN:
			info.UnknownNodes++
		case STATUS_GOOD:
			info.GoodNodes++
		case STATUS_BAD:
			info.BadNodes++
		}
	}

	return info
}
