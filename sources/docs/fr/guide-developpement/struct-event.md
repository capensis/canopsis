# Structure d'un évènement

**TODO :** à continuer avec cette doc https://git.canopsis.net/canopsis/old-deprecated-doc/blob/master/old_doc/en/developer_guide/specifications/event.md  

## Focus AMQP

-   Vhost: canopsis  
-   Routing key: `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`  
-   Exchange: canopsis.events  
-   Exchange Options: type: "topic", durable: true, auto_delete: false  
-   Content Type: "application/json"  

## Structure basique d'un évènement  

Voici la structure de base d'un événement, commune à tous les type d'événements.

```javascript
{
    'connector':        // Connector Type (gelf, nagios, snmp, ...)
    'connector_name':   // Connector Identifier (nagios1, nagios2, ...)
    'event_type':       // Event type (see below)
    'source_type':      // Source of event ('component', or 'resource')
    'component':        // Component's name
    'resource':         // Resource's name (only if source_type is 'resource')

    // The following is optional
    'hostgroups':       // Nagios hostgroups for component, default []
    'servicegroups':    // Nagios servicegroups for resource, default []
    'timestamp':        // UNIX timestamp for when the event  was emitted (optional: set by the server to now)

    'output':           // Message
    'long_output':      // Description
    'display_name':     // Name to display in Canopsis
    'tags':             // Tags for the event (optional, the server adds connector, connector_name, event_type, source_type, component and resource if present)

    'perf_data':        // Nagios formatted perfdata string
    'perf_data_array':  // array of metrics (see below)
}
```

## Ajout d'éléments et personnalisation

Aprés avoir définit la structure de base de l'événement, choississez ce que vous voulez ajouter à celui-ci et ajoutez les champs suivants.

### Event Check Structure

```javascript
{
    'event_type': 'check',

    'state':                // Check state (0 - OK, 1 - WARNING, 2 - CRITICAL, 3 - UNKNOWN), default is 0
    'state_type':           // Check state type (0 - SOFT, 1 - HARD), default is 1
    'status':               // 0 == Ok | 1 == En cours | 2 == Furtif | 3 == Bagot | 4 == Annule
    // The following is optional
    'scheduled':            // True if the check was scheduled, False otherwise

    'check_type':           // Nagios Check Type (host or service)
    'current_attempt':      // Attempt ID for the check
    'max_attempts':         // Max attempts before sending HARD state
    'execution_time':       // Check duration
    'latency':              // Check latency (time between schedule and execution)
    'command_name':         // Check command
}
```

### Event Log Structure

```javascript
{
    'event_type': 'log',

    'output':           // Becomes mandatory
    'long_output':      // Remains optional
    'display_name':     // Remains optional

    'level':            // Optional log level
    'facility':         // Optional log facility
}
```

### Event Acknowledgment Structure

```javascript
{
    'event_type': 'ack',

    'ref_rk':               // Routing Key of acknowledged event
    'author':               // Acknowledgment author
    'output':               // Acknowledgment comment
}
```

### Event Cancel Structure

```javascript
{
    'event_type': 'cancel',

    'ref_rk':               // Routing Key of event
    'author':               // author
    'output':               // comment
}
```

### Event Undo Cancel Structure

```javascript
{
    'event_type': 'uncancel',

    'ref_rk':               // Routing Key of event
    'author':               // author
    'output':               // comment
}
```

### Event Ackremove Structure

```javascript
{
    'event_type': 'ackremove',

    'ref_rk':               // Routing Key of event
    'author':               // author
    'output':               // comment
}
```

### Event Downtime Structure

```javascript
{
    'event_type': 'downtime',

    'author':               // Downtime author
    'output':               // Downtime comment
    'start':                // UNIX timestamp for downtime's start
    'end':                  // UNIX timestamp for downtime's end
    'duration':             // Downtime's duration
    'entry':                // Downtime's schedule date/time (as a UNIX timestamp)
    'fixed':                // Does the downtime starts at 'start' or at next check after 'start' ?
    'downtime_id':          // Downtime's identifier
}
```

### Event SNMP Structure

```javascript
{
    'event_type': 'trap',
    'snmp_severity':        // SNMP severity
    'snmp_state':           // SNMP state
    'snmp_oid':             // SNMP oid
}
```

### Event Calendar Structure

```javascript
{
    'event_type': 'calendar',
    'resource':                 // iCal event UID
    'start':                    // iCal event start UNIX timestamp
    'end':                      // iCal event end UNIX timestamp
    'all_day':                  // True or False
    'output':                   // iCal event title
}
```

### Event Perf Structure

Un évènement de type \'perf\' ne sera jamais sauvegarder dans une base de données, il est uniquement utilisé pour envoyer des perfdata :  

```javascript
{
    'event_type': 'perf',

    'perf_data':
    'perf_data_array':
}
```

See bellow for more informations about those fields.

### Event Statistics Counter Increment Structure

```javascript
{
    'event_type': 'statcounterinc',

    'stat_name':            // The name of the counter to increment
    'alarm':                // The alarm
    'entity':               // The entity which sent the event
}
```
Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.  

### Event Statistics Duration Structure

```javascript
{
    'event_type': 'statduration',

    'stat_name':            // The name of the duration
    'duration':             // The value of the duration (in seconds)
    'current_alarm':        // The alarm
    'current_entity':       // The entity which sent the event
}
```

Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON. 

### Event Statistics State Interval Structure

```javascript
{
    'event_type': 'statstateinterval',

    'stat_name':            // The name of the state
    'duration':             // The time spent in this state (in seconds)
    'state':                // The value of the state
    'alarm':                // The alarm
    'entity':               // The entity which sent the event
}
```

Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON. 

