VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse HEAD)

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)"

test:
	go test -v ./...

proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/*/*.proto
	protoc --go_out=. --go_opt=paths=source_relative pkg/*/*.proto

image:
	docker build -t p2p .