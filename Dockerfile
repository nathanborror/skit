FROM ubuntu:14.04
MAINTAINER Nathan Borror <nathan@nathanborror.com>

RUN apt-get update
RUN apt-get install build-essential golang git mercurial sqlite3

ENV GOPATH /go

ADD . /go/src/github.com/nathanborror/skit
RUN go install github.com/nathanborror/skit

EXPOSE 8080
ENTRYPOINT /go/bin/skit
