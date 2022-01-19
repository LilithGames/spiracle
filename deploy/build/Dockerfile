FROM golang:1.17.5-buster
WORKDIR /spiracle

RUN go install github.com/buraksezer/olric/cmd/olric-cli@latest
COPY bin/spiracle spiracle
