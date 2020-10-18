package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Sha256Sum computes SHA-256 sum of given bytes and returns hex representation of it.
func Sha256Sum(b []byte) string {
	checksum := sha256.Sum256(b)
	return hex.EncodeToString(checksum[:])
}

func CreateDir(dir string) error {
	return os.MkdirAll(dir, 0775)
}

func ReadFile(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func OpenFile(dir, fileName string) (*os.File, error) {
	filePath := filepath.Join(dir, fileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ReadFilePiece(reader io.ReaderAt, pieceSize, pieceNumber uint64) ([]byte, error) {
	b := make([]byte, pieceSize)
	n, err := reader.ReadAt(b, int64(pieceSize*pieceNumber))
	if err != nil && err != io.EOF {
		return nil, err
	}
	return b[:n], nil
}

func WriteFilePiece(writer io.WriterAt, pieceNumber uint64, piece []byte) error {
	size := uint64(len(piece))
	if size == 0 {
		return nil
	}
	_, err := writer.WriteAt(piece, int64(size*pieceNumber))
	if err != nil {
		return err
	}
	return nil
}

// AllocateZeroedFile creates a new file filled with zeroes.
func AllocateZeroedFile(filePath string, fileSize uint64) error {
	info, err := os.Stat(filePath)
	if info != nil {
		return os.ErrExist
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0664)
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
