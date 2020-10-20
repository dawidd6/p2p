package piece

import "io"

func Read(reader io.ReaderAt, pieceSize, pieceOffset int64) ([]byte, error) {
	b := make([]byte, pieceSize)

	n, err := reader.ReadAt(b, pieceOffset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b[:n], nil
}

func Write(writer io.WriterAt, pieceOffset int64, pieceData []byte) error {
	_, err := writer.WriteAt(pieceData, pieceOffset)
	if err != nil {
		return err
	}

	return nil
}
