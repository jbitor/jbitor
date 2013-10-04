package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

/*
	Provides a BValue type for encoding/decoding data as described here:
	http://www.bittorrent.org/beps/bep_0003.html
*/

type BValueType int

const (
	STRING BValueType = iota
	INTEGER
	LIST
	DICTIONARY
)

type BValue struct {
	t     BValueType
	value interface{}
}

func (bval BValue) WriteBencoded(buffer *bytes.Buffer) (err error) {
	switch bval.t {
	case STRING:
		buffer.WriteString(
			fmt.Sprintf("%v:%v", len(bval.value.(string)), bval.value.(string)))

	case INTEGER:
		buffer.WriteString(
			fmt.Sprintf("i%ve", bval.value.(int64)))

	case LIST:
		buffer.WriteString("l")
		for _, item := range bval.value.([]BValue) {
			var item_str string
			item_str, err = item.Bencode()
			if err != nil {
				return
			}
			buffer.WriteString(item_str)
		}
		buffer.WriteString("e")

	case DICTIONARY:
		// FIXME: keys must be sorted.

		buffer.WriteString("d")
		for key, item := range bval.value.(map[string]BValue) {
			buffer.WriteString(
				fmt.Sprintf("%v:%v", len(key), key))

			var item_str string
			item_str, err = item.Bencode()
			if err != nil {
				return
			}
			buffer.WriteString(item_str)
		}
		buffer.WriteString("e")

	default:
		err = errors.New(fmt.Sprintf("Illegal BValue.t: %v", bval.t))
	}

	return
}

func (bval BValue) Bencode() (str string, err error) {
	var buffer bytes.Buffer

	err = bval.WriteBencoded(&buffer)

	if err == nil {
		str = buffer.String()
	}

	return
}

func Bdecode(string) (bval *BValue, err error) {
	err = errors.New("Bdecode Not implemented")

	return
}

func Bencode(data interface{}) (bencoded string, err error) {
	val, err := NewBValue(data)

	if err == nil {
		bencoded, err = val.Bencode()
	}

	return
}

func NewBValue(data interface{}) (bval *BValue, err error) {
	// Uses reflection to convert an arbitrary object to a
	// bencoding.Value if possible.

	t := reflect.TypeOf(data)

	switch t.Kind() {
	case reflect.Int64:
		bval = &(BValue{t: INTEGER, value: data.(int64)})

	case reflect.String:
		bval = &(BValue{t: STRING, value: data.(string)})

	case reflect.Map:
		bvals := make(map[string]BValue)

		bval = &(BValue{t: DICTIONARY, value: bvals})

		for key, item := range data.(map[string]interface{}) {
			var item_bval *BValue

			item_bval, err = NewBValue(item)

			if err != nil {
				return
			}

			bvals[key] = *item_bval
		}

	case reflect.Array, reflect.Slice:
		length := len(data.([]interface{}))

		bvals := make([]BValue, length)

		bval = &(BValue{t: LIST, value: bvals})

		for index, item := range data.([]interface{}) {
			var item_bval *BValue

			item_bval, err = NewBValue(item)

			if err != nil {
				return
			}

			bvals[index] = *item_bval
		}

	default:
		err = errors.New(fmt.Sprintf("Invalid type for bencoding: %v", t))
	}

	return
}

func (bval BValue) String() string {
	// Returns a JSON-like representation of this bencoding.Value.

	switch bval.t {
	case STRING:
		return bval.value.(string)

	case INTEGER:
		return fmt.Sprintf("%v", bval.value.(int64))

	case LIST:
		var buffer bytes.Buffer
		buffer.WriteString("[")

		for index, item := range bval.value.([]BValue) {
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

		for key, item := range bval.value.(map[string]BValue) {
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
