package torrent

import "net"

// For now, this is just a stub holding data from the DHT.

type RemotePeer struct {
	address net.TCPAddr
}

type LocalPeer struct {
	port int
}
