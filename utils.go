package main

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

// ReadFilePiece reads only a specified portion of file.
func ReadFilePiece(filePath string, size, number int) ([]byte, error) {
	b := make([]byte, size)

	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	n, err := file.ReadAt(b, int64(size*number))
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b[:n], file.Close()
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
