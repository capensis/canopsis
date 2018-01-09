# Events

This document describes the structure of events in Canopsis.

## Contents

### Description

An event in Canopsis is the representation of asynchronously incoming
data, sent by a connector.

They are described by a **CEvent** data schema per type of event.

Each event contains typed informations (listed bellow), associated to
one or more entities.


### Logging events

Those events are used to log informations to Canopsis.

#### Comment

A `comment` event **MUST** contain:

- an author
- a message
- a criticality level
- a reference to the commented user event

And it **MAY** contain:

- a detailed message

#### Log

A `log` event **MUST** contain:

- a message
- a severity level

And it **MAY** contain:

- a detailed message

#### User

A `user` event **MUST** contain:

- an author
- a message
- a criticality level

And it **MAY** contain:

- a detailed message


### Supervising events

Those events are used to store changes in a supervision environment to
Canopsis.

#### Cancel

A `cancel` event **MUST** contain:

- an author
- a message
- a reference to the check event to cancel

#### ChangeState

A `changestate` event **MUST** contain:

- an author
- a state
- a message
- a reference to the check event to modify

#### Check

A `check` event **MUST** contain:

- a state
- a message

And it **MAY** contain:

- a state specification
- a detailed message

#### Downtime

A `downtime` event **MUST** contain:

- an author
- a message
- a period

#### Selector

A `selector` event **MUST** contain:

- a state
- a message

And it **MAY** contain:

- a displayed name
- a state specification
- a detailed message

#### Snooze

A `snooze` event **MUST** contain:

- an author
- a message

And it **MAY** contain:

- a period

#### Trap

A `trap` event **MUST** contain:

- a state
- a severity level
- an OID

#### Uncancel

An `uncancel` event **MUST** contain:

- an author
- a message
- a reference to the canceled check event


### Ticketing events

Those events are used to represent interactions with a CMDB.

#### Assocticket

A `declareticket` event **MUST** contain:

- an author
- a message
- a ticket ID
- a reference to the check event to assign the ticket to

#### Declareticket

A `declareticket` event **MUST** contain:

- an author
- a message
- a reference to the check event to create a ticket for


### Acknowledging events

Those events are used to manage supervising events.

#### Acknowledgment

An `ack` event **MUST** contain:

- an author
- a message
- a reference to the check event to acknowledge

#### Ackremove

An `ackremove` event **MUST** contain:

- an author
- a message
- a reference to the check event to *unacknowledge*


### Performance events

Those events are used to store metrics in Canopsis.

#### perf

A `perf` event **MUST** contain:

- one or more metric

**NB:** `perf` events are not stored and can be included in every other events.
