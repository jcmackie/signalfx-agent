package apiserver

import (
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/signalfx/signalfx-agent/internal/core/common/kubernetes"
	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/prometheusexporter"
	"io"
	"time"
)

func init() {
	monitors.Register(&monitorMetadata, func() interface{} { return &prometheusexporter.Monitor{} }, &Config{})
}

// Config is the config for this monitor and implements interface PrometheusConfig.
type Config struct {
	config.MonitorConfig
	// Configuration of the Kubernetes API client.
	KubernetesAPI *kubernetes.APIConfig `yaml:"kubernetesAPI" default:"{}"`
	// Path to the metrics endpoint on the exporter server, usually `/metrics` (the default).
	MetricPath string `yaml:"metricPath" default:"/metrics"`
}

// Validate k8s-specific configuration.
func (c *Config) Validate() error {
	return c.KubernetesAPI.Validate()
}

// NewPrometheusClient is a PrometheusConfig interface method implementation that creates the prometheus client.
func (c *Config) NewPrometheusClient() (*prometheusexporter.PrometheusClient, error) {
	k8sClient, err := kubernetes.MakeClient(c.KubernetesAPI)
	if err != nil {
		return nil, err
	}
	return  &prometheusexporter.PrometheusClient{
		GetMetricFamilies: func() (metricFamilies []*dto.MetricFamily, err error) {
			var body io.ReadCloser
			defer func() { if body != nil {body.Close()} }()
			if body, err = k8sClient.CoreV1().RESTClient().Get().RequestURI(c.MetricPath).Stream(); err != nil {
				return
			}
			decoder := expfmt.NewDecoder(body, expfmt.FmtText)
			metricFamilies = make([]*dto.MetricFamily, 0)
			for {
				var mf dto.MetricFamily
				if err = decoder.Decode(&mf); err != nil || err == io.EOF {
					return
				}
				metricFamilies = append(metricFamilies, &mf)
			}
			return
		},
	}, nil
}

// GetInterval is a PrometheusConfig interface method implementation for getting the configured monitor run interval.
func (c *Config) GetInterval() time.Duration {
	return time.Duration(c.IntervalSeconds)*time.Second
}
