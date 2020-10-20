package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type Hash struct {
	hash.Hash
}

func New() *Hash {
	return &Hash{sha256.New()}
}

func (hash *Hash) HexSum() string {
	return hex.EncodeToString(hash.Sum(nil))
}

// Compute computes SHA-256 sum of given bytes and returns hex representation of it.
func Compute(b []byte) string {
	checksum := sha256.Sum256(b)
	return hex.EncodeToString(checksum[:])
}
