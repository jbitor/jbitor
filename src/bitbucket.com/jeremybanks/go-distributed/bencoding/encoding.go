package bencoding

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

func (bval Value) WriteBencoded(buffer *bytes.Buffer) (err error) {
	switch bval.t {
	case STRING:
		value, ok := bval.value.(string)
		if !ok {
			panic("STRING BValue did not have a string value.")
		}

		buffer.WriteString(fmt.Sprintf("%v:%v", len(value), value))

	case INTEGER:
		value, ok := bval.value.(int64)
		if !ok {
			panic("INTEGER BValue did not have an int64 value.")
		}

		buffer.WriteString(fmt.Sprintf("i%ve", value))

	case LIST:
		value, ok := bval.value.([]*Value)
		if !ok {
			panic("LIST BValue did not have a []*Value value.")
		}

		buffer.WriteString("l")
		for _, item := range value {
			var item_str string
			item_str, err = item.Bencode()
			if err != nil {
				return
			}
			buffer.WriteString(item_str)
		}
		buffer.WriteString("e")

	case DICTIONARY:
		value, ok := bval.value.(map[string]*Value)
		if !ok {
			panic("LIST BValue did not have a []*Value value.")
		}

		keys := make([]string, len(value))

		i := 0
		for key, _ := range value {
			keys[i] = key
			i += 1
		}

		sort.Strings(keys)

		buffer.WriteString("d")
		for _, key := range keys {
			item := value[key]

			buffer.WriteString(fmt.Sprintf("%v:%v", len(key), key))

			var item_str string
			item_str, err = item.Bencode()
			if err != nil {
				return
			}
			buffer.WriteString(item_str)
		}
		buffer.WriteString("e")

	default:
		err = errors.New(fmt.Sprintf("Illegal Value.t: %v", bval.t))
	}

	return
}

func (bval Value) Bencode() (str string, err error) {
	var buffer bytes.Buffer

	// TODO: Short-circuit if Marshaller

	err = bval.WriteBencoded(&buffer)

	if err == nil {
		str = buffer.String()
	}

	return
}

func Bencode(data interface{}) (bencoded string, err error) {
	var bval *Value
	var bval_data Value
	var ok bool

	if bval, ok = data.(*Value); ok {

	} else if bval_data, ok = data.(Value); ok {
		bval = &bval_data
	} else {
		bval, err = NewValue(data)
	}

	if err == nil {
		bencoded, err = bval.Bencode()
	}

	return
}
