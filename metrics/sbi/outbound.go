package sbi

// Outbound
import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/utils"
)

func GetSbiOutboundMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	OutboundReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      OUTBOUND_REQ_COUNTER_NAME,
			Help:      OUTBOUND_REQ_COUNTER_DESC,
		},
		[]string{OUT_TARGET_SERVICE_NAME_LABEL, OUT_STATUS_CODE_LABEL, OUT_METHOD_LABEL},
	)

	metrics = append(metrics, OutboundReqCounter)

	OutboundRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      OUTBOUND_REQ_HIST_NAME,
			Help:      OUTBOUND_REQ_HISTOGRAM_DESC,
			Buckets: []float64{
				0.0001,
				0.0050,
				0.0100,
				0.0200,
				0.0250,
				0.0500,
			},
		},
		[]string{OUT_METHOD_LABEL, OUT_TARGET_SERVICE_NAME_LABEL, OUT_STATUS_CODE_LABEL},
	)

	metrics = append(metrics, OutboundRequestDuration)

	return metrics
}

func incrOutboundReqCounter(method string, serviceName string, statusCode int) {
	if IsSbiMetricsEnabled() {
		OutboundReqCounter.With(prometheus.Labels{
			OUT_TARGET_SERVICE_NAME_LABEL: serviceName,
			OUT_STATUS_CODE_LABEL:         utils.FormatStatus(statusCode),
			OUT_METHOD_LABEL:              method,
		}).Add(1)
	}
}

func incrOutboundReqDurationCounter(method string, serviceName string, statusCode int, duration float64) {
	if IsSbiMetricsEnabled() {
		OutboundRequestDuration.With(prometheus.Labels{
			OUT_TARGET_SERVICE_NAME_LABEL: serviceName,
			OUT_METHOD_LABEL:              method,
			OUT_STATUS_CODE_LABEL:         utils.FormatStatus(statusCode),
		}).Observe(duration)
	}
}

func SbiMetricHook(method string, serviceName string, statusCode int, duration float64) {
	incrOutboundReqCounter(method, serviceName, statusCode)
	incrOutboundReqDurationCounter(method, serviceName, statusCode, duration)
}
