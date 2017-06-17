package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"time"

	"github.com/czerwonk/atlas_exporter/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "0.5.2"

var (
	showVersion          = flag.Bool("version", false, "Print version information.")
	listenAddress        = flag.String("web.listen-address", ":9400", "Address on which to expose metrics and web interface.")
	metricsPath          = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	filterInvalidResults = flag.Bool("filter.invalid-results", true, "Exclude offline/incompatible probes")
	cacheTtl             = flag.Int("cache.ttl", 3600, "Cache time to live in seconds")
	cacheCleanUp         = flag.Int("cache.cleanup", 300, "Interval for cache clean up in seconds")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: atlas_exporter [ ... ]\n\nParameters:")
		fmt.Println()
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
	log.Infof("Starting atlas exporter (Version: %s)\n", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>RIPE Atlas Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>RIPE Atlas Exporter</h1>
			<h2>Example</h2>
			<p>Metrics for measurement with id 8809582:</p>
			<p><a href="` + *metricsPath + `?measurement_id=8809582">` + r.Host + *metricsPath + `?measurement_id=8809582</a></p>
			<h2>More Information</h2>
			<p><a href="https://github.com/czerwonk/atlas_exporter">github.com/czerwonk/atlas_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	log.Infof("Cache TTL: %v\n", time.Duration(*cacheTtl)*time.Second)
	log.Infof("Cache cleanup interval (seconds): %v\n", time.Duration(*cacheCleanUp)*time.Second)
	initCache()

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			log.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("measurement_id")

	if len(id) == 0 {
		return errors.New("Parameter measurement_id has to be defined")
	}

	metrics, err := getMeasurement(id)

	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		reg := prometheus.NewRegistry()
		reg.MustRegister(metric.NewMetricCollector(id, metrics))

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
	}

	return nil
}
