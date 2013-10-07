package torrentmetainfo

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"errors"
)

func (metainfo T) UnmarshalBencodingValue(bval *bencoding.Value) (err error) {
	val, ok := bval.Value.(map[string]*bencoding.Value)
	if !ok {
		return errors.New("Root not a dictionary")
	}

	metainfo.Name, ok = val["name"].Value.(string)
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

func (metainfo T) UnmarshalTorrentBencoding(encoded []byte) (err error) {
	// Loads metadata from .torrent file data.
	//
	// This is different from UnmarshalBencoding because that function
	// only takes the metainfo section of the torrent data, whereas this
	// takes the entire torrent file data.
	var torrentBval *bencoding.Value
	torrentBval, err = bencoding.Bdecode(encoded)
	if err == nil {
		bval, ok := torrentBval.Value.(map[string]interface{})["info"]

		if ok {
			err = metainfo.UnmarshalBencodingValue(bval.(*bencoding.Value))
		} else {
			err = errors.New("Couldn't find 'info' key in torrent data.")
		}
	}
	return
}
