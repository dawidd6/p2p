all: p2p p2pd p2p-tracker p2p-trackerd

proto:
	bash -c "rm -f pkg/proto/*.go"
	bash -c "protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/*.proto"
	bash -c "protoc --go_out=. --go_opt=paths=source_relative pkg/proto/*.proto"

proto-plugins:
	  GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	  GO111MODULE=off go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

p2p p2pd p2p-tracker p2p-trackerd:
	go build -o ./bin/$@ ./cmd/$@