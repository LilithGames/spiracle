FROM golang:1.16.4-buster
WORKDIR /spiracle

RUN go install github.com/buraksezer/olric/cmd/olric-cli@v0.4.0-beta.2
COPY bin/spiracle spiracle
