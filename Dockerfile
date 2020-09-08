FROM golang:1

COPY . /app

WORKDIR /app

RUN rake

ENTRYPOINT ["bin/p2p"]