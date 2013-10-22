package dht

import (
	"fmt"
	"github.com/jeremybanks/go-distributed/bencoding"
	"github.com/jeremybanks/go-distributed/torrent"
	"net"
	"time"
)

type RemoteNode struct {
	Id      torrent.BTID // may be unknown
	Address net.UDPAddr

	LastRequestTo    time.Time
	LastResponseTo   time.Time
	LastResponseFrom time.Time
	LastRequestFrom  time.Time

	ConsecutiveFailedQueries int

	Source *RemoteNode
	// If specified, identifies the RemoteNode which gave us this Node.
	// This could be used to identify a Source which provides us with many
	// bad nodes.
}

func RemoteNodeFromAddress(address net.UDPAddr) (remote *RemoteNode) {
	// Creates a RemoteNode with a known address but an unknown ID.
	// You may want to .SendPing() this node so that it learns its ID!
	remote = new(RemoteNode)
	remote.Id = ""
	remote.Address = address

	// default values are fine for other fields

	return remote
}

func RemoteNodeFromBencodingDict(dict bencoding.Dict) (remote *RemoteNode) {
	remote = new(RemoteNode)

	remote.Id = torrent.BTID(dict["Id"].(bencoding.String))
	remote.Address = decodeNodeAddress(dict["Address"].(bencoding.String))
	remote.LastRequestTo = time.Unix(int64(dict["LastRequestToSec"].(bencoding.Int)), 0)
	remote.LastResponseTo = time.Unix(int64(dict["LastResponseToSec"].(bencoding.Int)), 0)
	remote.LastRequestFrom = time.Unix(int64(dict["LastRequestFromSec"].(bencoding.Int)), 0)
	remote.LastResponseFrom = time.Unix(int64(dict["LastResponseFromSec"].(bencoding.Int)), 0)
	remote.ConsecutiveFailedQueries = int(dict["ConsecutiveFailedQueries"].(bencoding.Int))

	return remote
}

func (remote *RemoteNode) MarshalBencodingDict() (dict bencoding.Dict) {
	dict = bencoding.Dict{
		"Id": bencoding.String(remote.Id),

		"Address": encodeNodeAddress(remote.Address),

		"LastRequestToSec":    bencoding.Int(remote.LastRequestTo.Unix()),
		"LastResponseToSec":   bencoding.Int(remote.LastResponseTo.Unix()),
		"LastRequestFromSec":  bencoding.Int(remote.LastRequestFrom.Unix()),
		"LastResponseFromSec": bencoding.Int(remote.LastResponseFrom.Unix()),

		"ConsecutiveFailedQueries": bencoding.Int(remote.ConsecutiveFailedQueries),
	}

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