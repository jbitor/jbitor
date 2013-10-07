package torrentutils

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"bytes"
	"crypto/sha1"
	"io"
	"os"
)

type CreationOptions struct {
	Path           string
	PieceLength    int64
	ForceMultiFile bool
}

func GenerateTorrentMetaInfo(options CreationOptions) (bencoding.Dict, error) {
	file, err := os.Open(options.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	multiFile := fileInfo.IsDir() || options.ForceMultiFile

	if multiFile {
		panic("multiFile not implemented yet")
	}

	pieces := make([][]byte, 0)

	if fileInfo.Size() == 0 {
		hasher := sha1.New()
		pieces = append(pieces, hasher.Sum(nil))
	} else {
		offset := int64(0)

		for {
			pieceData := make([]byte, options.PieceLength)
			pieceSize, err := file.Read(pieceData)
			pieceData = pieceData[:pieceSize]

			if err != nil && err != io.EOF {
				return nil, err
			}
			if pieceSize == 0 {
				if offset != fileInfo.Size() {
					panic("Unexpected byte shortage")
				}
				break
			}

			hasher := sha1.New()
			hasher.Write(pieceData)
			pieces = append(pieces, hasher.Sum(nil))
			offset += int64(pieceSize)

			if err == io.EOF || int64(pieceSize) < options.PieceLength {
				if offset != fileInfo.Size() {
					panic("Unexpected byte shortage")
				}
				break
			}
		}
	}

	infoDict := bencoding.Dict{
		"name":         bencoding.String(fileInfo.Name()),
		"length":       bencoding.Int(fileInfo.Size()),
		"piece length": bencoding.Int(options.PieceLength),
		"pieces":       bencoding.String(bytes.Join(pieces, []byte{})),
	}

	if err != nil {
		panic(err)
	}

	return infoDict, nil
}
