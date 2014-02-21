## jbitor

web: https://github.com/jbitor/jbitor  
git: https://github.com/jbitor/jbitor.git

<a href="https://travis-ci.org/jbitor/jbitor/branches">
<img src="https://travis-ci.org/jbitor/jbitor.png?branch=master"
     alt="master branch build status on Travis CI" />
</a>

---

### Implemented Functionality

#### `./jbitorgtk`

A GTK GUI which:

- displays the DHT connection state.

#### `./jbitor torrent create TARGET > TARGET.torrent`

Generates a torrent file for the target file or directory. The piece size is currently hardcoded.

#### `./jbitor dht connect`

Maintains a (client-only) connection to the mainline BitTorrent DHT until terminated.

#### `./jbitor dht get-peers INFOHASH`

Uses the DHT to find BitTorrent peers for the torrent with the given infohash, and outputs their connection info.

#### `./jbitor json from-bencoding < FOO.benc > FOO.json`  <br />  `bin/cli json to-bencoding < FOO.json > FOO.benc`

Used to convert between equivalent JSON and Bencoding data. Data that does not have an equivalent representation in the other format will cause an error.

---

Use `./doc` to run `godoc` and open a browser pointing at `jbitor`.

---

For the GUI you may need to install something like:

    PKG_CONFIG_PATH=/opt/X11/lib/pkgconfig brew install gtk+
    PKG_CONFIG_PATH=/opt/X11/lib/pkgconfig go get github.com/mattn/go-gtk/gtk

---

Copyright 2013-2014 Jeremy Banks <<j@jeremybanks.ca>>.

Currently released under the GPLv3. Maybe BSD later.
