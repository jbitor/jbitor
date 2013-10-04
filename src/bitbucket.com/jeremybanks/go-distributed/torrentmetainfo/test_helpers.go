package torrentmetainfo

func makeTestTorrentWithoutFiles() (metainfo T) {
	metainfo.name = "hello.txt"
	metainfo.pieces = "12345678901234567890"
	metainfo.piece_length = 1024
	metainfo.length = 20
	metainfo.files = nil
	return
}

func makeTestTorrentWithFiles() (metainfo T) {
	panic("not implemented")
}
