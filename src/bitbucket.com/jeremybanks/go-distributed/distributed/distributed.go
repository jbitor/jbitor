package main

import (
	"bitbucket.com/jeremybanks/go-distributed/bencoding"
	"fmt"
)

func main() {
	msg, _ := bencoding.NewValue("I guess the tests passed!")
	fmt.Println(msg)
}
