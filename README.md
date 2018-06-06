# OpenStack-Exporter

Simple Prometheus exporter for OpenStack

## Deployment

First clone the repo

```git clone https://github.com/iGene/openstack-exporter.git```

A Dockerfile is provided in this repo, to build a OpenStack Exporter Docker image, simply run

```docker build -t openstack-exporter openstack-exporter/```

Copy the sample configuraion file and fill in

- OpenStack Username
- OpenStack Password
- OpenStack Keystone Endpoint

Launch the exporter using Docker

```docker run --name openstack-exporter -v $(pwd)/openstack.toml:/etc/openstack-exporter/openstack.toml -d -p 9183:9183 openstack-exporter```

Check if its working by

```curl localhost:9183/metrics```
