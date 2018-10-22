# Event filter API

The event filter API allows to manage the rules of the event filter of the go
engine `che`.

## List the rules

```
GET /api/v2/eventfilter/rules
```

Returns an array containing all the rules of the event filter.

```json
[
    {
        "_id": "6b90880a-c4f0-4a4d-8a51-de0c7e14581e",
        "type": "drop",
        "pattern": {...},
        "priority": 100,
    },
    ...
]
```

## Get a rule

```
GET /api/v2/eventfilter/rules/<rule_id>
```

Returns the rule whose id is `<rule_id>` or an error if there is no such
rule.

## Create a rule

```json
POST /api/v2/eventfilter/rules
{
    "type": "drop",
    "pattern": {...},
    "priority": 100,
}
```

Creates a rule and returns its id. An error is returned if the rule is
invalid.

## Delete a rule

```
DEL /api/v2/eventfilter/rules/<rule_id>
```

Deletes the rule whose id is `<rule_id>` or return an error if there is no
such rule.

## Edit a rule

```json
PUT /api/v2/eventfilter/rules/<rule_id>
{
    "type": "drop",
    "pattern": {...},
    "priority": 100,
}
```

Edits the rule whose id is `<rule_id>`, or returns an error if the rule is
invalid or if its id changed.
