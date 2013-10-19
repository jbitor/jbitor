package main

import (
	"encoding/hex"
	"github.com/jeremybanks/go-distributed/bencoding"
	"github.com/jeremybanks/go-distributed/dht"
	"io/ioutil"
	"net"
	"os"
	"time"
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

		save := func() {
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
		}

		go func() {
			for {
				time.Sleep(10 * time.Second)
				save()
			}
		}()

		defer save()
	} else {
		local = dht.NewLocalNode()
	}

	terminate := make(chan bool)
	terminated := make(chan error)

	go local.Run(terminate, terminated)

	knownNode := dht.RemoteNodeFromAddress(net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 6881,
	})

	knownNode = local.AddOrGetRemoteNode(knownNode)

	logger.Printf("Hello, I am %v.\n", local)

	infoHash, _ := hex.DecodeString("5497a53543938b77ef660939d3b32e02be7bc213")
	logger.Printf("Trying to GetPeers for %v.\n", infoHash)
	peersResult, nodesResult, errChan := local.GetPeers(knownNode, string(infoHash))

	select {
	case result := <-peersResult:
		logger.Printf("GetPeers got peers: %v\n", result)
	case result := <-nodesResult:
		logger.Printf("GetPeers got nodes: %v\n", result)
	case err := <-errChan:
		logger.Printf("GetPeers had error: %v\n", err)
	}

	if err := <-terminated; err != nil {
		logger.Fatalf("Error in running LocalNode: %v\n", err)
	}
}
