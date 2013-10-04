package bencoding

import (
	"errors"
)

func Bdecode(encoded []byte) (bval *Value, err error) {
	err = errors.New("Bdecode Not implemented")

	return
}

func BdecodeTo(encoded []byte, target *Bdecodable) (err error) {
	panic("BdecodeTo not implemented")

	// TODO: Short-circuit if Unmarshaller too

	// bval, err := Bdecode(str)

	// if err != nil {
	// 	return
	// }

	// err = target.InitFromBencodingValue(bval)

	return
}
