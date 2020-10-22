FROM golang:1-alpine AS builder
COPY . /app
WORKDIR /app
RUN apk -U add make git
RUN make

FROM alpine:3
COPY --from=builder /app/bin/* /bin/
CMD ["p2p"]
