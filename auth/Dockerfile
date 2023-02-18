FROM golang:1.19-alpine
WORKDIR /usr/src/goapp

COPY go.mod go.sum ./

RUN GOPROXY=https://goproxy.io,direct go mod download && GOPROXY=https://goproxy.io,direct go mod verify

COPY data/*.go ./data/
COPY *.go ./

RUN go build -o /docker-go-sha

CMD [ "/docker-go-sha" ]
