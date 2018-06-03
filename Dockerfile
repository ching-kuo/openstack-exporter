FROM golang:1.10.2-alpine3.7 AS build-env

WORKDIR /go/src/openstack-exporter
COPY . .

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go get -d -v ./...
RUN go build -o openstack-exporter -v .

FROM alpine

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates

COPY --from=build-env /go/src/openstack-exporter/openstack-exporter /usr/bin/

CMD ["openstack-exporter"]
