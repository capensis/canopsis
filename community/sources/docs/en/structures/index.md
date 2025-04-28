
# Event Structure

## Focus on AMQP

- Vhost: canopsis
- Routing key: `<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]`
- Exchange: canopsis.events
- Exchange Options: type: "topic", durable: true, auto_delete: false
- Content Type: "application/json"

## Basic Event Structure

Here is the basic structure of an [event](../../guide-utilisation/vocabulaire/index.md#evenement), common to all [event types](#list-of-event-types).

```javascript
{
    "event_type":       // Event type (see below) - value is a `string`
    "source_type":      // Source of event ("component" or "resource") - value is a `string`
    "connector":        // Connector type (gelf, nagios, snmp, ...) - value is a `string`
    "connector_name":   // Connector identifier (nagios1, nagios2, ...) - value is a `string`
    "component":        // Component name - value is a `string`
    "resource":         // Resource name (only if source_type is "resource") - value is a `string`

    // Optional fields

    "timestamp":        // UNIX timestamp when the event was emitted (optional: server sets to now if omitted) - integer `number` type
                        // Since release 4.5, only timestamps within 24 hours are kept and forwarded; others are dropped by the FIFO engine.
                        // Also since 4.5, `last_update_date` field inside created alarm is initialized from the event's timestamp field, and `creation_date` is initialized from time.Now()
    "output":           // Message - value is a `string`
    "long_output":      // Description - value is a `string`
}
```

## List of Event Types

Canopsis handles the following types of events:

| Event type | Function | Origin | Internal |
|:-----------|:---------|:-------|:---------|
| ack | Acknowledge an alarm | UI/Engine-axe | no |
| ackremove | Remove acknowledgment from the alarm | UI/Engine-axe | no |
| assocticket | Associate a ticket with an alarm | UI/Engine-axe | no |
| cancel | Cancel an alarm | UI/Engine-axe | no |
| check | Create alarm or change status/state | Various connectors | no |
| contextupdate | Update entity without chenging LastEventDate | Various connectors | no |
| comment | Comment an alarm | UI | no |
| declareticket | Declare a ticket on the alarm | UI | no |
| done | Mark an alarm as done | ??? | no |
| changestate | Change an alarm's state | UI/Engine-axe | no |
| snooze | Snooze an alarm | UI/Engine-axe | no |
| unsnooze | Unsnooze an alarm | Engine-axe | yes |
| uncancel | Uncancel an alarm | ??? | no |
| metaalarm | Create a metaalarm | Engine-correlation | yes |
| metaalarmupdated | Notify that a metaalarm is updated | Engine-correlation | yes |
| pbhenter | Notify that an alarm enters PBH state | Engine-pbehavior | yes |
| pbhleaveandenter | Notify that an alarm leaves and enters a PBH state | Engine-pbehavior | yes |
| pbhleave | Notify that an alarm leaves PBH state | Engine-pbehavior | yes |
| pbhcreate | Create a PBH behavior | Engine-axe | yes |
| resolve_done | Resolve "done" alarms | Engine-axe | yes |
| resolve_cancel | Resolve canceled alarms | Engine-axe | yes |
| resolve_close | Resolve closed alarms | Engine-axe | yes |
| resolve_deleted | Resolve alarm for deleted entity | Engine-che | yes |
| updatestatus | Update an alarm's status | Engine-axe | yes |
| manual_metaalarm_group, manual_metaalarm_ungroup, manual_metaalarm_update | Group/ungroup and update manual metaalarms | UI | yes |
| activate | Activate an alarm | Engine-axe | yes |
| run_delayed_scenario | Run a delayed action | Engine-action | yes |
| instructionstarted, instructionpaused, instructionresumed, instructioncompleted, instructionfailed, instructionaborted, autoinstructionstarted, autoinstructioncompleted, autoinstructionfailed, autoinstructionalreadyrunning | Add instruction execution statuses to alarm steps | Engine-remediation/API | yes |
| recomputeentityservice | Recompute the service context graph and state | API | yes |
| updateentityservice | Update the service cache in engines | API | yes |
| entityupdated | Notify engines that an entity is updated | Engine-che/API | yes |
| entitytoggled | Notify engines that an entity is enabled/disabled | API | yes |
| alarmskipped | Check if an alarm was skipped during service recompute | Engine-service | yes |
| junittestsuiteupdated | Notify that the test suite is updated | connector-junit | yes |
| junittestcaseupdated | Notify that the test case is updated | connector-junit | yes |
| noevents | Create an alarm for an entity by idle rule | Engine-axe | yes |

**Important Note:** Internal events are mainly for internal system communication. Sending them manually may disrupt Canopsis' consistent state!

## Public events format
#### Mandatory fields
Event is represented by a json object, which may be passed to the canopsis by API call or RabbitMQ message.

All events should have the following fields:

* `event_type` - Event type (see below) - should be `string`.
* `source_type` - Source of event ("component", or "resource") - should be `string`.
* `connector` - Connector Type (gelf, nagios, snmp, ...) - should be `string`.
* `connector_name` - Connector Identifier (nagios1, nagios2, ...) - should be `string`.
* `component` - Component's name - should be `string`.
* `resource` - Resource's name (only if source_type is "resource") - should be `string`.
* `timestamp` - UNIX timestamp for when the event  was emitted (optional: set by the server to now) - should be `integer`.

Some information about `timestamp` field:
Though it's an optional field, it's a good practice to set it every time where it's possible, otherwise it may cause some alarm history inconsistency.  
The `engine-fifo` checks if timestamp is in valid range, which is **`now() ± 24h`**. If it's not, then it sets it to `now()`.

### Check event

Additional fields:
* `state` -  Alarm's state - should be (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL). Mandatory.
* `output` - Alarm's message - should be `string`. Optional.
* `long_output` - Alarm's description - should be `string`. Optional.
* `close_delay` - timeout in seconds to resolve open alarm

**!Note:** though `output` and `long_output` is always truncated by values, which are set by `OutputLength` and `LongOutputLength` parameters in **canopsis alarm config**.

Example:
```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "check",
  "component": "test_component",
  "resource": "test_resource",
  "state": 1,
  "output": "test output",
  "long_output": "test long output",
  "timestamp": 1639403732
}
```

Also it may have various additional extra fields, which are used for an entity enrichment feature.

### Contextupdate event

It's intended to trigger event filter rules and provide enrichment data for entity. Same fields as in `check` type event can be used considering that event's processing finished in Engine-che. Anything that impacts respective alarm is useless, for example: "state", "close_delay".

Example:
```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "contextupdate",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "long_output": "test long output",
  "timestamp": 1639403732
}
```

Also it may have various additional extra fields, which are used for an entity enrichment feature.

### Ack event

Additional fields:
* `author` - Ack step author - should be `string`. Optional.
* `user_id` - Ack step user - should be `string`. Optional.
* `output` - Ack step message - should be `string`. Optional.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "ack",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "timestamp": 1639403732
}
```

### Ackremove event

Additional fields:
* `author` - Ackremove step author - should be `string`. Optional.
* `user_id` - Ackremove step user - should be `string`. Optional.
* `output` - Ackremove step message - should be `string`. Optional.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "ackremove",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "timestamp": 1639403732
}
```

