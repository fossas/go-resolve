FROM golang:1.10-alpine

RUN apk add git

ADD . /go/src/github.com/fossas/go-resolve
RUN go install github.com/fossas/go-resolve/...
