FROM golang:1.13.8-alpine3.11 as build-env

ENV GO111MODULE=on
ENV GOPROXY=http://goproxy.cn
ENV BUILDPATH=github.com/icowan/grpc-world
RUN mkdir -p /go/src/${BUILDPATH}
COPY ./ /go/src/${BUILDPATH}
RUN cd /go/src/${BUILDPATH} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM alpine:latest

COPY --from=build-env /go/bin/grpc-world /go/bin/grpc-world

WORKDIR /go/bin/
CMD ["/go/bin/grpc-world", "-http-addr", ":8080", "-grpc-addr", ":8081"]