# Engine-remediation

## Services interactions

A simple schema which only displays direct interactions with all databases, RMQ, external APIs, etc.

```mermaid
flowchart
    ER[engine-remediation]
    MDB[(MongoDB)]
    RMQ[(RabbitMQ)]
    TDB[(TimescaleDB)]
    TSC[Task scheduler client]
    ER ---|store executions| MDB
    ER ---|receive/send events| RMQ
    ER ---|store metrics| TDB
    ER ---|run jobs| TSC
```

## Detailed schemas

The following schemas display flows of events by each use-case of business logic.

### Run an auto instruction on an event.

```mermaid
flowchart
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    ER[engine-remediation]
    OE[other engines]
    MDB[(MongoDB)]
    TSC[Task scheduler client]
    EF -- 1 . Event --> ECH
    ECH -- 2 . Event --> EAX
    EAX -- 3 . Event --> OE
    EAX -- 3 . Event --> ER
    ER -.->|4 . Create executions| MDB
    ER -.->|5 . Run jobs| TSC
    ER -.->|6 . Fetch result| TSC
    ER -.->|7 . Update alarm| EAX
```

### Run a simplified manual instruction.

```mermaid
flowchart
    A[API]
    EAX[engine-axe]
    ER[engine-remediation]
    MDB[(MongoDB)]
    TSC[Task scheduler client]
    A -.->|1 . Run instruction| ER
    ER -.->|2 . Create executions| MDB
    ER -.->|3 . Run jobs| TSC
    ER -.->|4 . Fetch result| TSC
    ER -.->|5 . Update alarm| EAX
```

### Run a manual instruction.

```mermaid
flowchart
    A[API]
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    ER[engine-remediation]
    MDB[(MongoDB)]
    TSC[Task scheduler client]
    A -.->|1 . Create executions| MDB
    A -.->|2 . Update alarm| EF
    EF -.->|3 . Update alarm| ECH
    ECH -.->|4 . Update alarm| EAX
    A -.->|5 . Run internal job| ER
    ER -.->|6 . Run external job| TSC
    ER -.->|7 . Fetch external result| TSC
    ER -.->|8 . Update alarm| EAX
    A -.->|9 . Fetch internal result| MDB
```

### Run an auto instruction on a scenario emitted trigger.

```mermaid
flowchart
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    ER[engine-remediation]
    OE[other engines]
    EAC[engine-action]
    MDB[(MongoDB)]
    TSC[Task scheduler client]
    EF -- 1 . Event --> ECH
    ECH -- 2 . Event --> EAX
    EAX -- 3 . Event --> OE
    OE -- 4 . Event --> EAC
    EAC -.->|5 . Update alarm| EAX
    EAX -.->|6 . Emitted trigger| ER
    ER -.->|7 . Create executions| MDB
    ER -.->|8 . Run jobs| TSC
    ER -.->|9 . Fetch result| TSC
    ER -.->|10 . Update alarm| EAX
```

### Run a scenario on an instruction emitted trigger.

```mermaid
flowchart
    EF[engine-fifo]
    ECH[engine-che]
    EAX[engine-axe]
    ER[engine-remediation]
    OE[other engines]
    EAC[engine-action]
    EF -- 1 . Event --> ECH
    ECH -- 2 . Event --> EAX
    EAX -- 3 . Event --> ER
    ER -.->|4 . Update alarm| EAX
    EAX -- 5 . Emitted trigger --> EF
    EF -- 6 . Emitted trigger --> ECH
    ECH -- 7 . Emitted trigger --> EAX
    EAX -- 8 . Emitted trigger --> OE
    OE -- 9 . Emitted trigger --> EAC
    EAC -.->|10 . Update alarm| EAX
```
