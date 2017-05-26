package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/czerwonk/atlas_exporter/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const version string = "0.3.0"

var (
	showVersion          = flag.Bool("version", false, "Print version information.")
	listenAddress        = flag.String("web.listen-address", ":9400", "Address on which to expose metrics and web interface.")
	metricsPath          = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	filterInvalidResults = flag.Bool("filter.invalid-results", true, "Exclude offline/incompatible probes")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: atlas_exporter [ ... ]\n\nParameters:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("atlas_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Metric exporter for RIPE Atlas measurements")
	fmt.Println("This software uses Go bindings from the DNS-OARC project (https://github.com/DNS-OARC/ripeatlas)")
}

func startServer() {
	fmt.Printf("Starting atlas exporter (Version: %s)\n", version)
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	fmt.Printf("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("measurement_id")

	if len(id) == 0 {
		return errors.New("Parameter measurement_id has to be defined.")
	}

	metrics, err := getMeasurement(id)

	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		reg := prometheus.NewRegistry()
		reg.MustRegister(metric.NewMetricCollector(id, metrics))

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	}

	return nil
}
