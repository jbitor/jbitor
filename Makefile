bin/jbitor-web:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/dht
	go fmt github.com/jbitor/bittorrent
	go fmt github.com/jbitor/webclient
	go fmt github.com/jbitor/cli/jbitor-web
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/dht
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/webclient
	go install github.com/jbitor/cli/jbitor-web
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/dht
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/webclient
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-web
	#

bin/jbitor-get-peers:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/dht
	go fmt github.com/jbitor/bittorrent
	go fmt github.com/jbitor/cli/jbitor-get-peers
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/dht
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/cli/jbitor-get-peers
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/dht
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-get-peers
	#

bin/jbitor-json:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/cli/jbitor-json
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/cli/jbitor-json
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-json
	#

bin/jbitor-create:
	### Formatting...
	#
	go fmt github.com/jbitor/bencoding
	go fmt github.com/jbitor/bittorrent
	go fmt github.com/jbitor/cli/jbitor-create
	#
	### Installing packages...
	#
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/cli/jbitor-create
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-create
	#

PHONY:
