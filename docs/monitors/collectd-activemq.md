<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/activemq

Monitor Type: `collectd/activemq` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/activemq))

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Overview

SignalFx's integration with ActiveMQ wraps the [Collectd GenericJMX
monitor](https://docs.signalfx.com/en/latest/integrations/agent/monitors/collectd-genericjmx.html)
to monitor ActiveMQ.

Use this monitor to gather the following types of information from ActiveMQ:

* Broker (Totals per broker)
* Queue (Queue status)
* Topic (Topic status)

To monitor the age of messages inside ActiveMQ queues, see [ActiveMQ
message age listener](https://github.com/signalfx/integrations/tree/master/amq-message-age)[](sfx_link:amq-message-age).

<!--- SETUP --->
## Example config

```yaml
monitors:
 - type: collectd/activemq
   host: localhost
   # This is the remote JMX port for ActiveMQ
   port: 1099
```


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: collectd/activemq
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `host` | **yes** | `string` | Host to connect to -- JMX must be configured for remote access and accessible from the agent |
| `port` | **yes** | `integer` | JMX connection port (NOT the RMI port) on the application.  This correponds to the `com.sun.management.jmxremote.port` Java property that should be set on the JVM when running the application. |
| `name` | no | `string` |  |
| `serviceName` | no | `string` | This is how the service type is identified in the SignalFx UI so that you can get built-in content for it.  For custom JMX integrations, it can be set to whatever you like and metrics will get the special property `sf_hostHasService` set to this value. |
| `serviceURL` | no | `string` | The JMX connection string.  This is rendered as a Go template and has access to the other values in this config. NOTE: under normal circumstances it is not advised to set this string directly - setting the host and port as specified above is preferred. (**default:** `service:jmx:rmi:///jndi/rmi://{{.Host}}:{{.Port}}/jmxrmi`) |
| `instancePrefix` | no | `string` | Prefixes the generated plugin instance with prefix. If a second `instancePrefix` is specified in a referenced MBean block, the prefix specified in the Connection block will appear at the beginning of the plugin instance, and the prefix specified in the MBean block will be appended to it |
| `username` | no | `string` | Username to authenticate to the server |
| `password` | no | `string` | User password to authenticate to the server |
| `customDimensions` | no | `map of strings` | Takes in key-values pairs of custom dimensions at the connection level. |
| `mBeansToCollect` | no | `list of strings` | A list of the MBeans defined in `mBeanDefinitions` to actually collect. If not provided, then all defined MBeans will be collected. |
| `mBeansToOmit` | no | `list of strings` | A list of the MBeans to omit. This will come handy in cases where only a few MBeans need to omitted from the default list |
| `mBeanDefinitions` | no | `map of objects (see below)` | Specifies how to map JMX MBean values to metrics.  If using a specific service monitor such as cassandra, kafka, or activemq, they come pre-loaded with a set of mappings, and any that you add in this option will be merged with those.  See [collectd GenericJMX](https://collectd.org/documentation/manpages/collectd-java.5.shtml#genericjmx_plugin) for more details. |


The **nested** `mBeanDefinitions` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `objectName` | no | `string` | Sets the pattern which is used to retrieve MBeans from the MBeanServer. If more than one MBean is returned you should use the `instanceFrom` option to make the identifiers unique |
| `instancePrefix` | no | `string` | Prefixes the generated plugin instance with prefix |
| `instanceFrom` | no | `list of strings` | The object names used by JMX to identify MBeans include so called "properties" which are basically key-value-pairs. If the given object name is not unique and multiple MBeans are returned, the values of those properties usually differ. You can use this option to build the plugin instance from the appropriate property values. This option is optional and may be repeated to generate the plugin instance from multiple property values |
| `values` | no | `list of objects (see below)` | The `value` blocks map one or more attributes of an MBean to a value list in collectd. There must be at least one `value` block within each MBean block |
| `dimensions` | no | `list of strings` |  |


The **nested** `values` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `type` | no | `string` | Sets the data set used within collectd to handle the values of the MBean attribute |
| `table` | no | `bool` | Set this to true if the returned attribute is a composite type. If set to true, the keys within the composite type is appended to the type instance. (**default:** `false`) |
| `instancePrefix` | no | `string` | Works like the option of the same name directly beneath the MBean block, but sets the type instance instead |
| `instanceFrom` | no | `list of strings` | Works like the option of the same name directly beneath the MBean block, but sets the type instance instead |
| `attribute` | no | `string` | Sets the name of the attribute from which to read the value. You can access the keys of composite types by using a dot to concatenate the key name to the attribute name. For example: “attrib0.key42”. If `table` is set to true, path must point to a composite type, otherwise it must point to a numeric type. |
| `attributes` | no | `list of strings` | The plural form of the `attribute` config above.  Used to derive multiple metrics from a single MBean. |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - ***`counter.amq.TotalConnectionsCount`*** (*counter*)<br>    Total connections count per broker
 - ***`gauge.amq.TotalConsumerCount`*** (*gauge*)<br>    Total number of consumers subscribed to destinations on the broker
 - `gauge.amq.TotalDequeueCount` (*gauge*)<br>    Total number of messages that have been acknowledged from the broker.
 - ***`gauge.amq.TotalEnqueueCount`*** (*gauge*)<br>    Total number of messages that have been sent to the broker.
 - ***`gauge.amq.TotalMessageCount`*** (*gauge*)<br>    Total number of unacknowledged messages on the broker
 - ***`gauge.amq.TotalProducerCount`*** (*gauge*)<br>    Total number of message producers active on destinations on the broker
 - `gauge.amq.queue.AverageBlockedTime` (*gauge*)<br>    Average time (ms) that messages have spent blocked by Flow Control.
 - ***`gauge.amq.queue.AverageEnqueueTime`*** (*gauge*)<br>    Average time (ms) that messages have been held at this destination
 - `gauge.amq.queue.AverageMessageSize` (*gauge*)<br>    Average size of messages in this queue, in bytes.
 - `gauge.amq.queue.BlockedSends` (*gauge*)<br>    Number of messages blocked by Flow Control.
 - ***`gauge.amq.queue.ConsumerCount`*** (*gauge*)<br>    Number of consumers subscribed to this queue.
 - ***`gauge.amq.queue.DequeueCount`*** (*gauge*)<br>    Number of messages that have been acknowledged and removed from the queue.
 - ***`gauge.amq.queue.EnqueueCount`*** (*gauge*)<br>    Number of messages that have been sent to the queue.
 - ***`gauge.amq.queue.ExpiredCount`*** (*gauge*)<br>    Number of messages that have expired from the queue.
 - `gauge.amq.queue.ForwardCount` (*gauge*)<br>    Number of messages that have been forwarded from this queue to a networked broker.
 - ***`gauge.amq.queue.InFlightCount`*** (*gauge*)<br>    The number of messages that have been dispatched to consumers, but not acknowledged.
 - ***`gauge.amq.queue.ProducerCount`*** (*gauge*)<br>    Number of producers publishing to this queue
 - ***`gauge.amq.queue.QueueSize`*** (*gauge*)<br>    The number of messages in the queue that have yet to be consumed.
 - `gauge.amq.queue.TotalBlockedTime` (*gauge*)<br>    The total time (ms) that messages have spent blocked by Flow Control.
 - `gauge.amq.topic.AverageBlockedTime` (*gauge*)<br>    Average time (ms) that messages have been blocked by Flow Control.
 - ***`gauge.amq.topic.AverageEnqueueTime`*** (*gauge*)<br>    Average time (ms) that messages have been held at this destination.
 - `gauge.amq.topic.AverageMessageSize` (*gauge*)<br>    Average size of messages on this topic, in bytes.
 - `gauge.amq.topic.BlockedSends` (*gauge*)<br>    Number of messages blocked by Flow Control
 - ***`gauge.amq.topic.ConsumerCount`*** (*gauge*)<br>    The number of consumers subscribed to this topic
 - `gauge.amq.topic.DequeueCount` (*gauge*)<br>    Number of messages that have been acknowledged and removed from the topic.
 - ***`gauge.amq.topic.EnqueueCount`*** (*gauge*)<br>    The number of messages that have been sent to the topic.
 - ***`gauge.amq.topic.ExpiredCount`*** (*gauge*)<br>    The number of messages that have expired from this topic.
 - `gauge.amq.topic.ForwardCount` (*gauge*)<br>    The number of messages that have been forwarded from this topic to a networked broker.
 - ***`gauge.amq.topic.InFlightCount`*** (*gauge*)<br>    The number of messages that have been dispatched to consumers, but have not yet been acknowledged.
 - ***`gauge.amq.topic.ProducerCount`*** (*gauge*)<br>    Number of producers publishing to this topic.
 - ***`gauge.amq.topic.QueueSize`*** (*gauge*)<br>    Number of messages in the topic that have yet to be consumed.
 - `gauge.amq.topic.TotalBlockedTime` (*gauge*)<br>    The total time (ms) that messages have spent blocked by Flow Control

#### Group jvm
All of the following metrics are part of the `jvm` metric group. All of
the non-default metrics below can be turned on by adding `jvm` to the
monitor config option `extraGroups`:
 - ***`gauge.jvm.threads.count`*** (*gauge*)<br>    Number of JVM threads
 - ***`gauge.loaded_classes`*** (*gauge*)<br>    Number of classes loaded in the JVM
 - ***`invocations`*** (*cumulative*)<br>    Total number of garbage collection events
 - ***`jmx_memory.committed`*** (*gauge*)<br>    Amount of memory guaranteed to be available in bytes
 - ***`jmx_memory.init`*** (*gauge*)<br>    Amount of initial memory at startup in bytes
 - ***`jmx_memory.max`*** (*gauge*)<br>    Maximum amount of memory that can be used in bytes
 - ***`jmx_memory.used`*** (*gauge*)<br>    Current memory usage in bytes
 - ***`total_time_in_ms.collection_time`*** (*cumulative*)<br>    Amount of time spent garbage collecting in milliseconds

### Non-default metrics (version 4.7.0+)

**The following information applies to the agent version 4.7.0+ that has
`enableBuiltInFiltering: true` set on the top level of the agent config.**

To emit metrics that are not _default_, you can add those metrics in the
generic monitor-level `extraMetrics` config option.  Metrics that are derived
from specific configuration options that do not appear in the above list of
metrics do not need to be added to `extraMetrics`.

To see a list of metrics that will be emitted you can run `agent-status
monitors` after configuring this monitor in a running agent instance.

### Legacy non-default metrics (version < 4.7.0)

**The following information only applies to agent version older than 4.7.0. If
you have a newer agent and have set `enableBuiltInFiltering: true` at the top
level of your agent config, see the section above. See upgrade instructions in
[Old-style whitelist filtering](../legacy-filtering.md#old-style-whitelist-filtering).**

If you have a reference to the `whitelist.json` in your agent's top-level
`metricsToExclude` config option, and you want to emit metrics that are not in
that whitelist, then you need to add an item to the top-level
`metricsToInclude` config option to override that whitelist (see [Inclusion
filtering](../legacy-filtering.md#inclusion-filtering).  Or you can just
copy the whitelist.json, modify it, and reference that in `metricsToExclude`.



