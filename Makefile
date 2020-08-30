all: p2p p2pd p2p-tracker p2p-trackerd

proto:
	bash -c "protoc --go-grpc_out=. pkg/proto/*.proto"
	bash -c "protoc --go_out=. pkg/proto/*.proto"

p2p p2pd p2p-tracker p2p-trackerd:
	go build -o ./bin/$@ ./cmd/$@