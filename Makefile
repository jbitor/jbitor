bin/distributed:
	### Formatting...
	#
	go fmt github.com/jeremybanks/go-distributed/bencoding
	go fmt github.com/jeremybanks/go-distributed/dht
	go fmt github.com/jeremybanks/go-distributed/torrentutils
	go fmt github.com/jeremybanks/go-distributed/torrent
	go fmt github.com/jeremybanks/go-distributed/cli
	go fmt github.com/jeremybanks/go-distributed/cli/distributed
	#
	### Installing packages...
	#
	go install github.com/jeremybanks/go-distributed/bencoding
	go install github.com/jeremybanks/go-distributed/dht
	go install github.com/jeremybanks/go-distributed/torrentutils
	go install github.com/jeremybanks/go-distributed/torrent
	go install github.com/jeremybanks/go-distributed/cli
	go install github.com/jeremybanks/go-distributed/cli/distributed
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/bencoding
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/dht
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/torrentutils
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/torrent
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/cli
	go test -run=. -bench=NONE github.com/jeremybanks/go-distributed/cli/distributed
	#

bin/distributedgtk:
    ### Formatting Packages...
    #
	go fmt github.com/jeremybanks/go-distributed/gui
	go fmt github.com/jeremybanks/go-distributed/gui/distributedgtk
	#
    ### Installing Packages...
    #
	go install github.com/mattn/go-gtk/gtk
	go install github.com/jeremybanks/go-distributed/gui
	go install github.com/jeremybanks/go-distributed/gui/distributedgtk
	#
	### Testing...
	go test github.com/jeremybanks/go-distributed/gui
	go test github.com/jeremybanks/go-distributed/gui/distributedgtk
	#

PHONY:
