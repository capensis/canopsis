.. _TR__Alarm:

================
Alarm management
================

This document specifies the alarm management in Canopsis, and its
implementation.

.. contents::
   :depth: 3

References
==========

 - :ref:`FR::Alarm <FR__Alarm>`
 - :ref:`FR::Configuration <FR__Configurable>`
 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Storage <FR__Storage>`
 - :ref:`TR::Storage <TR__Storage>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Romain Hennuyer", "2017/04/04", "0.5", "Updating change state description", ""
   "Romain Hennuyer", "2017/03/28", "0.4", "Add steps comment", ""
   "Jean-Baptiste Braun", "2016/12/13", "0.3", "Add steps hard limit", ""
   "Jean-Baptiste Braun", "2016/12/13", "0.2", "Add flapping crop feature", ""
   "David Delassus", "2015/10/26", "0.1", "Document creation", ""

Contents
========

.. _TR__Alarm__Manager:

Alerts manager
--------------

An ``Alerts`` :ref:`configurable <FR__Configurable>` provides:

 - the ability to archive an :ref:`event <FR__Event>`
 - alarm cycle management operations:
    - create a new one
    - update existing one
    - get last one
    - find old alarms
    - tagging

.. _TR__Alarm__Cycle:

Alarm Cycle
-----------

Alarm cycles are persisted in a :ref:`timed storage <FR__Storage_Type>` with
the following informations:

 - data identifier: the entity identifier of the received event
 - value: set of alarm steps (see :ref:`data model <TR__Alarm__DataModel>` below)
 - timestamp: date/time of alarm appearance

.. _TR__Alarm__DataModel:

Alarm data model
----------------

The set of alarm steps will compute informations for an easier use:

.. csv-table::
   :header: Field, Description

   state, ``stateinc`` step or ``statedec`` step or ``changestate`` step or ``None``
   status, ``statusinc`` step or ``statusdec`` step or ``None``
   ack, ``ack`` step or ``None`` if unacknowledged
   canceled, ``cancel`` step or ``None`` if uncanceled
   ticket, ``declareticket`` step or ``assocticket`` step or ``None`` if no ticketing informations
   resolved, timestamp of resolution or ``None`` if alarm is still on going
   steps, array of steps (see data models bellow)
   tags, array of tag (as strings) used for filtering

.. _TR__Alarm__DataModel__StateInc:

Alarm step "state increase" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``stateinc``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new state

.. _TR__Alarm__DataModel__StateDec:

Alarm step "state decrease" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statedec``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new state

.. _TR__Alarm__DataModel__StatusInc:

Alarm step "status increase" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statusinc``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new status

.. _TR__Alarm__DataModel__StatusDec:

Alarm step "status decrease" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statusdec``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``{connector}.{connector_name}``
   m, step message, ``event["output"]``
   val, step associated value, new status

**NB:** if status decreases to ``OFF``, then the alarm value ``resolved`` is
set to this step timestamp.

.. _TR__Alarm__DataModel__Acknowledge:

Alarm step "acknowledge" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``ack``
   t, step timestamp, ``event["timestamp"]``
   a, step author, acknowledgment author
   m, step message, acknowledgment message

.. _TR__Alarm__DataModel__Unacknowledge:

Alarm step "unacknowledge" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``ackremove``
   t, step timestamp, ``event["timestamp"]``
   a, step author, acknowledgment removal author
   m, step message, acknowledgment removal message

**NB:** this step reset the alarm value ``ack`` to ``None``.

.. _TR__Alarm__DataModel__Cancel:

Alarm step "cancel" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``cancel``
   t, step timestamp, ``event["timestamp"]``
   a, step author, alarm canceling author
   m, step message, alarm canceling message

.. _TR__Alarm__DataModel__Comment:

Alarm step "comment" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``comment``
   t, step timestamp, ``event["timestamp"]``
   a, step author, comment author
   m, step message, comment message

.. _TR__Alarm__DataModel__Restore:

Alarm step "restore" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``uncancel``
   t, step timestamp, ``event["timestamp"]``
   a, step author, alarm restoring author
   m, step message, alarm restoring message

**NB:** this step reset the alarm value ``cancel`` to ``None``.

.. _TR__Alarm__DataModel__DeclareTicket:

Alarm step "declare ticket" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``declareticket``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ticket declaration author
   m, step message, ticket declaration message
   val, ticket number, ``None``

.. _TR__Alarm__DataModel__AssocTicket:

Alarm step "associate ticket" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``assocticket``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ticket association author
   m, step message, ticket association message
   val, ticket number, ``event["ticket"]``

.. _TR__Alarm__DataModel__ChangeState:

Alarm step "change state" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``changestate``
   t, step timestamp, ``event["timestamp"]``
   a, step author, state requalification author
   m, step message, state requalification message
   val, state, new state

Alarm step "snooze" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``snooze``
   t, step timestamp, ``event["timestamp"]``
   a, step author, ``event["author"]``
   m, step message, ``Auto snooze generated by rule <event_filter_rule>``
   val, timestamp, date until end of snooze

