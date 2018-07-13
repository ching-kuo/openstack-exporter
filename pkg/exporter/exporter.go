package exporter

import (
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/iGene/openstack-exporter/pkg/collector"
	"github.com/prometheus/client_golang/prometheus"
)

// OpenStackExporter Wraps all collectors in a single exporter to extract metrics
// and make sure it is thread safe.
type OpenStackExporter struct {
	sync.Mutex
	collectors []prometheus.Collector
}

// verify that the exporter implementation is correct
var _ prometheus.Collector = &OpenStackExporter{}

// NewOpenStackExporter creates an instance to OpenStackExporter and returns a
// reference to it.
func NewOpenStackExporter(provider *gophercloud.ProviderClient, region string) *OpenStackExporter {
	return &OpenStackExporter{
		collectors: []prometheus.Collector{
			collector.NewComputeCollector(provider, region),
			collector.NewBlockStorageCollector(provider, region),
			collector.NewNetworkCollector(provider, region),
            collector.NewClusterHealthCollector(provider, region),
		},
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics.
func (o *OpenStackExporter) Collect(ch chan<- prometheus.Metric) {
	// Only one Collect call in progress at a time.
	o.Lock()
	defer o.Unlock()

	for _, oo := range o.collectors {
		oo.Collect(ch)
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (o *OpenStackExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, oo := range o.collectors {
		oo.Describe(ch)
	}
}
