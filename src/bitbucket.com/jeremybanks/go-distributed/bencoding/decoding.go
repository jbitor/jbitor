package bencoding

import (
	"errors"
)

func Bdecode(str string) (bval *Value, err error) {
	err = errors.New("Bdecode Not implemented")

	return
}

func BdecodeTo(str string, target *Bdecodable) (err error) {
	panic("BdecodeTo not implemented")

	// TODO: Short-circuit if Unmarshaller too

	// bval, err := Bdecode(str)

	// if err != nil {
	// 	return
	// }

	// err = target.InitFromBencodingValue(bval)

	return
}
