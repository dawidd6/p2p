// Package torrent provides convenient functions for Torrent type
package torrent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dawidd6/p2p/pkg/hasher"
	"github.com/dawidd6/p2p/pkg/piece"
)

const FileExtension = "torrent.json"

var (
	WrongExtensionError = errors.New("wrong torrent file extension")
)

// Create makes a new torrent from file
func Create(file *os.File, pieceSize int64, trackerAddr string) (*Torrent, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileHasher := hasher.New()
	pieceHasher := hasher.New()
	pieceHashes := make([]string, 0)

	for {
		pieceData := make([]byte, pieceSize)

		n, err := file.Read(pieceData)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		pieceData = pieceData[:n]

		n, err = fileHasher.Write(pieceData)
		if err != nil {
			return nil, err
		}

		n, err = pieceHasher.Write(pieceData)
		if err != nil {
			return nil, err
		}

		pieceHashes = append(pieceHashes, pieceHasher.HexSum())
		pieceHasher.Reset()
	}

	return &Torrent{
		FileName:       file.Name(),
		FileHash:       fileHasher.HexSum(),
		FileSize:       info.Size(),
		PieceSize:      pieceSize,
		PieceHashes:    pieceHashes,
		TrackerAddress: trackerAddr,
	}, nil
}

// Load reads torrent from given file
func Load(filePath string) (*Torrent, error) {
	if !strings.HasSuffix(filePath, FileExtension) {
		return nil, WrongExtensionError
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	torrent := &Torrent{}

	err = json.Unmarshal(b, torrent)
	if err != nil {
		return nil, err
	}

	return torrent, nil
}

// Save writes given torrent to a file
func Save(torrent *Torrent, dir, fileName string) error {
	filePath := filepath.Join(dir, fmt.Sprintf("%s.%s", fileName, FileExtension))

	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, message, 0666)
}

// Verify checks if given file is indeed a torrent data file
func Verify(torrent *Torrent, f *os.File) error {
	hash := hasher.New()

	for i, pieceHash := range torrent.PieceHashes {
		pieceNumber := int64(i)
		pieceOffset := piece.Offset(torrent.PieceSize, pieceNumber)
		pieceData, err := piece.Read(f, torrent.PieceSize, pieceOffset)
		if err != nil {
			return err
		}

		err = hash.Verify(pieceData, pieceHash)
		if err != nil {
			return err
		}
	}

	return nil
}
