# OpenStack-Exporter
Simple Prometheus exporter for OpenStack.

## Build
First clone the repo, and should be execute `Makefile` using make tool:
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
Should change this part into 3 sections to selecting each section describing how to configure it using the below methods.

1. Configuration File:

```sh
$ cp openstack.toml.example openstack.toml
$ docker run -d -p 9183:9183  \
      -v $(pwd)/openstack.toml:/etc/openstack-exporter/openstack.toml \
      --name openstack-exporter \
      igene/openstack-exporter:v0.1.0 --config /etc/openstack-exporter/openstack.toml
```

2. Command line option:

```sh
$ docker run -d -p 9183:9183  \
      --name openstack-exporter \
      igene/openstack-exporter:v0.1.0 \
      --keystone-url=http://172.22.132.21/identity/v3 \
      --project-name=admin \
      --username=admin \
      --password=secret \
      --domain-name=default \
      --region-name=RegionOne
```

3. Environment variables:

```sh
$ cp -rp openrc.example openrc
$ docker run -d -p 9183:9183  \
      --env-file=openrc \
      --name openstack-exporter \
      igene/openstack-exporter:v0.1.0
```

Check if its working by:
```sh
$ curl localhost:9183/metrics
```
