FROM golang:1.6.2-alpine

RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/ErnestIO/aws-definition-mapper
WORKDIR /go/src/github.com/ErnestIO/aws-definition-mapper

RUN make deps && go install

ENTRYPOINT ./entrypoint.sh

