bin/jbitor:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/dht
	go fmt github.com/jbitor/bittorrent
	go fmt github.com/jbitor/webclient
	go fmt github.com/jbitor/cli
	go fmt github.com/jbitor/cli/jbitor
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/dht
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/webclient
	go install github.com/jbitor/cli
	go install github.com/jbitor/cli/jbitor
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/dht
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/webclient
	go test -run=. -bench=NONE github.com/jbitor/cli
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor
	#

PHONY:
