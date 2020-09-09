package torrent

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dawidd6/p2p/file"
)

func CreateTorrent(name string, filePaths []string) (*Torrent, error) {
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

		fmt.Println(len(fileContent))
	}

	return torrent, nil
}
