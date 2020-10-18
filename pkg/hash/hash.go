package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// Compute computes SHA-256 sum of given bytes and returns hex representation of it.
func Compute(b []byte) string {
	checksum := sha256.Sum256(b)
	return hex.EncodeToString(checksum[:])
}
