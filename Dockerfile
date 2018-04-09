FROM golang:1.10

ADD . /go/src/github.com/fossas/go-resolve
RUN go install github.com/fossas/go-resolve/...