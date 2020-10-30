package piece_test

import (
	"os"
	"testing"

	"github.com/dawidd6/p2p/pkg/piece"
	"github.com/stretchr/testify/assert"
)

func TestOffset(t *testing.T) {
	expected := 2048
	got := piece.Offset(1024, 2)
	assert.EqualValues(t, expected, got)
}

func TestWriteRead(t *testing.T) {
	file, err := os.Create("test.txt")
	assert.NoError(t, err)

	t.Cleanup(func() {
		err = os.Remove(file.Name())
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)
	})

	expected := []byte("hello there")
	offset := piece.Offset(1024, 2)

	err = piece.Write(file, offset, expected)
	assert.NoError(t, err)

	got, err := piece.Read(file, 1024, offset)
	assert.NoError(t, err)
	assert.EqualValues(t, len(expected), len(got))
	assert.EqualValues(t, expected, got)
}
