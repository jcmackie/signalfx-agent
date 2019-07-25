package prometheusexporter

import (
	"context"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	"github.com/signalfx/signalfx-agent/internal/utils"
)

func init() {
	monitors.Register(&monitorMetadata, func() interface{} { return &Monitor{} }, &Config{})
}

// Monitor for prometheus exporter metrics
type Monitor struct {
	Output types.Output
	// Optional set of metric names that will be sent by default, all other
	// metrics derived from the exporter being dropped.
	IncludedMetrics map[string]bool
	// Extra dimensions to add in addition to those specified in the config.
	ExtraDimensions map[string]string
	// If true, IncludedMetrics is ignored and everything is sent.
	SendAll bool
	cancel  func()
	client  *Client
}

// Configure the monitor and kick off volume metric syncing
func (m *Monitor) Configure(conf ConfigInterface) (err error) {
	logger := conf.GetLoggingEntry()
	if m.client, err = conf.NewClient(); err != nil {
		logger.WithError(err).Error("Could not create prometheus client")
		return
	}
	var ctx context.Context
	ctx, m.cancel = context.WithCancel(context.Background())
	utils.RunOnInterval(ctx, func() {
		var metricFamilies []*dto.MetricFamily
		if metricFamilies, err = m.client.GetMetricFamilies(); err != nil {
			logger.WithError(err).Error("Could not get prometheus metrics")
			return
		}
		dps := convertMetricFamilies(metricFamilies)
		now := time.Now()
		for i := range dps {
			dps[i].Timestamp = now
			m.Output.SendDatapoint(dps[i])
		}
	}, conf.GetInterval())
	return
}

// Shutdown stops the metric sync
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
