FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/tinrab/meower

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY util util
COPY mq mq
COPY schema schema
COPY meow-service meow-service
COPY pusher-service pusher-service
COPY storage-service storage-service

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
