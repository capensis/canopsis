.. include:: /includes/links.rst

.. _glossary:

Glossary
========

Here are the terminology and the different concepts used in Canopsis |resource| .

.. _component:

Component
---------

A component is a (virtual) machine which contains a set of .

An example of component is a server.

The association with a system machine is realized with a name. If two components have the same name, they have to get different |connector| names in order to be differentiated by Canopsis.

.. _resource:

Resource
--------

A resource is contained in a |component| and is identified by a unique name in the scope of one component.

An example of resource is a service.

.. _connector:

Connector
---------

A connector is the mean of Canopsis to access to a |component| .

.. _metric:

Metric
------

A metric is related measure that facilitates the quantification of some system characteristic.

For example, a service monitored by a supervisor can publish an event with a metric about its availability in time duration.

A metric is identified with a unique name in the scope of the triplet (component, resource, connector). Therefore, an event contains at least one metric.

A metric contains at least a value, and optionally an unit (meter, time, etc.), maximal/minimal values,
a type (GAUGE, COUNTER, etc.) and warning/critical and thresholds.

Go to |metricsPage| for more details.

.. _publisher:

Publisher
---------

A publisher is a |supervisor| which feeds Canopsis with events through |queues|.

nagios, icinga, graylog ..

.. _hypervisor:

Hypervisor
----------

An hypervisor is at the top of supervisors in order to supervise and to analyze a large, distributed and heterogeneous system. Its goal is to supervise supervisors and to provide better data analyses about monitored systems with historical and near real-time concerns.

.. _supervisor:

Supervisor
----------

A supervisor is a system which uses pollers in order to retrieve information from a system infrastructure such as services availability.

.. _engine:

Engine
------

An engine is an event processor.

Several |engines| exist in Canopsis, such as the |filter| engine which allows events to be processed by Canopsis engines, or the |consolidation| engine which calculates consolidation to do with input engines.

.. _event:

Event
-----

An event is a data information processed by Canopsis engines. An event can be published by a supervisor or by an engine in the case of engine processing results.

It is identified by the triplet (component, resource, connector). Therefore only the last event published with this triplet is saved in a database, unlike event processing results (perf data) which are all saved in a database.

.. _queue:

Queue
-----

A queue is an event target which can be processed by synchronous |engines|.

.. _perf_data:

Perf data
---------

Every event consumed by Canopsis is transformed into a perf data which is saved in a database for future processing such as data analysis.
