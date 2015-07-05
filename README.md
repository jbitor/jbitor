## jbitor

BitTorrent stuff. Unlikely to ever be a full client.

web: https://github.com/jbitor/jbitor  
git: https://github.com/jbitor/jbitor.git

##### godoc

- [bencoding](https://godoc.org/github.com/jbitor/bencoding)
- [bittorrent](https://godoc.org/github.com/jbitor/bittorrent)
- [cli](https://godoc.org/github.com/jbitor/cli)
- [dht](https://godoc.org/github.com/jbitor/dht)
- [webclient](https://godoc.org/github.com/jbitor/webclient)

---

### Provided Commands

#### `./jbitor-web`

Serves a web GUI at <http://127.0.0.1:8080/>.
Maintains a (client-only) connection to the mainline BitTorrent DHT until terminated.

#### `./jbitor-create TARGET > TARGET.torrent`

Generates a torrent file for the target file or directory. The piece size is currently hardcoded.

#### `./jbitor-get-peers INFOHASH`

Uses the DHT to find BitTorrent peers for the torrent with the given infohash, and outputs their connection info.

#### `./jbitor-json from-bencoding < FOO.benc > FOO.json`  <br />  `./jbitor-json to-bencoding < FOO.json > FOO.benc`

Used to convert between equivalent JSON and Bencoding data. Data that does not have an equivalent representation in the other format will cause an error.

---

`./test` tests and builds everything, and runs some simple things.

`./example` downloads an Ubuntu torrent file from peers found from the DHT.

`./doc` runs `godoc` and opens a browser pointing viewing `jbitor`'s docs.

---

Copyright 2013-2015 Jeremy Banks <<j@jeremybanks.ca>>.
