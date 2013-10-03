package bencoding

import (
	"bytes"
	"fmt"
)

type ValueType int

const (
	STRING ValueType = iota
	INTEGER
	LIST
	DICTIONARY
)

type Value struct {
	t            ValueType
	string_value string
	int_value    int64
	list_value   []Value
	map_value    map[string]Value
}

func (value Value) String() string {
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

func Encode() string {
	val := Value{t: LIST, list_value: []Value{
		Value{t: STRING, string_value: "Hello"},
		Value{t: STRING, string_value: "World"}}}

	return fmt.Sprintf("%v\n", val)
}
