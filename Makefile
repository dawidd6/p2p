all: p2p p2pd p2p-tracker p2p-trackerd

proto:
	bash -c "cd pkg/proto; protoc --go_out=plugins=grpc:. *.proto"

p2p p2pd p2p-tracker p2p-trackerd:
	go build -o ./bin/$@ ./cmd/$@