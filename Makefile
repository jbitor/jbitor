run: PHONY
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
	### Running...
	#
	bin/cli torrent make . > tmp/test.torrent
	#
	bin/cli dht helloworld
	#

PHONY:
