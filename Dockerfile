FROM ubuntu:14.04
MAINTAINER Nathan Borror <nathan@nathanborror.com>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update
RUN apt-get install -y build-essential golang git mercurial sqlite3

ENV GOPATH /go

RUN go get -v github.com/nathanborror/gommon/... # last change: 2014-08-09
ADD . /go/src/github.com/nathanborror/skit
RUN go get -v github.com/nathanborror/skit
WORKDIR /go/src/github.com/nathanborror/skit

EXPOSE 8080
ENTRYPOINT /go/bin/skit
