package torrent

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dawidd6/p2p/pkg/hash"
	"github.com/dawidd6/p2p/pkg/piece"

	"github.com/dawidd6/p2p/pkg/errors"
)

const FileExtension = "torrent.json"

func Create(file *os.File, pieceSize int64, trackerAddr string) (*Torrent, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileHash := hash.New()
	pieceHash := hash.New()
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

		n, err = fileHash.Write(pieceData)
		if err != nil {
			return nil, err
		}

		n, err = pieceHash.Write(pieceData)
		if err != nil {
			return nil, err
		}

		pieceHashes = append(pieceHashes, pieceHash.HexSum())
		pieceHash.Reset()
	}

	return &Torrent{
		FileName:       file.Name(),
		FileHash:       fileHash.HexSum(),
		FileSize:       info.Size(),
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

	return ioutil.WriteFile(filename, message, 0666)
}

func Verify(torrent *Torrent, f *os.File) error {
	for i, pieceHash := range torrent.PieceHashes {
		pieceNumber := int64(i)
		pieceOffset := torrent.PieceSize * pieceNumber
		pieceData, err := piece.Read(f, torrent.PieceSize, pieceOffset)
		if err != nil {
			return err
		}

		if hash.Compute(pieceData) != pieceHash {
			return errors.PieceChecksumMismatchError
		}
	}

	return nil
}
