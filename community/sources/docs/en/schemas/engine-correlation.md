# Engine-correlation

## Services interactions

A simple schema which only displays direct interactions with all databases, RMQ, external APIs, etc.

```mermaid
flowchart
    ECO[engine-correlation]
    MDB[(MongoDB)]
    RMQ[(RabbitMQ)]
    TDB[(TimescaleDB)]
    R[(Redis)]
    ECO ---|store counters to detect meta alarms| R
    ECO ---|store meta alarms| MDB
    ECO ---|receive/send events| RMQ
    ECO ---|store metrics| TDB
```

## Detailed schemas

The following schemas display flows of events by each use-case of business logic.

### Create a meta alarm on an event.

```mermaid
flowchart
    EF[engine-fifo]
    EC[engine-che]
    EAX[engine-axe]
    ECO[engine-correlation]
    EN[next engine]
    MDB[(MongoDB)]
    R[(Redis)]
    EF -- 1 . Event --> EC
    EC -- 2 . Event --> EAX
    EAX -- 3 . Event --> ECO
    ECO -.->|4 . Check counters| R
    ECO -- 5 . Meta alarm event --> EF
    EF -- 6 . Meta alarm event --> EC
    EC -- 7 . Meta alarm event --> EAX
    EAX -.->|8 . Create meta alarm| MDB
    ECO -- 5. Event --> EN
```

### Update a meta alarm on an event.

```mermaid
flowchart
    EF[engine-fifo]
    EC[engine-che]
    EAX[engine-axe]
    EN[next engine]
    MDB[(MongoDB)]
    EF -- 1 . Event --> EC
    EC -- 2 . Event --> EAX
    EAX -.->|3 . Update meta alarm| MDB
    EAX -- 4. Event --> EN
```

### Update a child on a meta alarm change (ack, comment, declare ticket, etc.).

```mermaid
flowchart
    EF[engine-fifo]
    EC[engine-che]
    EAX[engine-axe]
    EN[next engine]
    MDB[(MongoDB)]
    EF -- 1 . Meta alarm event --> EC
    EC -- 2 . Meta alarm event --> EAX
    EAX -.->|3 . Update meta alarm| MDB
    EAX -- 4. Child event --> EF
    EAX -- 4. Meta alarm event --> EN
    EF -- 5. Child event --> EC
    EC -- 6. Child event --> EAX
    EAX -.->|7 . Update child alarm| MDB
    EAX -- 8. Child event --> EN
```
