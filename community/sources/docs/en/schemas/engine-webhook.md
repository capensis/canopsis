# Engine-webhook

## Services interactions

A simple schema which only displays direct interactions with all databases, RMQ, external APIs, etc.

```mermaid
flowchart
    EW[engine-webhook]
    MDB[(MongoDB)]
    RMQ[(RabbitMQ)]
    TDB[(TimescaleDB)]
    NTC[Notification tool client]
    TTC[Ticketing tool client]
    EW ---|store requests| MDB
    EW ---|receive/send events| RMQ
    EW ---|store metrics| TDB
    EW ---|create notifications| NTC
    EW ---|create tickets| TTC
```

## Detailed schemas

The following schemas display flows of events by each use-case of business logic.

### Send a notification by a webhook.

```mermaid
flowchart
    EAC[engine-action]
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    EW[engine-webhook]
    OE[other engines]
    R[(Redis)]
    NTC[Notification tool client]
    EF -- 1 . Event --> ECH
    ECH -- 2 . Event --> EAX
    EAX -- 3 . Event --> OE
    OE -- 4 . Event --> EAC
    EAC -.->|5 . Store scenario executions| R
    EAC -.->|6 . Store request| MDB
    EAC -.->|7 . Run webhook| EW
    EW -.->|8 . Send HTTP request| NTC
    EW -.->|9 . Update alarm| EAX
    EAX -.->|10 . Result alarm| EAC
```

### Create a ticket by a webhook.

```mermaid
flowchart
    EAC[engine-action]
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    EW[engine-webhook]
    OE[other engines]
    MDB[(MongoDB)]
    R[(Redis)]
    TTC[Ticketing tool client]
    EF -- 1 . Event --> ECH
    ECH -- 2 . Event --> EAX
    EAX -- 3 . Event --> OE
    OE -- 4 . Event --> EAC
    EAC -.->|5 . Store scenario executions| R
    EAC -.->|6 . Store request| MDB
    EAC -.->|7 . Run webhook| EW
    EW -.->|8 . Send HTTP request| TTC
    EW -.->|9 . Update alarm| EAX
    EAX -.->|10 . Result alarm| EAC
```

### Create a ticket by a declare ticket rule.

```mermaid
flowchart
    A[API]
    EAX[engine-axe]
    EW[engine-webhook]
    TTC[Ticketing tool client]
    A -.->|1 . Store request| MDB
    A -.->|2 . Run declare ticket| EW
    EW -.->|3 . Send HTTP request| TTC
    EW -.->|4 . Update alarm| EAX
    A -.->|5 . Fetch result| MDB
```
