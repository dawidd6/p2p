package errors

import "errors"

var (
	WrongTorrentExtensionError   = errors.New("wrong torrent file extension")
	MetadataContextNotOkError    = errors.New("failed to get metadata from context")
	FileChecksumMismatchError    = errors.New("saved checksum of file in torrent file does not match actual checksum of the file")
	PieceChecksumMismatchError   = errors.New("saved checksum of piece in torrent file does not match actual checksum of the piece")
	TorrentChecksumMismatchError = errors.New("saved checksum of torrent does not match actual checksum of the torrent")
	TorrentNotFound              = errors.New("torrent not found")
)
