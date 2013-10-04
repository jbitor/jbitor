package bencoding

import (
	"errors"
	"fmt"
	"testing"
)

// Test helper function

func testEncodings(t *testing.T, encodings map[string]interface{}) {
	for target, input := range encodings {
		actual, err := Bencode(input)

		if err != nil {
			t.Error("Error while encoding", input, "-", err)
		} else if target != string(actual) {
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

func testDecodings(t *testing.T, decodings map[string]*Value) {
	for input, target := range decodings {
		actual, err := Bdecode([]byte(input))

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
		_, err := Bdecode([]byte(undecodable))

		if err == nil {
			t.Error("No error when attempting to decode undecodable", undecodable)
		}
	}
}

// Test helper types

type bencodableInt32 int32

func (self bencodableInt32) MarshalBencodingValue() (bval *Value, err error) {
	bval, err = NewValue(int64(self))

	if err != nil {
		panic(fmt.Sprintf("failed to create Value: %v", err))
	}

	return
}

func (self bencodableInt32) MarshalBencoding() (encoded []byte, err error) {
	return Bencode(self)
}

func (self bencodableInt32) UnmarshalBencodingValue(bval *Value) (err error) {
	value, ok := bval.Value.(int64)
	if !ok {
		err = errors.New(fmt.Sprintf("%s is not an int64", bval.Value))
	} else {
		self = bencodableInt32(value)
	}
	return
}

func (self bencodableInt32) UnmarshalBencoding(encoded []byte) (err error) {
	var bval *Value
	bval, err = Bdecode(encoded)
	if err == nil {
		err = self.UnmarshalBencodingValue(bval)
	}
	return
}
