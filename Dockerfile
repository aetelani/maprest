FROM golang:alpine

WORKDIR /go/src/github.com/aetelani/maprest

RUN go get -d -v github.com/aetelani/maprest
RUN go build

CMD ["./maprest"]
