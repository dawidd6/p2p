// Package torrent provides convenient functions for Torrent type
package torrent

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dawidd6/p2p/pkg/hasher"
	"github.com/dawidd6/p2p/pkg/piece"
)

const FileExtension = "torrent.json"

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

func File(dir, name string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.%s", name, FileExtension))
}

// Load reads torrent from a reader
func Read(reader io.Reader) (*Torrent, error) {
	b, err := ioutil.ReadAll(reader)
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

// Save writes given torrent to a writer
func Write(writer io.Writer, torrent *Torrent) error {
	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	_, err = writer.Write(message)
	if err != nil {
		return err
	}

	return nil
}

// Verify checks if given file is indeed a torrent data file
func Verify(reader io.ReaderAt, torrent *Torrent) error {
	hash := hasher.New()

	for i, pieceHash := range torrent.PieceHashes {
		pieceNumber := int64(i)
		pieceOffset := piece.Offset(torrent.PieceSize, pieceNumber)
		pieceData, err := piece.Read(reader, torrent.PieceSize, pieceOffset)
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
