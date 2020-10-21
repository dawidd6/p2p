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

type Hash struct {
	hash.Hash
}

func New() *Hash {
	return &Hash{sha256.New()}
}

func (hash *Hash) HexSum() string {
	return hex.EncodeToString(hash.Sum(nil))
}

func (hash *Hash) Verify(b []byte, h string) error {
	hash.Reset()
	hash.Write(b)

	if hash.HexSum() != h {
		return ChecksumMismatchError
	}

	return nil
}
