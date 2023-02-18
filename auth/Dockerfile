FROM golang:1.19-alpine
WORKDIR /usr/src/goapp

COPY go.mod ./
COPY go.sum ./
COPY *.go ./
COPY data/*.go ./data/

RUN go mod tidy -go=1.19

RUN go env -w GO111MODULE="auto"

RUN go build -o /docker-go-sha

CMD [ "/docker-go-sha" ]