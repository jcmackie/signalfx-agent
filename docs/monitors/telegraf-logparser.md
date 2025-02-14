<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# telegraf/logparser

Monitor Type: `telegraf/logparser` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/telegraf/monitors/telegraflogparser))

**Accepts Endpoints**: No

**Multiple Instances Allowed**: Yes

## Overview

This monitor is based on the Telegraf logparser plugin.
The monitor tails log files. More information about the Telegraf plugin
can be found [here](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/logparser).
All metrics emitted from this monitor will have the `plugin` dimension set to `telegraf-logparser`

Sample YAML configuration:

```yaml
 - type: telegraf/logparser
   files:
    - '$file'
   watchMethod: poll       # specify the file watch method ("inotify" or "poll")
   fromBeginning: true     # specify to read from the beginning
   measurementName: test-measurement # the metric name prefix
   patterns:
    - "%{COMMON_LOG_FORMAT}" # specifies the apache common log format
   timezone: UTC
```


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: telegraf/logparser
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `files` | **yes** | `list of strings` | Paths to files to be tailed |
| `watchMethod` | no | `string` | Method for watching changes to files ("ionotify" or "poll") (**default:** `poll`) |
| `fromBeginning` | no | `bool` | Whether to start tailing from the beginning of the file (**default:** `false`) |
| `measurementName` | no | `string` | Name of the measurement |
| `patterns` | no | `list of strings` | A list of patterns to match. |
| `namedPatterns` | no | `list of strings` | A list of named grok patterns to match. |
| `customPatterns` | no | `string` | Custom grok patterns. (`grok` only) |
| `customPatternFiles` | no | `list of strings` | List of paths to custom grok pattern files. |
| `timezone` | no | `string` | Specifies the timezone.  The default is UTC time.  Other options are `Local` for the local time on the machine, `UTC`, and `Canada/Eastern` (unix style timezones). |



The agent does not do any built-in filtering of metrics coming out of this
monitor.


