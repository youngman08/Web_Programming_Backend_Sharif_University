# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.io,direct go mod download && GOPROXY=https://goproxy.io,direct go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]

