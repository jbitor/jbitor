bin/cli:
	### Formatting...
	#
	go fmt github.com/jeremybanks/go-distributed/bencoding
	go fmt github.com/jeremybanks/go-distributed/dht
	go fmt github.com/jeremybanks/go-distributed/torrentutils
	go fmt github.com/jeremybanks/go-distributed/cli
	#
	### Installing packages...
	#
	go install github.com/jeremybanks/go-distributed/bencoding
	go install github.com/jeremybanks/go-distributed/dht
	go install github.com/jeremybanks/go-distributed/torrentutils
	go install github.com/jeremybanks/go-distributed/cli
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/bencoding
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/dht
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/torrentutils
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/cli
	#

run: bin/cli PHONY
	@bin/cli torrent make test-torrents/hello/ > test-torrents/hello.torrent
	@echo
	@bin/cli dht helloworld tmp/dht-node.benc
	@echo
	@bin/cli json from-bencoding < tmp/dht-node.benc | python -m json.tool > tmp/dht-node.json

PHONY:
