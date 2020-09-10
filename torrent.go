package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"path/filepath"
	"time"
)

const TORRENT_EXTENSION = ".torrent.json"

func CreateTorrent(name string, filePaths []string) (*Torrent, error) {
	if name == "" {
		name = filePaths[0]
	}

	torrent := &Torrent{
		Name:      name,
		Timestamp: time.Now().UTC().Unix(),
		Files:     make([]*File, 0),
	}

	for _, filePath := range filePaths {
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(fileContent)
		checksum := sha256.Sum256(fileContent)
		fileName := filepath.Base(filePath)
		pieces := make([]*Piece, 0)
		file := &File{
			Name:   fileName,
			Sha256: hex.EncodeToString(checksum[:]),
			Pieces: pieces,
		}

		for i := 1; i < math.MaxInt64; i++ {
			chunk := make([]byte, PIECE_LENGTH)

			_, err := reader.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			number := uint64(i)
			checksum := sha256.Sum256(chunk)
			piece := &Piece{
				Number: number,
				Sha256: hex.EncodeToString(checksum[:]),
			}

			file.Pieces = append(file.Pieces, piece)
		}

		torrent.Files = append(torrent.Files, file)
	}

	return torrent, nil
}

func LoadTorrent(filePath string) (*Torrent, error) {
	torrent := &Torrent{}

	message, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(message, torrent)
	if err != nil {
		return nil, err
	}

	return torrent, nil
}

func (torrent *Torrent) SaveTorrent() error {
	filename := fmt.Sprintf("%s%s", torrent.Name, TORRENT_EXTENSION)

	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}
