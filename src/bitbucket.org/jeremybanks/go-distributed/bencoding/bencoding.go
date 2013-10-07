package bencoding

import (
	"bytes"
	"fmt"
)

type Bencodable interface {
	WriteBencoded(*bytes.Buffer) error
}

type Int int64
type Str string
type List []Bencodable
type Dict map[Str]Bencodable

func Bencode(bval Bencodable) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := bval.WriteBencoded(buffer)
	if err == nil {
		return buffer.Bytes(), nil
	} else {
		return nil, err
	}
}

func (bint Int) WriteBencoded(buffer *bytes.Buffer) error {
	_, err := fmt.Fprintf(buffer, "i%ve", int64(bint))
	return err
}

func (bstr Str) WriteBencoded(buffer *bytes.Buffer) error {
	_, err := fmt.Fprintf(buffer, "%v:%v", len(bstr), string(bstr))
	return err
}

func (blist List) WriteBencoded(buffer *bytes.Buffer) error {
	_, err := buffer.Write([]byte("l"))
	if err != nil {
		return err
	}

	for _, item := range blist {
		err = item.WriteBencoded(buffer)
		if err != nil {
			return err
		}
	}

	_, err = buffer.Write([]byte("e"))
	if err != nil {
		return err
	}

	return nil
}

func (bdict Dict) WriteBencoded(buffer *bytes.Buffer) error {
	_, err := buffer.Write([]byte("d"))
	if err != nil {
		return err
	}

	for key, value := range bdict {
		err = key.WriteBencoded(buffer)
		if err != nil {
			return err
		}
		err = value.WriteBencoded(buffer)
		if err != nil {
			return err
		}
	}

	_, err = buffer.Write([]byte("e"))
	if err != nil {
		return err
	}

	return nil
}
