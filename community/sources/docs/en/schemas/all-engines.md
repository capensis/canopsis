# All engines

## General schemas

A schema which only displays direct interactions between engines and between engines and external services.

### Pro

```mermaid
flowchart
    subgraph Host1
        H1R1[Resource1]
        H1R2[Resource2]
        H1MT[Monitoring tool Centreon/Zabbix/Icinga2/etc.]
        H1C[Canopsis connector]
        H1R1 --> H1MT
        H1R2 --> H1MT
        H1MT --> H1C
    end
    subgraph Host2
        H2R[External source Database/Junit/Logstash/etc.]
        H2C[Canopsis connector]
        H2R --> H2C
    end
    H1C -- AMQP --> EF
    H2C -- AMQP --> EF
    subgraph Canopsis
        EAC[engine-action]
        EAX[engine-axe]
        ECH[engine-che]
        EF[engine-fifo]
        EPH[engine-pbehavior]
        ESE[engine-service]
        ECO[engine-correlation]
        EDI[engine-dynamic-infos]
        ER[engine-remediation]
        EW[engine-webhook]
        A[API]
        EF -- AMQP --> ECH
        ECH -- AMQP --> EAX
        EAX -- AMQP --> ECO
        ECO -- AMQP --> ESE
        ESE -- AMQP --> EDI
        EDI -- AMQP --> EAC
        A -- AMQP --> EF
        EPH <-. AMQP .-> EAX
        A -. AMQP .-> EPH
        ER <-. AMQP .-> EAX
        EW -. AMQP .-> EAX
        EAC <-. AMQP .-> EW
        EAC <-. AMQP .-> EAX
        ESE <-. AMQP .-> EAX
        A -. AMQP .-> EDI
        A -. AMQP .-> ER
        A -. AMQP .-> EW
    end
    subgraph Pilot computer
        UI[Browser UI]
        NTC[Notification tool client]
        TTC[Ticketing tool client]
        TSC[Task scheduler client]
    end
    A <-- HTTP --> UI
    subgraph Host3
        NT[Notification tool Mattermost/E-mail/Etc.]
    end
    EW -- HTTP --> NT
    NT --> NTC
    subgraph Host4
        TT[Ticketing tool Service now/etc.]
    end
    EW -- HTTP --> TT
    TT <--> TTC
    subgraph Host5
        TS[Task scheduler Rundeck/AWX/Jenkins/etc.]
    end
    ER -- HTTP --> TS
    TS <--> TSC
    subgraph Host6
        DS[Data source Service now/etc.]
    end
    EF -- HTTP --> DS
    ECH -- HTTP --> DS
```

### Community

```mermaid
flowchart
    subgraph Host1
        H1R1[Resource1]
        H1R2[Resource2]
        H1MT[Monitoring tool Centreon/Zabbix/Icinga2/etc.]
        H1C[Canopsis connector]
        H1R1 --> H1MT
        H1R2 --> H1MT
        H1MT --> H1C
    end
    subgraph Host2
        H2R[External source Database/Junit/Logstash/etc.]
        H2C[Canopsis connector]
        H2R --> H2C
    end
    H1C -- AMQP --> EF
    H2C -- AMQP --> EF
    subgraph Canopsis
        EAC[engine-action]
        EAX[engine-axe]
        ECH[engine-che]
        EF[engine-fifo]
        EP[engine-pbehavior]
        ESE[engine-service]
        A[API]
        EF -- AMQP --> ECH
        ECH -- AMQP --> EAX
        EAX -- AMQP --> ESE
        ESE -- AMQP --> EAC
        A -- AMQP --> EF
        EP <-. AMQP .-> EAX
        A -. AMQP .-> EP
        EAC <-. AMQP .-> EAX
        ESE <-. AMQP .-> EAX
    end
    subgraph Pilot computer
        UI[Browser UI]
    end
    A <-- HTTP --> UI
```

## Detailed schemas

- [engine-fifo](./engine-fifo.md)
- [engine-che](./engine-che.md)
- [engine-axe](./engine-axe.md)
- [engine-correlation](./engine-correlation.md)
- [engine-service](./engine-service.md)
- [engine-dynamic-infos](./engine-dynamic-infos.md)
- [engine-action](./engine-action.md)
- [engine-pbehavior](./engine-pbehavior.md)
- [engine-remediation](./engine-remediation.md)
- [engine-webhook](./engine-webhook.md)
