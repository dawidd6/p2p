VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse HEAD)

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)"

test:
	go test -v ./...

proto-plugins:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.0

proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/*/*.proto
	protoc --go_out=. --go_opt=paths=source_relative pkg/*/*.proto

image:
	docker build -t p2p .