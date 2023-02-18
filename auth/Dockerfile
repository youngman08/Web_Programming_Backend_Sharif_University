# FROM golang:1.16 as base
# FROM base as dev
# RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# WORKDIR /opt/app/api
# COPY ./auth .
# CMD ["air"]
# syntax=docker/dockerfile:1


# FROM golang:1.16-alpine
# WORKDIR /app
# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download
# COPY *.go ./
# RUN go env -w GO111MODULE=off
# RUN set GOPATH=/home/rezasoumi/go
# RUN go install
# RUN go build -o /docker-gs-ping
# EXPOSE 8090
# CMD [ "/docker-gs-ping" ]


FROM golang:1.19-alpine
WORKDIR /usr/src/goapp
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
# COPY ../../../../../usr/local/go/src/net/netip/* ./
COPY data/*.go ./data/
RUN go mod tidy -go=1.19
RUN go env -w GO111MODULE="auto"
RUN go build -o /docker-go-sha
CMD [ "/docker-go-sha" ]