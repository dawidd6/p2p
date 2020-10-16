package torrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/dawidd6/p2p/pkg/file"

	"github.com/dawidd6/p2p/pkg/errors"
	"github.com/dawidd6/p2p/pkg/piece"
	"github.com/dawidd6/p2p/pkg/utils"

	"google.golang.org/protobuf/proto"
)

const FileExtension = ".torrent.json"

func CreateTorrentFromDir(name, dir string) (*Torrent, error) {
	panic("TODO")
	return nil, nil // TODO	recursively add files
}

func CreateTorrentFromFiles(name string, filePaths []string) (*Torrent, error) {
	if name == "" {
		name = filePaths[0]
	}

	size := uint64(0)
	files := make([]*file.File, 0)

	for _, filePath := range filePaths {
		pieces := make([]*piece.Piece, 0)

		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(fileContent)

		for {
			chunk := make([]byte, piece.PieceLength)

			n, err := reader.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			p := &piece.Piece{
				Sha256: utils.Sha256Sum(chunk[:n]),
			}

			pieces = append(pieces, p)
		}

		size += uint64(len(fileContent))

		f := &file.File{
			Name:   filepath.Base(filePath),
			Sha256: utils.Sha256Sum(fileContent),
			Size:   uint64(len(fileContent)),
			Pieces: pieces,
		}

		files = append(files, f)
	}

	torrent := &Torrent{
		Name:      name,
		Size:      size,
		Timestamp: uint64(time.Now().UTC().Unix()),
		Files:     files,
		Trackers:  []string{"localhost:8889"}, // TODO customizable trackers urls
	}

	message, err := proto.Marshal(torrent)
	if err != nil {
		return nil, err
	}

	torrent.Sha256 = utils.Sha256Sum(message)

	return torrent, nil
}

func LoadTorrentFromFile(filePath string) (*Torrent, error) {
	if !strings.HasSuffix(filePath, FileExtension) {
		return nil, errors.WrongTorrentExtensionError
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return LoadTorrentFromBytes(b)
}

func LoadTorrentFromBytes(b []byte) (*Torrent, error) {
	torrent := &Torrent{}

	err := json.Unmarshal(b, torrent)
	if err != nil {
		return nil, err
	}

	return torrent, nil
}

func SaveTorrentToFile(torrent *Torrent) error {
	filename := fmt.Sprintf("%s%s", torrent.Name, FileExtension)

	message, err := json.MarshalIndent(torrent, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}

func VerifyFiles(torrent *Torrent, dir string) error {
	return utils.DoInDirectory(dir, func() error {
		for _, f := range torrent.Files {
			fileContent, err := ioutil.ReadFile(f.Name)
			if err != nil {
				return err
			}

			expected := f.Sha256
			got := utils.Sha256Sum(fileContent)

			if got != expected {
				fmt.Println("file:", f.Name)
				fmt.Println("got:", got)
				fmt.Println("expected:", expected)
				return errors.FileChecksumMismatchError
			}

			for i, p := range f.Pieces {
				chunk, err := utils.ReadFilePiece(f.Name, piece.PieceLength, int64(i))
				if err != nil {
					return err
				}

				expected := p.Sha256
				got := utils.Sha256Sum(chunk)

				if got != expected {
					fmt.Println("number:", i)
					fmt.Println("got:", got)
					fmt.Println("expected:", expected)
					return errors.PieceChecksumMismatchError
				}
			}
		}

		return nil
	})
}
