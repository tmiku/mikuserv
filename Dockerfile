FROM golang:1.21-alpine3.18

ADD . /go/src/mikuserv

WORKDIR "/go/src/mikuserv"

RUN go get github.com/resend/resend-go/v2
RUN go build

ENTRYPOINT "./mikuserv"
