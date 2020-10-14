package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256(t *testing.T) {
	want := "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"
	got := Sha256Sum([]byte("hello\n"))

	assert.Equal(t, got, want)
}

// dd bs=1 skip=16 count=16 if=test/fixtures/go.sum 2>/dev/null
func TestReadFilePiece(t *testing.T) {
	piece, err := ReadFilePiece("../../test/fixtures/go.sum", 16, 0)
	want := "442f9ad4aedc52dca7bf9219c7fb48a4906afb8877b21120cf8d68281704dfc8"
	got := Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, got, want)

	piece, err = ReadFilePiece("../../test/fixtures/go.sum", 16, 10)
	want = "019c4df04108da14cc6d0f3f628c7935c7e423c4b3839e559266f9178a74c62e"
	got = Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, got, want)
}

func TestDoInDirectory(t *testing.T) {
	err := DoInDirectory("", func() error {
		return nil
	})

	assert.NoError(t, err)
}
