Alarm management
================

This document specifies the alarm management in Canopsis, and its
implementation.

References
----------

> -   FR::Alarm &lt;FR\_\_Alarm&gt;
> -   FR::Configuration &lt;FR\_\_Configurable&gt;
> -   FR::Event &lt;FR\_\_Event&gt;
> -   FR::Storage &lt;FR\_\_Storage&gt;
> -   TR::Storage &lt;TR\_\_Storage&gt;

Updates
-------

Contents
--------

### Alerts manager

An `Alerts` configurable &lt;FR\_\_Configurable&gt; provides:

> -   the ability to archive an event &lt;FR\_\_Event&gt;
> -   
>
>     alarm cycle management operations:
>
>     :   -   create a new one
>         -   update existing one
>         -   get last one
>         -   find old alarms
>         -   tagging
>
### Alarm Cycle

Alarm cycles are persisted in a
timed storage &lt;FR\_\_Storage\_Type&gt; with the following
informations:

> -   data identifier: the entity identifier of the received event
> -   value: set of alarm steps (see
>     data model &lt;TR\_\_Alarm\_\_DataModel&gt; below)
> -   timestamp: date/time of alarm appearance

### Alarm data model

The set of alarm steps will compute informations for an easier use:

#### Alarm step "state increase" data model

#### Alarm step "state decrease" data model

#### Alarm step "status increase" data model

#### Alarm step "status decrease" data model

**NB:** if status decreases to `OFF`, then the alarm value `resolved` is
set to this step timestamp.

#### Alarm step "acknowledge" data model

#### Alarm step "unacknowledge" data model

**NB:** this step reset the alarm value `ack` to `None`.

#### Alarm step "cancel" data model

#### Alarm step "comment" data model

#### Alarm step "restore" data model

**NB:** this step reset the alarm value `cancel` to `None`.

#### Alarm step "declare ticket" data model

#### Alarm step "associate ticket" data model

#### Alarm step "change state" data model

#### Alarm step "snooze" data model

#### Alarm step "statecounter" data model

#### Alarm step "hardlimit" data model

Unit Tests
----------

### Get alarm history

`get_alarms([resolved], [tags], [exclude_tags], [timewindow]) -> alarms`:

> -   `resolved` (optional) as a `boolean`: get resolved alarms or
>     unresolved alarms
> -   `tags` (optional) as a `string` or a `list`: get alarms with
>     listed tags
> -   `exclude_tags` (optional) as a `string` or a `list`: get alarms
>     without listed tags
> -   `timewindow` (optional) as a
>     `canopsis.timeserie.timewindow.TimeWindow`: get alarms within time
>     interval
> -   `alarms` as a `cursor`: alarms that matched the previous
>     parameters

### Creating new alarm

`make_alarm(alarm_id, event)`:

> -   `alarm_id` as `string`: the entity id of the alarm
> -   `event` as `dict`: the event which produces the alarm

#### Case 1: new alarm

**Expected:** The alarm **MUST** be present in the configured timed
storage with all values set to `None`.

#### Case 2: existing alarm

**Expected:** The existing alarm **MUST** be left untouched, and no new
alarm should be created.

### Get last unresolved alarm

`get_current_alarm(alarm_id) -> alarm`:

> -   `alarm_id` as `string`: the entity id of the alarm
> -   `alarm` as a `dict`: the current unresolved alarm, or `None` if no
>     alarm found, or all of them are resolved

#### Case 1: there is an unresolved alarm

**Expected:** `alarm` **MUST NOT** be `None`, and should contains a
value described by the
alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;.

#### Case 2: there is no alarm, or no resolved ones

**Expected:** `alarm` **MUST** be `None`.

### Update existing alarm

`update_current_alarm(alarm, new_value, [tags])`:

> -   `alarm` as described by the
>     timed storage data model &lt;TR\_\_Storage\_\_DataModel\_\_Timed&gt;:
>     alarm to update
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: value to use
>     for the alarm
> -   `tags` (optional) as a `list` or a `string`: tags to add to the
>     alarm value

