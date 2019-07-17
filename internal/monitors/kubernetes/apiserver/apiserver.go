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

// Config for the kubernetes apiserver monitor. Config implements interface PrometheusConfig.
type Config struct {
	config.MonitorConfig
	// Configuration of the Kubernetes API client
	KubernetesAPI *kubernetes.APIConfig `yaml:"kubernetesAPI" default:"{}"`
	// Path to the metrics endpoint on the exporter server, usually `/metrics`
	// (the default).
	MetricPath string `yaml:"metricPath" default:"/metrics"`
}

// PrometheusConfig method implementation.
func (c *Config) NewPrometheusClient() (*prometheusexporter.PrometheusClient, error) {
	k8sClient, err := kubernetes.MakeClient(c.KubernetesAPI)
	if err != nil {
		return nil, err
	}
	return  &prometheusexporter.PrometheusClient{
		GetMetricFamilies: func() (metricFamilies []*dto.MetricFamily, err error) {
			var body io.ReadCloser
			body, err = k8sClient.CoreV1().RESTClient().Get().RequestURI(c.MetricPath).Stream()
			if err != nil {
				return
			}
			defer body.Close()
			decoder := expfmt.NewDecoder(body, expfmt.FmtText)
			metricFamilies = make([]*dto.MetricFamily, 0)
			for {
				var mf dto.MetricFamily
				err = decoder.Decode(&mf)
				if err == io.EOF {
					return metricFamilies, nil
				} else if err != nil {
					return nil, err
				}
				metricFamilies = append(metricFamilies, &mf)
			}
			return
		},
	}, nil
}

// PrometheusConfig method implementation.
func (c *Config) GetInterval() time.Duration {
	return time.Duration(c.IntervalSeconds)*time.Second
}
