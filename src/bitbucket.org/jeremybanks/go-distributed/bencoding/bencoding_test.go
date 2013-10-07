package bencoding

import (
	"reflect"
	"strconv"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	var value Bencodable
	for _, value = range []Bencodable{
		Int(-1),
		Int(0),
		Int(1),
		String(""),
		String("helloe worldee"),
		List{Int(-1), Int(-1), Int(-2), Int(0), Int(4), Int(5)},
		List{Int(0), String("Hello worlde")},
		Dict{"hello": List{String("world"), Int(2), String("you")}},
	} {
		encoded, err := Encode(value)
		if err != nil {
			t.Error("Error encoding", value, err)
			continue
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Error("After encoding", value, ", error decoding", strconv.Quote(string(encoded)), err)
			continue
		}

		if !reflect.DeepEqual(value, decoded) {
			t.Error("Value did not safely round-trip", value, decoded)
		}
	}
}
