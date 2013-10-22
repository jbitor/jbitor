package dht

import (
	"errors"
	"github.com/jeremybanks/go-distributed/bencoding"
	"github.com/jeremybanks/go-distributed/torrent"
	"io/ioutil"
	"os"
	"syscall"
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

	terminateOpenNode chan<- bool

	// The data file we read/write the node state from/to.
	openDataFile *os.File

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
	var (
		openDataFile   *os.File
		nodeData       []byte
		nodeDict       bencoding.Bencodable
		nodeDictAsDict bencoding.Dict
		ok             bool
		local          *localNodeClient
	)

	local = new(localNodeClient)

	openDataFile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	local.openDataFile = openDataFile

	err = syscall.Flock(int(openDataFile.Fd()), syscall.LOCK_EX) // block for exclusive lock of file
	if err != nil {
		return
	}

	nodeData, err = ioutil.ReadAll(local.openDataFile)
	if err != nil {
		logger.Printf("Unable to read existing DHT node file (%v). Creating a new one.\n", err)
		local.openNode = newLocalNode()
	} else if len(nodeData) == 0 {
		logger.Printf("Existing DHT node file was empty. Creating a new one.\n")
		local.openNode = newLocalNode()
	} else {
		nodeDict, err = bencoding.Decode(nodeData)
		if err != nil {
			openDataFile.Close()
			return
		}

		nodeDictAsDict, ok = nodeDict.(bencoding.Dict)
		if !ok {
			err = errors.New("Node data wasn't a dict.")
			logger.Printf("%v\n", err)
			openDataFile.Close()
			return
		}

		local.openNode = localNodeFromBencodingDict(nodeDictAsDict)
		logger.Printf("Loaded local node info from %v.\n", path)
	}

	terminateOpenNode := make(chan bool)
	local.terminateOpenNode = terminateOpenNode

	c = Client(local)

	go local.openNode.Run(terminateOpenNode)

	go func() {
		for local.openNode != nil {
			c.Save()
			time.Sleep(15 * time.Second)
		}
	}()

	return
}

func (c *localNodeClient) Close() (err error) {
	if c.openNode == nil {
		return errors.New("dht.Client is not open.")
	}

	err = c.Save()

	_ = c.openDataFile.Close()
	c.openDataFile = nil

	_ = c.openNode.Connection.Close()
	c.openNode = nil

	return
}

func (c *localNodeClient) Save() (err error) {
	var (
		nodeData []byte
	)

	if c.openNode == nil {
		err = errors.New("Client is closed.")
		return
	}

	nodeData, err = bencoding.Encode(c.openNode)
	if err != nil {
		return
	}

	err = c.openDataFile.Truncate(0)
	if err != nil {
		return
	}

	_, err = c.openDataFile.WriteAt(nodeData, 0)
	if err != nil {
		return
	}

	err = c.openDataFile.Sync()
	if err != nil {
		return
	}

	logger.Printf("Saved DHT client state.\n")

	return
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
