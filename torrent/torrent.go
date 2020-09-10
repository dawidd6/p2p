package torrent

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"path/filepath"
	"time"

	"github.com/dawidd6/p2p/piece"
	"github.com/dawidd6/p2p/proto"
)

const TORRENT_EXTENSION = "json"

type Torrent struct {
	*proto.Torrent
}

func CreateTorrent(name string, filePaths []string) (*Torrent, error) {
	if name == "" {
		return nil, errors.New("name of torrent can't be empty")
	}

	torrent := &Torrent{
		Torrent: &proto.Torrent{
			Name:      name,
			Timestamp: time.Now().UTC().Unix(),
			Files:     make([]*proto.File, 0),
		},
	}

	for _, filePath := range filePaths {
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(fileContent)
		sha := sha256.Sum256(fileContent)
		fileName := filepath.Base(filePath)
		pieces := make([]*proto.Piece, 0)
		fil := &proto.File{
			Name:   fileName,
			Sha256: sha[:],
			Pieces: pieces,
		}

		for i := 0; i < math.MaxInt64; i++ {
			chunk := make([]byte, piece.PIECE_LENGTH)

			_, err := reader.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			number := uint64(i)
			sha := sha256.Sum256(fileContent)
			pie := &proto.Piece{
				Number: number,
				Sha256: sha[:],
			}

			fil.Pieces = append(fil.Pieces, pie)
		}

		torrent.Files = append(torrent.Files, fil)
	}

	return torrent, nil
}

func Load(filePath string) (*Torrent, error) {
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

func (torrent *Torrent) Save() error {
	filename := fmt.Sprintf("%s.%s", torrent.Name, TORRENT_EXTENSION)

	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}
