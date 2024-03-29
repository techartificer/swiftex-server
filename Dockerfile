FROM golang:alpine AS builder

RUN apk add git openssh

ENV GOPATH=/go

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GO111MODULE=on
# RUNy go env -w GOPROXY="https://goproxy.io,direct"

WORKDIR ${GOPATH}/src/github.com/techartificer/swiftex-server

COPY go.mod go.sum ${GOPATH}/src/github.com/techartificer/swiftex-server/
RUN go mod download -x

COPY . ${GOPATH}/src/github.com/techartificer/swiftex-server

RUN go build -v -o swiftex
RUN mv swiftex /go/bin/swiftex

FROM alpine

WORKDIR /root

COPY --from=builder /go/bin/swiftex /usr/local/bin/swiftex
EXPOSE 4141

CMD [ "swiftex" ]

