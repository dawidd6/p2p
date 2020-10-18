package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CreateDirectory(dir string) error {
	return os.MkdirAll(dir, 0775)
}

func ReadAll(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}

func Open(dir, fileName string) (*os.File, error) {
	filePath := filepath.Join(dir, fileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	return file, nil
}
