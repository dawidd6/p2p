VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse HEAD)

DOCKER_NETWORK = p2p-network
DOCKER_P2P_IMAGE = dawidd6/p2p
DOCKER_DB_IMAGE = eqalpha/keydb

# Don't enable cgo, so built binaries are portable
export CGO_ENABLED = 0

# Build binaries
.PHONY: build
build:
	@mkdir -p bin
	go build -ldflags "-s -w -X github.com/dawidd6/p2p/pkg/version.Version=$(VERSION)" -o bin ./cmd/...

# Test code
.PHONY: test
test:
	go test -v -count=1 ./...

# Build Docker image of the project
.PHONY: image
image:
	docker build -t $(DOCKER_P2P_IMAGE) .

# Test integration of the system in containers
.PHONY: integration-test
integration-test:
	# Create network
	docker network inspect $(DOCKER_NETWORK) || docker network create $(DOCKER_NETWORK)
	# Create torrent
	docker run --rm --user="$(shell id -u):$(shell id -g)" --volume="$(PWD):$(PWD)" --workdir="$(PWD)" --tty $(DOCKER_P2P_IMAGE) \
		p2p create --piece-size=4096 --tracker-host=tracker go.sum
	# Start database
	docker run --rm --network=$(DOCKER_NETWORK) --name=db --detach $(DOCKER_DB_IMAGE)
	# Start tracker
	docker run --rm --network=$(DOCKER_NETWORK) --name=tracker --detach $(DOCKER_P2P_IMAGE) \
		p2ptrackerd --announce-interval=5s --host=0.0.0.0 --db-host=db
	# Start daemons
	docker run --rm --network=$(DOCKER_NETWORK) --name=daemon-seed-1 --volume="$(PWD)":"$(PWD)" --workdir="$(PWD)" --detach $(DOCKER_P2P_IMAGE) \
		p2pd --host=0.0.0.0 --port=8881 --seed-port=44441
	docker run --rm --network=$(DOCKER_NETWORK) --name=daemon-seed-2 --volume="$(PWD)":"$(PWD)" --workdir="$(PWD)" --detach $(DOCKER_P2P_IMAGE) \
		p2pd --host=0.0.0.0 --port=8882 --seed-port=44442
	docker run --rm --network=$(DOCKER_NETWORK) --name=daemon-leech-1 --workdir="$(PWD)" --detach $(DOCKER_P2P_IMAGE) \
		p2pd --host=0.0.0.0 --port=8883 --seed-port=44443
	docker run --rm --network=$(DOCKER_NETWORK) --name=daemon-leech-2 --workdir="$(PWD)" --detach $(DOCKER_P2P_IMAGE) \
		p2pd --host=0.0.0.0 --port=8884 --seed-port=44444
	# Add torrents
	docker run --rm --network=$(DOCKER_NETWORK) --volume="$(PWD):$(PWD):ro" --workdir="$(PWD)" --tty $(DOCKER_P2P_IMAGE) \
		p2p add --host=daemon-leech-1 --port=8883 go.sum.torrent.json
	docker run --rm --network=$(DOCKER_NETWORK) --volume="$(PWD):$(PWD):ro" --workdir="$(PWD)" --tty $(DOCKER_P2P_IMAGE) \
		p2p add --host=daemon-leech-2 --port=8884 go.sum.torrent.json
	docker run --rm --network=$(DOCKER_NETWORK) --volume="$(PWD):$(PWD):ro" --workdir="$(PWD)" --tty $(DOCKER_P2P_IMAGE) \
		p2p add --host=daemon-seed-1 --port=8881 go.sum.torrent.json
	docker run --rm --network=$(DOCKER_NETWORK) --volume="$(PWD):$(PWD):ro" --workdir="$(PWD)" --tty $(DOCKER_P2P_IMAGE) \
		p2p add --host=daemon-seed-2 --port=8882 go.sum.torrent.json
	# Wait some time
	sleep 10s
	# Get statuses
	docker run --rm --network=$(DOCKER_NETWORK) --tty $(DOCKER_P2P_IMAGE) \
		p2p status --host=daemon-seed-1 --port=8881
	docker run --rm --network=$(DOCKER_NETWORK) --tty $(DOCKER_P2P_IMAGE) \
		p2p status --host=daemon-seed-2 --port=8882
	docker run --rm --network=$(DOCKER_NETWORK) --tty $(DOCKER_P2P_IMAGE) \
		p2p status --host=daemon-leech-1 --port=8883
	docker run --rm --network=$(DOCKER_NETWORK) --tty $(DOCKER_P2P_IMAGE) \
		p2p status --host=daemon-leech-2 --port=8884
	# Print logs
	docker logs tracker
	docker logs daemon-seed-1
	docker logs daemon-leech-1
	docker logs daemon-leech-2
	# Stop all containers
	docker stop db tracker daemon-seed-1 daemon-leech-1 daemon-leech-2
	# Delete network
	docker network rm $(DOCKER_NETWORK)

# Generate code from Protocol Buffers definitions
.PHONY: proto
proto:
	protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/*/*.proto
	protoc --go_out=. --go_opt=paths=source_relative pkg/*/*.proto

