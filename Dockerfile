FROM golang:alpine

WORKDIR /go

# Prepare build cache
RUN apk add git
RUN go get -d -v github.com/gorilla/mux

# Setup app packages, build and run
CMD go get -d -v github.com/aetelani/maprest && go install -v ./... && maprest
