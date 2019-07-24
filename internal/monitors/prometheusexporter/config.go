package prometheusexporter

import (
	dto "github.com/prometheus/client_model/go"
	"time"
)

// PrometheusConfig is the interface for configuring the prometheus exporter monitor.
type PrometheusConfig interface {
	NewPrometheusClient() (*PrometheusClient, error)
	GetInterval() time.Duration
}

// PrometheusClient is the prometheus exporter monitor client for scraping prometheus metrics.
type PrometheusClient struct {
	GetMetricFamilies func() ([]*dto.MetricFamily, error)
}
