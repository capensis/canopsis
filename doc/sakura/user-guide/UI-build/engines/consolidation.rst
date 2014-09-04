.. _consolidation:

Consolidation
=============

.. |sla| raw:: html

   <font color="red">sla</font>
   
A `consolidation <#consolidation2>`__ permits to generate a new metric in
`aggregating <#aggregation>`__ a set of metrics and in crushing
aggregated values in one final metric.

An example of consolidation is drawn in the following figure.

|image1|

Four graphs are represented vertically, separated by an horizontal
dotted line: "metric 1",  "aggregation 1",  "aggregation 2" and "consolidation".

All graphs provide values where the abscissa axis is in time, and
ordinate in the same arbitrary unit.

The two intervals [t0; t1[ and [t1; t2[ are defined with the same
duration and the last interval [t2; now] is less longer than previous
ones.

Finally, there are `aggregation operations <#operation>`__ between the
"metric 1" and "aggregation 1" graphs and `consolidation
operations <#%20aggregation%20and%20consolidation>`__ between the both
aggregations graphs and the consolidation graph.

Let's see how `aggregation <#aggregation>`__ and `consolidation <#consolidation2>`__ works.

Aggregation
-----------

Description
~~~~~~~~~~~

An aggregation is a function which takes in parameter one metric, one
aggregation `operation <#operation>`__ and one aggregation interval.

The aggregation process starts to group metric values which are in the
same aggregation interval (in the previous figure, the interval [t0; t1[
is composed of five values).

Then the aggregation `operation <#operation>`__ is applied on each
aggregation interval and it generates a new metric where values are
positioned at the beginning of their aggregation interval (the result of
values between t0 and t1 is at t0).

Example
~~~~~~~

Let a metric related to a service |sla| ,
you can get the average sla per day in selecting the aggregation
operation mean and an aggregation interval of one day.

.. _consolidation2:

Consolidation
-------------

Description
~~~~~~~~~~~

A consolidation generates one metric per consolidation operation from at
least one aggregation. While consolidation metric component and resource
names are customizable, the metric name is the consolidation operation
name.

The consolidation operation is applied for every aggregation values at
the same position on the abscissa axis. It implies all aggregation must
have the same aggregation interval.

Example
~~~~~~~

Let five aggregations A1, A2, A3, A4 and A5 which are respectively mean
per day of |sla| of all services of one
server "S". You can get the general mean per day of S
|sla| in applying the consolidation
operation mean on all aggregations.

Engine configuration
~~~~~~~~~~~~~~~~~~~~

The engine consolidation is dedicated to calculate consolidations from
Canopsis input events.

In order the parameterize such consolidations, let's go to the view
consolidation from the consolidation menu.

|image2|

The view consolidation allows you to
edit/add/delete/search/import/export consolidations.

|image3|

In editing mode (after double clicking on an existing consolidation or
in pushing the button "add"), a wizard with different tabs allows you to
configure the consolidation.

The tab "options" permits to set consolidation, component and resource
names of the metrics to generate.

|image4|

The tab "consolidation" defines the aggregation interval and the
aggregation/consolidation operations.

|image5|

Finally, the tab "filter" permits to select metrics to aggregate and to
consolidate.

|image6|

When you have set up your consolidation, the engine will calculate the
consolidation and reveals you number of metrics which have been
aggregated and the consolidation time processing.

.. _operation:

Aggregation/Consolidation Operations
------------------------------------

Here are the list of available operations for aggregations and
consolidations


Aggregation only
~~~~~~~~~~~~~~~~

First
_____

Select the first value.

Last
_____

Select the last value.

Aggregation and consolidation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Min
___

Select the minimal value.

Max
___

Select the maximal value.

Mean
____

Select the mean value (sum(values) / count(values)).

Delta
_____

Select the delta value (max(values) - min(values))

Sum
___

Select the sum value.

.. |image1| image:: ../../../_static/images/consolidation/consolidation.png
.. |image2| image:: ../../../_static/images/consolidation/consolidation_menu.png
.. |image3| image:: ../../../_static/images/consolidation/consolidation_view.png
.. |image4| image:: ../../../_static/images/consolidation/consolidation_options.png
.. |image5| image:: ../../../_static/images/consolidation/consolidation_consolidation.png
.. |image6| image:: ../../../_static/images/consolidation/consolidation_filter.png
