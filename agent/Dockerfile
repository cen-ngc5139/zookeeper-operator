FROM golang:1.13-alpine as builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
RUN apk add --no-cache git
WORKDIR     /usr/src/zk-agent
COPY        . /usr/src/zk-agent
RUN         go build -v 

FROM        alpine:3.10
COPY        --from=builder /usr/src/zk-agent/zk-agent /usr/local/bin/zk-agent
ENTRYPOINT  ["/usr/local/bin/zk-agent"]
