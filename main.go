package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kballard/go-shellquote"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"

	kitlog "github.com/go-kit/kit/log"
)

const metricPrefix = "lvm_"

func main() {
	listenAddress := flag.String("web.listen-address", ":9845", "The address to listen on for HTTP requests")
	configFile := flag.String("web.config", "", "Path to config yaml file that can enable TLS or authentication")
	metricsPath := flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	disableExporterMetrics := flag.Bool("web.disable-exporter-metrics", false, "Exclude metrics about the exporter itself")
	cmd := flag.String("command", "/usr/sbin/lvm", "Path to the LVM binary")

	flag.Parse()

	cmdParts, err := shellquote.Split(*cmd)
	if err != nil {
		log.Fatalf("Parsing command failed: %v", err)
	}

	registry := prometheus.NewPedanticRegistry()

	lvmRegistry := prometheus.WrapRegistererWithPrefix(metricPrefix, registry)
	lvmRegistry.MustRegister(newCommandCollector(cmdParts))

	if !*disableExporterMetrics {
		registry.MustRegister(
			prometheus.NewBuildInfoCollector(),
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
			version.NewCollector("lvm_exporter"),
		)
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		MaxRequestsInFlight: 3,
	}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>LVM exporter</title></head>
			<body>
			<h1>LVM exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Listening on %q", *listenAddress)

	server := &http.Server{Addr: *listenAddress}
	logger := kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})

	if err := web.ListenAndServe(server, *configFile, logger); err != nil {
		log.Fatal(err)
	}
}
