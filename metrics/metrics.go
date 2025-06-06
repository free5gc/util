// Package metricsInfo sets and initializes Prometheus metricsInfo.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/utils"
)

// Init initializes all Prometheus metricsInfo
func Init(initMetrics InitMetrics) *prometheus.Registry {
	reg := prometheus.NewRegistry()

	globalLabels := prometheus.Labels{
		NF_TYPE_LABEL: initMetrics.GetNfName(),
	}

	prometheus.WrapRegistererWith(globalLabels, reg)
	// Uncomment to remove the basic go prometheus customCollectors.
	// wrappedReg.Unregister(customCollectors.NewProcessCollector(customCollectors.ProcessCollectorOpts{}))

	if len(initMetrics.GetCustomCollectors()) > 0 {
		utils.EnableBusinessMetrics()
		for _, customCollector := range initMetrics.GetCustomCollectors() {
			initMetric(customCollector, reg)
		}
	}

	return reg
}

func initMetric(metrics []prometheus.Collector, reg prometheus.Registerer) {
	for _, metric := range metrics {
		reg.MustRegister(metric)
	}
}
