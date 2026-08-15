// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/czerwonk/atlas_exporter/atlas"
	"github.com/czerwonk/atlas_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	_ "net/http/pprof"
)

const (
	generalTimeout = 60 * time.Second
	streamTimeout  = 5 * time.Minute
	version        = "1.0.6"
)

var (
	showVersion         = flag.Bool("version", false, "Print version information.")
	listenAddress       = flag.String("web.listen-address", ":9400", "Address on which to expose metrics and web interface.")
	metricsPath         = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	cacheTTL            = flag.Int("cache.ttl", 3600, "Cache time to live in seconds")
	cacheCleanUp        = flag.Int("cache.cleanup", 300, "Interval for cache clean up in seconds")
	configFile          = flag.String("config.file", "", "Path to congig file to use")
	timeout             = flag.Duration("timeout", generalTimeout, "Timeout")
	workerCount         = flag.Uint("worker.count", 8, "Number of go routines retrieving probe information")
	streaming           = flag.Bool("streaming", true, "Retrieve data by subscribing to Atlas Streaming API")
	streamingBufferSize = flag.Uint("streaming.buffer-size", 100, "Size of buffer to prevent locking socket.io go routines")
	streamingTimeout    = flag.Duration("streaming.timeout", streamTimeout, "When no update is received in this timespan a reconnect is initiated.")
	profiling           = flag.Bool("profiling", false, "Enables pprof endpoints")
	goMetrics           = flag.Bool("metrics.go", true, "Enables go runtime prometheus metrics")
	processMetrics      = flag.Bool("metrics.process", true, "Enables process runtime prometheus metrics")
	logLevel            = flag.String("log.level", "info", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]")
	tlsEnabled          = flag.Bool("tls.enabled", false, "Enables TLS")
	tlsCertChainPath    = flag.String("tls.cert-file", "", "Path to TLS cert file")
	tlsKeyPath          = flag.String("tls.key-file", "", "Path to TLS key file")
	cfg                 *config.Config
	strategy            atlas.Strategy
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

	setLogLevel(*logLevel)

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	err := loadConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	if *streaming {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		strategy = atlas.NewStreamingStrategy(ctx, cfg, *streamingBufferSize, *streamingTimeout)
	} else {
		strategy = atlas.NewRequestStrategy(cfg, *workerCount)
	}

	if !*profiling {
		http.DefaultServeMux = http.NewServeMux()
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

func loadConfig() error {
	if len(*configFile) == 0 {
		cfg = &config.Config{}
		return nil
	}

	b, err := os.ReadFile(*configFile)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}

	c, err := config.Load(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("could not parse config file: %v", err)
	}
	cfg = c

	return nil
}

func startServer() {
	log.Infof("Starting atlas exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>RIPE Atlas Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>RIPE Atlas Exporter</h1>
			<h2>Example</h2>
			<p>Metrics for measurement configured in configuration file:</p>
			<p><a href="` + *metricsPath + `">` + r.Host + *metricsPath + `</a></p>
			<p>Metrics for measurement with id 8809582:</p>
			<p><a href="` + *metricsPath + `?measurement_id=8809582">` + r.Host + *metricsPath + `?measurement_id=8809582</a></p>
			<h2>More information</h2>
			<p><a href="https://github.com/czerwonk/atlas_exporter">github.com/czerwonk/atlas_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	log.Infof("Cache TTL: %v", time.Duration(*cacheTTL)*time.Second)
	log.Infof("Cache cleanup interval: %v", time.Duration(*cacheCleanUp)*time.Second)
	atlas.InitCache(time.Duration(*cacheTTL)*time.Second, time.Duration(*cacheCleanUp)*time.Second)

	log.Infof("Listening for %s on %s (TLS: %v)", *metricsPath, *listenAddress, *tlsEnabled)
	if *tlsEnabled {
		log.Fatal(http.ListenAndServeTLS(*listenAddress, *tlsCertChainPath, *tlsKeyPath, nil))
		return
	}

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

	s := strategy

	ids := []string{}
	if len(id) > 0 {
		ids = append(ids, id)
		s = atlas.NewRequestStrategy(cfg, *workerCount)
	} else {
		ids = append(ids, cfg.MeasurementIDs()...)
	}

	if len(ids) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	measurements, err := s.MeasurementResults(ctx, ids)
	if err != nil {
		return err
	}

	reg := prometheus.NewRegistry()

	// add process metrics
	if *processMetrics {
		processCollector := collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
		reg.MustRegister(processCollector)
	}

	// add go collector metrics
	if *goMetrics {
		goCollector := collectors.NewGoCollector()
		reg.MustRegister(goCollector)
	}

	if len(measurements) > 0 {
		c := newCollector(measurements)
		reg.MustRegister(c)
	}

	l := log.New()
	l.Level = log.ErrorLevel

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      l,
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)

	return nil
}
