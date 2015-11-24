.. _TR__Statistics:

==========
Statistics
==========

This document is the canopsis statistics computation reference where all internal computed metrics are explained.

.. contents::
   :depth: 2

References
==========

- :ref:`FR::Statistics`
- :ref:`TR::Statistics`

Updates
=======


.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Eric Régnier", "2015/10/02", "1.0", "Initial functional requirements", ""

Contents
========

.. _TR__Statistics:

Statistics enhencement
----------------------

Description related to :ref:`desc <FR__Statistics>`.

Software architecture
>>>>>>>>>>>>>>>>>>>>>

Statistics metrics can be produced over many place into Canopsis. Here is next a description of canopsis internal metrics that are produced depending on some condition that have to be explained for better understanding of each of these indicators.

Mainly, there is a statistics manager that contains all statistics computation methods. This manager can be called from many location such as engines. It can aslo be used through the api webservice where a client may reach the procedure to generate some metrics.

Metrics are produced and sent to canopsis system thanks to an event. An event is viewed in Canopsis a context element (see context documentation for more information about this concept). This is why a metric a more than just a value, it is also an information within canopsis related to many other information that can be of various nature.

Computed metrics generate a trend type metric that can increase or decrease the value of a context element's metric. this way, Canopsis system is able to compute back many information through the time thanks to an aggregation system. This can be viewed as an evolving curve over time where past data are linked to a single Canopsis context element. This way, it is possible to retrieve an accurate value information for any existing time period.

Technical guide
>>>>>>>>>>>>>>>

The reason why some metrics may be produced over many code procedure is that the nature of the indicator may bary a lot and contextual information required to compute metrics value may not be available elsewhere. This is the case for alert count indicator that requires that the information of the previous event state is compared to the current passing event state. That's why such a computation have to be done in a code place where the previous event state is available. It technically can be from every where, but for performance purpose (avoid query database in many places) it is only available in the event store engine.

Automated tests
>>>>>>>>>>>>>>>

At the moment, statistics that are computed from the manager are covered by unit tests as long as the business code belongs to the manager.

Writting statistics business code within a manager is the best way to isolate specific code and keep control on how it can be tested.

Automated tests to run for this feature are located in the python sources folder : stats/test/stats.py this will test the stats manager.

Functional tests
>>>>>>>>>>>>>>>>

The section bellow have to explain how (what condition) statistics metrics are produced by the internal Canopsis system. These explaination have to be clear enough to make a tester (person) able to trigger the associated code computation.

Performance
>>>>>>>>>>>

For some performance issues, some metric may be computed asynchronously instead of from the event stream.

