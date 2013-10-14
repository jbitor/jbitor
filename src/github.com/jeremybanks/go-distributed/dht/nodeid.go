package dht

// Implements a BitTorrent DHT Node (peer), as described in
// http://www.bittorrent.org/beps/bep_0005.html.

import (
	"crypto/rand"
	"encoding/hex"
)

type NodeId string // of length 20

var UnknownNodeId NodeId = NodeId("") // empty string is unkown node ID

func GenerateNodeId() NodeId {
	// Generates a new random NodeID.
	// Panics if unable to generate random number.

	bytes := new([20]byte)
	n, err := rand.Read(bytes[:])

	if n < 20 {
		panic("too few bytes generated for some reason?")
	}

	if err != nil {
		panic(err)
	}

	return NodeId(bytes[:])
}

func (id NodeId) String() string {
	if len(id) > 0 {
		return hex.EncodeToString([]byte(id))
	} else {
		return "[unknown ID]"
	}
}
