// Configuration structures and helper logic for all agent components.
package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/signalfx/neo-agent/core/config/stores"
	"github.com/signalfx/neo-agent/core/filters"
	"github.com/signalfx/neo-agent/utils"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	// TODO: Consider whether to allow store configuration from the main config
	// file.  There is a major chicken/egg problem with this and reloading
	// stores is very tricky.  Right now, stores can only be configured via
	// envvars and I think it is best to keep it that way.
	//Stores              map[string]StoreConfig `yaml:"stores,omitempty" default:"{}"`
	SignalFxAccessToken string `yaml:"signalFxAccessToken"`
	// The ingest URL for SignalFx, without the path
	IngestURL string `yaml:"ingestUrl" default:"https://ingest.signalfx.com"`
	// The hostname that will be reported as the "host" dimension on metrics
	// for which host applies
	Hostname string `yaml:"hostname"`
	// How often to send metrics to SignalFx.  Monitors can't override this
	// individually.
	IntervalSeconds int `yaml:"intervalSeconds" default:"10"`
	// Dimensions that will be automatically added to all metrics reported
	GlobalDimensions map[string]string `yaml:"globalDimensions" default:"{}"`
	Observers        []ObserverConfig  `yaml:"observers" default:"[]"`
	Monitors         []MonitorConfig   `yaml:"monitors" default:"[]"`
	Writer           WriterConfig      `yaml:"writer" default:"{}"`
	Logging          LogConfig         `yaml:"logging" default:"{}"`
	// Configure the underlying collectd daemon
	Collectd         CollectdConfig `yaml:"collectd" default:"{}"`
	MetricsToExclude []MetricFilter `yaml:"metricsToExclude" default:"[]"`
	ProcFSPath       string         `yaml:"procFSPath" default:"/proc"`
	PythonEnabled    bool           `yaml:"pythonEnabled" default:"false"`
}

func (c *Config) setDefaultHostname() {
	fqdn := fqdn.Get()
	if fqdn == "unknown" {
		log.Info("Error getting fully qualified hostname")
	} else {
		log.Infof("Using hostname %s", fqdn)
		c.Hostname = fqdn
	}
}

func (c *Config) Initialize(metaStore *stores.MetaStore) (*Config, error) {
	c.overrideFromEnv()

	c.setDefaultHostname()

	if !c.validate() {
		return nil, fmt.Errorf("Configuration did not validate!")
	}

	c.propagateValuesDown(metaStore)
	idGenerator := newIdGenerator()
	for i := range c.Monitors {
		c.Monitors[i].EnsureID(idGenerator)
	}

	return c, nil
}

func (c *Config) IngestURLAsURL() *url.URL {
	if url, err := url.Parse(c.IngestURL); err == nil {
		return url
	}
	return nil
}

// Support overridding a few config options with envvars.  No need to allow
// everything to be overridden.
func (c *Config) overrideFromEnv() {
	c.SignalFxAccessToken = utils.FirstNonEmpty(c.SignalFxAccessToken, os.Getenv("SFX_ACCESS_TOKEN"))
	c.Hostname = utils.FirstNonEmpty(c.Hostname, os.Getenv("SFX_HOSTNAME"))
	c.IngestURL = utils.FirstNonEmpty(c.IngestURL, os.Getenv("SFX_INGEST_URL"))

	intervalSeconds, err := strconv.ParseInt(os.Getenv("SFX_INTERVAL_SECONDS"), 10, 32)
	if err != nil {
		c.IntervalSeconds = utils.FirstNonZero(c.IntervalSeconds, int(intervalSeconds))
	}
}

// Validate everything except for Observers and Monitors
func (c *Config) validate() bool {
	valid := true

	if c.SignalFxAccessToken == "" {
		log.Error("signalFxAccessToken must be set!")
		valid = false
	}
	if _, err := url.Parse(c.IngestURL); err != nil {
		log.WithFields(log.Fields{
			"ingestURL": c.IngestURL,
			"error":     err,
		}).Error("ingestURL is not a valid URL")
	}

	return valid
}

func (c *Config) makeFilterSet() *filters.FilterSet {
	fs := make([]*filters.Filter, 0)
	for _, mte := range c.MetricsToExclude {
		dims := mte.ConvertDimensionsMapForSliceValues()
		mte.ConvertMetricNameToSlice()
		fs = append(fs, filters.New(mte.MonitorType, mte.MetricNames, dims))
	}

	return &filters.FilterSet{
		Filters: fs,
	}
}

// Send values from the top of the config down to nested configs that might
// need them
func (c *Config) propagateValuesDown(metaStore *stores.MetaStore) {
	filterSet := c.makeFilterSet()

	ingestURL, err := url.Parse(c.IngestURL)
	if err != nil {
		panic("IngestURL was supposed to be validated already")
	}

	for i := range c.Monitors {
		c.Monitors[i].IngestURL = ingestURL
		c.Monitors[i].GlobalDimensions = c.GlobalDimensions
		c.Monitors[i].SignalFxAccessToken = c.SignalFxAccessToken
		c.Monitors[i].Hostname = c.Hostname
		c.Monitors[i].Filter = filterSet
		c.Monitors[i].ProcFSPath = c.ProcFSPath
		// Top level interval serves as a default
		c.Monitors[i].IntervalSeconds = utils.FirstNonZero(c.Monitors[i].IntervalSeconds, c.IntervalSeconds)
		c.Monitors[i].MetaStore = metaStore
	}

	c.Collectd.Hostname = c.Hostname
	c.Collectd.IntervalSeconds = c.IntervalSeconds
	c.Collectd.Filter = filterSet

	c.Writer.IngestURL = ingestURL
	c.Writer.Filter = filterSet
	c.Writer.SignalFxAccessToken = c.SignalFxAccessToken
}

type LogConfig struct {
	Level string `yaml:"level,omitempty" default:"info"`
	// TODO: Support log file output and other log targets
}

func (lc *LogConfig) LogrusLevel() *log.Level {
	if lc.Level != "" {
		level, err := log.ParseLevel(lc.Level)
		if err != nil {
			log.WithFields(log.Fields{
				"level": lc.Level,
			}).Error("Invalid log level")
			return nil
		}
		return &level
	}
	return nil
}

type CustomConfigurable interface {
	GetOtherConfig() map[string]interface{}
}

// Collectd high-level configurations
type CollectdConfig struct {
	DisableCollectd      bool   `yaml:"disableCollectd,omitempty" default:"false"`
	IntervalSeconds      int    `yaml:"intervalSeconds,omitempty" default:"10"`
	Timeout              int    `yaml:"timeout,omitempty" default:"40"`
	ReadThreads          int    `yaml:"readThreads,omitempty" default:"5"`
	WriteQueueLimitHigh  int    `yaml:"writeQueueLimitHigh,omitempty" default:"500000"`
	WriteQueueLimitLow   int    `yaml:"writeQueueLimitLow,omitempty" default:"400000"`
	CollectInternalStats bool   `yaml:"collectInternalStats,omitempty" default:"false"`
	LogLevel             string `yaml:"logLevel,omitempty" default:"notice"`
	// The following are propagated from the top-level config
	Hostname string             `yaml:"-"`
	Filter   *filters.FilterSet `yaml:"-"`
}

type StoreConfig struct {
	OtherConfig map[string]interface{} `yaml:",inline,omitempty" default:"{}"`
}

func (sc *StoreConfig) GetOtherConfig() map[string]interface{} {
	return sc.OtherConfig
}

var (
	// EnvReplacer replaces . and - with _
	EnvReplacer   = strings.NewReplacer(".", "_", "-", "_")
	configTimeout = 10 * time.Second
)

// Used to ensure unique IDs for monitors and observers
func newIdGenerator() func(string) int {
	ids := map[string]int{}

	return func(name string) int {
		ids[name] += 1
		return ids[name]
	}
}
