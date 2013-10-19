package dht

// Implements a BitTorrent DHT Node (peer), as described in
// http://www.bittorrent.org/beps/bep_0005.html.

import (
	"crypto/rand"
	"encoding/hex"
	weakrand "math/rand"
)

type NodeId string // of length 20

var UnknownNodeId NodeId = NodeId("") // empty string is unkown node ID

func GenerateNodeId() NodeId {
	// Securely generates a random NodeID.
	// Panics if unable to generate secure random number.

	bytes := new([20]byte)
	n, err := rand.Read(bytes[:])

	if n < 20 {
		panic("unable to generate 20 secure random bytes for node ID")
	}

	if err != nil {
		panic(err)
	}

	return NodeId(bytes[:])
}

func GenerateWeakNodeId() NodeId {
	// Generates a random NodeID using the weak/fast RNG.

	bytes := new([20]byte)

	for i := range bytes {
		bytes[i] = byte(weakrand.Int() & 0xFF)
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