Alarm step "statecounter" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``statedeccounter``
   t, step timestamp, timestamp of last status change
   a, step author, author of last status change
   m, step message, ````
   val, dict, keys are ``statedec``, ``stateinc``, ``state:0``, ``state:1``, ``state:2`` and ``state:3`` and values are counters

Alarm step "hardlimit" data model
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. csv-table::
   :header: Field, Description, Default Value

   _t, step type, ``hardlimit``
   t, step timestamp, hard limit timer
   a, step author, ``__canopsis__``
   m, step message, ``This alarm has reached an hard limit (<limit> steps recorded) : no more steps will be appended. Please cancel this alarm or extend hard limit value.``
   val, integer, limit at the time it has been reached

Unit Tests
==========

Get alarm history
-----------------

``get_alarms([resolved], [tags], [exclude_tags], [timewindow]) -> alarms``:

 * ``resolved`` (optional) as a ``boolean``: get resolved alarms or unresolved alarms
 * ``tags`` (optional) as a ``string`` or a ``list``: get alarms with listed tags
 * ``exclude_tags`` (optional) as a ``string`` or a ``list``: get alarms without listed tags
 * ``timewindow`` (optional) as a ``canopsis.timeserie.timewindow.TimeWindow``: get alarms within time interval
 * ``alarms`` as a ``cursor``: alarms that matched the previous parameters

Creating new alarm
------------------

``make_alarm(alarm_id, event)``:

 * ``alarm_id`` as ``string``: the entity id of the alarm
 * ``event`` as ``dict``: the event which produces the alarm

Case 1: new alarm
~~~~~~~~~~~~~~~~~

**Expected:** The alarm **MUST** be present in the configured timed storage
with all values set to ``None``.

Case 2: existing alarm
~~~~~~~~~~~~~~~~~~~~~~

**Expected:** The existing alarm **MUST** be left untouched, and no new alarm should be created.

Get last unresolved alarm
-------------------------

``get_current_alarm(alarm_id) -> alarm``:

 * ``alarm_id`` as ``string``: the entity id of the alarm
 * ``alarm`` as a ``dict``: the current unresolved alarm, or ``None`` if no alarm found, or all of them are resolved

Case 1: there is an unresolved alarm
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

**Expected:** ``alarm`` **MUST NOT** be ``None``, and should contains a value described by the :ref:`alarm data model <TR__Alarm__DataModel>`.

Case 2: there is no alarm, or no resolved ones
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

**Expected:** ``alarm`` **MUST** be ``None``.

Update existing alarm
---------------------

``update_current_alarm(alarm, new_value, [tags])``:

 * ``alarm`` as described by the :ref:`timed storage data model <TR__Storage__DataModel__Timed>`: alarm to update
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: value to use for the alarm
 * ``tags`` (optional) as a ``list`` or a ``string``: tags to add to the alarm value

Case 1: there is no alarm
~~~~~~~~~~~~~~~~~~~~~~~~~

**Expected:** A new document **SHOULD** be created.

Case 2: there is an existing alarm
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

**Expected:**

 - the alarm value **MUST** be replaced by ``new_value``
 - the ``tags`` **MUST** be added to the alarm value

Task: acknowledge
-----------------

