package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

/*
	Provides a BValue type for encoding/decoding data as described here:
	https://wiki.theory.org/BitTorrentSpecification#Bencoding
*/

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
	dict_value   map[string]BValue
}

func (bval BValue) WriteBencoded(buffer *bytes.Buffer) (err error) {
	switch bval.t {
	case STRING:
		buffer.WriteString(
			fmt.Sprintf("%v:%v", len(bval.string_value), bval.string_value))

	case INTEGER:
		buffer.WriteString(
			fmt.Sprintf("i%ve", bval.int_value))

	case LIST:
		buffer.WriteString("l")
		for _, item := range bval.list_value {
			var item_str string
			item_str, err = item.Bencode()
			if err != nil {
				return
			}
			buffer.WriteString(item_str)
		}
		buffer.WriteString("e")

	case DICTIONARY:
		for key, item := range bval.dict_value {
			buffer.WriteString(
				fmt.Sprintf("%v:%v", len(key), key))

			var item_str string
			item_str, err = item.Bencode()
			if err != nil {
				return
			}
			buffer.WriteString(item_str)
		}

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

	// http://blog.golang.org/laws-of-reflection

	t := reflect.TypeOf(data)

	switch t.Kind() {
	case reflect.Int64:
		bval = &(BValue{t: INTEGER, int_value: data.(int64)})
	case reflect.String:
		bval = &(BValue{t: STRING, string_value: data.(string)})
	case reflect.Map:
		bval = &(BValue{t: DICTIONARY, dict_value: make(map[string]BValue)})

		for key, item := range data.(map[string]interface{}) {
			var item_bval *BValue

			item_bval, err = NewBValue(item)

			if err != nil {
				return
			}

			bval.dict_value[key] = *item_bval
		}
	case reflect.Array, reflect.Slice:
		length := len(data.([]interface{}))

		bval = &(BValue{
			t:          LIST,
			list_value: make([]BValue, length),
		})

		for index, item := range bval.list_value {
			var item_bval *BValue

			item_bval, err = NewBValue(item)

			if err != nil {
				return
			}

			bval.list_value[index] = *item_bval
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
		return bval.string_value

	case INTEGER:
		return fmt.Sprintf("%v", bval.int_value)

	case LIST:
		var buffer bytes.Buffer
		buffer.WriteString("[")

		for index, item := range bval.list_value {
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

		for key, item := range bval.dict_value {
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
