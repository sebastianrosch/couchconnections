ARG GOLANG_VERSION=1.13.4

FROM golang:$GOLANG_VERSION as builder
ARG BUILD_VERSION=0.0.1

WORKDIR /go/src/github.com/sebastianrosch/couchconnections
COPY . .

RUN make build VERSION=$BUILD_VERSION

FROM debian
RUN apt-get update -y && apt-get -y install ca-certificates

WORKDIR /app
COPY --from=builder /go/src/github.com/sebastianrosch/couchconnections/couchconnections-api /couchconnections-api

ENTRYPOINT ["/couchconnections-api"]