``alerts.useraction.ack(manager, alarm, author, message, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to acknowledge
 * ``author`` as a ``string``: the acknowledgment author
 * ``message`` as a ``string``: the acknowledgment message
 * ``event`` as a ``dict``: the :ref:`acknowledgment event <FR__Event__Ack>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``ack`` **MUST** be set to the :ref:`acknowledge step <TR__Alarm__DataModel__Acknowledge>`
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: unacknowledge
-------------------

``alerts.useraction.ackremove(manager, alarm, author, message, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to unacknowledge
 * ``author`` as a ``string``: the acknowledgment removing author
 * ``message`` as a ``string``: the acknowledgment removing message
 * ``event`` as a ``dict``: the :ref:`acknowledgment removing event <FR__Event__Ackremove>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``ack`` **MUST** be set to ``None``
 - the :ref:`unacknowledge step <TR__Alarm__DataModel__Unacknowledge>` **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: Cancel
------------

``alerts.useraction.cancel(manager, alarm, author, message, event) -> new_value, status``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to cancel
 * ``author`` as a ``string``: the alarm canceling author
 * ``message`` as a ``string``: the alarm canceling message
 * ``event`` as a ``dict``: the :ref:`alarm canceling event <FR__Event__Cancel>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value
 * ``status`` as an ``int`` which will always be set to ``CANCELED`` (will trigger a change of status on the alarm)

**Expected:**

 - the alarm ``cancel`` **MUST** be set to :ref:`cancel step <TR__Alarm__DataModel__Cancel>`
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: comment
-----------------

``alerts.useraction.comment(manager, alarm, author, message, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to comment
 * ``author`` as a ``string``: the comment author
 * ``message`` as a ``string``: the comment message
 * ``event`` as a ``dict``: the :ref:`comment event <FR__Event__Comment>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``comment`` **MUST** be set to :ref:`comment step <TR__Alarm__DataModel__Comment>`
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: Restore
-------------

``alerts.useraction.uncancel(manager, alarm, author, message, event) -> new_value, status``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to restore
 * ``author`` as a ``string``: the alarm restoring author
 * ``message`` as a ``string``: the alarm restoring message
 * ``event`` as a ``dict``: the :ref:`alarm restoring event <FR__Event__Uncancel>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value
 * ``status`` as an ``int`` which will be set to the previous status or the actual status as it should have been without the *cancel* (will trigger a change of status on the alarm)

**Expected:**

 - the alarm ``cancel`` **MUST** be set to ``None``
 - the :ref:`cancel step <TR__Alarm__DataModel__Cancel>` **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: Declare ticket
--------------------

``alerts.useraction.declareticket(manager, alarm, author, message, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm used for ticket declaration
 * ``author`` as a ``string``: the ticket declaration author
 * ``message`` as a ``string``: the ticket declaration message
 * ``event`` as a ``dict``: the :ref:`ticket declaration event <FR__Event__Declareticket>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``ticket`` **MUST** be set to the :ref:`ticket declaration step <TR__Alarm__DataModel__Declareticket>`
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: Associate ticket
----------------------

``alerts.useraction.assocticket(manager, alarm, author, message, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm used for ticket association
 * ``author`` as a ``string``: the ticket association author
 * ``message`` as a ``string``: the ticket association message
 * ``event`` as a ``dict``: the :ref:`ticket association event <FR__Event__Assocticket>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``ticket`` **MUST** be set to the :ref:`ticket association step <TR__Alarm__DataModel__Assocticket>`
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``

Task: Change State
------------------

``alerts.useraction.changestate(manager, alarm, author, message, event) -> new_value`` (as same as ``alerts.useraction.keepstate``):

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to change
 * ``author`` as a ``string``: the change state author
 * ``message`` as a ``string``: the change state message
 * ``event`` as a ``dict``: the :ref:`change state event <FR__Event__Changestate>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``ticket`` **MUST** be set to the :ref:`change state step <TR__Alarm__DataModel__ChangeState>`
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm **MUST** be returned as ``new_value``
 - the alarm **MUST** be recognized as an unchangable state
 - the alarm **MUST** always update his state if the state is OK (0)

Task: State increase
--------------------

``alerts.systemaction.state_increase(manager, alarm, state, event) -> new_value, status``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to change
 * ``state`` as ``int``: the increased state
 * ``event`` as a ``dict``: the :ref:`check event <FR__Event__Check>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value
 * ``status`` as ``int``: the new computed status from state history

**Expected:**

 - the alarm ``state`` **MUST** be set to the :ref:`state increase step <TR__Alarm__DataModel__StateInc>` **only if** there was no :ref:`change state step <TR__Alarm__DataModel__ChangeState>` set
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm ``status`` **MUST** be computed accordingly to the functional tests

Task: State decrease
--------------------

``alerts.systemaction.state_decrease(manager, alarm, state, event) -> new_value, status``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to change
 * ``state`` as ``int``: the decreased state
 * ``event`` as a ``dict``: the :ref:`check event <FR__Event__Check>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value
 * ``status`` as ``int``: the new computed status from state history

**Expected:**

 - the alarm ``state`` **MUST** be set to the :ref:`state decrease step <TR__Alarm__DataModel__StateDec>` **only if** there was no :ref:`change state step <TR__Alarm__DataModel__ChangeState>` set
 - the step **MUST** be added to the ``steps`` set of the alarm
 - the alarm ``status`` **MUST** be computed accordingly to the functional tests

Task: Status increase
---------------------

``alerts.systemaction.status_increase(manager, alarm, status, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to change
 * ``status`` as ``int``: the increased status
 * ``event`` as a ``dict``: the :ref:`check event <FR__Event__Check>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``status`` **MUST** be set to the :ref:`status increase step <TR__Alarm__DataModel__StatusInc>`
 - the step **MUST** be added to the ``steps`` set of the alarm

Task: Status decrease
----------------------

``alerts.systemaction.status_decrease(manager, alarm, status, event) -> new_value``:

 * ``manager`` as an ``Alerts`` configurable: the task caller
 * ``alarm`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the alarm to change
 * ``status`` as ``int``: the decreased status
 * ``event`` as a ``dict``: the :ref:`check event <FR__Event__Check>`
 * ``new_value`` as described by the :ref:`alarm data model <TR__Alarm__DataModel>`: the new alarm value

**Expected:**

 - the alarm ``status`` **MUST** be set to the :ref:`status increase step <TR__Alarm__DataModel__StatusInc>`
 - the step **MUST** be added to the ``steps`` set of the alarm

Utility: Get previous step
--------------------------

``get_previous_step(alarm, steptypes, [ts]) -> step``:

Utility: Get last state
-----------------------

``get_last_state(alarm, [ts]) -> state``:

Utility: Get last status
------------------------

``get_last_status(alarm, [ts]) -> status``:

Utility: Is flapping ?
----------------------

``is_flapping(manager, alarm) -> result``:

Utility: Is stealthy ?
----------------------

``is_stealthy(manager, alarm) -> result``:

Utility: Is keeped state ?
----------------------

``is_keeped_state(alarm) -> result``:

