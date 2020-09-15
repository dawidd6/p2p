package main

const (
	WrongTorrentExtensionError   = "wrong torrent file extension, should be " + TorrentExtension
	MetadataContextNotOkError    = "failed to get metadata from context"
	FileChecksumMismatchError    = "saved checksum of file in torrent file does not match actual checksum of the file"
	PieceChecksumMismatchError   = "saved checksum of piece in torrent file does not match actual checksum of the piece"
	TorrentChecksumMismatchError = "saved checksum of torrent does not match actual checksum of the torrent"
)
