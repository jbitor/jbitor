package bencoding

import "testing"

func TestEncodeString(t *testing.T) {
	encoded, _ := Bencode("spam")
	if encoded != "4:spam" {
		t.Fail()
	}
}

func TestEncodeInteger(t *testing.T) {
	encoded, _ := Bencode(int64(3))
	if encoded != "i3e" {
		t.Error("encoded is", encoded)
	}

	encoded, _ = Bencode(int64(-3))
	if encoded != "i-3e" {
		t.Error("encoded is", encoded)
	}

	encoded, _ = Bencode(int64(0))
	if encoded != "i0e" {
		t.Error("encoded is", encoded)
	}
}

func TestEncodeList(t *testing.T) {
	encoded, _ := Bencode([]interface{}{"spam", "eggs"})
	if encoded != "l4:spam4:eggse" {
		t.Error("encoded is", encoded)
	}
}

func TestEncodeDictionary(t *testing.T) {
	encoded, _ := Bencode(map[string]interface{}{"cow": "moo", "spam": "eggs"})
	if encoded != "d3:cow3:moo4:spam4:eggse" {
		t.Error("encoded is", encoded)
	}

	encoded, _ = Bencode((map[string]interface{}{"spam": []interface{}{"a", "bee"}}))
	if encoded != "d4:spaml1:a1:bee" {
		t.Error("encoded is", encoded)
	}
}

// TODO: add tests asserting that an error happens
// also, test for error codes above
