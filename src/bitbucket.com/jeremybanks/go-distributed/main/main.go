package main

import (
	"bitbucket.com/jeremybanks/go-distributed/bencoding"
	"fmt"
)

func main() {
	val, err := bencoding.NewBValue(map[string]interface{}{
		"piece length": int64(1024),
		"pieces":       "\x00234567890123456789\xFF",
		"name":         "Test Data",
		"length":       int64(512),
		"misc": map[string]interface{}{
			"hello": "World!",
		},
	})

	if err != nil {
		fmt.Printf("Error creating value:\n\t%v\n", err)
		return
	}

	fmt.Printf("Value:\n\t%v\n\n", val)

	bencoded, err := val.Bencode()

	if err != nil {
		fmt.Printf("Error encoding:\n\t%v\n", err)
		return
	}

	fmt.Printf("Encoded:\n\t%v\n\n", bencoded)

	restored, err := bencoding.Bdecode(bencoded)

	if err != nil {
		fmt.Printf("Error decoding:\n\t%v\n", err)
		return
	}

	fmt.Printf("Decoded value:\n\t%v\n\n", restored)

}
