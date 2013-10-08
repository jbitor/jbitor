package main

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"bitbucket.org/jeremybanks/go-distributed/torrentutils"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v COMMAND\n", os.Args[0])
		os.Exit(1)
		return
	}

	command := os.Args[1]
	commandArgs := os.Args[2:]

	switch command {
	case "torrent":
		cmdTorrent(commandArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %v\n", command)
		os.Exit(1)
	}
}

func cmdTorrent(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %v torrent SUBCOMMAND\n", os.Args[0])
		os.Exit(1)
		return
	}

	subcommand := args[0]
	subcommandArgs := args[1:]

	switch subcommand {
	case "make":
		cmdTorrentMake(subcommandArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown torrent subcommand: %v\n", subcommand)
		os.Exit(1)
	}
}

func cmdTorrentMake(args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %v torrent make PATH\n", os.Args[0])
		os.Exit(1)
		return
	}

	infoDict, err := torrentutils.GenerateTorrentMetaInfo(torrentutils.CreationOptions{
		Path:           args[0],
		PieceLength:    524288,
		ForceMultiFile: false,
	})
	if err != nil {
		panic(err)
	}

	infoData, err := bencoding.Encode(infoDict)
	if err != nil {
		panic(err)
	}

	torrentDict := bencoding.Dict{
		"info":     infoDict,
		"announce": bencoding.String("http://localhost/"),
	}

	torrentData, err := bencoding.Encode(torrentDict)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encoding torrent data:", err)
		os.Exit(1)
		return
	}

	hasher := sha1.New()
	hasher.Write(infoData)
	hash := hasher.Sum(nil)
	infoHashHex := hex.EncodeToString(hash)

	fmt.Fprintf(os.Stderr, "Generated torrent btih=%v.\n", infoHashHex)
	os.Stderr.Sync()

	os.Stdout.Write(torrentData)
	os.Stdout.Sync()

	_ = infoHashHex
}
