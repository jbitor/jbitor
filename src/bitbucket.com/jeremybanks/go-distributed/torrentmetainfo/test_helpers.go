package torrentmetainfo

func makeTestTorrentWithoutFiles() (metainfo T) {
	metainfo.Name = "data.go"
	metainfo.Pieces = "\x8E\x66\xF6\xC8\x3A\xF4\x23\x52\xEF\xDF\x5A\x2C\x1D\x02\x16\x21\x22\xD7\x63\x94"
	metainfo.PieceLength = 32768
	metainfo.Length = 613
	metainfo.Files = nil
	return
}

func makeTestTorrentWithFiles() (metainfo T) {
	metainfo.Name = "test"
	metainfo.Length = 2629
	metainfo.Pieces = "\x0E\x35\xD4\x04\x61\xFB\x99\x77\x46\x5E\xAB\xB6\xA0\x9A\xC7\x84\x48\xF9\x69\x98"
	metainfo.PieceLength = 32768
	metainfo.Length = 3063
	metainfo.Files = &[]File{
		File{Offset: 0, Length: 2629, Path: []string{"data.go"}},
		File{Offset: 2629, Length: 434, Path: []string{"types.go"}},
	}
	return
}
