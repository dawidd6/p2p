package piece

import "io"

func Read(reader io.ReaderAt, pieceSize, pieceNumber uint64) ([]byte, error) {
	b := make([]byte, pieceSize)
	n, err := reader.ReadAt(b, int64(pieceSize*pieceNumber))
	if err != nil && err != io.EOF {
		return nil, err
	}
	return b[:n], nil
}

func Write(writer io.WriterAt, pieceNumber uint64, pieceData []byte) error {
	size := uint64(len(pieceData))
	if size == 0 {
		return nil
	}
	_, err := writer.WriteAt(pieceData, int64(size*pieceNumber))
	if err != nil {
		return err
	}
	return nil
}
