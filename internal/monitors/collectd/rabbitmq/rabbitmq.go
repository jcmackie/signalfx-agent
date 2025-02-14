package rabbitmq

import (
	"github.com/signalfx/golib/pointer"
	"github.com/signalfx/signalfx-agent/internal/core/config"

	"github.com/signalfx/signalfx-agent/internal/utils"

	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/collectd"
	"github.com/signalfx/signalfx-agent/internal/monitors/collectd/python"
	"github.com/signalfx/signalfx-agent/internal/monitors/pyrunner"
)

func init() {
	monitors.Register(&monitorMetadata, func() interface{} {
		return &Monitor{
			python.PyMonitor{
				MonitorCore: pyrunner.New("sfxcollectd"),
			},
		}
	}, &Config{})
}

// Config is the monitor-specific config with the generic config embedded
type Config struct {
	config.MonitorConfig `yaml:",inline" acceptsEndpoints:"true"`
	python.CommonConfig  `yaml:",inline"`
	pyConf               *python.Config
	Host                 string `yaml:"host" validate:"required"`
	Port                 uint16 `yaml:"port" validate:"required"`
	// The name of the particular RabbitMQ instance.  Can be a Go template
	// using other config options. This will be used as the `plugin_instance`
	// dimension.
	BrokerName         string `yaml:"brokerName" default:"{{.host}}-{{.port}}"`
	CollectChannels    *bool  `yaml:"collectChannels"`
	CollectConnections *bool  `yaml:"collectConnections"`
	CollectExchanges   *bool  `yaml:"collectExchanges"`
	CollectNodes       *bool  `yaml:"collectNodes"`
	CollectQueues      *bool  `yaml:"collectQueues"`
	HTTPTimeout        *int   `yaml:"httpTimeout"`
	VerbosityLevel     string `yaml:"verbosityLevel"`
	Username           string `yaml:"username" validate:"required"`
	Password           string `yaml:"password" validate:"required" neverLog:"true"`
}

// PythonConfig returns the embedded python.Config struct from the interface
func (c *Config) PythonConfig() *python.Config {
	c.pyConf.CommonConfig = c.CommonConfig
	return c.pyConf
}

// Monitor is the main type that represents the monitor
type Monitor struct {
	python.PyMonitor
}

// Configure configures and runs the plugin in python
func (m *Monitor) Configure(conf *Config) error {
	sendChannelMetrics := conf.CollectChannels
	sendConnectionMetrics := conf.CollectConnections
	sendExchangeMetrics := conf.CollectExchanges
	sendNodeMetrics := conf.CollectNodes
	sendQueueMetrics := conf.CollectQueues

	if m.Output.HasEnabledMetricInGroup(groupChannel) {
		sendChannelMetrics = pointer.Bool(true)
	}
	if m.Output.HasEnabledMetricInGroup(groupConnection) {
		sendConnectionMetrics = pointer.Bool(true)
	}
	if m.Output.HasEnabledMetricInGroup(groupExchange) {
		sendExchangeMetrics = pointer.Bool(true)
	}
	if m.Output.HasEnabledMetricInGroup(groupNode) {
		sendNodeMetrics = pointer.Bool(true)
	}
	if m.Output.HasEnabledMetricInGroup(groupQueue) {
		sendQueueMetrics = pointer.Bool(true)
	}

	conf.pyConf = &python.Config{
		MonitorConfig: conf.MonitorConfig,
		Host:          conf.Host,
		Port:          conf.Port,
		ModuleName:    "rabbitmq",
		ModulePaths:   []string{collectd.MakePythonPluginPath("rabbitmq")},
		TypesDBPaths:  []string{collectd.DefaultTypesDBPath()},
		PluginConfig: map[string]interface{}{
			"Host":               conf.Host,
			"Port":               conf.Port,
			"BrokerName":         conf.BrokerName,
			"Username":           conf.Username,
			"Password":           conf.Password,
			"CollectChannels":    sendChannelMetrics,
			"CollectConnections": sendConnectionMetrics,
			"CollectExchanges":   sendExchangeMetrics,
			"CollectNodes":       sendNodeMetrics,
			"CollectQueues":      sendQueueMetrics,
			"HTTPTimeout":        conf.HTTPTimeout,
			"VerbosityLevel":     conf.VerbosityLevel,
		},
	}

	// the python runner's templating system does not convert to map first
	// this requires TitleCase template values.  For BrokerName we accept
	// either upper or lower case values.  Converting the map to yaml
	// and explicitly rendering the BrokerName will change everything to lower case.
	mp, err := utils.ConvertToMapViaYAML(conf)
	if err != nil {
		return err
	}
	brokerName, err := collectd.RenderValue(conf.BrokerName, mp)
	if err != nil {
		return err
	}
	conf.pyConf.PluginConfig["BrokerName"] = brokerName

	return m.PyMonitor.Configure(conf)
}

// GetExtraMetrics returns additional metrics that should be allowed through.
func (c *Config) GetExtraMetrics() []string {
	var extraMetrics []string

	if c.CollectChannels != nil && *c.CollectChannels {
		extraMetrics = append(extraMetrics, groupMetricsMap[groupChannel]...)
	}

	if c.CollectConnections != nil && *c.CollectConnections {
		extraMetrics = append(extraMetrics, groupMetricsMap[groupConnection]...)
	}

	if c.CollectExchanges != nil && *c.CollectExchanges {
		extraMetrics = append(extraMetrics, groupMetricsMap[groupExchange]...)
	}

	if c.CollectNodes != nil && *c.CollectNodes {
		extraMetrics = append(extraMetrics, groupMetricsMap[groupNode]...)
	}

	if c.CollectQueues != nil && *c.CollectQueues {
		extraMetrics = append(extraMetrics, groupMetricsMap[groupQueue]...)
	}

	return extraMetrics
}
