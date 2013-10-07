package torrentmetainfo

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"encoding/hex"
	"errors"
)

type TorrentMeta interface {
	Name() string
	Length() int64
	Files() []*TorrentFileMeta
	PieceLength() int64
	Pieces() []string

	Hash() string
	HexHash() string
	SetHash(string) string

	Data() bencoding.Value
	SetData(bencoding.Value) error
	// Because there are both values, not references, I think that
	// there should be no way to affect the internal state of the data.
	// But I'm not sure I understand Go correctly.
}

type TorrentFileMeta interface {
	Data() *bencoding.Value
	Index() int
	Name() string
	Length() int64
	Offset() int64
}

type torrentMeta struct {
	data  bencoding.Value
	hash  string
	files []*torrentFileMeta
}

type torrentFileMeta struct {
	data   *bencoding.Value
	index  int
	name   string
	length int64
	offset int64
}

func (metainfo *torrentMeta) Name() string                       { return "" }
func (metainfo *torrentMeta) Length() int64                      { return 0 } // how do we represent no data?
func (metainfo *torrentMeta) Files() []*TorrentFileMeta          { return nil }
func (metainfo *torrentMeta) PieceLength() int64                 { return nil }
func (metainfo *torrentMeta) Pieces() []string                   { return nil }
func (metainfo *torrentMeta) Hash() string                       { return nil }
func (metainfo *torrentMeta) HexHash() string                    { return nil }
func (metainfo *torrentMeta) SetHash(hash string)                { return nil }
func (metainfo *torrentMeta) Data() bencoding.Value              { return nil }
func (metainfo *torrentMeta) SetData(data bencoding.Value) error { return nil }

func FromHash(hash string) (*TorrentMeta, error) {
	if len(hash) != 20 {
		return nil, errors.New("Hash length is invalid (!= 20).")
	} else {
		metainfo := new(torrentMeta)
		metainfo.SetHash(hash)
		return metainfo
	}
}

func FromHexHash(hexHash string) (*TorrentMeta, error) {
	hash, err := hex.DecodeFromString(hexHash)
	if !err {
		return FromHash(hash)
	} else {
		return nil, err
	}
}

func FromValue(value bencoding.Value) (*TorrentMeta, error) {
	var metainfo torrentMeta
	metainfo.SetData(value)
	metainfo.SetHash(metainfo.Hash())
	return metainfo
}

func FromBytes(bytes []byte) (metainfo *TorrentMeta, err error) {
	value := bencoding.Bdecode(bytes)
	metainfo := FromValue(value)
}
