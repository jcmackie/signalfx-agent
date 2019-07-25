package prometheusexporter

import (
	"time"

	dto "github.com/prometheus/client_model/go"
)

// ConfigInterface is the interface for configuring the prometheus exporter monitor.
type ConfigInterface interface {
	NewClient() (*Client, error)
	GetInterval() time.Duration
}

// Client is the prometheus exporter monitor client for scraping prometheus metrics.
type Client struct {
	GetMetricFamilies func() ([]*dto.MetricFamily, error)
}
