package sbi

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	OUTBOUND_REQ_COUNTER_NAME = "outbound_request_total"
	OUTBOUND_REQ_COUNTER_DESC = "Total number of SBI outbound requests attempted or sent by the NF"

	OUTBOUND_REQ_HIST_NAME      = "outbound_request_duration_seconds"
	OUTBOUND_REQ_HISTOGRAM_DESC = "Histogram of request latencies"
)

const (
	INBOUND_REQ_COUNTER_NAME = "inbound_request_total"
	INBOUND_REQ_COUNTER_DESC = "Total number of SBI inbound requests received by the NF"

	INBOUND_REQ_HIST_NAME = "inbound_request_duration_seconds"
	INBOUND_REQ_HIST_DESC = "Histogram of request latencies"
)

const (
	SUBSYSTEM_NAME = "sbi"
)

var (
	OutboundReqCounter      *prometheus.CounterVec
	OutboundRequestDuration *prometheus.HistogramVec
	InboundReqCounter       *prometheus.CounterVec
	InboundRequestDuration  *prometheus.HistogramVec
)

// Labels names for the outbound sbi metrics
const (
	OUT_TARGET_SERVICE_NAME_LABEL = "target_service_name"
	OUT_STATUS_CODE_LABEL         = "status_code"
	OUT_METHOD_LABEL              = "method"
)

// Labels names for the inbound sbi metrics
const (
	IN_STATUS_CODE_LABEL  = "status_code"
	IN_METHOD_LABEL       = "method"
	IN_CAUSE_LABEL        = "cause"
	IN_PATH_LABEL         = "path"
	IN_PB_DETAILS_CTX_STR = "problem"
)

type OutboundMetricBasicInfo struct {
	StatusCode        int     `json:"status_code"`
	TargetServiceName string  `json:"target_service_name"`
	Method            string  `json:"method"`
	Duration          float64 `json:"duration"`
}

var sbiMetricsEnabled bool

func IsSbiMetricsEnabled() bool {
	return sbiMetricsEnabled
}

func EnableSbiMetrics() {
	sbiMetricsEnabled = true
}
