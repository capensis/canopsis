.. _TR__Statistics:

=============================
Canopsis Statistics Generator
=============================

This document specifies how the Statistics Generator must be implemented.

.. contents::
   :depth: 3

References
==========

 - :ref:`FR::Statistics <FR__Statistics>`
 - :ref:`FR::Configurable <FR__Configurable>`
 - :ref:`FR::Event <FR__Event>`
 - :ref:`FR::Engine <FR__Engine>`
 - :ref:`FR::Webservice <FR__Webservice>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/15", "0.1", "Document creation", ""

Contents
========

.. _TR__Statistics__MetricProducer:

Metric Producer
---------------

A :ref:`configurable registry <FR__Configurable__Registry>` will provide:

 - a :ref:`filter configuration <FR__Statistics__Configuration>` storage
 - a method returning ``True`` or ``False`` wether an :ref:`event <FR__Event>` match the filter or not
 - a method to create an :ref:`aggregated serie <FR__Serie>` from a :ref:`basic metric <FR__Statistics__Desc>`

.. _TR__Statistics__UserMetricProducer:

User Metric Producer
--------------------

A subclass of the :ref:`MetricProducer <TR__Statistics__MetricProducer>` configurable,
providing:

 - a method per :ref:`resource <FR__Statistics__User>`:
    - parameters:
        - the username as a parameter
        - the events involved
    - returns: one or more :ref:`metrics <FR__Event__Perf>` to publish

.. _TR__Statistics__EventMetricProducer:

Event Metric Producer
---------------------

A subclass of the :ref:`MetricProducer <TR__Statistics__MetricProducer>` configurable,
providing:

 - a method per :ref:`resource <FR__Statistics__Event>`:
    - parameters:
        - the events involved
    - returns: one or more :ref:`metrics <FR__Event__Perf>` to publish

.. _TR__Statistics__Engine:

Statistics engine
-----------------

An :ref:`engine <FR__Engine>` will listen for the following events:

 - :ref:`check events <FR__Event__Check>`
 - :ref:`ack events <FR__Event__Check>`

This engine will be composed of the :ref:`UserMetricProducer <TR__Statistics__UserMetricProducer>` and the :ref:`EventMetricProducer <TR__Statistics__EventMetricProducer>`, which will be used to:

 - periodically fetch filters from storage
 - cache check events with an OK state, and ack events, to know the involved events when needed
 - if the event match one filter, call methods from both :ref:`MetricProducer <TR__Statistics__MetricProducer>`
 - publish all the returned events

*NB: after benchmarking, we may have to dispatch those calls in specific engines.*

.. _TR__Statistics__Webservice:

Statistics webservice
---------------------

A :ref:`webservice <FR__Webservice>` will provide two routes:

 - one to call when a Canopsis UI is opened, which will increment a counter for the user
 - one to call when a Canopsis UI is closed, which will decrement a counter for the user

When the counter is incremented, and was ``0``, the current timestamp will be registered.
When the counter fall back to ``0``, the difference between the registered timestamp and the current timestamp will be published as a metric.

Unit tests
==========

.. warning::

   **TODO:** unit test of each MetricProducer methods

.. warning::

   **TODO:** unit test of each type of UserMetricProducer method

.. warning::

   **TODO:** unit test of each type of EventMetricProducer method

.. warning::

   **TODO:** unit test for each route of the webservice
