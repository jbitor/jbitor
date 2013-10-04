package bencoding

import (
	"errors"
)

func Bdecode(str string) (bval *BValue, err error) {
	err = errors.New("Bdecode Not implemented")

	return
}

// func BdecodeTo(str string, target *BDecodable) (err error) {
//  bval, err := Bdecode(str)

//  if err != nil {
//      return
//  }

//  err = target.initFromBValue(bval)

//  return
// }
