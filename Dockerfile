FROM golang:1.10-alpine AS build-env
LABEL maintainer="Gene Kuo<igene@igene.tw>"

ENV GOPATH /go
WORKDIR $GOPATH/src/github.com/iGene/openstack-exporter

RUN apk add --no-cache git make g++ && \
  go get -u github.com/golang/dep/cmd/dep

COPY . .
RUN make && \
  mv openstack-exporter /tmp/openstack-exporter

# Run stage
FROM alpine

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*
COPY --from=build-env /tmp/openstack-exporter /bin/openstack-exporter

ENTRYPOINT ["openstack-exporter"]
