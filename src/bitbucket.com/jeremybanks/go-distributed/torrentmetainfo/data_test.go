package torrentmetainfo

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	metainfo := makeTestTorrentWithoutFiles()
	target := ""
	actual := metainfo.String()

	if target != actual {
		t.Error(fmt.Sprintf("String representation was:\n%v\nexpected:\n%v\n.", actual, target))
	}
}
