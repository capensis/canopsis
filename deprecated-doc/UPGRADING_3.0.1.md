# Upgrading to Canopsis 3.0.1

## Removed engines

The following Canopsis engines need to be removed (through systemd or Docker):
- `linklist`
- `selector`
- `perfdata`
- `context`

*Warning:* `context` is removed, but `context-graph` MUST be kept.

You also need to remove any existing `linklist` task from MongoDB:
```bash
mongodb:PRIMARY> db.getCollection('object').remove({'task':'tasklinklist'})
```

## Added engines

`perfdata` MUST be replaced by the `metric` engine. Enable it through systemd or Docker.
