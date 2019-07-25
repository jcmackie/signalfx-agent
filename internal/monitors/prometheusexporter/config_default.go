package prometheusexporter

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/sirupsen/logrus"
)

// Config is the default config for this monitor and implements ConfigInterface.
type Config struct {
	config.MonitorConfig `yaml:",inline" acceptsEndpoints:"true"`

	// Host of the exporter
	Host string `yaml:"host" validate:"required"`
	// Port of the exporter
	Port uint16 `yaml:"port" validate:"required"`

	// Basic Auth username to use on each request, if any.
	Username string `yaml:"username"`
	// Basic Auth password to use on each request, if any.
	Password string `yaml:"password" neverLog:"true"`

	// If true, the agent will connect to the exporter using HTTPS instead of
	// plain HTTP.
	UseHTTPS bool `yaml:"useHTTPS"`
	// If useHTTPS is true and this option is also true, the exporter's TLS
	// cert will not be verified.
	SkipVerify bool `yaml:"skipVerify"`

	// Path to the metrics endpoint on the exporter server, usually `/metrics`
	// (the default).
	MetricPath string `yaml:"metricPath" default:"/metrics"`

	// Send all the metrics that come out of the Prometheus exporter without
	// any filtering.  This option has no effect when using the prometheus
	// exporter monitor directly since there is no built-in filtering, only
	// when embedding it in other monitors.
	SendAllMetrics bool `yaml:"sendAllMetrics"`
}

// NewClient is a ConfigInterface method implementation that creates the default prometheus client.
func (c *Config) NewClient() (*Client, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipVerify},
		},
	}
	var scheme string
	if c.UseHTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host := c.Host
	// Handle IPv6 addresses properly
	if strings.ContainsAny(host, ":") {
		host = "[" + host + "]"
	}
	url := fmt.Sprintf("%s://%s:%d%s", scheme, host, c.Port, c.MetricPath)
	return &Client{
		GetMetricFamilies: func() (metricFamilies []*dto.MetricFamily, err error) {
			var req *http.Request
			var resp *http.Response
			// Prometheus 2.0 deprecated protobuf and now only does the text format.
			if req, err = http.NewRequest("GET", url, nil); err != nil {
				return
			}
			if c.Username != "" {
				req.SetBasicAuth(c.Username, c.Password)
			}
			if resp, err = httpClient.Do(req); err != nil {
				return
			}
			if resp.StatusCode != 200 {
				err = fmt.Errorf("prometheus exporter at %s returned status %d", resp.Request.URL.String(), resp.StatusCode)
				return
			}
			defer resp.Body.Close()
			decoder := expfmt.NewDecoder(resp.Body, expfmt.ResponseFormat(resp.Header))
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

func (c *Config) GetExtraMetrics() []string {
	// Maintain backwards compatibility with the config flag that existing
	// prior to the new filtering mechanism.
	if c.SendAllMetrics {
		return []string{"*"}
	}
	return nil
}

var _ config.ExtraMetrics = &Config{}
