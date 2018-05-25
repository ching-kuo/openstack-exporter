package collectors

import (
    "log"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

type blockStorageCollector struct{
    provider gophercloud.ProviderClient

    TotalVolumeSize prometheus.Gauge

    TotalVolumeNumber prometheus.Gauge
}

func NewBlockStorageCollector(provider gophercloud.ProviderClient) *blockStorageCollector{
    return &blockStorageCollector{
        provider: provider,

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

func (b *blockStorageCollector) collect() error{
    region := gophercloud.EndpointOpts{Region: "RegionOne"}
    blockStorageClient, err := openstack.NewBlockStorageV3(&b.provider,region)
    if err != nil {
        return err
    }

    opts := volumes.ListOpts{
        AllTenants: true,
    }

    var size float64 = 0
    var number float64 = 0

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
        number += 1
    }
    b.TotalVolumeSize.Set(size)
    b.TotalVolumeNumber.Set(number)

    return nil
}

func (b *blockStorageCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range b.collectorList() {
		metric.Describe(ch)
	}
}

func (b *blockStorageCollector) Collect(ch chan<- prometheus.Metric) {

	if err := b.collect(); err != nil {
		log.Println("failed collecting compute metrics:", err)
	}

	for _, metric := range b.collectorList() {
		metric.Collect(ch)
	}

}
