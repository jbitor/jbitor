package dht

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"fmt"
	"net"
	"time"
)

type RemoteNode struct {
	Id      NodeId // may be unknown
	Address net.UDPAddr

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

func RemoteNodeFromBencodingDict(dict bencoding.Dict) (remote *RemoteNode) {
	panic("not implemented")
}

func (remote *RemoteNode) MarshalBencodingDict() (dict bencoding.Dict) {
	dict = bencoding.Dict{
		"id": bencoding.String(remote.Id),

		// XXX: Ideally this would use proper compact representation, but for
		// now maybe just use [[ address, port ]...]
		"address": bencoding.List{
			bencoding.String(remote.Address.IP.String()),
			bencoding.Int(remote.Address.Port),
		},

		"lastRequestToSec":    bencoding.Int(remote.LastRequestTo.Unix()),
		"lastResponseToSec":   bencoding.Int(remote.LastResponseTo.Unix()),
		"lastRequestFromSec":  bencoding.Int(remote.LastRequestFrom.Unix()),
		"lastResponseFromSec": bencoding.Int(remote.LastResponseFrom.Unix()),
	}
	/*

		LastRequestTo    time.Time
		LastResponseTo   time.Time
		LastResponseFrom time.Time
		LastRequestFrom  time.Time

		ConsecutiveFailedQueries int
	*/

	return dict
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
