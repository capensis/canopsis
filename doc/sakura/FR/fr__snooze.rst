.. _FR__Snooze:

=============
Snoozed alarm
=============

This document describes the snooze behaviour feature of Canopsis.

.. contents::
   :depth: 3

References
==========

List of referenced functional requirements:

 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Context <FR__Context>`
 - :ref:`FR::Alarm <FR__Alarm>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Florent Demeulenaere", "2016/06/22", "0.2", "Add snoozed description", ""
   "Florent Demeulenaere", "2016/06/21", "0.1", "Document creation", ""

Contents
========

.. _FR__Snooze__Desc:

Description
-----------

This feature aims to add the snoozed behaviour to the alarm cycle. This
behaviour could be used in two different ways.

The snoozed behaviour at the alarm creation
-------------------------------------------

When an alarm is created because of a check KO, this alarm **CAN** be snoozed
if the check KO matches a filter.  It means that the alarm is not displayed
since its creation but only after a defined time period if the alarm is still
opened.

If the alarm state is closed before the end of the defined time period then the
alarm **MUST** be completely ignored for any reports.

This behavior **MUST** be optional and configurable by group of entities.

The snoozed behaviour triggered by the user
-------------------------------------------

When an alarm is displayed on the widget list, the user **MUST** be able to
snooze it.

If an alarm is snoozed by a user, the alarm is switched to a ``snooze`` state
(not visible by default) during a time period defined by the user.

During this time period, the alarm is not visible. If the alarm state switch to
OK during the snooze period, this alarm **MUST** will not be shown again.

Multiple snoozes
----------------

When an alarm is already snoozed (whether automatically or manually), an user
**MUST** be able to snooze it again. The alarm will be displayed again after
the last configured snoozed delay.

The actual snooze period used for further computation will start at the date
of the first snooze event and end at the last configured date.

SLA side effect
---------------

When SLA reports are computed, the user **CAN** configure if snooze time has to
be deduced from alarms.

Functional tests
----------------

- When a check KO creating an alarm matches a given filter, the alarm is
  snoozed during a custom time period before being displayed on the UI.

- The user is able to snooze an alarm during a custom time period.
