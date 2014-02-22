## jbitor

web: https://github.com/jbitor/jbitor  
git: https://github.com/jbitor/jbitor.git

<a href="https://travis-ci.org/jbitor/jbitor/branches">
<img src="https://travis-ci.org/jbitor/jbitor.png?branch=master"
     alt="master branch build status on Travis CI" />
</a>

---

### Provided Commands

#### `./jbitor-web`

Serves a web GUI at <http://127.0.0.1:47935/>.
Maintains a (client-only) connection to the mainline BitTorrent DHT until terminated.

#### `./jbitor-create TARGET > TARGET.torrent`

Generates a torrent file for the target file or directory. The piece size is currently hardcoded.

#### `./jbitor-get-peers INFOHASH`

Uses the DHT to find BitTorrent peers for the torrent with the given infohash, and outputs their connection info.

#### `./jbitor-json from-bencoding < FOO.benc > FOO.json`  <br />  `./jbitor-json to-bencoding < FOO.json > FOO.benc`

Used to convert between equivalent JSON and Bencoding data. Data that does not have an equivalent representation in the other format will cause an error.

---

`./test` tests and builds everything, and runs some simple things.

`./doc` runs `godoc` and opens a browser pointing viewing `jbitor`'s docs.
(Requires `godoc`, you may need to `go get code.google.com/p/go.tools/cmd/godoc`.)

---

Copyright 2013-2014 Jeremy Banks <<j@jeremybanks.ca>>.

Currently released under the GPLv3. Maybe BSD later.
