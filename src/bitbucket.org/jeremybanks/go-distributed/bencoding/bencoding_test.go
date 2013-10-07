package bencoding

import (
	"reflect"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	testCases := map[string]Bencodable{
		"i-1e":   Int(-1),
		"i0e":    Int(0),
		"i1e":    Int(1),
		"i1023e": Int(1023),
		"li1ei2ed2:abl1:ci4e1:de2:aai1eeledee": List{
			Int(1), Int(2), Dict{
				"ab": List{String("c"), Int(4), String("d")},
				"aa": Int(1),
			}, List{}, Dict{},
		},
	}

	for originalEncodedStr, originalDecoded := range testCases {
		originalEncoded := []byte(originalEncodedStr)

		encoded, err := Encode(originalDecoded)
		if err != nil {
			t.Error("Error encoding", originalDecoded, err)
		} else {
			if !reflect.DeepEqual(encoded, originalEncoded) {
				t.Error(
					"Encoded value", strconv.Quote(string(encoded)),
					"does not equal expected value", strconv.Quote(string(originalEncoded)))
			}
		}

		decoded, err := Decode(originalEncoded)
		if err != nil {
			t.Error("Error decoding", strconv.Quote(string(originalEncoded)), err)
		} else {
			if !reflect.DeepEqual(decoded, originalDecoded) {
				t.Error(
					"Decoded value", decoded,
					"does not equal expected value", originalDecoded)
			}
		}
	}
}
