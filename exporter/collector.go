package exporter

import (
	"fmt"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// Describes all of the metric types that we're exporting.
// Declaring this (along with [exporter.Collect]) implements a [prometheus.Collector]
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.info
	ch <- e.metrics.temp
	ch <- e.metrics.resets
	ch <- e.metrics.outputPower
	ch <- e.metrics.outputPowerLimit
	ch <- e.metrics.outputProbeTemp
	ch <- e.metrics.outputProbeHumidity
	ch <- e.metrics.outputAlarmEnabled
	ch <- e.metrics.outputAlarmHigh
	ch <- e.metrics.outputAlarmLow
	ch <- e.metrics.outputRamping
	ch <- e.metrics.outputRampEnd
	ch <- e.metrics.outputError
}

// Polls a Herpstat SpyderWeb, then sends the relevant data back to Prometheus via a channel.
// Declaring this (along with [exporter.Describe]) implements a [prometheus.Collector].
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	level.Debug(logger).Log("msg", fmt.Sprintf("%s was called", *webTelemetryPath))

	if !e.herpstat.poll() {
		level.Warn(logger).Log("msg", "Returning previously cached data.")
	}

	ch <- newCounterMetric(e.metrics.info, 1, e.herpstat.info.system.infoLabelValues()...)
	ch <- newGaugeMetric(e.metrics.temp, e.herpstat.info.system.Temp, e.herpstat.info.system.labelValues()...)
	ch <- newGaugeMetric(e.metrics.safetyRelay, e.herpstat.info.system.safetyrelay(), e.herpstat.info.system.safetyRelayLabelValues()...)
	ch <- newGaugeMetric(e.metrics.resets, e.herpstat.info.system.PowerResets, e.herpstat.info.system.labelValues()...)

	for i := range *e.herpstat.info.outputs {
		output := &(*e.herpstat.info.outputs)[i]
		systemName := e.herpstat.info.system.Name

		ch <- newCounterMetric(e.metrics.outputInfo, 1, output.infoLabelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputPower, output.Power, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputPowerLimit, output.PowerLimit, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputProbeTemp, output.ProbeTemp, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputProbeHumidity, output.ProbeHumidity, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputAlarmEnabled, output.AlarmEnabled, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputAlarmHigh, output.AlarmHigh, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputAlarmLow, output.AlarmLow, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputRamping, output.ramping(), output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputRampEnd, output.RampEnd, output.labelValues(&systemName)...)
		ch <- newGaugeMetric(e.metrics.outputError, output.ErrorCode, output.errorLabelValues(&systemName)...)
	}
}

// create a new [prometheus.CounterValue] metric
// We aren't manually counting anything ourselves - these will be used to keep track of metadata.
// Their metrics will have a value of 1 and their labels will have all of the relevant info.
func newCounterMetric(desc *prometheus.Desc, value float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(desc, prometheus.CounterValue, value, labelValues...)
}

// creates a new [prometheus.GaugeValue] metric
func newGaugeMetric(desc *prometheus.Desc, value float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, labelValues...)
}
