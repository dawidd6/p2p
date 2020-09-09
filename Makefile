PREFIX ?= /usr/local
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse HEAD)

build:
	go build -ldflags "-X main.version=$(VERSION)"

install:
	GOBIN=$(PREFIX)/bin go install -ldflags "-X main.version=$(VERSION)"

proto-plugins:
	GO111MODULE=off go get google.golang.org/protobuf/cmd/protoc-gen-go
	GO111MODULE=off go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

proto: proto-plugins
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative */*.proto
	protoc --go_out=. --go_opt=paths=source_relative */*.proto