#### Case 1: there is no alarm

**Expected:** A new document **SHOULD** be created.

#### Case 2: there is an existing alarm

**Expected:**

> -   the alarm value **MUST** be replaced by `new_value`
> -   the `tags` **MUST** be added to the alarm value

### Task: acknowledge

`alerts.useraction.ack(manager, alarm, author, message, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     acknowledge
> -   `author` as a `string`: the acknowledgment author
> -   `message` as a `string`: the acknowledgment message
> -   `event` as a `dict`: the
>     acknowledgment event &lt;FR\_\_Event\_\_Ack&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `ack` **MUST** be set to the
>     acknowledge step &lt;TR\_\_Alarm\_\_DataModel\_\_Acknowledge&gt;
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: unacknowledge

`alerts.useraction.ackremove(manager, alarm, author, message, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     unacknowledge
> -   `author` as a `string`: the acknowledgment removing author
> -   `message` as a `string`: the acknowledgment removing message
> -   `event` as a `dict`: the
>     acknowledgment removing event &lt;FR\_\_Event\_\_Ackremove&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `ack` **MUST** be set to `None`
> -   the
>     unacknowledge step &lt;TR\_\_Alarm\_\_DataModel\_\_Unacknowledge&gt;
>     **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: Cancel

`alerts.useraction.cancel(manager, alarm, author, message, event) -> new_value, status`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     cancel
> -   `author` as a `string`: the alarm canceling author
> -   `message` as a `string`: the alarm canceling message
> -   `event` as a `dict`: the
>     alarm canceling event &lt;FR\_\_Event\_\_Cancel&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value
> -   `status` as an `int` which will always be set to `CANCELED` (will
>     trigger a change of status on the alarm)

**Expected:**

> -   the alarm `cancel` **MUST** be set to
>     cancel step &lt;TR\_\_Alarm\_\_DataModel\_\_Cancel&gt;
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: comment

`alerts.useraction.comment(manager, alarm, author, message, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     comment
> -   `author` as a `string`: the comment author
> -   `message` as a `string`: the comment message
> -   `event` as a `dict`: the
>     comment event &lt;FR\_\_Event\_\_Comment&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `comment` **MUST** be set to
>     comment step &lt;TR\_\_Alarm\_\_DataModel\_\_Comment&gt;
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: Restore

`alerts.useraction.uncancel(manager, alarm, author, message, event) -> new_value, status`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     restore
> -   `author` as a `string`: the alarm restoring author
> -   `message` as a `string`: the alarm restoring message
> -   `event` as a `dict`: the
>     alarm restoring event &lt;FR\_\_Event\_\_Uncancel&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value
> -   `status` as an `int` which will be set to the previous status or
>     the actual status as it should have been without the *cancel*
>     (will trigger a change of status on the alarm)

**Expected:**

> -   the alarm `cancel` **MUST** be set to `None`
> -   the cancel step &lt;TR\_\_Alarm\_\_DataModel\_\_Cancel&gt;
>     **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: Declare ticket

`alerts.useraction.declareticket(manager, alarm, author, message, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm used
>     for ticket declaration
> -   `author` as a `string`: the ticket declaration author
> -   `message` as a `string`: the ticket declaration message
> -   `event` as a `dict`: the
>     ticket declaration event &lt;FR\_\_Event\_\_Declareticket&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `ticket` **MUST** be set to the
>     ticket declaration step &lt;TR\_\_Alarm\_\_DataModel\_\_Declareticket&gt;
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: Associate ticket

`alerts.useraction.assocticket(manager, alarm, author, message, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm used
>     for ticket association
> -   `author` as a `string`: the ticket association author
> -   `message` as a `string`: the ticket association message
> -   `event` as a `dict`: the
>     ticket association event &lt;FR\_\_Event\_\_Assocticket&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `ticket` **MUST** be set to the
>     ticket association step &lt;TR\_\_Alarm\_\_DataModel\_\_Assocticket&gt;
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`

### Task: Change State

`alerts.useraction.changestate(manager, alarm, author, message, event) -> new_value`
(as same as `alerts.useraction.keepstate`):

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     change
> -   `author` as a `string`: the change state author
> -   `message` as a `string`: the change state message
> -   `event` as a `dict`: the
>     change state event &lt;FR\_\_Event\_\_Changestate&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `ticket` **MUST** be set to the
>     change state step &lt;TR\_\_Alarm\_\_DataModel\_\_ChangeState&gt;
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm **MUST** be returned as `new_value`
> -   the alarm **MUST** be recognized as an unchangable state
> -   the alarm **MUST** always update his state if the state is OK (0)

### Task: State increase

`alerts.systemaction.state_increase(manager, alarm, state, event) -> new_value, status`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     change
> -   `state` as `int`: the increased state
> -   `event` as a `dict`: the check event &lt;FR\_\_Event\_\_Check&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value
> -   `status` as `int`: the new computed status from state history

**Expected:**

> -   the alarm `state` **MUST** be set to the
>     state increase step &lt;TR\_\_Alarm\_\_DataModel\_\_StateInc&gt;
>     **only if** there was no
>     change state step &lt;TR\_\_Alarm\_\_DataModel\_\_ChangeState&gt;
>     set
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm `status` **MUST** be computed accordingly to the
>     functional tests

### Task: State decrease

`alerts.systemaction.state_decrease(manager, alarm, state, event) -> new_value, status`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     change
> -   `state` as `int`: the decreased state
> -   `event` as a `dict`: the check event &lt;FR\_\_Event\_\_Check&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value
> -   `status` as `int`: the new computed status from state history

**Expected:**

> -   the alarm `state` **MUST** be set to the
>     state decrease step &lt;TR\_\_Alarm\_\_DataModel\_\_StateDec&gt;
>     **only if** there was no
>     change state step &lt;TR\_\_Alarm\_\_DataModel\_\_ChangeState&gt;
>     set
> -   the step **MUST** be added to the `steps` set of the alarm
> -   the alarm `status` **MUST** be computed accordingly to the
>     functional tests

### Task: Status increase

`alerts.systemaction.status_increase(manager, alarm, status, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     change
> -   `status` as `int`: the increased status
> -   `event` as a `dict`: the check event &lt;FR\_\_Event\_\_Check&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `status` **MUST** be set to the
>     status increase step &lt;TR\_\_Alarm\_\_DataModel\_\_StatusInc&gt;
> -   the step **MUST** be added to the `steps` set of the alarm

### Task: Status decrease

`alerts.systemaction.status_decrease(manager, alarm, status, event) -> new_value`:

> -   `manager` as an `Alerts` configurable: the task caller
> -   `alarm` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the alarm to
>     change
> -   `status` as `int`: the decreased status
> -   `event` as a `dict`: the check event &lt;FR\_\_Event\_\_Check&gt;
> -   `new_value` as described by the
>     alarm data model &lt;TR\_\_Alarm\_\_DataModel&gt;: the new alarm
>     value

**Expected:**

> -   the alarm `status` **MUST** be set to the
>     status increase step &lt;TR\_\_Alarm\_\_DataModel\_\_StatusInc&gt;
> -   the step **MUST** be added to the `steps` set of the alarm

### Utility: Get previous step

`get_previous_step(alarm, steptypes, [ts]) -> step`:

### Utility: Get last state

`get_last_state(alarm, [ts]) -> state`:

### Utility: Get last status

`get_last_status(alarm, [ts]) -> status`:

### Utility: Is flapping ?

`is_flapping(manager, alarm) -> result`:

### Utility: Is stealthy ?

`is_stealthy(manager, alarm) -> result`:

Utility: Is keeped state ? ----------------------

`is_keeped_state(alarm) -> result`:
