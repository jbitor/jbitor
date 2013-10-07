package torrentmetainfo

import (
	"fmt"
	"testing"
)

func TestHexHashWithoutFiles(t *testing.T) {
	metainfo := makeTestTorrentWithoutFiles()
	target := "c339b4b17d653a1ed6b2eed15bc6dd19b828ea5c"
	actual, err := metainfo.HexHash()

	if target != actual {
		t.Log(fmt.Sprintf("error was: %v", err))
		t.Error(fmt.Sprintf("Hash was:\n%v\nexpected:\n%v\n.", actual, target))
	}
}

func TestHexHashWithFiles(t *testing.T) {
	metainfo := makeTestTorrentWithFiles()
	target := "0ed4af3bdd5b66ce072ae4b7d31a9369217b0d46"
	actual, err := metainfo.HexHash()

	if target != actual {
		t.Log(fmt.Sprintf("error was: %v", err))
		t.Error(fmt.Sprintf("Hash was:\n%v\nexpected:\n%v\n.", actual, target))
	}
}

func TestStringWithoutFiles(t *testing.T) {
	metainfo := makeTestTorrentWithoutFiles()
	target := "<torrentmetainfo.T with .HexHash() = c339b4b17d653a1ed6b2eed15bc6dd19b828ea5c>"
	actual := metainfo.String()

	if target != actual {
		t.Error(fmt.Sprintf("String representation was:\n%v\nexpected:\n%v\n.", actual, target))
	}
}

func TestStringWithFiles(t *testing.T) {
	metainfo := makeTestTorrentWithFiles()
	target := "<torrentmetainfo.T with .HexHash() = 0ed4af3bdd5b66ce072ae4b7d31a9369217b0d46>"
	actual := metainfo.String()

	if target != actual {
		t.Error(fmt.Sprintf("String representation was:\n%v\nexpected:\n%v\n.", actual, target))
	}
}
