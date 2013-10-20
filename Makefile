bin/distributed:
	### Formatting...
	#
	go fmt github.com/jeremybanks/go-distributed/bencoding
	go fmt github.com/jeremybanks/go-distributed/dht
	go fmt github.com/jeremybanks/go-distributed/torrentutils
	go fmt github.com/jeremybanks/go-distributed/torrent
	go fmt github.com/jeremybanks/go-distributed/cli
	go fmt github.com/jeremybanks/go-distributed/cli/distributed
	#
	### Installing packages...
	#
	go install github.com/jeremybanks/go-distributed/bencoding
	go install github.com/jeremybanks/go-distributed/dht
	go install github.com/jeremybanks/go-distributed/torrentutils
	go install github.com/jeremybanks/go-distributed/torrent
	go install github.com/jeremybanks/go-distributed/cli
	go install github.com/jeremybanks/go-distributed/cli/distributed
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/bencoding
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/dht
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/torrentutils
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/torrent
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/cli
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/cli/distributed
	#

hello: bin/distributed PHONY
	@bin/distributed torrent make test-torrents/hello/ > test-torrents/hello.torrent
	@echo
	@bin/distributed dht helloworld tmp/dht-node.benc
	@echo
	@bin/distributed json from-bencoding < tmp/dht-node.benc | python -m json.tool > tmp/dht-node.json

PHONY:
