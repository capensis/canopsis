.. _FR__Snooze:

=============
Snoozed alarm
=============

This document describes the snooze status feature of Canopsis.

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

This feature aims to add the snoozed status to the alarm cycle. This status could be used in two different ways.

The snoozed status at the alarm creation
----------------------------------------

When an alarm is created, this alarm **COULD** first of all be snoozed (for example is the alarm match the given filter). It means that the alarm is not displayed since its creation but only after a defined time period if the state is always "NOK".

If the alarm state switches to OK before the end of the defined time period then the alarm **MUST** be completely removed.

This behavior **MUST** be optional and configurable.

The snoozed status triggered by the user
----------------------------------------

When an alarm is displayed on the widget list, the user **MUST** be able to snooze it.

If an alarm is snoozed by a user, the alarm status is switched to "snoozed" during a time period defined by the user.

During this time period, the alarm is not visible. If the alarm state switch to OK during the snooze period, this alarm **MUST** be completely removed.

Functional tests
----------------

- When an alarm match a given filter, this alarm is snooze before displaying on the UI.

- The user is able to snooze an alarm during a time period set by the user.