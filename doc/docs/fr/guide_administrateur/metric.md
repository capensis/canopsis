# Metric

The metric engine is an engine that writes the performance date sent with the
events (in `perf_data` and `perf_data_array`) in the influxdb database.

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
