package dht

import (
	"fmt"
	"net"
)

type RemoteNode struct {
	Id      NodeId
	Address net.UDPAddr
}

func RemoteNodeFromAddress(address net.UDPAddr) (remote *RemoteNode) {
	// Creates a RemoteNode with a known address but an unknown ID.
	// You may want to .Ping() this node so that it learns its ID!
	remote = new(RemoteNode)
	remote.Id = UnknownNodeId
	remote.Address = address
	return remote
}

func GenerateFakeRemoteNode() (remote *RemoteNode) {
	remote = new(RemoteNode)
	remote.Id = GenerateNodeId()
	remote.Address = net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}
	return remote
}

func (remote *RemoteNode) String() string {
	return fmt.Sprintf("<RemoteNode %v at %v>", remote.Id, remote.Address)
}
