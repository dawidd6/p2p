// Package state includes torrent state struct
package state

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/dawidd6/p2p/pkg/torrent"
)

const FileExtension = ".state" + torrent.FileExtension

// Read reads state from a reader
func Read(reader io.Reader) (*State, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	state := &State{}

	err = json.Unmarshal(b, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// Write writes state torrent to a writer
func Write(writer io.Writer, state *State) error {
	message, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	_, err = writer.Write(message)
	if err != nil {
		return err
	}

	return nil
}
