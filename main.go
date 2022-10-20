package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kballard/go-shellquote"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"gopkg.in/alecthomas/kingpin.v2"

	kitlog "github.com/go-kit/kit/log"
)

const metricPrefix = "lvm_"

func main() {
	showVersion := kingpin.Flag("version", "Output version information and exit").Bool()
	listenAddress := kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests").Default(":9845").String()
	configFile := kingpin.Flag("web.config", "Path to config yaml file that can enable TLS or authentication").String()
	metricsPath := kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	disableExporterMetrics := kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself").Bool()
	cmd := kingpin.Flag("command", "Path to the LVM binary").Default("/usr/sbin/lvm").String()

	kingpin.Parse()

	if *showVersion {
		fmt.Println(version.Print("prometheus-lvm-exporter"))
		return
	}

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
