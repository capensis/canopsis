# Engine-pbehavior

## Services interactions

A simple schema which only displays direct interactions with all databases, RMQ, external APIs, etc.

```mermaid
flowchart
    EPH[engine-pbehavior]
    MDB[(MongoDB)]
    RMQ[(RabbitMQ)]
    TDB[(TimescaleDB)]
    R[(Redis)]
    EPH ---|fetch pbehaviors| MDB
    EPH ---|store computer pbehavior intervals| R
    EPH ---|receive/send events| RMQ
    EPH ---|store metrics| TDB
```

## Detailed schemas

The following schemas display flows of events by each use-case of business logic.

### Create/update/remove a pbehavior.

```mermaid
flowchart
    A[API]
    EF[engine-fifo]
    EPH[engine-pbehavior]
    R[(Redis)]
    A -.->|1 . Compute event| EPH
    EPH -.->|2 . Update computed intervals| R
    EPH -- 3 . Pbhenter/pbhleave events --> EF
```

### Update alarms in periodical process 

```mermaid
flowchart
    EF[engine-fifo]
    EPH[engine-pbehavior]
    R[(Redis)]
    EPH -.->|1 . Update computed intervals| R
    EPH -- 2 . Pbhenter/pbhleave events --> EF
```

### Scenarios

See [engine-action](./engine-actoin.md).
