.. _TR__Serie:

=====
Serie
=====

This document specifies how series features must be implemented.

.. contents::
   :depth: 3

References
==========

 - :ref:`FR::Serie <FR__Serie>`
 - :ref:`FR::Metric <FR__Metric>`
 - :ref:`TR::Metric <TR__Metric>`
 - :ref:`FR::Configurable <FR__Configurable>`
 - :ref:`FR::Context <FR__Context>`
 - :ref:`FR::Schema <FR__Schema>`
 - :ref:`FR::Engine <FR__Engine>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/15", "0.1", "Document creation", ""

Contents
========

.. _TR__Serie__Frontend:

Frontend Serie
--------------

A *SerieAdapter* **MUST** provide a method to:

  - fetch :ref:`perfdata <FR__Metric__PerfData>` according to the selected :ref:`metrics <FR__Metric>`
  - aggregate those perfdata using the :ref:`serie aggregation <FR__Serie__Aggregation>`
  - consolidate the aggregated points using the :ref:`serie formula <FR__Serie__Formula>`

This *SerieAdapter* **SHOULD** be used by a widget which wants to visualize the series.

In this case, the formula will be a *sand-boxed* JavaScript expression, where the
:ref:`operators <FR__Serie__Operators>` are methods from the *SerieAdapter* which will returns an array of points
according to the :ref:`regular expression <FR__Serie__Selection>`.

.. _TR__Serie__Backend:

Backend Serie
-------------

A *Serie* :ref:`configurable registry <FR__Configurable__Registry>` **MUST** provide:

 - a method to fetch metrics according to a regular expression and, eventually, a set of already selected metrics (to avoid an operation to the :ref:`Context <FR__Context>`)
 - a method to fetch perfdata using the :ref:`PerfData configurable <TR__Metric__PerfData>`, returning an array of :ref:`timeseries <TR__Metric__TimeSerie>`
 - a method to consolidate timeseries, using the serie formula

In this case, the formula will be a *sand-boxed* Python expression, where the operators are functions with the following prototype:

 - regular expression to select the correct timeseries (using the *Serie* configurable registry) as a parameter
 - returning a single value (according to the operator)

.. _TR__Serie__Engine:

Serie engine
------------

A serie :ref:`engine <FR__Engine>` will provide:

 - a method which periodically check which serie to calculate, and send them on the engine's queue
 - a method to consume the queue, and calculate the serie

The :ref:`data schema <FR__Schema__Data>` of the :ref:`serie configuration <FR__Serie__Configuration>` will provide:

 - ``last_computation`` as an ``integer``, which will be used with the ``aggregation_interval`` to know if a serie must be calculated or not

Unit tests
==========

SerieAdapter
------------

.. warning::

   **TODO:** unit test for the fetch method

.. warning::

   **TODO:** unit test for the aggregate method

.. warning::

   **TODO:** unit test for the consolidation method

Serie configurable registry
---------------------------

.. warning::

   **TODO:** unit test for the metric fetching method

.. warning::

   **TODO:** unit test for the perfdata fetching method

.. warning::

   **TODO:** unit test for the consolidate method

Serie Engine
------------

.. warning::

   **TODO:** unit test for the beat_processing

.. warning::

   **TODO:** unit test for the event_processing
