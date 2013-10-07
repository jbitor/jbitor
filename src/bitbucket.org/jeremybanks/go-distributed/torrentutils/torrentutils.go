package torrentutils

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"bytes"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type pieceHasher struct {
	PieceLength int64
	PieceHashes [][]byte
	PieceBuffer []byte
}

func newPieceHasher(pieceLength int64) pieceHasher {
	return pieceHasher{
		PieceLength: pieceLength,
		PieceHashes: [][]byte{},
		PieceBuffer: []byte{},
	}
}

type pieceHasherWrapperWriter struct {
	hasher *pieceHasher
}

func (wrapper pieceHasherWrapperWriter) Write(data []byte) (int, error) {
	return wrapper.Write(data)
}

func (self *pieceHasher) Writer() io.Writer {
	return pieceHasherWrapperWriter{hasher: self}
}

func (self *pieceHasher) Write(data []byte) (int, error) {
	written := 0
	self.PieceBuffer = bytes.Join([][]byte{self.PieceBuffer, data}, []byte(""))

	for int64(len(self.PieceBuffer)) >= self.PieceLength {
		piece := self.PieceBuffer[:self.PieceLength]

		hasher := sha1.New()
		hasher.Write(piece)
		written += len(piece)
		self.PieceHashes = append(self.PieceHashes, hasher.Sum(nil))

		self.PieceBuffer = self.PieceBuffer[self.PieceLength:]
	}

	return written, nil
}

func (self *pieceHasher) Pieces() []byte {
	piecesData := bytes.Join(self.PieceHashes, []byte(""))
	if len(self.PieceBuffer) > 0 {
		hasher := sha1.New()
		hasher.Write(self.PieceBuffer)
		piecesData = bytes.Join([][]byte{piecesData, hasher.Sum(nil)}, []byte(""))
	}
	return piecesData
}

type CreationOptions struct {
	Path           string
	PieceLength    int64
	ForceMultiFile bool
}

func GenerateTorrentMetaInfo(options CreationOptions) (bencoding.Dict, error) {
	fileInfo, err := os.Stat(options.Path)
	if err != nil {
		return nil, err
	}

	multiFile := fileInfo.IsDir() || options.ForceMultiFile

	pieces := make([]byte, 0)

	fileList := bencoding.List{}

	pieceHasher := newPieceHasher(options.PieceLength)

	if multiFile {
		err := filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(options.Path, path)
			if err != nil {
				return err
			}

			pathList := bencoding.List{}

			for _, component := range filepath.SplitList(relPath) {
				pathList = append(pathList, bencoding.String(component))
			}

			fileDict := bencoding.Dict{
				"path":   pathList,
				"length": bencoding.Int(info.Size()),
			}

			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			pieceHasher.Write(data)

			fileList = append(fileList, fileDict)

			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		if fileInfo.Size() > 0 {
			file, err := os.Open(options.Path)
			defer file.Close()

			if err != nil {
				return nil, err
			}

			io.Copy(pieceHasher.Writer(), file)

			file.Close()
		}
	}

	pieces = pieceHasher.Pieces()

	infoDict := bencoding.Dict{
		"name":         bencoding.String(fileInfo.Name()),
		"piece length": bencoding.Int(options.PieceLength),
		"pieces":       bencoding.String(pieces),
	}

	if multiFile {
		infoDict["files"] = fileList
	} else {
		infoDict["length"] = bencoding.Int(fileInfo.Size())
	}

	if err != nil {
		panic(err)
	}

	return infoDict, nil
}
