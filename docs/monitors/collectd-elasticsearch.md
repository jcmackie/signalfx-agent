<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/elasticsearch

Monitor Type: `collectd/elasticsearch` ([Source](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/elasticsearch))

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Overview

Monitors ElasticSearch instances. We strongly recommend using the
[elasticsearch](./elasticsearch.md) monitor instead, as it will
scale much better.

See https://github.com/signalfx/collectd-elasticsearch and
https://github.com/signalfx/integrations/tree/master/elasticsearch


## Configuration

To activate this monitor in the Smart Agent, add the following to your
agent config:

```
monitors:  # All monitor config goes under this key
 - type: collectd/elasticsearch
   ...  # Additional config
```

**For a list of monitor options that are common to all monitors, see [Common
Configuration](../monitor-config.md#common-configuration).**


| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `pythonBinary` | no | `string` | Path to a python binary that should be used to execute the Python code. If not set, a built-in runtime will be used.  Can include arguments to the binary as well. |
| `host` | **yes** | `string` |  |
| `port` | **yes** | `integer` |  |
| `additionalMetrics` | no | `list of strings` | AdditionalMetrics to report on |
| `cluster` | no | `string` | Cluster name to which the node belongs. This is an optional config that will override the cluster name fetched from a node and will be used to populate the plugin_instance dimension |
| `detailedMetrics` | no | `bool` | DetailedMetrics turns on additional metric time series (**default:** `true`) |
| `enableClusterHealth` | no | `bool` | EnableClusterHealth enables reporting on the cluster health (**default:** `true`) |
| `enableIndexStats` | no | `bool` | EnableIndexStats reports metrics about indexes (**default:** `true`) |
| `indexes` | no | `list of strings` | Indexes to report on (**default:** `[_all]`) |
| `indexInterval` | no | `unsigned integer` | IndexInterval is an interval in seconds at which the plugin will report index stats. It must be greater than or equal, and divisible by the Interval configuration (**default:** `300`) |
| `indexStatsMasterOnly` | no | `bool` | IndexStatsMasterOnly sends index stats from the master only (**default:** `false`) |
| `indexSummaryOnly` | no | `bool` |  (**default:** `false`) |
| `password` | no | `string` | Password used to access elasticsearch stats api |
| `protocol` | no | `string` | Protocol used to connect: http or https |
| `threadPools` | no | `list of strings` | ThreadPools to report on (**default:** `[search index]`) |
| `username` | no | `string` | Username used to access elasticsearch stats api |
| `version` | no | `string` |  |


## Metrics

These are the metrics available for this monitor.
Metrics that are categorized as
[container/host](https://docs.signalfx.com/en/latest/admin-guide/usage.html#about-custom-bundled-and-high-resolution-metrics)
(*default*) are ***in bold and italics*** in the list below.


 - `bytes.indices.primaries.merges.total-size` (*gauge*)<br>
 - `bytes.indices.total.merges.total-size` (*gauge*)<br>
 - `counter.http.total_open` (*cumulative*)<br>
 - `counter.indices.cache.field.eviction` (*cumulative*)<br>
 - `counter.indices.cache.filter.cache-count` (*cumulative*)<br>
 - `counter.indices.cache.filter.evictions` (*cumulative*)<br>
 - `counter.indices.cache.filter.hit-count` (*cumulative*)<br>
 - `counter.indices.cache.filter.miss-count` (*cumulative*)<br>
 - `counter.indices.cache.filter.total-count` (*cumulative*)<br>
 - `counter.indices.flush.time` (*cumulative*)<br>
 - `counter.indices.flush.total` (*cumulative*)<br>
 - `counter.indices.get.exists-time` (*cumulative*)<br>
 - `counter.indices.get.exists-total` (*cumulative*)<br>
 - `counter.indices.get.missing-time` (*cumulative*)<br>
 - `counter.indices.get.missing-total` (*cumulative*)<br>
 - `counter.indices.get.time` (*cumulative*)<br>
 - ***`counter.indices.get.total`*** (*cumulative*)<br>    The total number of get requests since node startup
 - `counter.indices.indexing.delete-time` (*cumulative*)<br>
 - `counter.indices.indexing.delete-total` (*cumulative*)<br>
 - `counter.indices.indexing.index-time` (*cumulative*)<br>
 - ***`counter.indices.indexing.index-total`*** (*cumulative*)<br>    The total number of index requests since node startup
 - `counter.indices.merges.time` (*cumulative*)<br>
 - ***`counter.indices.merges.total`*** (*cumulative*)<br>    Total number of merges since node startup
 - `counter.indices.merges.total-size` (*cumulative*)<br>
 - `counter.indices.primaries.fielddata.evictions` (*cumulative*)<br>
 - `counter.indices.primaries.flush.total` (*cumulative*)<br>
 - `counter.indices.primaries.flush.total-time` (*cumulative*)<br>
 - `counter.indices.primaries.get.exists-time` (*cumulative*)<br>
 - `counter.indices.primaries.get.exists-total` (*cumulative*)<br>
 - `counter.indices.primaries.get.missing-time` (*cumulative*)<br>
 - `counter.indices.primaries.get.missing-total` (*cumulative*)<br>
 - `counter.indices.primaries.get.time` (*cumulative*)<br>
 - `counter.indices.primaries.indexing.delete-time` (*cumulative*)<br>
 - `counter.indices.primaries.indexing.delete-total` (*cumulative*)<br>
 - `counter.indices.primaries.indexing.index-time` (*cumulative*)<br>
 - `counter.indices.primaries.indexing.index-total` (*cumulative*)<br>
 - `counter.indices.primaries.merges.total` (*cumulative*)<br>
 - `counter.indices.primaries.merges.total-docs` (*cumulative*)<br>
 - `counter.indices.primaries.merges.total-time` (*cumulative*)<br>
 - `counter.indices.primaries.refresh.total` (*cumulative*)<br>
 - `counter.indices.primaries.refresh.total-time` (*cumulative*)<br>
 - `counter.indices.primaries.search.fetch-time` (*cumulative*)<br>
 - `counter.indices.primaries.search.fetch-total` (*cumulative*)<br>
 - `counter.indices.primaries.search.query-time` (*cumulative*)<br>
 - `counter.indices.primaries.search.query-total` (*cumulative*)<br>
 - `counter.indices.primaries.segments.count` (*cumulative*)<br>
 - `counter.indices.primaries.translog.operations` (*cumulative*)<br>
 - `counter.indices.primaries.warmer.total` (*cumulative*)<br>
 - `counter.indices.primaries.warmer.total.primaries.warmer.total-time` (*cumulative*)<br>
 - `counter.indices.refresh.time` (*cumulative*)<br>
 - `counter.indices.refresh.total` (*cumulative*)<br>
 - `counter.indices.search.fetch-time` (*cumulative*)<br>
 - `counter.indices.search.fetch-total` (*cumulative*)<br>
 - ***`counter.indices.search.query-time`*** (*cumulative*)<br>    Total time spent in search queries (milliseconds)
 - ***`counter.indices.search.query-total`*** (*cumulative*)<br>    The total number of search requests since node startup
 - `counter.indices.search.scroll-time` (*cumulative*)<br>
 - `counter.indices.search.scroll.total` (*cumulative*)<br>
 - `counter.indices.total.fielddata.evictions` (*cumulative*)<br>
 - `counter.indices.total.flush.periodic` (*cumulative*)<br>
 - `counter.indices.total.get.exists-time` (*cumulative*)<br>
 - `counter.indices.total.get.exists-total` (*cumulative*)<br>
 - `counter.indices.total.get.missing-time` (*cumulative*)<br>
 - `counter.indices.total.get.missing-total` (*cumulative*)<br>
 - `counter.indices.total.get.time` (*cumulative*)<br>
 - `counter.indices.total.get.total` (*cumulative*)<br>
 - `counter.indices.total.indexing.delete-time` (*cumulative*)<br>
 - `counter.indices.total.indexing.delete-total` (*cumulative*)<br>
 - `counter.indices.total.indexing.index-time` (*cumulative*)<br>
 - ***`counter.indices.total.indexing.index-total`*** (*cumulative*)<br>    The total number of index requests per cluster
 - ***`counter.indices.total.merges.total`*** (*cumulative*)<br>    Total number of merges per cluster
 - `counter.indices.total.merges.total-docs` (*cumulative*)<br>
 - `counter.indices.total.merges.total-time` (*cumulative*)<br>
 - `counter.indices.total.search.fetch-total` (*cumulative*)<br>
 - `counter.indices.total.search.query-time` (*cumulative*)<br>
 - ***`counter.indices.total.search.query-total`*** (*cumulative*)<br>    The total number of search requests per cluster
 - `counter.indices.total.translog.earliest_last_modified_age` (*cumulative*)<br>
 - `counter.indices.total.translog.uncommitted_operations` (*cumulative*)<br>
 - `counter.indices.total.translog.uncommitted_size_in_bytes` (*cumulative*)<br>
 - `counter.jvm.gc.count` (*cumulative*)<br>
 - `counter.jvm.gc.old-count` (*cumulative*)<br>
 - `counter.jvm.gc.old-time` (*cumulative*)<br>
 - ***`counter.jvm.gc.time`*** (*cumulative*)<br>    Total garbage collection time (milliseconds)
 - `counter.jvm.uptime` (*cumulative*)<br>
 - `counter.thread_pool.completed` (*cumulative*)<br>
 - ***`counter.thread_pool.rejected`*** (*cumulative*)<br>    Number of rejected thread pool requests
 - `counter.transport.rx.count` (*cumulative*)<br>
 - `counter.transport.rx.size` (*cumulative*)<br>
 - `counter.transport.tx.count` (*cumulative*)<br>
 - `counter.transport.tx.size` (*cumulative*)<br>
 - ***`gauge.cluster.active-primary-shards`*** (*gauge*)<br>    The number of active primary shards
 - ***`gauge.cluster.active-shards`*** (*gauge*)<br>    The number of active shards
 - `gauge.cluster.initializing-shards` (*gauge*)<br>    The number of currently initializing shards
 - ***`gauge.cluster.number-of-data_nodes`*** (*gauge*)<br>    The current number of data nodes in the cluster
 - ***`gauge.cluster.number-of-nodes`*** (*gauge*)<br>    Total number of nodes in the cluster
 - ***`gauge.cluster.relocating-shards`*** (*gauge*)<br>    The number of shards that are currently being relocated
 - `gauge.cluster.status` (*gauge*)<br>    The health status of the cluster
 - ***`gauge.cluster.unassigned-shards`*** (*gauge*)<br>    The number of shards that are currently unassigned
 - `gauge.http.current_open` (*gauge*)<br>
 - ***`gauge.indices.cache.field.size`*** (*gauge*)<br>    Field data size (bytes)
 - ***`gauge.indices.cache.filter.size`*** (*gauge*)<br>    Filter cache size (bytes)
 - ***`gauge.indices.docs.count`*** (*gauge*)<br>    Number of documents on this node
 - ***`gauge.indices.docs.deleted`*** (*gauge*)<br>    Number of deleted documents on this node
 - `gauge.indices.get.current` (*gauge*)<br>
 - `gauge.indices.indexing.delete-current` (*gauge*)<br>
 - `gauge.indices.indexing.index-current` (*gauge*)<br>
 - ***`gauge.indices.merges.current`*** (*gauge*)<br>    Number of active merges
 - `gauge.indices.merges.current-docs` (*gauge*)<br>
 - `gauge.indices.merges.current-size` (*gauge*)<br>
 - `gauge.indices.merges.total-docs` (*gauge*)<br>
 - `gauge.indices.primaries.completion.size` (*gauge*)<br>
 - `gauge.indices.primaries.docs.count` (*gauge*)<br>
 - `gauge.indices.primaries.docs.deleted` (*gauge*)<br>
 - `gauge.indices.primaries.fielddata.memory-size` (*gauge*)<br>
 - `gauge.indices.primaries.flush.periodic` (*gauge*)<br>
 - `gauge.indices.primaries.get.current` (*gauge*)<br>
 - `gauge.indices.primaries.indexing.delete-current` (*gauge*)<br>
 - `gauge.indices.primaries.indexing.index-current` (*gauge*)<br>
 - `gauge.indices.primaries.merges.current` (*gauge*)<br>
 - `gauge.indices.primaries.merges.current-docs` (*gauge*)<br>
 - `gauge.indices.primaries.merges.current-size` (*gauge*)<br>
 - `gauge.indices.primaries.search.fetch-current` (*gauge*)<br>
 - `gauge.indices.primaries.search.open-contexts` (*gauge*)<br>
 - `gauge.indices.primaries.search.query-current` (*gauge*)<br>
 - `gauge.indices.primaries.segments.index-writer-memory` (*gauge*)<br>
 - `gauge.indices.primaries.segments.memory` (*gauge*)<br>
 - `gauge.indices.primaries.segments.version-map-memory` (*gauge*)<br>
 - `gauge.indices.primaries.store.size` (*gauge*)<br>
 - `gauge.indices.primaries.translog.earliest_last_modified_age` (*gauge*)<br>
 - `gauge.indices.primaries.translog.size` (*gauge*)<br>
 - `gauge.indices.primaries.translog.uncommitted_operations` (*gauge*)<br>
 - `gauge.indices.primaries.translog.uncommitted_size_in_bytes` (*gauge*)<br>
 - `gauge.indices.primaries.warmer.current` (*gauge*)<br>
 - `gauge.indices.search.fetch-current` (*gauge*)<br>
 - `gauge.indices.search.open-contexts` (*gauge*)<br>
 - `gauge.indices.search.query-current` (*gauge*)<br>
 - `gauge.indices.search.scroll.current` (*gauge*)<br>
 - ***`gauge.indices.segments.count`*** (*gauge*)<br>    Number of segments on this node
 - `gauge.indices.segments.index-writer-size` (*gauge*)<br>
 - `gauge.indices.segments.size` (*gauge*)<br>
 - `gauge.indices.store.size` (*gauge*)<br>
 - ***`gauge.indices.total.docs.count`*** (*gauge*)<br>    Number of documents in the cluster
 - `gauge.indices.total.docs.deleted` (*gauge*)<br>
 - ***`gauge.indices.total.fielddata.memory-size`*** (*gauge*)<br>    Field data size (bytes)
 - ***`gauge.indices.total.filter-cache.memory-size`*** (*gauge*)<br>    Filter cache size (bytes)
 - `gauge.indices.total.get.current` (*gauge*)<br>
 - `gauge.indices.total.indexing.delete-current` (*gauge*)<br>
 - `gauge.indices.total.indexing.index-current` (*gauge*)<br>
 - `gauge.indices.total.merges.current` (*gauge*)<br>
 - `gauge.indices.total.merges.current-docs` (*gauge*)<br>
 - `gauge.indices.total.merges.current-size` (*gauge*)<br>
 - `gauge.indices.total.search.open-contexts` (*gauge*)<br>
 - `gauge.indices.total.search.query-current` (*gauge*)<br>
 - `gauge.indices.total.store.size` (*gauge*)<br>
 - `gauge.indices.translog.uncommitted_operations` (*gauge*)<br>
 - `gauge.indices.translog.uncommitted_size_in_bytes` (*gauge*)<br>
 - ***`gauge.jvm.mem.heap-committed`*** (*gauge*)<br>    Total heap committed by the process (bytes)
 - ***`gauge.jvm.mem.heap-used`*** (*gauge*)<br>    Total heap used (bytes)
 - `gauge.jvm.mem.non-heap-committed` (*gauge*)<br>
 - `gauge.jvm.mem.non-heap-used` (*gauge*)<br>
 - `gauge.jvm.mem.pools.old.max_in_bytes` (*gauge*)<br>
 - `gauge.jvm.mem.pools.old.used_in_bytes` (*gauge*)<br>
 - `gauge.jvm.mem.pools.young.max_in_bytes` (*gauge*)<br>
 - `gauge.jvm.mem.pools.young.used_in_bytes` (*gauge*)<br>
 - `gauge.jvm.threads.count` (*gauge*)<br>
 - `gauge.jvm.threads.peak` (*gauge*)<br>
 - `gauge.process.cpu.percent` (*gauge*)<br>
 - ***`gauge.process.open_file_descriptors`*** (*gauge*)<br>    Number of currently open file descriptors
 - `gauge.thread_pool.active` (*gauge*)<br>    Number of active threads
 - `gauge.thread_pool.largest` (*gauge*)<br>    Highest active threads in thread pool
 - `gauge.thread_pool.queue` (*gauge*)<br>    Number of Tasks in thread pool
 - `gauge.thread_pool.threads` (*gauge*)<br>    Number of Threads in thread pool
 - `gauge.transport.server_open` (*gauge*)<br>
 - `percent.jvm.mem.heap-used-percent` (*gauge*)<br>

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



