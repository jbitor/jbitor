## go-distributed: peer-to-peer systems in Go

web: https://github.com/jeremybanks/go-distributed  
git: https://github.com/jeremybanks/go-distributed.git

---

### Implemented Functionality

#### `./run torrent create TARGET > TARGET.torrent`

Generates a torrent file for the target file or directory. The piece size is currently hardcoded.

#### `/run dht connect PATH.benc`

Maintains a (client-only) connection to the mainline BitTorrent DHT until terminated. State will be persisted in a bencoded file at the specified path. This will maintain a large list of healthy nodes, bootstrapped from the common bootstrap nodes.

#### `/run dht get-peers INFOHASH`

Uses the DHT to find BitTorrent peers for the torrent with the given infohash, and outputs their connection info.

#### `/run json from-bencoding < FOO.benc > FOO.json`  <br />  `bin/cli json to-bencoding < FOO.json > FOO.benc`

Used to convert between equivalent JSON and Bencoding data. DaFindta that does not have an equivalent representation in the other format will cause an error.

---

Use `./doc` to run `godoc` and open a browser pointing at `go-distributed`.

---

Copyright 2013 Jeremy Banks <<jeremy@jeremybanks.ca>>.

Currently released under the GPLv3. Maybe BSD later.
