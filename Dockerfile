FROM golang:1.24-alpine3.19

ADD . /go/src/mikuserv

WORKDIR "/go/src/mikuserv"

RUN go get github.com/resend/resend-go/v2
RUN go build

ENTRYPOINT "./mikuserv"
