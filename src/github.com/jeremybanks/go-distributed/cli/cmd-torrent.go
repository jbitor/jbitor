package main

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/jeremybanks/go-distributed/bencoding"
	//	"github.com/jeremybanks/go-distributed/dht"
	"github.com/jeremybanks/go-distributed/torrentutils"
	"os"
)

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
	case "find_peers":
		cmdTorrentFindPeers(subcommandArgs)
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
		logger.Fatalf("Error generating torrent: %v\n", err)
		return
	}

	infoData, err := bencoding.Encode(infoDict)
	if err != nil {
		logger.Fatalf("Error encoding torrent infodict (for hashing): %v\n", err)
		return
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
		logger.Fatalf("Error encoding torrent data: %v\n", err)
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

func cmdTorrentFindPeers(args []string) {
	if len(args) != 1 {
		logger.Fatalf("Usage: %v torrent find_peers INFOHASH\n", os.Args[0])
		return
	}

	hexInfoHash := args[0]

	infoHash, err := hex.DecodeString(hexInfoHash)

	if err != nil {
		logger.Fatalf("Specified infohash is not a valid hex byte string [%v].\n", err)
		return
	}

	if len(infoHash) != 20 {
		logger.Fatalf("Infohash length is wrong (%v instead of 20).\n", len(infoHash))
		return
	}

	//	dhtClient := dht.OpenClient(".dht-peer")
	//	defer dhtClient.Close()

}
