package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "herpstat"

var (
	systemLabelNames            = []string{"name"}
	systemSafetyRelayLabelNames = []string{"name", "relay"}
	systemInfoLabelNames        = []string{"name", "ip", "mac", "firmware", "outputs"}

	outputLabelNames      = []string{"system", "id"}
	outputInfoLabelNames  = []string{"system", "id", "name", "mode"}
	outputErrorLabelNames = []string{"system", "id", "error"}
)

type metrics struct {
	info                *prometheus.Desc
	temp                *prometheus.Desc
	resets              *prometheus.Desc
	safetyRelay         *prometheus.Desc
	outputInfo          *prometheus.Desc
	outputPower         *prometheus.Desc
	outputPowerLimit    *prometheus.Desc
	outputProbeTemp     *prometheus.Desc
	outputProbeHumidity *prometheus.Desc
	outputAlarmEnabled  *prometheus.Desc
	outputAlarmHigh     *prometheus.Desc
	outputAlarmLow      *prometheus.Desc
	outputRamping       *prometheus.Desc
	outputRampEnd       *prometheus.Desc
	outputError         *prometheus.Desc
}

// newOutputMetric is a convenience wrapper for [exporter.newMetric] that creates a new Prometheus desecriptor for
// an [exporter.output] metric with a given name, description, and optional labels.
// [exporter.outputLabelNames] is used as a default if no labels are given.
func newOutputMetric(name, description string, labels ...string) *prometheus.Desc {
	if labels == nil {
		labels = outputLabelNames
	}

	return newMetric("output", name, description, labels...)
}

// newSystemMetric is a convenience wrapper for [exporter.newMetric] that creates a new Prometheus desecriptor for
// an [exporter.system] metric with a given name, description, and optional labels.
// [exporter.systemLabelNames] is used as a default if no labels are given.
func newSystemMetric(name, description string, labels ...string) *prometheus.Desc {
	if labels == nil {
		labels = systemLabelNames
	}

	return newMetric("system", name, description, labels...)
}

// newMetric is a convenience wrapper for [prometheus.NewDesc] to create a new [prometheus.Desc] descriptor
// with a given susbsytem, name, description, and labels.
func newMetric(subsystem, name, description string, labels ...string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, name),
		description,
		labels,
		nil,
	)
}

func newMetrics() *metrics {
	return &metrics{
		info: newSystemMetric("info",
			"Information about the Herpstat system itself.",
			systemInfoLabelNames...,
		),
		temp: newSystemMetric("temp",
			"Current internal temperature.",
		),
		resets: newSystemMetric("reset_total",
			"Number of times Herpstat has reset",
		),
		safetyRelay: newSystemMetric("safetyrelay",
			"Safety relay status.",
			systemSafetyRelayLabelNames...,
		),
		outputInfo: newOutputMetric("info",
			"metadata about the output",
			outputInfoLabelNames...,
		),
		outputPower: newOutputMetric("power",
			"Current output power level.",
		),
		outputPowerLimit: newOutputMetric("power_limit",
			"Current output power limit.",
		),
		outputProbeTemp: newOutputMetric("probe_temperature",
			"Current probe temperature.",
		),
		outputProbeHumidity: newOutputMetric("probe_humidity",
			"Current probe humidity level.",
		),
		outputAlarmEnabled: newOutputMetric("alarm_enabled",
			"Output alarm enabled.",
		),
		outputAlarmHigh: newOutputMetric("alarm_high",
			"Output Alarm high value.",
		),
		outputAlarmLow: newOutputMetric("alarm_low",
			"Output Alarm low value.",
		),
		outputRamping: newOutputMetric("ramping",
			"Currently ramping?",
		),
		outputRampEnd: newOutputMetric("ramp_end",
			"Ramp end value.",
		),
		outputError: newOutputMetric("error",
			"Error Code.",
			outputErrorLabelNames...,
		),
	}
}
