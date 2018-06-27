VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

GOOS ?= $(shell go env GOOS)

ORG := github.com
OWNER := igenetw
REPOPATH ?= $(ORG)/$(OWNER)/openstack-exporter

.PHONY: build
build: openstack-exporter

.PHONY: openstack-exporter
openstack-exporter: depend
	GOOS=$(GOOS) go build -a -o $@ cmd/main.go

.PHONY: build_image
build_image:
	docker build -t $(OWNER)/openstack-exporter:$(VERSION) .

.PHONY: depend
depend:
	@dep ensure

.PHONY: clean
clean:
	@rm -rf openstack-exporter
