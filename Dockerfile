FROM golang:1-alpine AS builder
COPY . /app
WORKDIR /app
RUN apk -U add make git
RUN make

FROM alpine
COPY --from=builder /app/bin /
CMD ["p2p"]
