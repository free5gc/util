// utils.go

package utils

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const (
	SuccessMetric = "successful"
	FailureMetric = "failure"
)

type MetricTypeEnabled string

const (
	SBI  MetricTypeEnabled = "sbi"
	NAS  MetricTypeEnabled = "nas"
	NGAP MetricTypeEnabled = "ngap"
)

var businessMetricsEnabled bool

func IsBusinessMetricsEnabled() bool {
	return businessMetricsEnabled
}

func EnableBusinessMetrics() {
	businessMetricsEnabled = true
}

func getValueFromCounter(c prometheus.Counter) (float64, error) {
	m := &dto.Metric{}
	if err := c.Write(m); err != nil {
		return 0, err
	}
	return m.GetCounter().GetValue(), nil
}

func GetCounterVecValue(counterName string, counter *prometheus.CounterVec, labels prometheus.Labels) (float64, error) {
	foundCounter, err := counter.GetMetricWith(labels)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve the %s counter, %v", counterName, err.Error())
	}

	counterValue, err := getValueFromCounter(foundCounter)
	if err != nil {
		return 0, fmt.Errorf("failed to get %s counter value, %v", counterName, err.Error())
	}

	return counterValue, nil
}

func GetStatus(metricStatusSuccess *bool) string {
	if metricStatusSuccess != nil && *metricStatusSuccess {
		return SuccessMetric
	}
	return FailureMetric
}

func FormatStatus(statusCode int) string {
	code := http.StatusInternalServerError
	if statusCode != 0 {
		code = statusCode
	}

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

// readStringPtr return the value of the string pointer if non-nil. Returns an empty string otherwise
func ReadStringPtr(strPtr *string) string {
	if strPtr == nil {
		return ""
	}
	return *strPtr
}
