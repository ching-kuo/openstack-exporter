FROM golang:1.10.2-alpine3.7

WORKDIR /go/src/openstack-exporter
COPY . .

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["openstack-exporter"]
