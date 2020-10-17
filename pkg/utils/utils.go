package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Sha256Sum computes SHA-256 sum of given bytes and returns hex representation of it.
func Sha256Sum(b []byte) string {
	checksum := sha256.Sum256(b)
	return hex.EncodeToString(checksum[:])
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	if info == nil {
		return false
	}

	return true
}

// ReadFilePiece reads only a specified portion of file.
func ReadFilePiece(filePath string, pieceSize, pieceNumber uint64) ([]byte, error) {
	b := make([]byte, pieceSize)

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0664)
	if err != nil {
		return nil, err
	}

	n, err := file.ReadAt(b, int64(pieceSize*pieceNumber))
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b[:n], file.Close()
}

// WriteFilePiece writes only a specified portion of file.
func WriteFilePiece(filePath string, pieceNumber uint64, piece []byte) error {
	size := uint64(len(piece))

	file, err := os.OpenFile(filePath, os.O_WRONLY, 0664)
	if err != nil {
		return err
	}

	_, err = file.WriteAt(piece, int64(size*pieceNumber))
	if err != nil {
		return err
	}

	return file.Close()
}

// AllocateZeroedFile creates a new file filled with zeroes.
func AllocateZeroedFile(filePath string, fileSize uint64) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = file.Truncate(int64(fileSize))
	if err != nil {
		return err
	}

	return file.Close()
}

// DoInDirectory changes current working directory to `dir`,
// executes `f` functions and returns to previous working directory.
func DoInDirectory(dir string, f func() error) error {
	// Just execute given function if `dir` is empty.
	if dir == "" {
		return f()
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	err = f()
	if err != nil {
		return err
	}

	return os.Chdir(cwd)
}
