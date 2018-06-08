# OpenStack-Exporter
Simple Prometheus exporter for OpenStack.

## Build
First clone the repo, and excute `Makefile` using make tool:
```sh
$ git clone https://github.com/iGene/openstack-exporter.git $GOPATH/src/github.com/iGene/openstack-exporter
$ cd $GOPATH/src/github.com/iGene/openstack-exporter
$ make
```

A Dockerfile is provided in this repo, to build a OpenStack Exporter Docker image, simply run:
```sh
$ make build_image
```

## Deployment
Copy the sample configuraion file and fill in:

- OpenStack Username
- OpenStack Password
- OpenStack Keystone Endpoint

Launch the exporter using Docker:
```sh
$ docker run -d -p 9183:9183  \
      -v $(pwd)/openstack.toml:/etc/openstack-exporter/openstack.toml \
      --name openstack-exporter \
      openstack-exporter:v0.1.0
```

Check if its working by:
```sh
$ curl localhost:9183/metrics
```
