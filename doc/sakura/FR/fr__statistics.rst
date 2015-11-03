.. _FR__Statistics:

=============================
Canopsis Statistics Generator
=============================

This document describes the statistics generation process.

.. contents::
   :depth: 3

References
==========

List of referenced functional requirements:

 - :ref:`FR::Metric <FR__Metric>`
 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Serie <FR__Serie>`
 - :ref:`FR::Schema <FR__Schema>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/15", "0.1", "Document creation", ""

Contents
========

.. _FR__Statistics__Desc:

Description
-----------

The statistics generator allows the user to create :ref:`metrics <FR__Metric>`,
generating :ref:`perfomance data <FR__Metric__PerfData>` about Canopsis usage,
availability, and performance.

It will provide:

 - basic metrics, that **MAY** be expanded using a filter, calculated from :ref:`events <FR__Event>`
 - complex metrics, calculated from basic metrics using an :ref:`aggregated serie <FR__Serie>`

.. _FR__Statistics__User:

User statistics
---------------

Metrics
~~~~~~~

.. csv-table::
   :header: "Component", "Resource", "Metric", "Metric Type", "Description"

   ``$username``, alarm_ack, total, COUNTER, Total number of acknowledged alarms
   ``$username``, alarm_ack, ``$filter_name``, COUNTER, Number of acknowledged alarms matching filter
   ``$username``, alarm_ack_delay, last, GAUGE, Delay between alarm apparition and its acknowledgment
   ``$username``, alarm_ack_delay, average, GAUGE
   ``$username``, alarm_ack_delay, max, GAUGE
   ``$username``, alarm_ack_delay, min, GAUGE
   ``$username``, alarm_ack_solved, last, GAUGE, Delay between alarm resolution and its acknowledgment
   ``$username``, alarm_ack_solved, average, GAUGE
   ``$username``, alarm_ack_solved, max, GAUGE
   ``$username``, alarm_ack_solved, min, GAUGE
   ``$username``, alarm_solved, last, GAUGE, Delay between alarm apparition and its resolution
   ``$username``, alarm_solved, average, GAUGE
   ``$username``, alarm_solved, max, GAUGE
   ``$username``, alarm_solved, min, GAUGE
   ``$username``, session_duration, last, GAUGE, Delay between user's login and logout
   ``$username``, session_duration, average, GAUGE
   ``$username``, session_duration, max, GAUGE
   ``$username``, session_duration, min, GAUGE

Computation
~~~~~~~~~~~

 * ``alarm_ack``: will be incremented if the event matches the configured filter
 * ``alarm_ack_delay``: will be the difference between the timestamps of the alarm event and its acknowledgment
 * ``alarm_ack_solved``: will be the difference between the timestamps of the ack and the alarm resolution event
 * ``alarm_solved``: will be the difference between the timestamps of the alarm event and its resolution event
 * ``session_duration``: will be the difference between the timestamps of the first login and the last Canopsis page closed
 * ``average``, ``max`` and ``min`` will be calculated from ``last`` using an aggregated serie

.. _FR__Statistics__Event:

Event statistics
----------------

Metrics
~~~~~~~

.. csv-table::
   :header: "Component", "Resource", "Metric", "Metric Type", "Description"

   canopsis, alarm, total, COUNTER, Total number of alarms
   canopsis, alarm, ``$filter_name``, COUNTER, Number of alarms matching filter
   canopsis, alarm_ack, total, COUNTER, Total number of acknowledged alarms
   canopsis, alarm_ack, ``$filter_name``, COUNTER, Number of acknowledged alarms matching filter
   canopsis, alarm_ack_solved, total, COUNTER, Total number of acknowledged resolved alarms
   canopsis, alarm_ack_solved, ``$filter_name``, COUNTER, Number of acknowledged resolved alarms matching filter
   canopsis, alarm_solved, total, COUNTER, Total number of resolved alarms
   canopsis, alarm_solved, ``$filter_name``, COUNTER, Number of resolved alarms matching filter
   canopsis, alarm_ack_delay, sum, COUNTER, Cumulative delay between alarm apparition and its acknowledgment
   canopsis, alarm_ack_delay, last, GAUGE, Delay between alarm apparition and its acknowledgment
   canopsis, alarm_ack_delay, average, GAUGE
   canopsis, alarm_ack_delay, max, GAUGE
   canopsis, alarm_ack_delay, min, GAUGE
   canopsis, alarm_ack_solved, sum, COUNTER, Cumulative delay between alarm resolution and its acknowledgment
   canopsis, alarm_ack_solved, last, GAUGE, Delay between alarm resolution and its acknowledgment
   canopsis, alarm_ack_solved, average, GAUGE
   canopsis, alarm_ack_solved, max, GAUGE
   canopsis, alarm_ack_solved, min, GAUGE
   canopsis, alarm_solved, sum, COUNTER, Cumulative delay between alarm apparition and its resolution
   canopsis, alarm_solved, last, GAUGE, Delay between alarm apparition and its resolution
   canopsis, alarm_solved, average, GAUGE
   canopsis, alarm_solved, max, GAUGE
   canopsis, alarm_solved, min, GAUGE

Computation
~~~~~~~~~~~

 * ``alarm``: will be incremented on each :ref:`check event <FR__Event__Check>` in a non OK state, if it matches the configured filter
 * ``alarm_ack``: will be incremented for each :ref:`ack event <FR__Event__Ack>`, if it matches the configured filter
 * ``alarm_ack_solved``: will be incremented for each check event (previously acknowledged) in an OK state, if it matches the configured filter
 * ``alarm_solved``: will be incremented for each check event in an OK state, if it matches the configured filter
 * ``alarm_ack_delay.sum``: will be incremented by the difference of timestamps between the alarm event and its acknowledgment
 * ``alarm_ack_solved.sum``: will be incremented by the difference of timestamps between the acknowledgment and the alarm resolution event
 * ``alarm_solved.sum``: will be incremented by the difference of timestamps between an alarm event and its resolution event
 * ``*.last``: will be published with the incremental value of ``*.sum``
 * ``average``, ``max`` and ``min`` will be calculated from ``last`` with an aggregated serie

.. _FR__Statistics__Configuration:

Configuration
-------------

In order to create the event filters (to produce complex metrics), we will need 
at least 5 collections of filters:

 - ``user.alarm_ack``
 - ``canopsis.alarm``
 - ``canopsis.alarm_ack``
 - ``canopsis.alarm_ack_solved``
 - ``canopsis.alarm_solved``

For this, we will need a :ref:`data schema <FR__Schema__Data>` providing the following informations:

 - a ``filter_name`` as a ``string``
 - a ``filter_type`` as a ``string`` which represent each collection
 - a ``filter`` representing the event filter to use for event matching

A **Statistics** view will be available, providing:

 - a configuration tab: with a listing of all filters, allowing CRUD operations
 - a view tab: with graphs visualizing the produced metrics

Functional Tests
================

.. warning::

   **TODO:** listing of expected result per metric

.. warning::

   **TODO:** listing of expected metrics according to filters
