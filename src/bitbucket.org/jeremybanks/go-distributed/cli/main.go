package main

import (
	"os"
)

func main() {
	if len(os.Args) == 1 {
		logger.Fatalf("Usage: %v COMMAND\n", os.Args[0])
		return
	}

	command := os.Args[1]
	commandArgs := os.Args[2:]

	switch command {
	case "torrent":
		cmdTorrent(commandArgs)
	case "dht":
		cmdDht(commandArgs)
	case "json":
		cmdJson(commandArgs)
	default:
		logger.Fatalf("Unknown command: %v\n", command)
		return
	}
}
