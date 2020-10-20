package torrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dawidd6/p2p/pkg/file"
	"github.com/dawidd6/p2p/pkg/hash"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/errors"
)

const FileExtension = "torrent.json"

func Create(filePath string, pieceSize uint64, trackerAddr string) (*Torrent, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(fileContent)
	pieceHashes := make([]string, 0)

	for {
		pieceData := make([]byte, pieceSize)

		n, err := reader.Read(pieceData)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		pieceData = pieceData[:n]
		pieceHash := hash.Compute(pieceData)
		pieceHashes = append(pieceHashes, pieceHash)
	}

	return &Torrent{
		FileName:       filepath.Base(filePath),
		FileHash:       hash.Compute(fileContent),
		FileSize:       uint64(len(fileContent)),
		PieceSize:      pieceSize,
		PieceHashes:    pieceHashes,
		TrackerAddress: trackerAddr,
	}, nil
}

func Load(filePath string) (*Torrent, error) {
	if !strings.HasSuffix(filePath, FileExtension) {
		return nil, errors.WrongTorrentExtensionError
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

func Save(torrent *Torrent) error {
	filename := fmt.Sprintf("%s.%s", torrent.FileName, FileExtension)

	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}

func Verify(torrent *Torrent, dir string) error {
	f, err := file.Open(dir, torrent.FileName)
	if err != nil {
		return err
	}

	fileContent, err := file.ReadAll(f)
	if err != nil {
		return err
	}

	if hash.Compute(fileContent) != torrent.FileHash {
		return errors.FileChecksumMismatchError
	}

	for i, pieceHash := range torrent.PieceHashes {
		pieceData, err := piece.Read(f, torrent.PieceSize, uint64(i))
		if err != nil {
			return err
		}

		if hash.Compute(pieceData) != pieceHash {
			return errors.PieceChecksumMismatchError
		}
	}

	return nil
}
