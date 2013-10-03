package bencoding

import (
	"bytes"
	"errors"
	"fmt"
)

type BValueType int

const (
	STRING BValueType = iota
	INTEGER
	LIST
	DICTIONARY
)

type BValue struct {
	t            BValueType
	string_value string
	int_value    int64
	list_value   []BValue
	map_value    map[string]BValue
}

func (value BValue) String() string {
	// Returns a JSON-like representation of this bencoding.Value.

	switch value.t {
	case STRING:
		return value.string_value

	case INTEGER:
		return fmt.Sprintf("%v", value.int_value)

	case LIST:
		var buffer bytes.Buffer
		buffer.WriteString("[")

		for index, value := range value.list_value {
			if index > 0 {
				buffer.WriteString(", ")
			}

			buffer.WriteString(value.String())
		}

		buffer.WriteString("]")
		return buffer.String()

	case DICTIONARY:
		return "{ ... }"
	}

	return "<invalid .t for bencoding.Value>"
}

func NewBValue(data interface{}) (value *BValue, err error) {
	// Uses reflection to convert an arbitrary object to a
	// bencoding.Value if possible.

	return nil, errors.New("Not implemented")
}

func NewBValueString(val string) *BValue {
	return &BValue{t: STRING, string_value: val}
}

func Bencode(data interface{}) (bencoded string, err error) {
	val, err := NewBValue(data)

	if err == nil {
		bencoded = val.Bencode()
	}

	return
}

func (value BValue) Bencode() string {
	return ""
}

// func Bdecode(data []byte) BValue {
// 	// 	
// }

// func Encode() string {
// 	val := Value{t: LIST, list_value: []Value{
// 		Value{t: STRING, string_value: "Hello"},
// 		Value{t: STRING, string_value: "World"}}}

// 	return fmt.Sprintf("%v\n", val)
// }
