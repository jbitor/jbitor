## go-distributed: peer-to-peer systems in Go

web: https://github.com/jeremybanks/go-distributed  
git: https://github.com/jeremybanks/go-distributed.git

---

### Implemented Functionality

#### `bin/cli torrent create TARGET > TARGET.torrent`

Generates a torrent file for the target file or directory. The piece size is currently hardcoded.

#### `bin/cli json from-bencoding < FOO.benc > FOO.json`  <br />  `bin/cli json to-bencoding < FOO.json > FOO.benc`

Used to convert between equvialent JSON and Bencoding data. Data that does not have an equivalent representation in the other format will cause an error.

#### `bin/cli dht hellworld PATH.benc`

Uses an existing node at `localhost:6118` to establish a mainline DHT peer list, persisted at the specified path.

---

Copyright 2013 Jeremy Banks <<jeremy@jeremybanks.ca>>
