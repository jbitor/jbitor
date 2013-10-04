package bencoding

import (
	"fmt"
	"testing"
)

// Test helpers

type bencodableInt32 int32

func (self bencodableInt32) BValue() *BValue {
	val, err := NewBValue(int64(self))

	if err != nil {
		panic(fmt.Sprintf("failed to create BValue: %v", err))
	}

	return val
}

func testEncodings(t *testing.T, encodings map[string]interface{}) {
	for target, input := range encodings {
		actual, err := Bencode(input)

		if err != nil {
			t.Error("Error while encoding", input, "-", err)
		} else if target != actual {
			t.Error("Encoding", input, "produced", actual, "instead of", target)
		}
	}
}

func testUnencodables(t *testing.T, unencodables []interface{}) {
	for _, unencodable := range unencodables {
		_, err := Bencode(unencodable)

		if err == nil {
			t.Error("No error when attempting to encode unencodable", unencodable)
		}
	}
}

func testDecodings(t *testing.T, decodings map[string]*BValue) {
	for input, target := range decodings {
		actual, err := Bdecode(input)

		if err != nil {
			t.Error("Error while decoding", input, "-", err)
		} else if target != actual {
			// XXX: I'm not sure if you can compare pointers like this...
			// Is there any way to check for recursive equality of pointerful structs?
			t.Error("Decoding", input, "produced", actual, "instead of", target)
		}
	}
}

func testUndecodables(t *testing.T, undecodables []string) {
	for _, undecodable := range undecodables {
		_, err := Bdecode(undecodable)

		if err == nil {
			t.Error("No error when attempting to decode undecodable", undecodable)
		}
	}
}

// Test cases

func TestIntegerEncoding(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"i3e":   int64(3),
		"i-3e":  int64(-3),
		"i6e":   int64(6),
		"i0e":   int64(0),
		"i16e":  BValue{t: INTEGER, value: int64(16)},  // already a BValue
		"i17e":  &BValue{t: INTEGER, value: int64(17)}, // already a *BValue
		"i-99e": bencodableInt32(-99),
	})

	testUnencodables(t, []interface{}{
		int32(99),  // wrong integer type
		uint64(99), // wrong integer type
	})

	testUndecodables(t, []string{
		"i-0e", // non-canonical encoding
		"i01e", // non-canonical encoding
	})
}

func TestListEncoding(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"l4:spam4:eggse": []interface{}{"spam", "eggs"},
	})
}

func TestDictionaryEncoding(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"d3:cow3:moo4:spam4:eggse": map[string]interface{}{"cow": "moo", "spam": "eggs"},
		"d4:spaml1:a1:bee":         map[string]interface{}{"spam": []interface{}{"a", "b"}},
	})

	testUndecodables(t, []string{
		("d1:bi1e1:ai2ee"), // keys out of order
	})
}

func TestEncodeUnverifiedExample(t *testing.T) {
	testEncodings(t, map[string]interface{}{
		"d6:lengthi512e4:miscd5:hello6:World!e4:name9:Test Data12:piece lengthi1024e6:pieces20:\x00234567890123456789\xFFe": map[string]interface{}{
			"piece length": int64(1024),
			"pieces":       "\x00234567890123456789\xFF",
			"name":         "Test Data",
			"length":       int64(512),
			"misc": map[string]interface{}{
				"hello": "World!",
			},
		},
	})
}
