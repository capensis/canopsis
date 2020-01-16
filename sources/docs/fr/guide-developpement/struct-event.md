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
    "event_type":       // Event type (see below) - value field is `string` type
    "source_type":      // Source of event ("component", or "resource") - value field is `string` type
    "connector":        // Connector Type (gelf, nagios, snmp, ...) - value field is `string` type
    "connector_name":   // Connector Identifier (nagios1, nagios2, ...) - value field is `string` type
    "component":        // Component's name - value field is `string` type
    "resource":         // Resource's name (only if source_type is "resource") - value field is `string` type

    // /!\ The following is optional /!\

    "timestamp":        // UNIX timestamp for when the event  was emitted (optional: set by the server to now) - value field is `number` type
    "output":           // Message - value field is `string` type
    "long_output":      // Description - value field is `string` type
}
```

## Ajout d'éléments et personnalisation

Aprés avoir défini la structure de base de l'évènement, choississez ce que vous voulez ajouter à celui-ci et ajoutez les champs suivants.

### Event Check Structure

```javascript
{
    "event_type": "check",  // mandatory - value field is `string` type

    "state":                // Check state (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL), default is 0 - value field is `number` type
}
```

### Event Acknowledgment Structure

```javascript
{
    "event_type": "ack",    // mandatory - value field is `string` type

    "author":               // Acknowledgment author, optional - value field is `string` type
    "output":               // Acknowledgment comment, optional - value field is `string` type
}
```

### Event Declareticket Structure

```javascript
{
    "event_type": "declareticket",    // mandatory - value field is `string` type

    "author":               // Declareticket author, optional - value field is `string` type
    "output":               // Declareticket comment, optional - value field is `string` type
}
```

### Event Assocticket Structure

```javascript
{
    "event_type": "assocticket",    // mandatory - value field is `string` type

    "author":               // Assocticket author, optional - value field is `string` type
    "ticket":               // Assocticket number, optional - value field is `string` type
    "output":               // Assocticket comment, optional - value field is `string` type
}
```

### Event Snooze Structure

```javascript
{
  "event_type": "snooze",   // mandatory - value field is `string` type

  "duration":         // snooze duration, in seconds - value field is `number` type
  "author":           // snooze author, optional - value field is `string` type
  "output":           // snooze comment, optional - value field is `string` type
}
```

### Event Changestate Structure

```javascript
{
  "event_type": "changestate",   // mandatory

  "author":           // changestate author, optional
  "output":           // changestate comment, optional
}
```

### Event Cancel Structure

```javascript
{
    "event_type": "cancel",     // mandatory - value field is `string` type

    "author":               // author, optional - value field is `string` type
    "output":               // comment, optional - value field is `string` type
}
```

### Event Undo Cancel Structure

```javascript
{
    "event_type": "uncancel",   // mandatory - value field is `string` type

    "author":               // author, optional - value field is `string` type
    "output":               // comment, optional - value field is `string` type
}
```


### Event Ackremove Structure

```javascript
{
    "event_type": "ackremove",  // mandatory - value field is `string` type

    "author":               // author, optional - value field is `string` type
    "output":               // comment, optional - value field is `string` type
}
```

### Event SNMP Structure

```javascript
{
    "event_type": "trap",  // mandatory - value field is `string` type

    "snmp_severity":        // SNMP severity, mandatory - value field is `string` type
    "snmp_state":           // SNMP state, mandatory - value field is `string` type
    "snmp_oid":             // SNMP oid, mandatory - value field is `string` type
}
```

### Event Statistics Counter Increment Structure

```javascript
{
    "event_type": "statcounterinc",     // mandatory

    "stat_name":            // The name of the counter to increment, mandatory
    "alarm":                // The alarm, mandatory
    "entity":               // The entity which sent the event, mandatory
}
```
Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.

### Event Statistics Duration Structure

```javascript
{
    "event_type": "statduration",   // mandatory

    "stat_name":            // The name of the duration, mandatory
    "duration":             // The value of the duration (in seconds), mandatory
    "current_alarm":        // The alarm, mandatory
    "current_entity":       // The entity which sent the event, mandatory
}
```

Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.

### Event Statistics State Interval Structure

```javascript
{
    "event_type": "statstateinterval",      // mandatory

    "stat_name":            // The name of the state, mandatory
    "duration":             // The time spent in this state (in seconds), mandatory
    "state":                // The value of the state, mandatory
    "alarm":                // The alarm, mandatory
    "entity":               // The entity which sent the event, mandatory
}
```

Le champ `alarm` devrait contenir la valeur de l'alarme sous forme d'objet JSON.
Le champ `entity` devrait contenir l'entité sous forme d'objet JSON.


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
snooze | Utilisé pour placer un Snooze sur une alarme |
declareticket | Utilisé pour déclarer un ticket |
assocticket | Utilisé pour associer un ticket |
changestate | Utilisé pour changer et verrouiller la criticité d'une alarme |
