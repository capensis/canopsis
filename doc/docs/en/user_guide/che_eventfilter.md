# Event filter

This page documents the event filter implemented in the go engine `che`. The
python engine `event-filter` is documented
[here](../../en/user_guide/event_filter.md).


## Rules

A rule is a JSON document containing the following parameters:

 - `type` (required): the type of rule (`enrichment`, `drop`, `break`).
 - `pattern` (optional): the rule will be applied to the events matched by
   this pattern. If the pattern is not specified, the rule will be applied to
   all the events.
 - `priority` (optional, 0 by default): the priority of the rule (the rules
   are executed in ascending order of priority).
 - `enabled` (optional, `true` by default): `false` to disable the rule.

For example, the following rule drops the events whose resource is
`invalid_resource`:

```json
{
    "type": "drop",
    "pattern": {
        "resource": "invalid_resource"
    },
    "priority": 10
}
```


### Execution of the rules

When an event is received by the `che` engine, the rules are executed in
ascending order of priority. If the event is matched by the pattern of an
enabled rule, this rule is applied.

 - If the outcome of the rule is `pass`, the event passes to the next rule.
 - If the outcome of the rule is `break`, the event breaks out of the event
   filter, and the remaining rules are not applied.
 - If the outcome of the rule is `drop`, the event is deleted. The deleted
   event is logged.

The outcome of a rule of type `break` is always `break`, and the outcome of a
rule of type `drop` is always `drop`.

If the event is invalid at the end of the execution of the rules, it is
deleted. These events are logged. An event is valid if:

 - its `source_type` field is `component`, and its `component` field is defined;
   *or*
 - ts `source_type` field is `resource`, and its `component` and `resource`
   fields are defined.

If the `debug` field of an event is `true`, the processing of the event by the
event filter is traced. This field can be set with an enrichment rule.

### Patterns

The pattern of a rule selects the events to which the rule is applied. A
pattern is a dictionary containing the values of some of the fields of an
event. For example:

```json
"pattern": {
    "component": "component_name",
    "resource": "resource_name"
}
```

A rule containing this pattern is applied to the events whose component is
`component_name` and whose resource is `resource_name`.

Alternatively, a dictionnary containing `operator: value` pairs can be provided
instead of a value. The available operators are:

 - `>=`, `>`, `<`, `<=`: compare the value of a field to a numerical value.
 - `regex_match`: match the value of a field with a regular expression.

For example, the following pattern matches the events whose state is between 1
and 3 and whose output is matched by a regular expression.

```json
"pattern": {
    "state": {">=": 1, "<=": 3},
    "output": {"regex_match": "Warning: CPU Load is critical (.*)"}
}
```

If a rule is applied after the enrichment with the entity corresponding to the
event, the entity corresponding to the event is available in the
`current_entity` field. It is the possible to define a dictionnary in
`current_entity` to filter on the entity's fields. For example, the following
pattern selects events whose entity is enabled and does not have a
`service_description` information:

```json
"pattern": {
    "current_entity": {
        "enabled": true,
        "infos": {
            "service_description": null
        }
    }
}
```


## Examples

### Dropping events

The following rule drops the event whose resource is `invalid_resource`:

```json
{
    "type": "drop",
    "pattern": {
        "resource": "invalid_resource"
    },
    "priority": 10
}
```

The following rule drops the events whose state is major or critical on
resources whose name starts with "cpu-":

```json
{
    "type": "drop",
    "pattern": {
        "state": {">=": 2},
        "resource": {"regex_match": "cpu-.*"}
    },
    "priority": 10
}
```

### Breaking out of the event filter

The following rule breaks the events of type `pbehavior` out of the event
filter:

```json
{
    "type": "break",
    "pattern": {
        "event_type": "pbehavior"
    },
    "priority": 0
}
```
