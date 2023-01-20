# Structure des évènements

## Focus AMQP

-   Vhost: canopsis
-   Routing key: `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`
-   Exchange: canopsis.events
-   Exchange Options: type: "topic", durable: true, auto_delete: false
-   Content Type: "application/json"

## Structure basique d'un évènement

Voici la structure de base d'un [évènement](../../guide-utilisation/vocabulaire/index.md#evenement), commune à tous les [types d'évènements](#liste-des-types-devenements).

```javascript
{
    "event_type":       // Event type (see below) - value field is `string` type
    "source_type":      // Source of event ("component", or "resource") - value field is `string` type
    "connector":        // Connector Type (gelf, nagios, snmp, ...) - value field is `string` type
    "connector_name":   // Connector Identifier (nagios1, nagios2, ...) - value field is `string` type
    "component":        // Component's name - value field is `string` type
    "resource":         // Resource's name (only if source_type is "resource") - value field is `string` type

    // /!\ The following is optional /!\

    "timestamp":        // UNIX timestamp for when the event  was emitted (optional: set by the server to now) - value field is an integer `number` type
    "output":           // Message - value field is `string` type
    "long_output":      // Description - value field is `string` type
}
```

## Liste des types d'évènements

Certains de ces événements déclenchent également un [trigger](../../guide-administration/architecture-interne/triggers.md).

Type              | Description                                                                                                                                |
------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
ack               | Acquitte une alarme                                                                                                                        |
ackremove         | Supprimer l'acquittement d'une alarme. (champ ack supprimé)                                                                                |
assocticket       | Associe un ticket                                                                                                                          |
declareticket     | Envoie un trigger declareticket                                                                                                            |
cancel            | Annule un évènement et met son statut dans un état "cancel", supprime également l'acquittement de l'évènement référent, le cas échéant.    |
changestate       | Change et verrouille la criticité d'une alarme                                                                                             |
check             | Envoie le résultat d'un check (depuis Nagios, Icinga,...)                                                                                  |
comment           | Ajoute un commentaire sur une alarme                                                                                                           |
snmp              | Envoyé par le connecteur [`snmp2canopsis`](../../interconnexions/Supervision/SNMPtrap.md) au moteur [`snmp`](../../guide-administration/moteurs/moteur-snmp.md) |
snooze            | Place un Snooze sur une alarme                                                                                                             |
statcounterinc    | Incrémente un compteur dans l'engine statistics                                                                                            |
statduration      | Ajoute une durée dans l'engine statistics                                                                                                  |
statstateinterval | Ajoute un état d'intervalle dans l'engine statistics                                                                                       |
uncancel          | Annule un évènement sur l'alarme et restaure son statut précédent                                                                          |
updatewatcher     | Déclenche la mise à jour de l'état d'un watcher (interne)                                                                                  |

## Ajout d'éléments et personnalisation

Après avoir défini la structure de base de l'[évènement](../../guide-utilisation/vocabulaire/index.md#evenement), choisissez le type d'événement que vous souhaitez envoyer et ajoutez les champs correspondants.

### Event Acknowledgment Structure

```javascript
{
    "event_type": "ack",    // mandatory - value field is `string` type

    "author":               // Acknowledgment author, optional - value field is `string` type
    "output":               // Acknowledgment comment, optional - value field is `string` type
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

### Event Assocticket Structure

```javascript
{
    "event_type": "assocticket",    // mandatory - value field is `string` type

    "author":               // Assocticket author, optional - value field is `string` type
    "ticket":               // Assocticket number, optional - value field is `string` type
    "output":               // Assocticket comment, optional - value field is `string` type
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

### Event Cancel Structure

```javascript
{
    "event_type": "cancel",     // mandatory - value field is `string` type

    "author":               // author, optional - value field is `string` type
    "output":               // comment, optional - value field is `string` type
}
```

### Event Changestate Structure

```javascript
{
  "event_type": "changestate",   // mandatory
  "state":                       // state that will be locked for the alarm (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL), default is 0, mandatory - value field is an integer `number` type

  "author":           // changestate author, optional
  "output":           // changestate comment, optional
}
```

### Event Check Structure

```javascript
{
    "event_type": "check",  // mandatory - value field is `string` type

    "state":                // Check state (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL), default is 0 - value field is an integer `number` type
}
```

### Event Comment Structure

```javascript
{
    "event_type": "comment", // mandatory - value field is `string` type

    "author":                // comment author
    "output":                // comment content
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

### Event Snooze Structure

```javascript
{
  "event_type": "snooze",   // mandatory - value field is `string` type

  "duration":         // snooze duration, in seconds - value field is an integer `number` type
  "author":           // snooze author, optional - value field is `string` type
  "output":           // snooze comment, optional - value field is `string` type
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

### Event Undo Cancel Structure

```javascript
{
    "event_type": "uncancel",   // mandatory - value field is `string` type

    "author":               // author, optional - value field is `string` type
    "output":               // comment, optional - value field is `string` type
}
```

### Event Updatewatcher Structure

```javascript
{
    "event_type": "updatewatcher",   // mandatory

    "connector": "watcher",          // fixed value
    "connector_name": "watcher",     // fixed value
    "source_type": "component",      // fixed value
    "component"                      // component value is the watcher id, mandatory
}
```
