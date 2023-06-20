# Engine-che

## Services interactions

A simple schema which only displays direct interactions with all databases, RMQ, external APIs, etc.

```mermaid
flowchart
    ECH[engine-che]
    MDB[(MongoDB)]
    RMQ[(RabbitMQ)]
    TDB[(TimescaleDB)]
    DS[Data source Service now/etc.]
    ECH ---|fetch event filter rules and fetch external data| MDB
    ECH ---|receive/send events| RMQ
    ECH ---|store metrics| TDB
    ECH ---|fetch external data| DS
```

## Detailed schemas

The following schemas display flows of events by each use-case of business logic.

### Create/update an entity on an event.

```mermaid
flowchart
    C[Canopsis connector]
    EF[engine-fifo]
    ECH[engine-che]
    EN[next engine]
    MDB[(MongoDB)]
    C -- 1 . Event --> EF
    EF -- 2 . Event --> ECH
    ECH -. 3 . Store entity .-> MDB
    ECH -- 4 . Event --> EN
```

### Event enrichment flow

```mermaid
flowchart
    C[Canopsis connector]
    EF[engine-fifo]
    ECH[engine-che]
    EN[next engine]
    MDB[(MongoDB)]
    DS[Data source Service now/etc.]
    C -- 1 . Event --> EF
    EF -- 2 . Event --> ECH
    ECH -. 3 . Fetch external data .-> DS
    ECH -. 4 . Store entity .-> MDB
    ECH -- 5 . Enriched event --> EN
```
