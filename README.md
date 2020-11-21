# P2P

![CI](https://github.com/dawidd6/p2p/workflows/CI/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dawidd6/p2p)](https://pkg.go.dev/github.com/dawidd6/p2p)

A gRPC based P2P file sharing system, written in Go.

## Building

Needs Go **1.15**.

```shell script
make
```

## Install

There are 2 supported methods of installing this project right now:

### Go

> Installing via `go get` could give you an unstable version of the project.

```shell script
go get github.com/dawidd6/p2p/cmd/...
```

### Docker

> `latest` Docker image tag corresponds to the latest git tag in this repo.
> Specific image tags corresponding to git tags are available too.

```shell script
docker pull dawidd6/p2p
```

## Usage

- `p2p` program to control the daemon.
- `p2pd` program to start the daemon.
- `p2ptrackerd` program to start the tracker.
