package exporter

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-kit/log/term"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

const (
	defaultListenAddress    = ":10010"
	defaultWebTelemetryPath = "/metrics"
	httpReadTimeout         = 12 * time.Second
)

var (
	debug = kingpin.Flag(
		"debug",
		"Enable debug logging. It's very noisy!",
	).Default("false").Bool()
	herpstatAddress = kingpin.Flag(
		"herpstat.address",
		"Your Herpstat SpyderWeb's address.",
	).Required().PlaceHolder("1.2.3.4").String()
	webDisableExporterMetrics = kingpin.Flag(
		"web.disable-exporter-metrics",
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	).Default("true").Bool()
	webTelemetryPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default(defaultWebTelemetryPath).String()

	webFlags = kingpinflag.AddFlags(kingpin.CommandLine, defaultListenAddress)

	logger log.Logger
)

// Exporter is consumed via [prometheus.MustRegister] and is used to generate and collect metrics
type Exporter struct {
	herpstat *herpstat
	metrics  *metrics
}

// Run function starts an HTTP server that listens on [exporter.webTelemetryPath] and exposes
// Prometheus metrics.
func Run() {
	kingpin.CommandLine.DefaultEnvars()
	kingpin.Parse()

	logger = term.NewColorLogger(os.Stdout, log.NewLogfmtLogger, logColors)
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	if *debug {
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "caller", log.DefaultCaller)
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	level.Info(logger).Log("msg", "Starting Herpstat SpyderWeb Exporter")
	level.Info(logger).Log("msg", "Herpstat URL", "url", fmt.Sprintf(rawstatusURL, *herpstatAddress))

	// create a new, clean prometheus registry without any exporter metrics
	registry := prometheus.NewRegistry()
	registry.MustRegister(&Exporter{
		herpstat: newHerpstat(),
		metrics:  newMetrics(),
	})

	// add the exporter metrics if requested
	if *debug || !*webDisableExporterMetrics {
		registry.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)
	}

	if *debug {
		registry.MustRegister(collectors.NewBuildInfoCollector())
	}

	http.Handle(*webTelemetryPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{ReadTimeout: httpReadTimeout}
	if err := web.ListenAndServe(server, webFlags, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
