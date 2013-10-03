package main

import (
	"bitbucket.com/jeremybanks/go-distributed/bencoding"
	"fmt"
)

func main() {
	val, err := bencoding.NewBValue(map[string]interface{}{
		"piece length": 1024,
		"pieces":       "\x00234567890123456789\xFF",
		"name":         "Test Data",
		"length":       512,
		"misc": map[string]interface{}{
			"hello": "World!",
		},
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n", val)
	}
}
