PREFIX ?= /usr/local
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse HEAD)
PACKAGE ?= github.com/dawidd6/p2p
BIN ?= bin

build:
	mkdir -p bin
	go build -o $(BIN) -ldflags "-s -w -X $(PACKAGE)/version.Version=$(VERSION)" ./cmd/...

install:
	GOBIN=$(PREFIX)/bin go install -ldflags "-s -w -X $(PACKAGE)/version.Version=$(VERSION)" ./cmd/...

proto-plugins:
	GO111MODULE=off go get google.golang.org/protobuf/cmd/protoc-gen-go
	GO111MODULE=off go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

proto: proto-plugins
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative */*.proto
	protoc --go_out=. --go_opt=paths=source_relative */*.proto

image:
	docker build -t p2p .