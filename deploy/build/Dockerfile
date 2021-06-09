FROM golang:1.16.4-buster
WORKDIR /spiracle

RUN go install github.com/buraksezer/olric/cmd/olric-cli@latest
COPY bin/spiracle spiracle
