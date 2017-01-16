.. _FR__Alarm:

================
Alarm management
================

This document describes the alarm management feature of Canopsis.

.. contents::
   :depth: 3

References
==========

List of referenced functional requirements:

 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Context <FR__Context>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jean-Baptiste Braun", "2016/12/13", "0.3", "Add steps hard limit", ""
   "Jean-Baptiste Braun", "2016/12/13", "0.2", "Add flapping crop feature", ""
   "David Delassus", "2015/10/22", "0.1", "Document creation", ""

Contents
========

.. _FR__Alarm__Desc:

Description
-----------

An alarm is described by a set of :ref:`events <FR__Event>`:

 - a :ref:`check event <FR__Event__Check>` with a non-OK state, starting the alarm
 - a check event, with an OK state, ending the alarm
 - eventually one or more check events with distinct non-OK state
 - eventually one or more :ref:`acknowledgment events <FR__Event__Ack>`
 - eventually a :ref:`ticket declaration event <FR__Event__Declareticket>` or a :ref:`ticket association event <FR__Event__Assocticket>`
 - eventually a :ref:`changing state event <FR__Event__Changestate>`
 - eventually an :ref:`alarm canceling event <FR__Event__Cancel>`
 - eventually an :ref:`alarm restoring event <FR__Event__Uncancel>`
 - eventually an :ref:`alarm snooze event <FR__Event__Snooze>`

This set of events is called, in Canopsis, an *alarm cycle* and is associated to
a :ref:`contextual entity <FR__Context__Entity>`.

.. _FR__Alarm__Step:

Alarm step
----------

The *alarm cycle* may only be ended after a status set to 0 if the period of
potential flapping has elapsed. It **MUST** have a list of steps. Each
*alarm step*, **CAN** carry one information among those below:

 - the alarm is *on going*
 - the alarm is *flapping*
 - the alarm is *stealthy*
 - the alarm has been acknowledged
 - the alarm has been associated to a ticket
 - a ticket has been declared for the alarm
 - the alarm has been canceled
 - the alarm has been restored from its canceled state
 - the alarm state escalated
 - the alarm state decreased
 - the alarm has been snoozed
 - the number of steps has reached an hard limit

Each step **MUST** be historized in its corresponding *alarm cycle*.
And once the alarm ended, the cycle **MUST** be closed, and archived.

.. _FR__Alarm__Step__Supervised:

Supervised step
~~~~~~~~~~~~~~~

A *supervised step* is an alarm step emitted by the incoming Canopsis events.

.. _FR__Alarm__Step__User:

User step
~~~~~~~~~

A *user step* is an alarm step emitted by a Canopsis user. It **MUST** provides
the author of this step change, and a message explaining why the change has been
made.

Flapping steps
~~~~~~~~~~~~~~

If an alarm is flapping and never stop being in this status, a large amount of
state escalations and state decreases will be recorded. To prevent this from
happening and alarms to be too voluminous, some state changes **CAN** be
*cropped*.

When a change of state is cropped, 2 step counters are incremented :

 - a counter that counts the number of state escalations or decreases
 - a counter that counts the number of state values : off (0), minor (1), major
   (2), critical (3)

Those kind of steps **MUST** be inserted in any order just after the last
change of status or incremented when flapping steps are cropped. There are 6
different counter types.

State changes are cropped only if the number of state changes since the last
change of status is superior to ``flapping_persistant_steps`` property
(configurable).

Hard limit
~~~~~~~~~~

`Hard limit` is a watchdog feature preventing an alarm to grow too big, in case
of `flapping steps cropping` not being enough. It simply controls the maximum
number of steps that an alarm can have.

If a step should be appended whereas the has limit has been reached, it
**MUST** be discarded and can *not* be retrieved.

The only step allowed to be considered when an alarm has reached its hard limit
is an alarm cancelation.

`Hard limit` number of steps to be kept is configurable. This value **CAN** be
updated at any moment, and alarms that had been frozen **MUST** continue to
record steps if this limit has been extended.

Functional test
---------------

Alarms cycle **MUST** respect the following schema:

.. figure:: ../_static/images/alarm/cycle.svg

Status **MUST** be computed in respect of the following schema:

.. figure:: ../_static/images/alarm/status.svg
