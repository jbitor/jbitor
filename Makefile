bin/cli:
	### Formatting...
	#
	go fmt bitbucket.org/jeremybanks/go-distributed/bencoding
	go fmt bitbucket.org/jeremybanks/go-distributed/dht
	go fmt bitbucket.org/jeremybanks/go-distributed/torrentutils
	go fmt bitbucket.org/jeremybanks/go-distributed/cli
	#
	### Installing packages...
	#
	go install bitbucket.org/jeremybanks/go-distributed/bencoding
	go install bitbucket.org/jeremybanks/go-distributed/dht
	go install bitbucket.org/jeremybanks/go-distributed/torrentutils
	go install bitbucket.org/jeremybanks/go-distributed/cli
	#
	### Testing...
	#
	go test -run=. -bench=NONE bitbucket.org/jeremybanks/go-distributed/bencoding
	go test -run=. -bench=NONE bitbucket.org/jeremybanks/go-distributed/dht
	go test -run=. -bench=NONE bitbucket.org/jeremybanks/go-distributed/torrentutils
	go test -run=. -bench=NONE bitbucket.org/jeremybanks/go-distributed/cli
	#

run: bin/cli PHONY
	@bin/cli torrent make test-torrents/hello/ > test-torrents/hello.torrent
	@echo
	@bin/cli json from-bencoding < test-torrents/hello.torrent | tee test-torrents/hello.torrent.json
	@echo
	@echo
	@bin/cli json to-bencoding < test-torrents/hello.torrent.json
	@echo
	@echo
	@bin/cli dht helloworld tmp/dht-node.benc

PHONY:
