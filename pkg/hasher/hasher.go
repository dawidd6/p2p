// Package hasher holds a hash.Hash convenient wrapper
package hasher

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
type Hasher struct {
	hash.Hash
}

// New returns a new instance of Hash
func New() *Hasher {
	return &Hasher{sha256.New()}
}

// HexSum calculates a checksum and returns a hex representation of it
func (hash *Hasher) HexSum() string {
	return hex.EncodeToString(hash.Sum(nil))
}

// Verify checks if given bytes' checksum is equal to the one provided
func (hash *Hasher) Verify(b []byte, h string) error {
	hash.Reset()
	hash.Write(b)

	if hash.HexSum() != h {
		return ChecksumMismatchError
	}

	return nil
}
