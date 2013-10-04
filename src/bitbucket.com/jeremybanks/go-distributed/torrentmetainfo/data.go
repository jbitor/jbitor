package torrentmetainfo

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func (metainfo T) Hash() (hash []byte, err error) {
	hasher := sha1.New()
	var bencoded []byte
	bencoded, err = metainfo.MarshalBencoding()
	if err == nil {
		hasher.Write(bencoded)
		hash = hasher.Sum(nil)
	}
	return
}

func (metainfo T) HexHash() (hexhash string, err error) {
	var hash []byte
	hash, err = metainfo.Hash()
	if err == nil {
		hexhash = hex.EncodeToString(hash)
	}
	return
}

func (metainfo T) String() string {
	hexhash, _ := metainfo.HexHash()
	return fmt.Sprintf("<torrentmetainfo.T with .HexHash() = %v>", hexhash)
}
