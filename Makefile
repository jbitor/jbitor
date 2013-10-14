run: bin/cli PHONY
	@bin/cli torrent make test-torrents/hello/ > test-torrents/hello.torrent
	@echo
	@bin/cli json from-bencoding < test-torrents/hello.torrent
	@echo
	@echo
	@bin/cli dht helloworld tmp/dht-node.benc

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

PHONY:
