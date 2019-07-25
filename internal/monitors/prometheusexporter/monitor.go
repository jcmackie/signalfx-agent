package prometheusexporter

import (
	"context"
	"io"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/prometheus/common/expfmt"
	"github.com/signalfx/golib/datapoint"
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
	SendAll      bool
	cancel       func()
	client       *Client
	loggingEntry *logrus.Entry
}

// Configure the monitor and kick off volume metric syncing
func (m *Monitor) Configure(conf ConfigInterface) (err error) {
	if m.loggingEntry == nil {
		m.loggingEntry = logrus.WithFields(logrus.Fields{"monitorType": conf.GetMonitorType()})
	}
	if m.client, err = conf.NewClient(); err != nil {
		m.loggingEntry.WithError(err).Error("Could not create prometheus client")
		return
	}
	var ctx context.Context
	ctx, m.cancel = context.WithCancel(context.Background())
	utils.RunOnInterval(ctx, func() {
		var bodyReader io.ReadCloser
		var format expfmt.Format
		defer func() {
			if bodyReader != nil {
				bodyReader.Close()
			}
		}()
		if bodyReader, format, err = m.client.GetBodyReader(); err != nil {
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
	return err
}

// Shutdown stops the metric sync
func (m *Monitor) Shutdown() {
	if m.cancel != nil {
		m.cancel()
	}
}
