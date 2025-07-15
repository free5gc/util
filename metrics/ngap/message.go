package ngap

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/ngap/ngapType"
	utils "github.com/free5gc/util/metrics/utils"
)

var (
	MsgRcvCounter  *prometheus.CounterVec
	MsgSentCounter *prometheus.CounterVec
)

func GetNgapHandlerMetrics(namespace string) []prometheus.Collector {
	var collectors []prometheus.Collector

	MsgRcvCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      MSG_RCV_COUNTER_NAME,
			Help:      MSG_RCV_COUNTER_DESC,
		},
		[]string{NAME_LABEL, STATUS_LABEL, CAUSE_LABEL},
	)

	collectors = append(collectors, MsgRcvCounter)

	MsgSentCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: SUBSYSTEM_NAME,
			Name:      MSG_SENT_COUNTER_NAME,
			Help:      MSG_SENT_COUNTER_DESC,
		},
		[]string{NAME_LABEL, STATUS_LABEL, CAUSE_LABEL},
	)

	collectors = append(collectors, MsgSentCounter)

	return collectors
}

func IncrMetricsRcvMsg(msgType string, metricStatusSuccess *bool, syntaxCause *ngapType.Cause) {
	if IsNgapMetricsEnabled() {
		msgCause := ""

		if syntaxCause != nil && syntaxCause.Present != 0 {
			msgCause = GetCauseErrorStr(syntaxCause)
		}

		MsgRcvCounter.With(prometheus.Labels{
			NAME_LABEL:   msgType,
			STATUS_LABEL: utils.GetStatus(metricStatusSuccess),
			CAUSE_LABEL:  msgCause,
		}).Add(1)
	}
}

func IncrMetricsSentMsg(msgType string, metricStatusSuccess *bool, syntaxCause ngapType.Cause, otherCause *string) {
	if IsNgapMetricsEnabled() {
		msgCause := ""
		causeErrStr := GetCauseErrorStr(&syntaxCause)

		switch {
		case causeErrStr != UNKNOWN_NGAP_TYPE_CAUSE_ERR:
			msgCause = causeErrStr
		case otherCause != nil:
			msgCause = *otherCause
		}

		MsgSentCounter.With(prometheus.Labels{
			NAME_LABEL:   msgType,
			STATUS_LABEL: utils.GetStatus(metricStatusSuccess),
			CAUSE_LABEL:  msgCause,
		}).Add(1)
	}
}
