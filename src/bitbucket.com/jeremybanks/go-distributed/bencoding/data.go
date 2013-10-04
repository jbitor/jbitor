package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func NewValue(data interface{}) (bval *Value, err error) {
	// Uses reflection to convert structures of int64, string,
	// map[string]interface{} and list[interface{}] to Values.
	// Also supports any BEncodable type.

	if value, ok := data.(int64); ok {
		bval = &(Value{t: INTEGER, value: value})
	} else if value, ok := data.(string); ok {
		bval = &(Value{t: STRING, value: value})
	} else if value, ok := data.(map[string]interface{}); ok {
		bvals := make(map[string]*Value)

		bval = &(Value{t: DICTIONARY, value: bvals})

		for key, item := range value {
			var item_bval *Value

			item_bval, err = NewValue(item)

			if err != nil {
				return
			}

			bvals[key] = item_bval
		}
	} else if value, ok := data.([]interface{}); ok {
		bvals := make([]*Value, len(value))

		bval = &(Value{t: LIST, value: bvals})

		for index, item := range value {
			var item_bval *Value

			item_bval, err = NewValue(item)

			if err != nil {
				return
			}

			bvals[index] = item_bval
		}
	} else if value, ok := data.(*Bencodable); ok {
		bval, err = (*value).MarshalBencodingValue()
	} else if value, ok := data.(Bencodable); ok {
		bval, err = value.MarshalBencodingValue()
	} else {
		err = errors.New(fmt.Sprintf("Invalid type for bencoding: %v", reflect.TypeOf(data)))
	}

	return
}

// When converting to and from JSON, byte values are mapped directly
// to unicode codepoints from 0 to 255. It will be impossible to
// unmarshal JSON that uses values outside of this range.

func (bval *Value) MarshalJSON() ([]byte, error) {
	panic("not implemented")
}

func (bval *Value) UnmarshalJSON([]byte) error {
	panic("not implemented")
}

func (bval *Value) String() string {
	// Returns a JSON-like representation of this bencoding.Value.

	switch bval.t {
	case STRING:
		return bval.value.(string)

	case INTEGER:
		return fmt.Sprintf("%v", bval.value.(int64))

	case LIST:
		var buffer bytes.Buffer
		buffer.WriteString("[")

		for index, item := range bval.value.([]*Value) {
			if index > 0 {
				buffer.WriteString(", ")
			}

			buffer.WriteString(item.String())
		}

		buffer.WriteString("]")
		return buffer.String()

	case DICTIONARY:
		var buffer bytes.Buffer
		buffer.WriteString("{")

		first := true

		for key, item := range bval.value.(map[string]*Value) {
			if first {
				first = false
			} else {
				buffer.WriteString(", ")
			}

			buffer.WriteString(key)

			buffer.WriteString(": ")

			buffer.WriteString(item.String())
		}

		buffer.WriteString("}")
		return buffer.String()
	}

	return fmt.Sprintf("<illegal Value.t: %v>", bval.t)
}
