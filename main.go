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

	kitlog "github.com/go-kit/log"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const metricPrefix = "lvm_"

func main() {
	showVersion := kingpin.Flag("version", "Output version information and exit").Bool()
	webConfig := webflag.AddFlags(kingpin.CommandLine, ":9845")
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

	server := &http.Server{}
	logger := kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})

	if err := web.ListenAndServe(server, webConfig, logger); err != nil {
		log.Fatal(err)
	}
}
