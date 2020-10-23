// Package hash holds a hash.Hash convenient wrapper
package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
)

var (
	ChecksumMismatchError = errors.New("checksum mismatch")
)

// Hash wraps hash.Hash type
type Hash struct {
	hash.Hash
}

// New returns a new instance of Hash
func New() *Hash {
	return &Hash{sha256.New()}
}

// HexSum calculates a checksum and returns a hex representation of it
func (hash *Hash) HexSum() string {
	return hex.EncodeToString(hash.Sum(nil))
}

// Verify checks if given bytes' checksum is equal to the one provided
func (hash *Hash) Verify(b []byte, h string) error {
	hash.Reset()
	hash.Write(b)

	if hash.HexSum() != h {
		return ChecksumMismatchError
	}

	return nil
}
