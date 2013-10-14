package main

import (
	"os"
)

func cmdJson(args []string) {
	if len(args) == 0 {
		logger.Fatalf("Usage: %v json SUBCOMMAND\n", os.Args[0])
		return
	}

	subcommand := args[0]
	subcommandArgs := args[1:]

	switch subcommand {
	case "from-bencoding":
		cmdJsonFromBencoding(subcommandArgs)
	case "to-bencoding":
		cmdJsonToBencoding(subcommandArgs)
	default:
		logger.Fatalf("Unknown torrent subcommand: %v\n", subcommand)
		return
	}

}

func cmdJsonFromBencoding(args []string) {
	logger.Fatalf("json from-bencoding not implemented")
	return
}

func cmdJsonToBencoding(args []string) {
	logger.Fatalf("json to-bencoding not implemented")
	return
}
