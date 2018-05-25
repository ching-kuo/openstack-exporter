package main

import (
    "fmt"
    "net/http"
    "log"
    "sync"
    "./collectors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/spf13/viper"
)

type OpenStackExporter struct {

    mu         sync.Mutex
    collectors []prometheus.Collector

}

var _ prometheus.Collector = &OpenStackExporter{}

func NewOpenStackExporter(provider *gophercloud.ProviderClient) *OpenStackExporter {
    return &OpenStackExporter{
        collectors: []prometheus.Collector{
            collectors.NewComputeCollector(*provider),
            collectors.NewBlockStorageCollector(*provider),
        },
    }
}

func (o *OpenStackExporter) Collect(ch chan<- prometheus.Metric) {
	o.mu.Lock()
	defer o.mu.Unlock()

	for _, oo := range o.collectors {
		oo.Collect(ch)
	}
}

func (o *OpenStackExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, oo := range o.collectors {
		oo.Describe(ch)
	}
}

func main(){

    viper.SetConfigName("openstack")
    viper.AddConfigPath(".")
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