### Assocticket event

Additional fields:
* `author` - Assocticket step author - should be `string`. Optional.
* `user_id` - Assocticket step user - should be `string`. Optional.
* `ticket` - Ticket number - should be `string`. Optional.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "assocticket",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "ticket": "123",
  "timestamp": 1639403732
}
```

### Declareticket event

No additional fields

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "declareticket",
  "component": "test_component",
  "resource": "test_resource",
  "timestamp": 1639403732
}
```

### Cancel event

Additional fields:
* `author` - Cancel step author - should be `string`. Optional.
* `user_id` - Cancel step user - should be `string`. Optional.
* `output` - Cancel step message - should be `string`. Optional.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "cancel",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "timestamp": 1639403732
}
```

### Uncancel event

Additional fields:
* `author` - Uncancel step author - should be `string`. Optional.
* `user_id` - Uncancel step user - should be `string`. Optional.
* `output` - Uncancel step message - should be `string`. Optional.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "uncancel",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "timestamp": 1639403732
}
```

### Comment event

Additional fields:
* `author` - Comment step author - should be `string`. Optional.
* `user_id` - Comment step user - should be `string`. Optional.
* `output` - Comment step message - should be `string`. Optional.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "comment",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "timestamp": 1639403732
}
```

### Snooze event

Additional fields:
* `author` - Snooze step author - should be `string`. Optional.
* `output` - Snooze step message - should be `string`. Optional.
* `user_id` - Snooze step user - should be `string`. Optional.
* `duration` - Snooze duration in seconds - should be `integer`. Mandatory.

Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "snooze",
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "duration": 120,
  "timestamp": 1639403732
}
```

### Changestate event

Additional fields:
* `author` - Changestate step author - should be `string`. Optional.
* `output` - Changestate step message - should be `string`. Optional.
* `user_id` - Changestate step user - should be `string`. Optional.
* `state` -  Alarm's state - should be (0 - INFO, 1 - MINOR, 2 - MAJOR, 3 - CRITICAL). Mandatory.


Example:

```json
{
  "connector": "test_connector",
  "connector_name": "test_connectorname",
  "source_type": "resource",
  "event_type": "changestate",
  "state": 3,
  "component": "test_component",
  "resource": "test_resource",
  "output": "test output",
  "author": "test author",
  "timestamp": 1639403732
}
```

### Fields author and user_id

If `author` or `user_id` are not set, the system sets `author` as authorized user's name and `user_id` as authorized user's id by default.

