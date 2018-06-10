package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/iGene/openstack-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	listenAddr  string
	metricsPath string
	configPath  string
	endpoint    string
	domain      string
	user        string
	password    string
	project     string
	region      string
)

func init() {
	pflag.StringVar(&listenAddr, "listen", ":9183", "<address>:<port> to listen on.")
	pflag.StringVar(&metricsPath, "telemetry-path", "/metrics", "Path under which to expose metrics.")
	pflag.StringVar(&configPath, "config", "", "Load the OpenStack config from path.")
	pflag.StringVar(&endpoint, "keystone-url", os.Getenv("OS_AUTH_URL"), "URL for the OpenStack Keystone API.")
	pflag.StringVar(&domain, "domain-name", os.Getenv("OS_DOMAIN_NAME"), "Domain name for the OpenStack Keystone.")
	pflag.StringVar(&user, "username", os.Getenv("OS_USERNAME"), "User for the OpenStack Keystone.")
	pflag.StringVar(&password, "password", os.Getenv("OS_PASSWORD"), "Password for the OpenStack Keystone.")
	pflag.StringVar(&project, "project-name", os.Getenv("OS_PROJECT_NAME"), "Project Name for the OpenStack Keystone.")
	pflag.StringVar(&region, "region-name", os.Getenv("OS_REGION_NAME"), "Region Name for the OpenStack Keystone.")
	pflag.Parse()
}

func checkFlags() {
	flags := []string{endpoint, domain, user, password, project}
	for _, f := range flags {
		if len(f) == 0 {
			fmt.Fprintf(os.Stderr, "Fatal error missing some flags:\n")
			pflag.PrintDefaults()
			os.Exit(1)
		}
	}
}

func loadOpenStackConfig(path string) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	endpoint = viper.GetString("global.endpoint")
	user = viper.GetString("global.username")
	password = viper.GetString("global.password")
	domain = viper.GetString("global.domain")
	project = viper.GetString("global.project")
	region = viper.GetString("global.region")
}

func main() {
	if len(configPath) != 0 {
		loadOpenStackConfig(configPath)
	}

	checkFlags()
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: endpoint,
		Username:         user,
		Password:         password,
		DomainName:       domain,
		TenantName:       project,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatalf("Fatal error autenticating: %s \n", err)
	}

	exporter := exporter.NewOpenStackExporter(provider, region)
	prometheus.MustRegister(exporter)
	if err != nil {
		log.Fatalf("Cannot export cluster")
	}

	http.Handle(metricsPath, promhttp.Handler())
	log.Printf("Starting exporter...")
	log.Fatalf("ListenAndServe error: %v", http.ListenAndServe(listenAddr, nil))
}
