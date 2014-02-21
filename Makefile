bin/jbitor:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/jbitor/dht
	go fmt github.com/jbitor/jbitor/utils
	go fmt github.com/jbitor/jbitor/torrent
	go fmt github.com/jbitor/jbitor/cli
	go fmt github.com/jbitor/jbitor/cli/jbitor
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/jbitor/dht
	go install github.com/jbitor/jbitor/utils
	go install github.com/jbitor/jbitor/torrent
	go install github.com/jbitor/jbitor/cli
	go install github.com/jbitor/jbitor/cli/jbitor
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/jbitor/dht
	go test -run=. -bench=NONE github.com/jbitor/jbitor/utils
	go test -run=. -bench=NONE github.com/jbitor/jbitor/torrent
	go test -run=. -bench=NONE github.com/jbitor/jbitor/cli
	go test -run=. -bench=NONE github.com/jbitor/jbitor/cli/jbitor
	#

PHONY:
