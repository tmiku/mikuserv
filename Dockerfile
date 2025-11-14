FROM golang:1.24-alpine3.22
WORKDIR "/go/src/mikuserv"

# copy dependencies first for docker cache
COPY go.mod go.sum ./
RUN go mod download

# add repo files to /go/src/mikuserv
COPY . .

RUN go build

ENTRYPOINT "./mikuserv"
