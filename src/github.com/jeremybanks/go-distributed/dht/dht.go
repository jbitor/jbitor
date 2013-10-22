// Package dht provides an implementation of a "mainline" BitTorrent
// Distributed Hash Table (DHT) node, as specified in BEP 5 (http://www.bittorrent.org/beps/bep_0005.html),
// and a higher-level client interface for querying the DHT.
package dht

import "net"

var bootstrapNodes = []RemoteNode{
	{
		Address: net.UDPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 6881,
		},
	},
	// XXX: These should have some flag indicating that they're only for
	// bootstrapping, to ensure they aren't used for any other purpose
	// and when there are no other known good nodes.
	{
		Address: net.UDPAddr{
			IP:   net.IPv4(67, 215, 242, 139),
			Port: 6881,
		},
	},
}
