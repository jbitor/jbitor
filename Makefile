bin/jbitor-web:
	### Formatting...
	#
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/dht
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/webclient
	goimports -w src/github.com/jbitor/cli/jbitor-web
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
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/dht
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/cli/jbitor-get-peers
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
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/cli/jbitor-json
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
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/cli/jbitor-create
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
