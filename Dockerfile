FROM golang:1-alpine AS builder
COPY . /app
WORKDIR /app
RUN apk -U add make git
RUN make

FROM alpine:3
COPY --from=builder /app/p2p /bin/p2p
USER p2p
ENTRYPOINT ["p2p"]
