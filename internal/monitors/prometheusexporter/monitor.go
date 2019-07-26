package prometheusexporter

import (
	"context"
	"github.com/prometheus/common/expfmt"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/signalfx-agent/internal/utils"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
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
	SendAll      bool
	ctx          context.Context
	cancel       func()
	client       *Client
	loggingEntry *logrus.Entry
	configErr    error
	mux          sync.Mutex
}

// Configure the monitor and kick off volume metric syncing
func (m *Monitor) Configure(conf ConfigInterface) error {
	if m.configureIfNilSync(conf); m.configErr == nil {
		m.fetchAndSendAsync(conf)
	}
	return m.configErr
}

func (m *Monitor) configureIfNilSync(conf ConfigInterface) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.cancel == nil {
		m.ctx, m.cancel = context.WithCancel(context.Background())
	}
	if m.loggingEntry == nil {
		m.loggingEntry = logrus.WithFields(logrus.Fields{"monitorType": conf.GetMonitorType()})
	}
	if m.client == nil {
		if m.client, m.configErr = conf.NewClient(); m.configErr != nil {
			m.loggingEntry.WithError(m.configErr).Error("Could not create prometheus client")
		}
	}
}

func (m *Monitor) fetchAndSendAsync(conf ConfigInterface) {
	utils.RunOnInterval(m.ctx, func() {
		bodyReader, format, err := m.client.GetBodyReader()
		defer func() {
			if bodyReader != nil {
				bodyReader.Close()
			}
		}()
		if err != nil {
			m.loggingEntry.WithError(err).Error("Could not get prometheus metrics")
			return
		}
		decoder := expfmt.NewDecoder(bodyReader, format)
		var dps []*datapoint.Datapoint
		if dps, err = decodeMetrics(decoder); err != nil {
			m.loggingEntry.WithError(err).Error("Could not decode prometheus metrics from response body")
			return
		}
		now := time.Now()
		for i := range dps {
			dps[i].Timestamp = now
			m.Output.SendDatapoint(dps[i])
		}
	}, conf.GetInterval())
}

// Shutdown stops the metric sync
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
