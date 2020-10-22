FROM golang:1-alpine AS builder
COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0
RUN apk -U add make git
RUN make build test

FROM alpine:3
COPY --from=builder /app/bin/* /bin/
CMD ["p2p"]
