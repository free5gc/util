package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/free5gc/util/metrics/utils"
)

const (
	NF_TYPE_LABEL = "nf_type"
)

type InitMetrics struct {
	metricsInfo      Metrics
	nfName           string
	metricsEnabled   map[utils.MetricTypeEnabled]bool
	customCollectors map[utils.MetricTypeEnabled][]prometheus.Collector
}

func NewInitMetrics(
	_metrics Metrics,
	_nfName string,
	_metricsEnabled map[utils.MetricTypeEnabled]bool,
	_customCollectors map[utils.MetricTypeEnabled][]prometheus.Collector,
) InitMetrics {
	return InitMetrics{
		metricsInfo:      _metrics,
		nfName:           _nfName,
		metricsEnabled:   _metricsEnabled,
		customCollectors: _customCollectors,
	}
}

func (initMetrics InitMetrics) GetMetricsInfo() Metrics {
	return initMetrics.metricsInfo
}

func (initMetrics InitMetrics) GetNfName() string {
	return initMetrics.nfName
}

func (initMetrics InitMetrics) GetMetricsEnabled() map[utils.MetricTypeEnabled]bool {
	return initMetrics.metricsEnabled
}

func (initMetrics InitMetrics) GetCustomCollectors() map[utils.MetricTypeEnabled][]prometheus.Collector {
	return initMetrics.customCollectors
}

type Metrics struct {
	Scheme      string `yaml:"scheme"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty"`
	Tls         Tls    `yaml:"tls,omitempty"`
	Namespace   string `yaml:"namespace"`
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}
