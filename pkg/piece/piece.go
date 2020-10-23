// Package piece provides convenient functions for piece reading/writing
package piece

import "io"

// Size is a default piece size
const Size = 256 * 1024 // 256 kB

// Offset calculates a file offset where piece should be read at
func Offset(pieceSize, pieceNumber int64) int64 {
	return pieceSize * pieceNumber
}

// Read retrieves a specified piece from file
func Read(reader io.ReaderAt, pieceSize, pieceOffset int64) ([]byte, error) {
	b := make([]byte, pieceSize)

	n, err := reader.ReadAt(b, pieceOffset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b[:n], nil
}

// Write saves a specified piece to file
func Write(writer io.WriterAt, pieceOffset int64, pieceData []byte) error {
	_, err := writer.WriteAt(pieceData, pieceOffset)
	if err != nil {
		return err
	}

	return nil
}
