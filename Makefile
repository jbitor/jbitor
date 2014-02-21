bin/jbitor:
	### Formatting...
	#
	go fmt github.com/jeremybanks/go-jbitor/bencoding
	go fmt github.com/jeremybanks/go-jbitor/dht
	go fmt github.com/jeremybanks/go-jbitor/utils
	go fmt github.com/jeremybanks/go-jbitor/torrent
	go fmt github.com/jeremybanks/go-jbitor/cli
	go fmt github.com/jeremybanks/go-jbitor/cli/jbitor
	#
	### Installing packages...
	#
	go install github.com/jeremybanks/go-jbitor/bencoding
	go install github.com/jeremybanks/go-jbitor/dht
	go install github.com/jeremybanks/go-jbitor/utils
	go install github.com/jeremybanks/go-jbitor/torrent
	go install github.com/jeremybanks/go-jbitor/cli
	go install github.com/jeremybanks/go-jbitor/cli/jbitor
	#
	### Testing...
	#
	go test -run=. -bench=NONE github.com/jeremybanks/go-jbitor/bencoding
	go test -run=. -bench=NONE github.com/jeremybanks/go-jbitor/dht
	go test -run=. -bench=NONE github.com/jeremybanks/go-jbitor/utils
	go test -run=. -bench=NONE github.com/jeremybanks/go-jbitor/torrent
	go test -run=. -bench=NONE github.com/jeremybanks/go-jbitor/cli
	go test -run=. -bench=NONE github.com/jeremybanks/go-jbitor/cli/jbitor
	#

bin/jbitorgtk:
    ### Formatting Packages...
    #
	go fmt github.com/jeremybanks/go-jbitor/gtkgui
	go fmt github.com/jeremybanks/go-jbitor/gtkgui/jbitorgtk
	#
    ### Installing Packages...
    #
	go install github.com/mattn/go-gtk/gtk
	go install github.com/jeremybanks/go-jbitor/gtkgui
	go install github.com/jeremybanks/go-jbitor/gtkgui/jbitorgtk
	#
	### Testing...
	go test github.com/jeremybanks/go-jbitor/gtkgui
	go test github.com/jeremybanks/go-jbitor/gtkgui/jbitorgtk
	#

PHONY:
