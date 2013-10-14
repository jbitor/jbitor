package dht

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[     DHT     ] ", log.Lshortfile)
}
