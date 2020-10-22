FROM golang:1-alpine AS builder
COPY . /app
WORKDIR /app
RUN apk -U add make git
RUN make

FROM alpine:3
COPY --from=builder /app/bin /
RUN groupadd -r p2p && useradd -r -g p2p p2p
USER p2p
CMD ["p2p"]
