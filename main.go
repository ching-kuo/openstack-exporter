package main

import (
    "fmt"
    "net/http"
    "log"
    "sync"

    "github.com/iGene/openstack-exporter/collectors"
    "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/spf13/viper"
)

// Wraps all collectors in a single exporter to extract metrics and make sure it is
// thread safe.

type OpenStackExporter struct {

    mu         sync.Mutex
    collectors []prometheus.Collector

}

// verify that the exporter implementation is correct

var _ prometheus.Collector = &OpenStackExporter{}

// NewOpenStackExporter creates an instance to OpenStackExporter and returns a 
// reference to it.

func NewOpenStackExporter(provider *gophercloud.ProviderClient) *OpenStackExporter {
    return &OpenStackExporter{
        collectors: []prometheus.Collector{
            collectors.NewComputeCollector(*provider),
            collectors.NewBlockStorageCollector(*provider),
            collectors.NewNetworkCollector(*provider),
        },
    }
}

// Collect is called by the Prometheus registry when collecting
// metrics.

func (o *OpenStackExporter) Collect(ch chan<- prometheus.Metric) {
	o.mu.Lock()
	defer o.mu.Unlock()

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

func main(){

    viper.SetConfigName("openstack")
    viper.AddConfigPath("/etc/openstack-exporter/")
    err := viper.ReadInConfig()
    viper.SetConfigType("toml")
    if err != nil { // Handle errors reading the config file
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
    } 
    endpoint := viper.GetString("global.endpoint")
    username := viper.GetString("global.username")
    password := viper.GetString("global.password")

    opts := gophercloud.AuthOptions{
        IdentityEndpoint: endpoint,
        Username: username,
        Password: password,
        DomainName: "default",
    }

    provider, err := openstack.AuthenticatedClient(opts)

    if err != nil {
	    panic(fmt.Errorf("Fatal error autenticating: %s \n", err))
    }

    prometheus.MustRegister(NewOpenStackExporter(provider))
    if err != nil {
        log.Fatalf("cannot export cluster")
    }

    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe(":9183", nil))
}
