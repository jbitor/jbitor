package torrent

import (
	"github.com/jeremybanks/go-distributed/bencoding"
	"net"
)

// For now, this is just a stub holding data from the DHT.

type RemotePeer struct {
	Address net.TCPAddr
}

type LocalPeer struct {
	Port int
}

func DecodePeerAddress(encoded bencoding.String) (addr net.TCPAddr) {
	return net.TCPAddr{
		IP:   net.IPv4(encoded[0], encoded[1], encoded[2], encoded[3]),
		Port: int(encoded[4])<<8 + int(encoded[5]),
	}
}
