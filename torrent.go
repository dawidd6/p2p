package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

const TorrentExtension = ".torrent.json"

func CreateTorrentFromDir(name, dir string) (*Torrent, error) {
	panic("TODO")
	return nil, nil // TODO	recursively add files
}

func CreateTorrentFromFiles(name string, filePaths []string) (*Torrent, error) {
	if name == "" {
		name = filePaths[0]
	}

	files := make([]*File, 0)

	for _, filePath := range filePaths {
		pieces := make([]*Piece, 0)

		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(fileContent)

		for {
			chunk := make([]byte, PieceLength)

			n, err := reader.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			piece := &Piece{
				Sha256: Sha256Sum(chunk[:n]),
			}

			pieces = append(pieces, piece)
		}

		file := &File{
			Name:   filepath.Base(filePath),
			Sha256: Sha256Sum(fileContent),
			Size:   uint64(len(fileContent)),
			Pieces: pieces,
		}

		files = append(files, file)
	}

	torrent := &Torrent{
		Name:      name,
		Timestamp: uint64(time.Now().UTC().Unix()),
		Files:     files,
		Trackers:  []string{"localhost:8889"}, // TODO customizable trackers urls
	}

	message, err := proto.Marshal(torrent)
	if err != nil {
		return nil, err
	}

	torrent.Sha256 = Sha256Sum(message)

	return torrent, nil
}

func LoadTorrentFromFile(filePath string) (*Torrent, error) {
	if !strings.HasSuffix(filePath, TorrentExtension) {
		return nil, errors.New(WrongTorrentExtensionError)
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

func (x *Torrent) SaveTorrentToFile() error {
	filename := fmt.Sprintf("%s%s", x.Name, TorrentExtension)

	message, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, message, 0644)
}

func (x *Torrent) VerifyFiles(dir string) error {
	return DoInDirectory(dir, func() error {
		for _, file := range x.Files {
			fileContent, err := ioutil.ReadFile(file.Name)
			if err != nil {
				return err
			}

			expected := file.Sha256
			got := Sha256Sum(fileContent)

			if got != expected {
				fmt.Println("file:", file.Name)
				fmt.Println("got:", got)
				fmt.Println("expected:", expected)
				return errors.New(FileChecksumMismatchError)
			}

			for i, piece := range file.Pieces {
				chunk, err := ReadFilePiece(file.Name, PieceLength, i)
				if err != nil {
					return err
				}

				expected := piece.Sha256
				got := Sha256Sum(chunk)

				if got != expected {
					fmt.Println("number:", i)
					fmt.Println("got:", got)
					fmt.Println("expected:", expected)
					return errors.New(PieceChecksumMismatchError)
				}
			}
		}

		return nil
	})
}
