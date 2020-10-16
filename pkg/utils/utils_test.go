package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const filePath = "test_file.txt"

func TestSha256(t *testing.T) {
	expected := "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"
	actual := Sha256Sum([]byte("hello\n"))

	assert.Equal(t, expected, actual)
}

func TestAllocateZeroedFile(t *testing.T) {
	size := int64(1024) // 1kB
	err := AllocateZeroedFile(filePath, size)

	assert.NoError(t, err)
	assert.FileExists(t, filePath)
}

func TestWriteFilePiece(t *testing.T) {
	err := WriteFilePiece(filePath, 0, []byte("1 line\n"))

	assert.NoError(t, err)

	err = WriteFilePiece(filePath, 2, []byte("2 line\n"))

	assert.NoError(t, err)
}

// dd bs=1 skip=16 count=16 if=test/fixtures/go.sum 2>/dev/null | sha256sum
func TestReadFilePiece(t *testing.T) {
	piece, err := ReadFilePiece(filePath, 7, 0)
	expected := "09b78b56926d5c5838ea538cc5bbe81f68195f4f7e39d07effb942fca5a33529"
	actual := Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, expected, actual)

	piece, err = ReadFilePiece(filePath, 7, 1)
	expected = "837885c8f8091aeaeb9ec3c3f85a6ff470a415e610b8ba3e49f9b33c9cf9d619"
	actual = Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, expected, actual)

	piece, err = ReadFilePiece(filePath, 7, 2)
	expected = "5cf363398ada33352093ed6598d40602871998ed57f0836d5a0835d462116b76"
	actual = Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, expected, actual)

	t.Cleanup(func() {
		err := os.Remove(filePath)
		assert.NoError(t, err)
		assert.NoFileExists(t, filePath)
	})
}

func TestDoInDirectory(t *testing.T) {
	err := DoInDirectory("", func() error {
		return nil
	})

	assert.NoError(t, err)
}
