# Statsng

The statsng engine is an engine that receives statistic events from the other
engines and uses them to compute statistics.

## Configuration

The statsng engine's configuration is stored in the `etc/statsng/engine.conf`
file, and has the following structure:

```
[ENGINE]
send_events = True
entity_tags = parent_service
pbehavior_tags = Maintenance,Pause
```

### Send events

The statistic events are only sent if the `send_events` option is set to
`True`. The statsng engine will not compute any statistics if it is set to
`False`.

### Entity Tags

The `entity_tags` option is a comma separated list of entity information ids.
These ids will be saved with the statistics, allowing them to be used in
filters in the [stats API](../developer_guide/apis/v2/stats.md).
