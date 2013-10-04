package torrentmetainfo

import (
	"bitbucket.com/jeremybanks/go-distributed/bencoding"
	"errors"
)

func (metainfo T) UnmarshalBencodingValue(bval *bencoding.Value) (err error) {
	val, ok := bval.Value.(map[string]*bencoding.Value)
	if !ok {
		return errors.New("Root not a dictionary")
	}

	metainfo.name, ok = val["name"].Value.(string)
	if !ok {
		return errors.New("name not a string")
	}

	_, has_files := val["files"]

	_ = has_files

	panic("UnmarshalBencodingValue not implemented")
}

func (metainfo T) UnmarshalBencoding(encoded []byte) (err error) {
	var bval *bencoding.Value
	bval, err = bencoding.Bdecode(encoded)
	if err == nil {
		err = metainfo.UnmarshalBencodingValue(bval)
	}
	return
}
