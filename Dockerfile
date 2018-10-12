FROM golang:alpine as builder
LABEL maintainer="Roald Nefs <info@roaldnefs.com>"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

COPY . /go/src/github.com/roaldnefs/tanuki

RUN set -x \
       && apk add --no-cache --virtual .build-deps \
               git \
       && cd /go/src/github.com/roaldnefs/tanuki \
       && go get -t -v ./... \
       && go build \
       && mv tanuki /usr/bin/tanuki \
       && rm -rf /go \
       && echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/tanuki /usr/bin/tanuki

ENTRYPOINT [ "tanuki" ]
CMD [ "--help" ]
