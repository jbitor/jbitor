package main

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"bitbucket.org/jeremybanks/go-distributed/dht"
	"bitbucket.org/jeremybanks/go-distributed/torrentutils"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net"
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
	default:
		logger.Fatalf("Unknown command: %v\n", command)
		return
	}
}

func cmdTorrent(args []string) {
	if len(args) == 0 {
		logger.Fatalf("Usage: %v torrent SUBCOMMAND\n", os.Args[0])
		return
	}

	subcommand := args[0]
	subcommandArgs := args[1:]

	switch subcommand {
	case "make":
		cmdTorrentMake(subcommandArgs)
	default:
		logger.Fatalf("Unknown torrent subcommand: %v\n", subcommand)
		return
	}
}

func cmdTorrentMake(args []string) {
	if len(args) != 1 {
		logger.Fatalf("Usage: %v torrent make PATH\n", os.Args[0])
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
		logger.Fatalf("Error encoding torrent data:", err)
		return
	}

	hasher := sha1.New()
	hasher.Write(infoData)
	hash := hasher.Sum(nil)
	infoHashHex := hex.EncodeToString(hash)

	logger.Printf("Generated torrent btih=%v.\n", infoHashHex)

	os.Stdout.Write(torrentData)
	os.Stdout.Sync()

	_ = infoHashHex
}

func cmdDht(args []string) {
	if len(args) == 0 {
		logger.Fatalf("Usage: %v dht SUBCOMMAND\n", os.Args[0])
		return
	}

	subcommand := args[0]
	subcommandArgs := args[1:]

	switch subcommand {
	case "helloworld":
		cmdDhtHelloWorld(subcommandArgs)
	default:
		logger.Fatalf("Unknown dht subcommand: %v\n", subcommand)
		return
	}
}

func cmdDhtHelloWorld(args []string) {
	var node *dht.LocalNode

	if len(args) > 0 {
		path := args[0]

		// XXX: We shold open and lock the file, so that instances can't clobber.
		nodeData, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Printf("Unable to read existing DHT node (%v). Creating a new one.\n", err)
			node = dht.NewLocalNode()
		} else {
			nodeDict, err := bencoding.Decode(nodeData)
			if err != nil {
				logger.Fatalf("%v\n", err)
				return
			}

			nodeDictAsDict, ok := nodeDict.(bencoding.Dict)
			if !ok {
				logger.Fatalf("\n")
				return
			}

			node = dht.LocalNodeFromBencodingDict(nodeDictAsDict)
		}

		defer func() {
			// save node
			nodeData, err := bencoding.Encode(node)

			if err != nil {
				logger.Fatalf("Error encoding local node state: %v\n", err)
				return
			}

			logger.Printf("Saving LocalNode state to %v.\n", path)
			// XXX: These flags give it 0 peermissions!
			err = ioutil.WriteFile(path, nodeData, 0644)

			if err != nil {
				logger.Fatalf("Error writing local node state: %v\n", err)
				return
			}
		}()
	} else {
		node = dht.NewLocalNode()
	}

	terminated := make(chan error)
	go node.Run(terminated)

	knownNode := dht.RemoteNodeFromAddress(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6881,
	})

	knownNode = node.AddOrGetRemoteNode(knownNode)

	logger.Printf("Hello, I am %v.\n", node)
	logger.Printf("I know of %v.\n", node.Nodes)

	logger.Printf("I am attempting to ping a DHT node at localhost:6881.\n")
	pingResult, pingErr := node.Ping(knownNode)

	logger.Printf("Ping initiated\n")

	nodeId, _ := hex.DecodeString("b7271d0b5577918ee92b1b5378d89e56ad08ba80")
	findResult, findErr := node.FindNode(knownNode, dht.NodeId(nodeId))

	for i := 0; i < 2; i++ {
		select {
		case result := <-pingResult:
			logger.Printf("got ping result: %v\n", *result)
		case result := <-pingErr:
			logger.Printf("got ping error: %v\n", result)
		case result := <-findResult:
			logger.Printf("got find result: %v\n", result)
		case result := <-findErr:
			logger.Printf("got find error: %v\n", result)
		}
	}

	logger.Printf("I know of %v.\n", node.Nodes)

}
