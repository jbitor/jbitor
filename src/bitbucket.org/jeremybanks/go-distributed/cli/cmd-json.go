package main

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"encoding/json"
	"io/ioutil"
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
	if len(args) != 0 {
		logger.Fatalf("Usage: %v json from-bencoding < FOO.torrent > FOO.torrent.json\n", os.Args[0])
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		logger.Fatalf("Error reading stdin: %v\n", err)
		return
	}

	decoded, err := bencoding.Decode(data)
	if err != nil {
		logger.Fatalf("Error bdecoding stdin: %v\n", err)
		return
	}

	jsonable, err := decoded.ToJsonable()
	if err != nil {
		logger.Fatalf("Error converting bencoded value to jsonable: %v\n", err)
	}

	jsoned, err := json.Marshal(jsonable)
	if err != nil {
		logger.Fatalf("Error json-encoding data: %v\n", err)
		return
	}

	os.Stdout.Write(jsoned)
	os.Stdout.Sync()
}

func cmdJsonToBencoding(args []string) {
	if len(args) != 0 {
		logger.Fatalf("Usage: %v json to-bencoding < FOO.torrent.json > FOO.torrent\n", os.Args[0])
		return
	}

	logger.Fatalf("json to-bencoding not implemented")
	return
}