package torrentmetainfo

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
)

func (metainfo T) MarshalBencodingValue() (bval *bencoding.Value, err error) {
	val := map[string]interface{}{
		"name":         metainfo.Name,
		"pieces":       metainfo.Pieces,
		"piece length": metainfo.PieceLength,
	}

	if metainfo.Files == nil {
		val["length"] = metainfo.Length
	} else {
		// Convert from []string to []interface{}:
		files := make([]interface{}, len(*metainfo.Files))

		for i, file := range *metainfo.Files {
			path := make([]interface{}, len(file.Path))

			for index, pathPart := range file.Path {
				path[index] = pathPart
			}

			files[i] = map[string]interface{}{
				"length": file.Length,
				"path":   path,
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
