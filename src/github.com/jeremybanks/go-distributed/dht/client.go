package dht

/*
import (
	"github.com/jeremybanks/go-distributed/torrent"
	"errors"
)

type Client interface {
	IsOpen() bool
	Close() (err error)

	GoodNodes() int
	// The number of "good" nodes known to the client.

	Connect(goodNodes int) (err error)
	// Blocks until the Client knows of `goodNodes` good nodes, or the client
	// is closed.

	GetPeers(infoHash []byte) (peers []*torrent.RemotePeer, err error)
	AnnouncePeer(infoHash []byte, local *torrent.LocalPeer) (err error)
}

type localNodeClient struct {
	openNode LocalNode
	closeChan chan<- bool
	errChan <-chan error
}

func OpenClient(path string) (client *Client, err error) {
	// TODO: use an exclusive open to lock path + ".lock"
	client = new(localNodeClient)

	closeChan := make(chan bool)

	client.closeChan = closeChan

	go client.node.Run()

	return client, nil
}

func (client *localNodeClient) IsOpen() bool {
	return client.openNode != nil
}

func (client *localNodeClient) Close() (err error) {
	if client.openNode == nil {
		return errors.New("dht.Client is not open.")
	}

	client.closeChan <- true
	client.openNode = nil
}

func (client *localNodeClient) run(closeChan <-chan bool) {
	for {
		select {
		case _ = <-closeChan:
			break
		default:
		}


	}


}
*/
