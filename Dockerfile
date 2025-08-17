FROM golang:alpine AS builder

WORKDIR /building

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o nas-net-tool .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /building/nas-net-tool /app/nas-net-tool

ENTRYPOINT ["/app/nas-net-tool"]