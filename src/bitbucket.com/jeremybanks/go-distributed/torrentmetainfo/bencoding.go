package torrentmetainfo

import (
	"bitbucket.com/jeremybanks/go-distributed/bencoding"
)

func (metainfo T) MarshalBencodingValue() (bval *bencoding.Value, err error) {
	val := map[string]interface{}{
		"name":         metainfo.name,
		"pieces":       metainfo.pieces,
		"piece length": metainfo.piece_length,
	}

	if metainfo.files == nil {
		val["length"] = metainfo.length
	} else {
		files := make([]interface{}, len(*metainfo.files))

		for i, file := range *metainfo.files {
			files[i] = map[string]interface{}{
				"length": file.length,
				"path":   file.path,
			}
		}

		val["files"] = files
	}

	bval, err = bencoding.NewValue(val)
	return
}

func (metainfo T) MarshalBencoding() (encoded []byte, err error) {
	return bencoding.Bencode(metainfo)
}
