FROM golang:alpine

WORKDIR /go/src/github.com/aetelani/maprest

RUN apk add git

# Setup packages
RUN go get -d -v github.com/aetelani/maprest
RUN go install -v ./...

CMD ["./maprest"]
