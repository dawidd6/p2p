package torrent

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"path/filepath"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/dawidd6/p2p/file"
	"github.com/dawidd6/p2p/piece"
)

const TORRENT_EXTENSION = "torrent"

func CreateTorrent(name string, filePaths []string) (*Torrent, error) {
	if name == "" {
		return nil, errors.New("name of torrent can't be empty")
	}

	torrent := &Torrent{
		Name:      name,
		Timestamp: time.Now().UTC().Unix(),
		Files:     make([]*file.File, 0),
	}

	for _, filePath := range filePaths {
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(fileContent)
		sha := sha256.Sum256(fileContent)
		fileName := filepath.Base(filePath)
		pieces := make([]*piece.Piece, 0)
		fil := &file.File{
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
			pie := &piece.Piece{
				Number: number,
				Sha256: sha[:],
			}

			fil.Pieces = append(fil.Pieces, pie)
		}

		torrent.Files = append(torrent.Files, fil)
	}

	return torrent, nil
}

func Load(filePaths []string) ([]*Torrent, error) {
	torrents := make([]*Torrent, 0)

	for _, filePath := range filePaths {
		torrent := &Torrent{}

		message, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		err = proto.Unmarshal(message, torrent)
		if err != nil {
			return nil, err
		}

		torrents = append(torrents, torrent)
	}

	return torrents, nil
}

func (torrent *Torrent) Save() error {
	filename := fmt.Sprintf("%s.%s", torrent.Name, TORRENT_EXTENSION)

	message, err := proto.Marshal(torrent)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}
