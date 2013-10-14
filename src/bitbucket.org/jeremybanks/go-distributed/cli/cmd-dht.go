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
	var local *dht.LocalNode

	if len(args) > 0 {
		path := args[0]

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		defer file.Close()

		if err != nil {
			logger.Fatalf("Unable to open file for DHT node data (%v).\n", err)
			return
		} else {
			nodeData, err := ioutil.ReadAll(file)

			if err != nil {
				logger.Printf("Unable to read existing DHT node file (%v). Creating a new one.\n", err)
				local = dht.NewLocalNode()
			} else if len(nodeData) == 0 {
				logger.Printf("Existing DHT node file was empty. Creating a new one.\n")
				local = dht.NewLocalNode()
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

				local = dht.LocalNodeFromBencodingDict(nodeDictAsDict)
				logger.Printf("Loaded local node info from %v.\n", path)
			}
		}

		defer func() {
			// save LocalNode
			nodeData, err := bencoding.Encode(local)

			if err != nil {
				logger.Fatalf("Error encoding local node state: %v\n", err)
				return
			}

			file.Truncate(0)
			file.WriteAt(nodeData, 0)
			file.Sync()

			if err != nil {
				logger.Fatalf("Error writing local node state: %v\n", err)
				return
			}

			logger.Printf("Saved LocalNode state to %v.\n", path)
		}()
	} else {
		local = dht.NewLocalNode()
	}

	terminated := make(chan error)
	go local.Run(terminated)

	knownNode := dht.RemoteNodeFromAddress(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6881,
	})

	knownNode = local.AddOrGetRemoteNode(knownNode)

	logger.Printf("Hello, I am %v.\n", local)
	logger.Printf("I know of %v.\n", local.Nodes)

	//	logger.Printf("I am attempting to ping a DHT node at localhost:6881.\n")
	//	pingResult, pingErr := local.SendPing(knownNode)
	//
	//	logger.Printf("SendPing initiated\n")

	nodeId, _ := hex.DecodeString("b7271d0b5577918ee92b1b5378d89e56ad08ba80")
	logger.Printf("Attempting to FindNode(%v)...", dht.NodeId(nodeId))

	findResult, findErr := local.FindNode(knownNode, dht.NodeId(nodeId))

	for len(local.OutstandingQueries) > 0 {
		select {
		case result := <-findResult:
			logger.Printf("FindNode result: %v\n", result)
		case result := <-findErr:
			logger.Printf("FindNode error: %v\n", result)
		}
	}

	logger.Printf("I know of %v.\n", local.Nodes)
}
