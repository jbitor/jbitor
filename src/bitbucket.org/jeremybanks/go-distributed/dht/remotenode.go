package dht

import (
	"fmt"
	"net"
	"time"
)

type RemoteNode struct {
	Id      NodeId
	Address net.UDPAddr

	OutstandingQueries map[string]Query

	LastRequestTo    time.Time
	LastResponseTo   time.Time
	LastResponseFrom time.Time
	LastRequestFrom  time.Time

	ConsecutiveFailedQueries int
}

func RemoteNodeFromAddress(address net.UDPAddr) (remote *RemoteNode) {
	// Creates a RemoteNode with a known address but an unknown ID.
	// You may want to .Ping() this node so that it learns its ID!
	remote = new(RemoteNode)
	remote.Id = UnknownNodeId
	remote.Address = address
	remote.ConsecutiveFailedQueries = 0
	return remote
}

func GenerateFakeRemoteNode() (remote *RemoteNode) {
	remote = new(RemoteNode)
	remote.Id = GenerateNodeId()
	remote.Address = net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1234}
	return remote
}

func (remote *RemoteNode) String() string {
	return fmt.Sprintf("<RemoteNode %v (%v) at %v>", remote.Id, remote.Status(), remote.Address)
}

func (remote *RemoteNode) Status() RemoteNodeStatus {
	switch {
	case remote.ConsecutiveFailedQueries >= 3:
		return STATUS_BAD
	case remote.LastResponseFrom.IsZero():
		return STATUS_UNKNOWN
	case time.Since(remote.LastResponseFrom) < time.Minute*15:
		return STATUS_GOOD
	case time.Since(remote.LastRequestFrom) < time.Minute*15:
		return STATUS_GOOD
	default:
		return STATUS_UNKNOWN
	}
}

type RemoteNodeStatus int

const (
	STATUS_UNKNOWN RemoteNodeStatus = iota
	STATUS_GOOD
	STATUS_BAD
)

func (status RemoteNodeStatus) String() string {
	switch status {
	case STATUS_UNKNOWN:
		return "STATUS_UNKNOWN"
	case STATUS_GOOD:
		return "STATUS_GOOD"
	case STATUS_BAD:
		return "STATUS_BAD"
	default:
		panic("invalid status value")
	}
}
