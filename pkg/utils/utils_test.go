package utils

import (
	"io/ioutil"
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
	size := uint64(7 * 4) // 1kB
	err := AllocateZeroedFile(filePath, size)

	assert.NoError(t, err)
	assert.FileExists(t, filePath)

	b, err := ioutil.ReadFile(filePath)

	assert.NoError(t, err)
	assert.EqualValues(t, size, len(b))
}

func TestWriteFilePiece(t *testing.T) {
	err := WriteFilePiece(filePath, 0, []byte("1 line\n"))

	assert.NoError(t, err)

	err = WriteFilePiece(filePath, 2, []byte("2 line\n"))

	assert.NoError(t, err)

	err = WriteFilePiece(filePath, 3, []byte("3 line\n"))

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
	assert.Equal(t, string(piece), "1 line\n")

	piece, err = ReadFilePiece(filePath, 7, 1)
	expected = "837885c8f8091aeaeb9ec3c3f85a6ff470a415e610b8ba3e49f9b33c9cf9d619"
	actual = Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, expected, actual)
	assert.Equal(t, piece, []byte{0, 0, 0, 0, 0, 0, 0})

	piece, err = ReadFilePiece(filePath, 7, 2)
	expected = "5cf363398ada33352093ed6598d40602871998ed57f0836d5a0835d462116b76"
	actual = Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, expected, actual)
	assert.Equal(t, string(piece), "2 line\n")

	piece, err = ReadFilePiece(filePath, 7, 3)
	expected = "5935306a339fd86a6031afe0ec31cde2020f0cfc6c658c0716a73c6e79f1730d"
	actual = Sha256Sum(piece)

	assert.NoError(t, err)
	assert.NotNil(t, piece)
	assert.Equal(t, expected, actual)
	assert.Equal(t, string(piece), "3 line\n")

}

func TestReadFilePieces(t *testing.T) {
	pieces, err := ReadFilePieces(filePath, 7)

	assert.NoError(t, err)
	assert.NotNil(t, pieces)

	for i, piece := range pieces {
		err := WriteFilePiece(filePath, uint64(i), piece)
		assert.NoError(t, err)
	}
}

func TestDoInDirectory(t *testing.T) {
	err := DoInDirectory("", func() error {
		return nil
	})

	assert.NoError(t, err)

	t.Cleanup(func() {
		err := os.Remove(filePath)
		assert.NoError(t, err)
		assert.NoFileExists(t, filePath)
	})
}

func TestIntegration(t *testing.T) {
	originalFilePath := "test.orig"
	newFilePath := "test"
	data := []byte("1234567890\nabcdef\n")
	pieceSize := 5 // TODO this number needs to be a divide of len(data) !!!
	t.Log(len(data))

	err := ioutil.WriteFile(originalFilePath, data, os.ModePerm)
	assert.NoError(t, err)
	assert.FileExists(t, originalFilePath)

	err = AllocateZeroedFile(newFilePath, uint64(len(data)))
	assert.NoError(t, err)
	assert.FileExists(t, newFilePath)

	for i := 0; i < len(data); i++ {
		piece, err := ReadFilePiece(originalFilePath, uint64(pieceSize), uint64(i))
		assert.NoError(t, err)
		assert.NotNil(t, piece)

		err = WriteFilePiece(newFilePath, uint64(i), piece)
		assert.NoError(t, err)
	}

	originalFileHash, err := GetFileHash(originalFilePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, originalFileHash)

	newFileHash, err := GetFileHash(newFilePath)
	assert.NoError(t, err)
	assert.NotEmpty(t, newFileHash)

	assert.Equal(t, originalFileHash, newFileHash)

	t.Cleanup(func() {
		err := os.Remove(originalFilePath)
		assert.NoError(t, err)
		assert.NoFileExists(t, filePath)

		err = os.Remove(newFilePath)
		assert.NoError(t, err)
		assert.NoFileExists(t, filePath)
	})
}
