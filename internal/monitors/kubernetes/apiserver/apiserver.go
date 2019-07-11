package apiserver

import (
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/prometheusexporter"
)

func init() {
	monitors.Register(&monitorMetadata, func() interface{} { return &prometheusexporter.Monitor{} }, &prometheusexporter.Config{})
}

