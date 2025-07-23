// Package metrics sets and initializes Prometheus metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/nas"
	"github.com/free5gc/util/metrics/ngap"
	"github.com/free5gc/util/metrics/sbi"
	"github.com/free5gc/util/metrics/utils"
)

// Init initializes all Prometheus metrics
func Init(initMetrics InitMetrics) *prometheus.Registry {
	reg := prometheus.NewRegistry()

	globalLabels := prometheus.Labels{
		NF_TYPE_LABEL: initMetrics.GetNfName(),
	}

	wrappedReg := prometheus.WrapRegistererWith(globalLabels, reg)
	// Uncomment to remove the basic go prometheus collectors.
	// wrappedReg.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	if len(initMetrics.GetCustomCollectors()) > 0 {
		utils.EnableBusinessMetrics()
		for _, customCollector := range initMetrics.GetCustomCollectors() {
			initMetric(customCollector, reg)
		}
	}

	if initMetrics.GetMetricsEnabled()[utils.SBI] {
		addSBIToRegistry(initMetrics.GetMetricsInfo().Namespace, wrappedReg)
	}

	if initMetrics.GetMetricsEnabled()[utils.NAS] {
		addNASToRegistry(initMetrics.GetMetricsInfo().Namespace, wrappedReg)
	}

	if initMetrics.GetMetricsEnabled()[utils.NGAP] {
		addNGAPToRegistry(initMetrics.GetMetricsInfo().Namespace, wrappedReg)
	}

	return reg
}

func addSBIToRegistry(namespace string, reg prometheus.Registerer) {
	var sbiMetrics []prometheus.Collector

	sbiMetrics = append(sbiMetrics, sbi.GetSbiOutboundMetrics(namespace)...)
	sbiMetrics = append(sbiMetrics, sbi.GetSbiInboundMetrics(namespace)...)

	initMetric(sbiMetrics, reg)
	sbi.EnableSbiMetrics()
}

func addNASToRegistry(namespace string, reg prometheus.Registerer) {
	var nasMetrics []prometheus.Collector

	nasMetrics = append(nasMetrics, nas.GetNasHandlerMetrics(namespace)...)

	initMetric(nasMetrics, reg)
	nas.EnableNasMetrics()
}

func addNGAPToRegistry(namespace string, reg prometheus.Registerer) {
	var ngapMetrics []prometheus.Collector

	ngapMetrics = append(ngapMetrics, ngap.GetNgapHandlerMetrics(namespace)...)

	initMetric(ngapMetrics, reg)
	ngap.EnableNgapMetrics()
}

func initMetric(metrics []prometheus.Collector, reg prometheus.Registerer) {
	for _, metric := range metrics {
		reg.MustRegister(metric)
	}
}
