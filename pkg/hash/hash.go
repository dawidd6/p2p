package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"

	"github.com/dawidd6/p2p/pkg/errors"
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

	_, err := hash.Write(b)
	if err != nil {
		return err
	}

	if hash.HexSum() != h {
		return errors.ChecksumMismatchError
	}

	return nil
}
