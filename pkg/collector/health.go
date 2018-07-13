package collector

import (
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
    "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/services"
	"github.com/prometheus/client_golang/prometheus"
)

// clusterHealthCollector collects statistics about hypervisor and Nova in an OpenStack
// cluster.
type clusterHealthCollector struct {
	provider *gophercloud.ProviderClient

	region string

    HealthStatus prometheus.Gauge
}

// NewClusterHealthCollector creates an instance of clusterHealthCollector.
func NewClusterHealthCollector(provider *gophercloud.ProviderClient, region string) *clusterHealthCollector {
	return &clusterHealthCollector{
		provider: provider,
		region:   region,
		HealthStatus: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "openstack_cluster_status",
				Help: "Health status of Cluster",
			},
		),
	}
}

func (c *clusterHealthCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		c.HealthStatus,
	}
}

func (c *clusterHealthCollector) collect() error {
	region := gophercloud.EndpointOpts{Region: c.region}
	computeClient, err := openstack.NewComputeV2(c.provider, region)
	if err != nil {
		return err
	}

    var down []string
    var up []string
    var found bool
	allPages, err := services.List(computeClient).AllPages()
	if err != nil {
		return err
	}

	allServices, err := services.ExtractServices(allPages)
	if err != nil {
		return err
	}

	for _, service := range allServices {
        if service.State == "down" {
            for _, s := range down {
                if service.State == s {
                    found = true
                    break
                }
            }
            if !found {
                down = append(down, service.Binary)
                found = false
            }
        }
        if service.State == "up" {
            for _, s := range up {
                if service.State == s {
                    found = true
                    break
                }
            }
            if !found {
                up = append(up, service.Binary)
                found = false
            }
        }
	}
    for _, down_s := range down {
        for _, up_s := range up{
            if up_s == down_s {
                c.HealthStatus.Set(2)
            }
        }
    }
    if len(down) != 0 {
        c.HealthStatus.Set(1)
    } else {
        c.HealthStatus.Set(0)
    }

	return nil
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by clusterHealthCollector.
func (c *clusterHealthCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.collectorList() {
		metric.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics.
func (c *clusterHealthCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Println("failed collecting cluster health metrics:", err)
	}

	for _, metric := range c.collectorList() {
		metric.Collect(ch)
	}
}
