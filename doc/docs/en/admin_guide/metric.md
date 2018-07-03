# Metric

The metric engine is an engine that writes the performance date sent with the
events (in `perf_data` and `perf_data_array`) in the influxdb database.

## Data model

The performance data are saved in an influxdb measurement with the same name as
the metric. This measurement can have three fields :

 - `value`: the value of the metric
 - `warn`: a warning threshold (may be `null`)
 - `crit`: a critical threshold (may be `null`)

It also has the tags `connector`, `connector_name`, `component` and `resource`.
You can add entity informations to the tags using the `tags` option in the
configuration.

## Configuration

The metric engine's configuration is stored in the `etc/metric/engine.conf`
file, and has the following structure:

```
[ENGINE]
tags = parent_service
```

### Tags

The `tags` option is a comma separated list of entity information ids.
These ids will be saved with the performance data.

Each `<information_id>` in this list will be a tag of the measurements created
by the metric engine. The value of this tag is the value stored in
`<entity>.infos.<information_id>`.
