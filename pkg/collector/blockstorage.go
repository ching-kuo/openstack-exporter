package collector

import (
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/prometheus/client_golang/prometheus"
)

// blockStorageCollector collects statistics about Cinder in an OpenStack Cluster.
type blockStorageCollector struct {
	provider gophercloud.ProviderClient

	region string

	TotalVolumeSize prometheus.Gauge

	TotalVolumeNumber prometheus.Gauge
}

// NewBlockStorageCollector creates an instance of blockStorageCollector.
func NewBlockStorageCollector(provider gophercloud.ProviderClient, region string) *blockStorageCollector {
	return &blockStorageCollector{
		provider: provider,
		region:   region,
		TotalVolumeSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "openstack_total_volume_size",
				Help: "Number of total size of volumes in GB",
			},
		),
		TotalVolumeNumber: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "openstack_total_volume_number",
				Help: "Number of total volume",
			},
		),
	}
}

func (b *blockStorageCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		b.TotalVolumeNumber,
		b.TotalVolumeSize,
	}
}

func (b *blockStorageCollector) collect() error {
	region := gophercloud.EndpointOpts{Region: b.region}
	blockStorageClient, err := openstack.NewBlockStorageV3(&b.provider, region)
	if err != nil {
		return err
	}

	opts := volumes.ListOpts{
		AllTenants: true,
	}

	var size float64
	var number float64

	allPages, err := volumes.List(blockStorageClient, opts).AllPages()
	if err != nil {
		return err
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return err
	}

	for _, volume := range allVolumes {
		size += float64(volume.Size)
		number++
	}
	b.TotalVolumeSize.Set(size)
	b.TotalVolumeNumber.Set(number)

	return nil
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by blockStorageCollector.
func (b *blockStorageCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range b.collectorList() {
		metric.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics.
func (b *blockStorageCollector) Collect(ch chan<- prometheus.Metric) {

	if err := b.collect(); err != nil {
		log.Println("failed collecting block storage metrics:", err)
	}

	for _, metric := range b.collectorList() {
		metric.Collect(ch)
	}

}
