package collector

import (
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
	"github.com/prometheus/client_golang/prometheus"
)

// computeCollector collects statistics about hypervisor and Nova in an OpenStack
// cluster.

type computeCollector struct {
	provider gophercloud.ProviderClient

	TotalRunningVMs prometheus.Gauge

	TotalMemoryMBUsed prometheus.Gauge

	TotalVCPUsUsed prometheus.Gauge
}

// NewComputeCollector creates an instance of computeCollector.
func NewComputeCollector(provider gophercloud.ProviderClient) *computeCollector {
	return &computeCollector{
		provider: provider,

		TotalRunningVMs: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "openstack_total_running_vms",
				Help: "Number of total vms running",
			},
		),

		TotalMemoryMBUsed: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "openstack_total_memory_mb_used",
				Help: "Number of total memory used in MB",
			},
		),

		TotalVCPUsUsed: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "openstack_total_vcpus_used",
				Help: "Number of total VCPU used",
			},
		),
	}
}

func (c *computeCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		c.TotalRunningVMs,
		c.TotalMemoryMBUsed,
		c.TotalVCPUsUsed,
	}
}

func (c *computeCollector) collect() error {
	region := gophercloud.EndpointOpts{Region: "RegionOne"}
	computeClient, err := openstack.NewComputeV2(&c.provider, region)
	if err != nil {
		return err
	}

	var v float64
	var m float64
	var cpu float64
	allPages, err := hypervisors.List(computeClient).AllPages()
	if err != nil {
		return err
	}

	allHypervisors, err := hypervisors.ExtractHypervisors(allPages)
	if err != nil {
		return err
	}

	for _, hypervisor := range allHypervisors {
		v += float64(hypervisor.RunningVMs)
		m += float64(hypervisor.MemoryMBUsed)
		cpu += float64(hypervisor.VCPUsUsed)

	}
	c.TotalRunningVMs.Set(v)
	c.TotalMemoryMBUsed.Set(m)
	c.TotalVCPUsUsed.Set(cpu)

	return nil
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by computeCollector.
func (c *computeCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.collectorList() {
		metric.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics.
func (c *computeCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Println("failed collecting compute metrics:", err)
	}

	for _, metric := range c.collectorList() {
		metric.Collect(ch)
	}

}
