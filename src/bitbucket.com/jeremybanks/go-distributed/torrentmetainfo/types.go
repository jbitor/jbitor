package torrentmetainfo

type T struct {
	Name        string
	Pieces      string
	PieceLength int64
	Length      int64
	Files       *[]File // nil if single-file torrent
}

type File struct {
	Offset int64
	Length int64
	Path   []string
}
