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

func RemoteNodeFromBencodingDict(dict bencoding.Dict) (remote *RemoteNode) {
	remote = new(RemoteNode)

	remote.Id = NodeId(dict["Id"].(bencoding.String))
	remote.Address = net.UDPAddr{
		IP:   net.ParseIP(string(dict["Address"].(bencoding.List)[0].(bencoding.String))),
		Port: int(dict["Address"].(bencoding.List)[1].(bencoding.Int)),
	}
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

		// XXX: Ideally this would use proper compact representation, but for
		// now maybe just use [[ address, port ]...]
		"Address": bencoding.List{
			bencoding.String(remote.Address.IP.String()),
			bencoding.Int(remote.Address.Port),
		},

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
