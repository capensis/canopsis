# Structure d'un évènement

## Focus AMQP

-   Vhost: canopsis
-   Routing key: `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`
-   Exchange: canopsis.events
-   Exchange Options: type: "topic", durable: true, auto_delete: false
-   Content Type: "application/json"

## Structure basique d'un évènement

Voici la structure de base d'un évènement, commune à tous les type d'évènements.

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

}
```

## Ajout d'éléments et personnalisation

Aprés avoir défini la structure de base de l'évènement, choississez ce que vous voulez ajouter à celui-ci et ajoutez les champs suivants.

### Event Check Structure

```javascript
{
    'event_type': 'check',

    'state':                // Check state (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL), default is 0

    // /!\ The following is optional /!\

    'scheduled':            // True if the check was scheduled, False otherwise

    'check_type':           // Nagios Check Type (host or service)
    'current_attempt':      // Attempt ID for the check
    'max_attempts':         // Max attempts before sending HARD state
    'execution_time':       // Check duration
    'latency':              // Check latency (time between schedule and execution)
    'command_name':         // Check command

    'state_type':           // ONLY FOR INFO, USED ONLY INTERNALLY BY CANOPSIS, DO NOT TRY TO FEED AN EVENT WITH THIS : Check state type (0 - SOFT, 1 - HARD), default is 1
    'status':               // ONLY FOR INFO, USED ONLY INTERNALLY BY CANOPSIS, DO NOT TRY TO FEED AN EVENT WITH THIS : 0 == Ok | 1 == En cours | 2 == Furtif | 3 == Bagot | 4 == Annule

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

### Event SNMP Structure

```javascript
{
    'event_type': 'trap',
    'snmp_severity':        // SNMP severity, mandatory
    'snmp_state':           // SNMP state, mandatory
    'snmp_oid':             // SNMP oid, mandatory
}
```

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

### Basic Alert Structure

Un alarme est le résultat de l'analyse des évènements. Elle historise et résume les changements d'état, les actions utilisateurs (acquittement, mise en pause, etc.).
Dans MongoDB, il contient les champs suivants.

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
statcounterinc | Utilisé pour incrémenter un compteur dans l'engine statistics |
statduration | Utilisé pour ajouter une durée dans l'engine statistics |
statstateinterval | Utilisé pour ajouter un état d'intervalle dans l'engine statistics |
ack | Utilisé pour acquitter une alerte |
cancel | Utilisé pour cancel un évènement et mettre son statut dans un état "cancel", supprime également l'acquittement de l'évènement référent, le cas échéant.  |
uncancel | Utilisé pour annuler un évènement. le statut précédent est restauré et accusé de réception aussi, le cas échéant.  |
ackremove | Utilisé pour supprimer un accusé de réception d'un évènement. (champ ack supprimé et collection ack mise à jour) |
