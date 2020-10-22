FROM golang:1-alpine AS builder
COPY . /app
WORKDIR /app
RUN apk -U add make git
RUN make

FROM alpine:3
COPY --from=builder /app/bin /
RUN addgroup -S p2p && adduser -S p2p -G p2p
USER p2p
CMD ["p2p"]
