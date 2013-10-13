package main

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"bitbucket.org/jeremybanks/go-distributed/dht"
	"bitbucket.org/jeremybanks/go-distributed/torrentutils"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
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
	case "dht":
		cmdDht(commandArgs)
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

	path := args[0]

	infoDict, err := torrentutils.GenerateTorrentMetaInfo(torrentutils.CreationOptions{
		Path:           path,
		PieceLength:    32768,
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
		"info": infoDict,
		"nodes": bencoding.List{
			bencoding.List{
				bencoding.String("127.0.0.1"),
				bencoding.Int(6881),
			},
		},
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

func cmdDht(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %v dht SUBCOMMAND\n", os.Args[0])
		os.Exit(1)
		return
	}

	subcommand := args[0]
	subcommandArgs := args[1:]

	switch subcommand {
	case "helloworld":
		cmdDhtHelloWorld(subcommandArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown dht subcommand: %v\n", subcommand)
		os.Exit(1)
	}
}

func cmdDhtHelloWorld(args []string) {
	node := dht.NewLocalNode()

	terminated := make(chan error)
	go node.Run(terminated)

	knownNode := dht.RemoteNodeFromAddress(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6881,
	})

	knownNode = node.AddOrGetRemoteNode(knownNode)

	fmt.Printf("Hello, I am %v.\n", node)
	fmt.Printf("I know of %v.\n", node.Nodes)

	fmt.Printf("\nI am attempting to ping a DHT node at localhost:6881.\n")
	pingResult, pingErr := node.Ping(knownNode)

	fmt.Printf("Ping initiated\n")

	select {
	case result := <-pingResult:
		fmt.Printf("got ping result: %v\n", *result)
	case result := <-pingErr:
		fmt.Printf("got ping error: %v\n", result)
	}

	fmt.Printf("I know of %v.\n", node.Nodes)

	nodeId, _ := hex.DecodeString("b7271d0b5577918ee92b1b5378d89e56ad08ba80")
	findResult, findErr := node.FindNode(knownNode, dht.NodeId(nodeId))

	select {
	case result := <-findResult:
		fmt.Printf("got find result: %v\n", result)
	case result := <-findErr:
		fmt.Printf("got find error: %v\n", result)
	}

}
