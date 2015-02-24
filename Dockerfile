# golang image where workspace (GOPATH) configured at /go.
FROM golang:latest

ADD ./keys/cert.pem cert.pem
ADD ./keys/key.pem key.pem

RUN go get github.com/hectorj2f/search_networking/server
RUN go install github.com/hectorj2f/search_networking/server

CMD ["./bin/server"]


EXPOSE 3333
