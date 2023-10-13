# Engine-axe

## Services interactions

A simple schema which only displays direct interactions with all databases, RMQ, external APIs, etc.

```mermaid
flowchart
    EAX[engine-axe]
    MDB[(MongoDB)]
    RMQ[(RabbitMQ)]
    TDB[(TimescaleDB)]
    EAX ---|store alarms| MDB
    EAX ---|receive/send events| RMQ
    EAX ---|store metrics| TDB
```

## Detailed schemas

The following schemas display flows of events by each use-case of business logic.

### Create/update an alarm on an event.

```mermaid
flowchart
    C[Canopsis connector]
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    EN[next engine]
    MDB[(MongoDB)]
    C -- 1 . Event --> EF
    EF -- 2 . Event --> ECH
    ECH -- 3 . Event --> EAX
    EAX -.->|4 . Store alarm| MDB
    EAX -- 5 . Event --> EN
```

### Create/update/remove a service.

```mermaid
flowchart
    A[API]
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    MDB[(MongoDB)]
    A -- 1 . Compute event --> EF
    EF -- 2 . Compute event --> ECH
    ECH -.->|3 . Update service dependencies| MDB
    ECH -- 4 . Compute event --> EAX
    EAX -.->|5 . Update service counters| MDB
    EAX -- 6 . Check event for service --> EF
```

### Update a service on an alarm change.

```mermaid
flowchart
    C[Canopsis connector]
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    MDB[(MongoDB)]
    C -- 1 . Event --> EF
    EF -- 2 . Event --> ECH
    ECH -.->|3 . Update resource dependencies| MDB
    ECH -- 4 . Event --> EAX
    EAX -.->|5 . Update service counters| MDB
    EAX -- 6 . Check event for service --> EF
```

### Scenarios

See [engine-action](./engine-action.md).

### Declare ticket rules

See [engine-webhook](./engine-webhook.md).

### Instructions

See [engine-remediation](./engine-remediation.md).

### Meta alarms

See [engine-correlation](./engine-correlation.md).
