package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

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
