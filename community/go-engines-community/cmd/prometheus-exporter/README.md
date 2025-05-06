# Prometheus Exporter for Canopsis

This exporter exposes internal metrics from Canopsis in prometheus-compatible format.

## Overview

- **Exporter Path**: `/metrics`
- **Port**: `9180` by default. You can change it using the `-port` flag.
- **Custom Metrics**:
  - `canopsis_eventfilter_errors`: Number of event filter errors
  - `canopsis_opened_alarms`: Number of opened alarms
  - `canopsis_resolved_alarms`: Number of resolved alarms
  - `canopsis_active_entities`: Number of active entities
  - `canopsis_disabled_entities`: Number of disabled entities
  - `canopsis_user_connections`: Number of user connections
  - `canopsis_enabled_users`: Number of enabled users
  - `canopsis_event_filters`: Number of event filters
  - `canopsis_active_pbehavior`: Number of active pbehaviors
  - `canopsis_meta_alarms_rules`: Number of meta alarm rules
  - `canopsis_dynamic_infos_rules`: Number of dynamic infos rules
  - `canopsis_engine_status{engine_name="<name>"}`: Engine status (1 = running, 0 = stopped)
  - `canopsis_last_exploitation_mod_time{type="<label>"}`: Last modification timestamp of exploitation elements

## Exporter Flags

The exporter supports the following command-line flags:

| Flag                     | Default | Description                                                  |
|--------------------------|---------|--------------------------------------------------------------|
| `-version`               | false   | Show the version information and exit                        |
| `-port`                  | 9180    | Port on which to run the exporter server                     |
| `-d`                     | false   | Enable debug logging                                         |
| `-updateMetricsInterval` | 10s     | Time interval between metric updates                         |

## Environment Variables

The exporter requires the following environment variables to be defined at runtime:

- `CPS_REDIS_URL`: The connection URI for Redis.
- `CPS_MONGO_URL`: The MongoDB connection URI.

### MongoDB Read Preference

To minimize load on the primary MongoDB node, it is recommended to use the `readPreference=secondary` option in the connection string. This allows the exporter to read from secondary members of the replica set without impacting the primary.

**Example URI:**

```
mongodb://cpsmongo:canopsis@localhost:27017,localhost:27018,localhost:27019/canopsis?replicaSet=cps&readPreference=secondary
```

## Prometheus Setup

To allow Prometheus to collect metrics from this exporter, include the following configuration in your `prometheus.yml` file. Replace `your-exporter-host` with the actual hostname or IP address where your Canopsis exporter is running.

```yaml
scrape_configs:
  - job_name: 'canopsis_exporter'
    static_configs:
      - targets: ['your-exporter-host:9180']
```
> **Note**: The Prometheus `scrape_interval` should be greater than or equal to the exporter's `-updateMetricsInterval` value.  
> If Prometheus scrapes more frequently than the exporter updates, it may collect outdated or duplicated metric values.
