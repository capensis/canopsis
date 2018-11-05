# Structure d'un évènement

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

    // /!\ The following is optional /!\
    
    'hostgroups':       // Nagios hostgroups for component, default []
    'servicegroups':    // Nagios servicegroups for resource, default []
    'timestamp':        // UNIX timestamp for when the event  was emitted (optional: set by the server to now)

    'output':           // Message
    'long_output':      // Description
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
    
    // /!\ The following is optional /!\
    
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
    'event_type': 'ack',    // mandatory

    'ref_rk':               // Routing Key of acknowledged event, mandatory
    'author':               // Acknowledgment author, mandatory
    'output':               // Acknowledgment comment, mandatory
 
```

### Event Cancel Structure

```javascript
{
    'event_type': 'cancel',     // mandatory

    'ref_rk':               // Routing Key of event, mandatory
    'author':               // author, mandatory
    'output':               // comment, mandatory
}
```

### Event Undo Cancel Structure

```javascript
{
    'event_type': 'uncancel',   // mandatory

    'ref_rk':               // Routing Key of event, mandatory
    'author':               // author, mandatory
    'output':               // comment, mandatory
}
```

### Event Ackremove Structure

```javascript
{
    'event_type': 'ackremove',  // mandatory

    'ref_rk':               // Routing Key of event, mandatory
    'author':               // author, mandatory
    'output':               // comment, mandatory
}
```

### Event Downtime Structure

```javascript
{
    'event_type': 'downtime',   // mandatory

    'author':               // Downtime author, mandatory
    'output':               // Downtime comment, mandatory
    'start':                // UNIX timestamp for downtime's start, mandatory
    'end':                  // UNIX timestamp for downtime's end, mandatory
    'duration':             // Downtime's duration, mandatory
    'entry':                // Downtime's schedule date/time (as a UNIX timestamp), mandatory
    'fixed':                // Does the downtime starts at 'start' or at next check after 'start' ?, mandatory
    'downtime_id':          // Downtime's identifier, mandatory
}
```

### Event SNMP Structure

```javascript
{
    'event_type': 'trap',
    'snmp_severity':        // SNMP severity, mandatory
    'snmp_state':           // SNMP state, mandatory
    'snmp_oid':             // SNMP oid, mandatory
}
```

### Event Perf Structure

Un évènement de type \'perf\' ne sera jamais sauvegarder dans une base de données, il est uniquement utilisé pour envoyer des perfdata :

```javascript
{
    'event_type': 'perf',   // mandatory

    'perf_data':            // mandatory
    'perf_data_array':      // mandatory
}
```

See bellow for more informations about those fields.

### Event Statistics Counter Increment Structure

```javascript
{
    'event_type': 'statcounterinc',     // mandatory

    'stat_name':            // The name of the counter to increment, mandatory
    'alarm':                // The alarm, mandatory
    'entity':               // The entity which sent the event, mandatory
}
```
Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.

### Event Statistics Duration Structure

```javascript
{
    'event_type': 'statduration',   // mandatory

    'stat_name':            // The name of the duration, mandatory
    'duration':             // The value of the duration (in seconds), mandatory
    'current_alarm':        // The alarm, mandatory
    'current_entity':       // The entity which sent the event, mandatory
}
```

Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.

### Event Statistics State Interval Structure

```javascript
{
    'event_type': 'statstateinterval',      // mandatory

    'stat_name':            // The name of the state, mandatory
    'duration':             // The time spent in this state (in seconds), mandatory
    'state':                // The value of the state, mandatory
    'alarm':                // The alarm, mandatory
    'entity':               // The entity which sent the event, mandatory
}
```

Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.

## Metrology

Pour envoyer des perfdata vers Canopsis, vous avez juste besoin de spécifier l'un des champs suivant :

```javascript
{
    'perf_data':        // Performance data ("Nagios format":http://nagiosplug.sourceforge.net/developer-guidelines.html#AEN201)
    'perf_data_array':  // Array of performance data with metric's type ('GAUGE', 'DERIVE', 'COUNTER', 'ABSOLUTE'), Ex:
    [
        {'metric': 'shortterm', 'value': 0.25, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
        {'metric': 'midterm',   'value': 0.16, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
        {'metric': 'longterm',  'value': 0.12, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' }
    ]
}
```

### Basic Alert Structure

Un alarme est le résultat de l'analyse des évènements. Elle historise et résume les changements d'état, les actions utilisateurs (acquittement, mise en pause, etc.).
Dans MongoDB, il contient les champs suivant.

```javascript
{
    '_id':          // MongoDB document ID
    'event_id':     // Event identifier (the routing key)
}
```


## List of event types

Type | Description |
-----|-------------|
check | Utilisé pour envoyer le résultat d'un check (depuis Nagios, Icinga,...)  |
comment | Utilisé pour envoyer un commentaire |
perf | Utilisé pour envoyer seulement des perfdata |
selector | Envoyé par l'engine selector |
sla |  Envoyé par l'engine selector sla |
statcounterinc | Utilisé pour incrémenter un compteur dans l'engine statistics |
statduration | Utilisé pour ajouter une durée dans l'engine statistics |
statstateinterval | Utilisé pour ajouter une état d'intervale dans l'engine statistics |
trap | Utilisé pour envoyer des traps SNMP|
user | Utilisé par l'utilisateur pour evoyer des informations |
ack | Utilisé pour acquitter une alerte |
downtime | Utilisé pour programmer un downtime |
cancel | Utilisé pour cancel un évènement et mettre son statut dans un état "cancel", supprime également l'acquittement de l'événement référent, le cas échéant.  |
uncancel | Utilisé pour annuler un événement. le statut précédent est restauré et accusé de réception aussi, le cas échéant.  |
ackremove | Utilisé pour supprimer un accusé de réception d'un événement. (champ ack supprimé et collection ack mise à jour) |