package torrent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	weakrand "math/rand"
)

// A BTID is a 20-byte string, used to identify Torrents (their "info hash")
// and DHT nodes (their node ID).
type BTID string

// Verify() returns an error if the BTID is of the wrong length.
// (Unknown/empty values are not valid.)
func (id BTID) Verify() (err error) {
	if len(id) != 20 {
		return errors.New(fmt.Sprintf("BTID has length %v, should be 20.", len(id)))
	} else {
		return nil
	}
}

// String() represents a BTID in hexadecimal.
func (id BTID) String() (str string) {
	if id != "" {
		return hex.EncodeToString([]byte(id))
	} else {
		return "unknown"
	}
}

// BTIDFromHex() returns a valid BTID for a given hex string, or an error.
func BTIDFromHex(hexStr string) (id BTID, err error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	id = BTID(bytes)

	err = id.Verify()
	if err != nil {
		return "", err
	}

	return id, nil
}

// SecureRandomBTID returns a BTID instantiated using a cryptographically
// secure random number generator, or an error if it was not possible to
// generate sufficient secure bytes.
func SecureRandomBTID() (id BTID, err error) {
	bytes := new([20]byte)
	n, err := rand.Read(bytes[:])

	if err != nil {
		return "", err
	}

	if n < 20 {
		return "", errors.New("Failed to generate 20 secure random bytes for BTID")
	}

	return BTID(bytes[:]), nil
}

// WeakRandomBTID returns a BTID instantiated using the standard random
// number generator. (Seeding the RNG is the responsibility of the user.)
func WeakRandomBTID() (id BTID) {
	bytes := new([20]byte)

	for i := range bytes {
		bytes[i] = byte(weakrand.Int() & 0xFF)
	}

	return BTID(bytes[:])
}
