package main

import (
	"bitbucket.org/jeremybanks/go-distributed/bencoding"
	"bitbucket.org/jeremybanks/go-distributed/dht"
	"encoding/hex"
	"io/ioutil"
	"net"
	"os"
)

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

		file, err := os.OpenFile(path, os.O_RDWR, 0644)
		defer file.Close()

		if err != nil {
			logger.Fatalf("Unable to open file for DHT node data (%v).\n", err)
			return
		} else {
			nodeData, err := ioutil.ReadAll(file)

			if err != nil {
				logger.Printf("Unable to read existing DHT node file (%v). Creating a new one.\n", err)
				node = dht.NewLocalNode()
			} else {
				nodeDict, err := bencoding.Decode(nodeData)
				if err != nil {
					logger.Fatalf("%v\n", err)
					return
				}

				nodeDictAsDict, ok := nodeDict.(bencoding.Dict)
				if !ok {
					logger.Fatalf("Node data wasn't a dict?\n")
					return
				}

				node = dht.LocalNodeFromBencodingDict(nodeDictAsDict)
				logger.Printf("Loaded local node info from %v.\n", path)
			}
		}

		defer func() {
			// save node
			nodeData, err := bencoding.Encode(node)

			if err != nil {
				logger.Fatalf("Error encoding local node state: %v\n", err)
				return
			}

			logger.Printf("Saving LocalNode state to %v.\n", path)
			file.Truncate(0)
			file.WriteAt(nodeData, 0)
			file.Sync()

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
