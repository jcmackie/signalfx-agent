package apiserver

import (
	"io"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/signalfx/signalfx-agent/internal/core/common/kubernetes"
	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/prometheusexporter"
	"github.com/sirupsen/logrus"
)

func init() {
	monitors.Register(&monitorMetadata, func() interface{} { return &prometheusexporter.Monitor{} }, &Config{})
}

// Config is the config for this monitor and implements ConfigInterface.
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

// NewClient is a ConfigInterface method implementation that creates the prometheus client.
func (c *Config) NewClient() (*prometheusexporter.Client, error) {
	k8sClient, err := kubernetes.MakeClient(c.KubernetesAPI)
	if err != nil {
		return nil, err
	}
	return &prometheusexporter.Client{
		GetMetricFamilies: func() (metricFamilies []*dto.MetricFamily, err error) {
			var body io.ReadCloser
			defer func() {
				if body != nil {
					body.Close()
				}
			}()
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
		},
	}, nil
}

// GetInterval is a ConfigInterface method implementation for getting the configured monitor run interval.
func (c *Config) GetInterval() time.Duration {
	return time.Duration(c.IntervalSeconds) * time.Second
}

var loggingEntry *logrus.Entry

// GetLoggingEntry is a ConfigInterface method implementation for getting the logging entry.
func (c *Config) GetLoggingEntry() *logrus.Entry {
	if loggingEntry == nil {
		loggingEntry = logrus.WithFields(logrus.Fields{"monitorType": monitorType})
	}
	return loggingEntry
}
