.. _FR__Meta_alarm:

==========
Meta alarm
==========

This document describes the meta alarm feature of Canopsis.

.. contents::
   :depth: 3

References
==========

List of referenced functional requirements:

 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Context <FR__Context>`
 - :ref:`FR::Alarm <FR__Alarm>`
 - :ref:`FR::Selector <FR__Selector>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Florent Demeulenaere", "2016/06/22", "0.1", "Document creation", ""

Contents
========

.. _FR__Meta_alarm__Desc:

Description
-----------

This feature aims to add the meta alarm in Canopsis. This alarm **COULD** be created when one or more event are triggered in a time period defined by the user.

Meta alarm
----------

A selection of one or several event **COULD** be triggered in a defined period of time in order to create a meta alarm in case of this entire selection is triggered.

Functional tests
----------------

- The user can make a selection of event(s) and defined a period of time and the description of the meta alarm that will be created.

- When all conditions are filled to created the meta alarm, this alarm is created and visible in Canopsis.