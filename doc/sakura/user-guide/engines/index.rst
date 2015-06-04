Engines
=======

.. toctree::
   :maxdepth: 2

   consolidation
   datacleaner
   event-filter
   event-requalification
   selector
   linklist
   snmp

.. include:: ../../includes/links.rst

Engines information
===================

.. |sla| raw:: html

   <font color="red">SLA</font>

.. |Topology| raw:: html

   <font color="red">Topology</font>

Canopsis process system events thanks to a set of engines.

Description
------------

There are two types of engines, asynchronous and synchronous.

A synchronous engine belongs to a chain of engines which consume events
from |queues|.

An asynchronous engine is independent from an execution chain of engines
and process only events stored in a database.

Architecture
------------

Here is an architecture view of engines (designed with
`CACOO <https://cacoo.com/>`__).

|image1|

In this view, two chain of engines and flow of
|event| ,
|perf_data| and Canopsis
configuration data are represented.

The event chain is the entry point for all poller/Canopsis events. Its
goal is to consume events and transform them into perfdata.

The alert chain manages all Canopsis alerts at the end of event
processing from the event chain.

Let's see roles for all engines (hyperlinks denote engines which are
configurable by the user).

Synchronous engines
-------------------

Cleaner
~~~~~~~

The Cleaner aims to clean events sent to Canopsis.

|event-filter|
~~~~~~~~~~~~~~~~~~~~~

The filter aims to check if an event has to be destroyed or processed.

|derogation|
~~~~~~~~~~~~~~~~~~~

The derogation aims to transform event fields from a poller to Canopsis
fields.

Tag
~~~

The tag aims to fill the tag fields of events.

Perfstore2
~~~~~~~~~~

The perfstore2 aims to save perfdata in the random access database.

Event store
~~~~~~~~~~~

The event store saves consumed events in the persistent database or send
an alert corresponding to an inconsistent event.

Alerts
~~~~~~

The alerts cleaner aims to clean alerts.

Alert counter
~~~~~~~~~~~~~

The alert counter count number of alerts by alert type.

|Topology|
~~~~~~~~~~~~~~~~~

The topology check rules with input alerts.

|selector|
~~~~~~~~~~~~~~~~~

The selector calculates a worst state about input alerts.

Asynchronous engines
--------------------

CollectDGW
~~~~~~~~~~

|sla|
~~~~~~~~~~~~

Calculate SLA of system services.

|consolidation|
~~~~~~~~~~~~~~~~~~~~~~

Do consolidation/aggregation on perfdata.

Perfstore2\_rotate
~~~~~~~~~~~~~~~~~~

Switch perfdata from the random access database to the persistent
database.

Engines are presented below as they appear in the default configuration. On some architectures it can be relevant to tweak their configuration and duplicate some engines.

Events queue
------------

+----------------+------------------------------------------------+------+------+
| Engine name    | Description                                    | Work | Beat |
+================+================================================+======+======+
| cleaner        | Clean events to ensure they won't cause errors | YES  | NO   |
|                | Same engine as the one in Alerts queue         |      |      |
+----------------+------------------------------------------------+------+------+
| event_filter   | Event firewall                                 | YES  | YES  |
+----------------+------------------------------------------------+------+------+
| derogation     |                                                | YES  | YES  |
+----------------+------------------------------------------------+------+------+
| tag            | Add tags to events                             | YES  | YES  |
+----------------+------------------------------------------------+------+------+
| perfstore2     | Store events' metrics in redis                 | YES  | NO   |
+----------------+------------------------------------------------+------+------+
| eventstore     | Store events in mongo                          | YES  | NO   |
+----------------+------------------------------------------------+------+------+


Alerts queue
------------


+----------------+------------------------------------------------+------+------+
| Engine name    | Description                                    | Work | Beat |
+================+================================================+======+======+
| cleaner        | Clean events to ensure they won't cause errors | YES  | NO   |
+----------------+------------------------------------------------+------+------+
| alertcounter   |                                                | YES  | NO   |
+----------------+------------------------------------------------+------+------+
| topology       |                                                | YES  | YES  |
+----------------+------------------------------------------------+------+------+
| selector       |                                                | YES  | YES  |
+----------------+------------------------------------------------+------+------+

Others
------

+-------------------+------------------------------------------------+------+------+
| Engine name       | Description                                    | Work | Beat |
+===================+================================================+======+======+
| collectdgw        |                                                | NO   | NO   |
+-------------------+------------------------------------------------+------+------+
| sla               |                                                | NO   | YES  |
+-------------------+------------------------------------------------+------+------+
| consolidation     |                                                | NO   | YES  |
+-------------------+------------------------------------------------+------+------+
| perfstore2_rotate | Move metrics and perfdatas from redis to mongo | NO   | YES  |
+-------------------+------------------------------------------------+------+------+


.. |image1| image:: ../../_static/images/engine/engines_map.png
