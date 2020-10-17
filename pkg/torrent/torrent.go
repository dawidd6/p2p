package torrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/dawidd6/p2p/pkg/errors"
	"github.com/dawidd6/p2p/pkg/utils"
)

const FileExtension = ".torrent.json"

func Create(filePath string, pieceSize uint64) (*Torrent, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(fileContent)
	pieceHashes := make([]string, 0)

	for {
		piece := make([]byte, pieceSize)

		n, err := reader.Read(piece)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		pieceHash := utils.Sha256Sum(piece[:n])
		pieceHashes = append(pieceHashes, pieceHash)
	}

	t := &Torrent{

		FileName:    filepath.Base(filePath),
		FileHash:    utils.Sha256Sum(fileContent),
		FileSize:    uint64(len(fileContent)),
		PieceSize:   pieceSize,
		PieceHashes: pieceHashes,
		TrackerUrls: []string{"localhost:8889"}, // TODO customizable trackers urls
	}

	return t, nil
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
	filename := fmt.Sprintf("%s%s", torrent.FileName, FileExtension)

	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}

func MissingPieces() {

}

func Verify(torrent *Torrent, dir string) error {
	return utils.DoInDirectory(dir, func() error {
		filePath := filepath.Join(dir, torrent.FileName)

		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		expected := torrent.FileHash
		got := utils.Sha256Sum(fileContent)

		if got != expected {
			return errors.FileChecksumMismatchError
		}

		for i, pieceHash := range torrent.PieceHashes {
			chunk, err := utils.ReadFilePiece(filePath, torrent.PieceSize, uint64(i))
			if err != nil {
				return err
			}

			expected := pieceHash
			got := utils.Sha256Sum(chunk)

			if got != expected {
				return errors.PieceChecksumMismatchError
			}
		}

		return nil
	})
}
