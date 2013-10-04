package main

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"fmt"
)

func main() {
	msg, _ := bencoding.NewValue("I guess the tests passed!")
	fmt.Println(msg)
}
