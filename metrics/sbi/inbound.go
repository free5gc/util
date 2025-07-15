package sbi

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/utils"
)

func GetSbiInboundMetrics(namespace string) []prometheus.Collector {
	var metrics []prometheus.Collector

	InboundReqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      INBOUND_REQ_COUNTER_NAME,
			Help:      INBOUND_REQ_COUNTER_DESC,
		},
		[]string{IN_STATUS_CODE_LABEL, IN_METHOD_LABEL, IN_CAUSE_LABEL, IN_PATH_LABEL},
	)

	metrics = append(metrics, InboundReqCounter)

	InboundRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      INBOUND_REQ_HIST_NAME,
			Help:      INBOUND_REQ_HIST_DESC,
			Buckets: []float64{
				0.0001,
				0.0050,
				0.0100,
				0.0200,
				0.0250,
				0.0500,
			},
		},
		[]string{IN_METHOD_LABEL, IN_PATH_LABEL, IN_STATUS_CODE_LABEL},
	)

	metrics = append(metrics, InboundRequestDuration)

	return metrics
}

func IncrInboundReqCounter(method string, path string, statusCode int, cause string) {
	if IsSbiMetricsEnabled() {
		InboundReqCounter.With(prometheus.Labels{
			IN_STATUS_CODE_LABEL: utils.FormatStatus(statusCode),
			IN_METHOD_LABEL:      method,
			IN_CAUSE_LABEL:       cause,
			IN_PATH_LABEL:        path,
		}).Add(1)
	}
}

func IncrInboundReqDurationCounter(method string, path string, statusCode int, duration float64) {
	if IsSbiMetricsEnabled() {
		InboundRequestDuration.With(prometheus.Labels{
			IN_PATH_LABEL:        path,
			IN_METHOD_LABEL:      method,
			IN_STATUS_CODE_LABEL: utils.FormatStatus(statusCode),
		}).Observe(duration)
	}
}
