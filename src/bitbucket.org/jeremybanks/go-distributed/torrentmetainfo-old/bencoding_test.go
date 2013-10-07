package torrentmetainfo

import (
	"fmt"
	"testing"
)

func TestBencodeWithoutFiles(t *testing.T) {
	metainfo := makeTestTorrentWithoutFiles()
	target := ("d6:lengthi613e4:name7:data.go12:piece lengthi32768e6:" +
		"pieces20:\x8E\x66\xF6\xC8\x3A\xF4\x23\x52\xEF\xDF\x5A" +
		"\x2C\x1D\x02\x16\x21\x22\xD7\x63\x94e")
	actual, err := metainfo.MarshalBencoding()

	if target != string(actual) {
		t.Log(fmt.Sprintf("error was: %v", err))
		t.Error(fmt.Sprintf("Bencoding was:\n%v\nexpected:\n%v\n.", string(actual), target))
	}
}

func TestBencodeWithFiles(t *testing.T) {
	metainfo := makeTestTorrentWithFiles()
	target := ("d5:filesld6:lengthi2629e4:pathl7:data.goeed6:length" +
		"i434e4:pathl8:types.goeee4:name4:test12:piece length" +
		"i32768e6:pieces20:\x0E\x35\xD4\x04\x61\xFB\x99\x77\x46" +
		"\x5E\xAB\xB6\xA0\x9A\xC7\x84\x48\xF9\x69\x98e")
	actual, err := metainfo.MarshalBencoding()

	if target != string(actual) {
		t.Log(fmt.Sprintf("error was: %v", err))
		t.Error(fmt.Sprintf("Bencoding was:\n%v\nexpected:\n%v\n.", string(actual), target))
	}
}
