all: p2p p2pd p2p-tracker p2p-trackerd

proto:
	bash -c "protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/*.proto"
	bash -c "protoc --go_out=. --go_opt=paths=source_relative pkg/proto/*.proto"

p2p p2pd p2p-tracker p2p-trackerd:
	go build -o ./bin/$@ ./cmd/$@