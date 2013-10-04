package bencoding

import "testing"

func TestIntegerDecoding(t *testing.T) {
	testUndecodables(t, []string{
		"i-0e", // non-canonical encoding
		"i01e", // non-canonical encoding
	})
}

func TestDictionaryDecoding(t *testing.T) {
	testUndecodables(t, []string{
		("d1:bi1e1:ai2ee"), // keys out of order
	})
}
