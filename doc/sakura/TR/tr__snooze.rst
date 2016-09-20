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

Notification
------------

Snoozes will be notified with a new type of event : ``snooze``. This event will
be sent by the ``alert`` engine for auto-triggered snoozes.

Storage
-------

In alarm steps, snooze will be a new type of step, recording both the snooze
start time and end time. A snooze property (default value is null) will be
added to alarms to record the last snooze step.

To retrieve non-snoozed alarms, a mongo filter filtering snooze=null or
snooze_end < now could be applied.

If a snoozed alarm is snoozed again :

 - the second snooze will be appended to steps
 - the snooze property end time will be updated with the second snooze end time

Configuration
-------------

As a first step, a global snooze time for all entities will be configurable in
etc/schema.d/crecord.statusmanagement.json. When contextv2 will be available,
automatic snooze time by entity will be used instead.
