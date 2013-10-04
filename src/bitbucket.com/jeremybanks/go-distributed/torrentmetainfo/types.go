package torrentmetainfo

type T struct {
	name         string
	pieces       string
	piece_length int64
	length       int64
	files        *[]File // nil if single-file torrent
}

type File struct {
	offset int64
	length int64
	path   []string
}
