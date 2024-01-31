FROM golang:1.21-alpine3.18

ADD . /go/src/mikuserv

WORKDIR "/go/src/mikuserv"

RUN go build

ENTRYPOINT "./mikuserv"