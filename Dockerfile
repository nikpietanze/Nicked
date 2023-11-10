FROM golang

ADD . /go/src/nicked

WORKDIR /go/src/nicked

RUN go mod tidy
RUN go build

ENTRYPOINT ./go/src/nicked/nicked.io

EXPOSE 80 443 8080
