FROM golang:latest

ENV GOPATH=/go

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

WORKDIR ${GOPATH}/src/github.com/techartificer/swiftex-server
COPY . ${GOPATH}/src/github.com/techartificer/swiftex-server

RUN go mod download

RUN go build -v -o swiftex
RUN mv swiftex /go/bin/swiftex

EXPOSE 4141

CMD [ "swiftex" ]



