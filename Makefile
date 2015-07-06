# TODO: This is disgusting. Kill it with fire.

all:
	export GOPATH=$$PWD
	echo $$GOPATH
	### Formatting...
	#
	goimports -w src/github.com/jbitor/cli/loggerconfig
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/bittorrent/dht
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/webclient
	goimports -w src/github.com/jbitor/cli/jbitor-web
	goimports -w src/github.com/jbitor/cli/jbitor-get-peers
	goimports -w src/github.com/jbitor/cli/jbitor-json
	goimports -w src/github.com/jbitor/cli/jbitor-create
	goimports -w src/github.com/jbitor/cli/jbitor-get-info
	#
	### Installing packages...
	#
	go install github.com/op/go-logging
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/bittorrent/dht
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/webclient
	go install github.com/jbitor/cli/jbitor-web
	go install github.com/jbitor/cli/jbitor-get-peers
	go install github.com/jbitor/cli/jbitor-json
	go install github.com/jbitor/cli/jbitor-create
	go install github.com/jbitor/cli/jbitor-get-info

bin/jbitor-web:
	export GOPATH=$$PWD
	### Formatting...
	#
	goimports -w src/github.com/jbitor/cli/loggerconfig
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/bittorrent/dht
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/webclient
	goimports -w src/github.com/jbitor/cli/jbitor-web
	#
	### Installing packages...
	#
	go install github.com/op/go-logging
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/bittorrent/dht
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/webclient
	go install github.com/jbitor/cli/jbitor-web
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/bittorrent/dht
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/webclient
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-web
	#

bin/jbitor-get-peers:
	export GOPATH=$$PWD
	### Formatting...
	#
	goimports -w src/github.com/jbitor/cli/loggerconfig
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/bittorrent/dht
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/cli/jbitor-get-peers
	#
	### Installing packages...
	#
	go install github.com/op/go-logging
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/bittorrent/dht
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/cli/jbitor-get-peers
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/bittorrent/dht
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-get-peers
	#

bin/jbitor-json:
	export GOPATH=$$PWD
	### Formatting...
	#
	goimports -w src/github.com/jbitor/cli/loggerconfig
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/cli/jbitor-json
	#
	### Installing packages...
	#
	go install github.com/op/go-logging
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/cli/jbitor-json
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-json
	#

bin/jbitor-create:
	export GOPATH=$$PWD
	### Formatting...
	#
	goimports -w src/github.com/jbitor/cli/loggerconfig
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/cli/jbitor-create
	#
	### Installing packages...
	#
	go install github.com/op/go-logging
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

bin/jbitor-get-info:
	export GOPATH=$$PWD
	### Formatting...
	#
	goimports -w src/github.com/jbitor/cli/loggerconfig
	goimports -w src/github.com/jbitor/bencoding
	goimports -w src/github.com/jbitor/bittorrent
	goimports -w src/github.com/jbitor/bittorrent/dht
	goimports -w src/github.com/jbitor/cli/jbitor-get-info
	#
	### Installing packages...
	#
	go install github.com/op/go-logging
	go install github.com/jbitor/bencoding
	go install github.com/jbitor/bittorrent
	go install github.com/jbitor/bittorrent/dht
	go install github.com/jbitor/cli/jbitor-get-info
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jbitor/bencoding
	go test -run=. -bench=NONE github.com/jbitor/bittorrent
	go test -run=. -bench=NONE github.com/jbitor/bittorrent/dht
	go test -run=. -bench=NONE github.com/jbitor/cli/jbitor-get-info
	#

PHONY:
