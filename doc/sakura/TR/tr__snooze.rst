.. _TR__Snooze:

=============
Snoozed alarm
=============

This document describes the snooze behaviour of Canopsis, and its
implementation.

.. contents::
   :depth: 3

References
==========

 - :ref:`FR::Alarm <FR__Alarm>`
 - :ref:`FR::Snooze <FR__Snooze>`
 - :ref:`FR::Snooze event <FR__Event__Snooze>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jean-Baptiste Braun", "2016/09/16", "0.1", "Proposition of implementation", ""

Contents
========

Snooze behaviour will be implemented in ``*_alarm`` collections. It won't be
recorded in events / events_log.

Snooze information
------------------

Snoozes will be notified with a new type of event : ``snooze``.

Auto-triggered snoozes are assumed by the ``alert`` engine if this type
of snooze is configured and if the alarm is just created.

Storage
-------

In alarm steps, snooze will be a new type of step, recording both the snooze
start time and end time. A snooze property (default value is null) will be
added to alarms to record the last snooze step.

To retrieve non-snoozed alarms, a mongo filter filtering snooze=null or
snooze_end < now is applied by default.

If a snoozed alarm is snoozed again :

 - The last snooze will be appended to steps
 - The new snooze will overwrite ``snooze`` field

Consequences :

 - A snoozed alarm that is snoozed again will stop being snoozed at the end of
   the most recent snooze
 - If the second snooze has a duration of 0 seconds, this behaviour achieve an
   'un-snooze' function

Configuration
-------------

As a first step, two global configuration directives are configurable in
etc/schema.d/crecord.statusmanagement.json :

 - auto_snooze : Do alarms have to be automatically snoozed after creation ?
 - snooze_default_time : Default snooze duration, either used in context of
   automatic snoozes or if a snooze event is received wi no duration property.

When contextv2 will be available, those properties will be configurable by
(group of) entity instead.
