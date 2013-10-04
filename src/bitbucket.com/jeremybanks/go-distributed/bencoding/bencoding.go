package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
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

type Bencodable interface {
	BValue() *BValue
}

// type BDecodable interface {
// 	initFromBValue(*BValue) error
// }

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
		for _, item := range bval.value.([]*BValue) {
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

		keys := make([]string, len(bval.value.(map[string]*BValue)))

		i := 0
		for key, _ := range bval.value.(map[string]*BValue) {
			keys[i] = key
			i += 1
		}

		sort.Strings(keys)

		buffer.WriteString("d")
		for _, key := range keys {
			item := bval.value.(map[string]*BValue)[key]

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

func Bdecode(str string) (bval *BValue, err error) {
	err = errors.New("Bdecode Not implemented")

	return
}

// func BdecodeTo(str string, target *BDecodable) (err error) {
// 	bval, err := Bdecode(str)

// 	if err != nil {
// 		return
// 	}

// 	err = target.initFromBValue(bval)

// 	return
// }

func Bencode(data interface{}) (bencoded string, err error) {
	var bval *BValue
	var bval_data BValue
	var ok bool

	if bval, ok = data.(*BValue); ok {

	} else if bval_data, ok = data.(BValue); ok {
		bval = &bval_data
	} else {
		bval, err = NewBValue(data)
	}

	if err == nil {
		bencoded, err = bval.Bencode()
	}

	return
}

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
