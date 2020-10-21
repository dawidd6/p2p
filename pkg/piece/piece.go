package piece

import "io"

const Size = 256 * 1024 // 256 kB

func Offset(pieceSize, pieceNumber int64) int64 {
	return pieceSize * pieceNumber
}

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
