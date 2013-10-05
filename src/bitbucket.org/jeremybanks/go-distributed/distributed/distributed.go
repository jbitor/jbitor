package main

import (
	"bitbucket.org/jeremybanks/go-distributed/torrentmetainfo"
	"fmt"
	"io/ioutil"
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
	case "info":
		cmdTorrentInfo(subcommandArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown torrent subcommand: %v\n", subcommand)
		os.Exit(1)
	}
}

func cmdTorrentInfo(args []string) {
	var data []byte
	var err error

	if len(args) == 1 {
		filename := args[0]
		data, err = ioutil.ReadFile(filename)
	} else if len(args) == 0 {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %v torrent info [FILE]\n", os.Args[0])
		os.Exit(1)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
		return
	}

	var metainfo torrentmetainfo.T
	err = metainfo.UnmarshalTorrentBencoding(data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error intepreting torrent data: %v\n", err)
	}

	hash, _ := metainfo.HexHash()
	fmt.Println("infohash:", hash)
	fmt.Println("    name:", metainfo.Name)
	fmt.Println("  length:", metainfo.Length)
	fmt.Println("   files:", metainfo.Files)
}
