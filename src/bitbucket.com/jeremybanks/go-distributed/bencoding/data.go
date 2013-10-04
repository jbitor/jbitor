package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func NewBValue(data interface{}) (bval *BValue, err error) {
	// Uses reflection to convert structures of int64, string,
	// map[string]interface{} and list[interface{}] to BValues.
	// Also supports any BEncodable type.

	if value, ok := data.(int64); ok {
		bval = &(BValue{t: INTEGER, value: value})
	} else if value, ok := data.(string); ok {
		bval = &(BValue{t: STRING, value: value})
	} else if value, ok := data.(map[string]interface{}); ok {
		bvals := make(map[string]*BValue)

		bval = &(BValue{t: DICTIONARY, value: bvals})

		for key, item := range value {
			var item_bval *BValue

			item_bval, err = NewBValue(item)

			if err != nil {
				return
			}

			bvals[key] = item_bval
		}
	} else if value, ok := data.([]interface{}); ok {
		length := len(data.([]interface{}))

		bvals := make([]*BValue, length)

		bval = &(BValue{t: LIST, value: bvals})

		for index, item := range value {
			var item_bval *BValue

			item_bval, err = NewBValue(item)

			if err != nil {
				return
			}

			bvals[index] = item_bval
		}
	} else if value, ok := data.(*Bencodable); ok {
		bval = (*value).BValue()
	} else if value, ok := data.(Bencodable); ok {
		bval = value.BValue()
	} else {
		err = errors.New(fmt.Sprintf("Invalid type for bencoding: %v", reflect.TypeOf(data)))
	}

	return
}

func (bval *BValue) String() string {
	// Returns a JSON-like representation of this bencoding.Value.

	switch bval.t {
	case STRING:
		return bval.value.(string)

	case INTEGER:
		return fmt.Sprintf("%v", bval.value.(int64))

	case LIST:
		var buffer bytes.Buffer
		buffer.WriteString("[")

		for index, item := range bval.value.([]*BValue) {
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

		for key, item := range bval.value.(map[string]*BValue) {
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

	return fmt.Sprintf("<illegal BValue.t: %v>", bval.t)
}
