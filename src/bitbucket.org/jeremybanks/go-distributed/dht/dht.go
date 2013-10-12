package dht

// Implements a BitTorrent DHT Node (peer), as described in
// http://www.bittorrent.org/beps/bep_0005.html.

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

type localNode struct {
	Id [20]byte
}

type remoteNode struct {
	Id [20]byte
}

func NewLocalNode() (*localNode, error) {
	node := new(localNode)

	n, err := rand.Read(node.Id[:20])

	if n < 20 {
		return nil, errors.New("too few bytes generated for some reason?")
	}

	if err != nil {
		return nil, err
	}

	return node, nil
}

func (node *localNode) String() string {
	return fmt.Sprintf("<localNode id=%v>",
		hex.EncodeToString(node.Id[:]))

}
