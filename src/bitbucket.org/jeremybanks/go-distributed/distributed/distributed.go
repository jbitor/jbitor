package main

import (
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
	case "info":
		cmdTorrentInfo(subcommandArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown torrent subcommand: %v\n", subcommand)
		os.Exit(1)
	}
}

func cmdTorrentInfo(args []string) {
	panic("torrent info subcommand not implemented")
}
