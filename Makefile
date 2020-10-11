VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse HEAD)
PACKAGE ?=
BIN ?= bin

build:
	mkdir -p bin
	go build -o $(BIN) -ldflags "-s -w -X github.com/dawidd6/p2p/pkg/version.Version=$(VERSION)" ./cmd/...

proto-plugins:
	GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	GO111MODULE=off go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/*.proto
	protoc --go_out=. --go_opt=paths=source_relative pkg/proto/*.proto

image:
	docker build -t p2p .