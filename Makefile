bin/jbitor:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/dht
	go fmt github.com/jbitor/jbitor/utils
	go fmt github.com/jbitor/torrent
	go fmt github.com/jbitor/cli
	go fmt github.com/jbitor/cli/jbitor
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/dht
	go install github.com/jbitor/jbitor/utils
	go install github.com/jbitor/torrent
	go install github.com/jbitor/cli
	go install github.com/jbitor/cli/jbitor
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/dht
	go test -run=. -bench=NONE github.com/jbitor/jbitor/utils
	go test -run=. -bench=NONE github.com/jbitor/torrent
	go test -run=. -bench=NONE github.com/jbitor/cli
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor
	#

PHONY:
