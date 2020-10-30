package hasher_test

import (
	"testing"

	"github.com/dawidd6/p2p/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestHasher(t *testing.T) {
	data := []byte("hello there")
	hash := hasher.New()
	n, err := hash.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	checksum := hash.HexSum()
	err = hash.Verify(data, checksum)
	assert.NoError(t, err)
}
