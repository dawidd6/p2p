PREFIX ?= /usr/local

build:
	go build

install:
	GOBIN=$(PREFIX)/bin go install

proto-plugins:
	GO111MODULE=off go get google.golang.org/protobuf/cmd/protoc-gen-go
	GO111MODULE=off go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

proto: proto-plugins
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/#{proto}/#{proto}.proto
	protoc --go_out=. --go_opt=paths=source_relative pkg/#{proto}/#{proto}.proto
