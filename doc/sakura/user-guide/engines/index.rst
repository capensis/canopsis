.. _user-engines:

Engines usage
=============

This document describes what is an engine, and what are the available engines.

For informations about how to configure a specific engine, take a look at:

.. toctree::
   :maxdepth: 1

   event-filter
   linklist
   selector
   datacleaner
   perfdata
   context


Introduction
------------

Events sent by connectors to Canopsis are processed with the help of engines.

An engine has multiple roles:

 * consuming an event : in order to process it, and then route it to the next engine(s)
 * performing a periodic task : called the "beat", this task will be executed at regular interval
 * consuming a dispatched record : called when records from database are available

Each engines is defined by a set of procedure, used to perform the above listed
tasks.

Description
-----------

A single engine is associated to one or more AMQP queue :

 * a queue ``Engine_<engine's name>`` to consume events
 * a queue ``Dispatcher_<engine's name>`` to consume records
 * ...

The beat is executed at regular interval in a parallel thread.

Engine Listing
--------------

+--------------------+----------------+-----------------+------+----------------------------------+
| Engine's name      | Event Consumer | Record Consumer | Beat | Description                      |
+--------------------+----------------+-----------------+------+----------------------------------+
| cleaner_events     | YES            | NO              | NO   | Drop invalid events              |
+--------------------+----------------+-----------------+------+----------------------------------+
| cleaner_alerts     | YES            | NO              | NO   | Drop invalid events              |
+--------------------+----------------+-----------------+------+----------------------------------+
| event_filter       | YES            | NO              | YES  | Apply filter rules               |
+--------------------+----------------+-----------------+------+----------------------------------+
| downtime           | YES            | NO              | YES  | Handle downtimes                 |
+--------------------+----------------+-----------------+------+----------------------------------+
| acknowledgement    | YES            | NO              | YES  | Acknowledge events               |
+--------------------+----------------+-----------------+------+----------------------------------+
| cancel             | YES            | NO              | NO   | Cancel alarms                    |
+--------------------+----------------+-----------------+------+----------------------------------+
| ticket             | YES            | NO              | NO   | Ticketing management             |
+--------------------+----------------+-----------------+------+----------------------------------+
| tag                | YES            | NO              | YES  | Add tags to event                |
+--------------------+----------------+-----------------+------+----------------------------------+
| perfdata           | YES            | NO              | NO   | Store perfdata from event        |
+--------------------+----------------+-----------------+------+----------------------------------+
| eventstore         | YES            | NO              | YES  | Store event in history           |
+--------------------+----------------+-----------------+------+----------------------------------+
| context            | YES            | NO              | YES  | Store contextual data from event |
+--------------------+----------------+-----------------+------+----------------------------------+
| topology           | YES            | NO              | NO   | Topology refresh and management  |
+--------------------+----------------+-----------------+------+----------------------------------+
| linklist           | YES            | NO              | YES  | Manage list of URLs for context  |
+--------------------+----------------+-----------------+------+----------------------------------+
| stats              | NO             | NO              | YES  | Calculate alarm stats            |
+--------------------+----------------+-----------------+------+----------------------------------+
| selector           | NO             | YES             | NO   | Selector refresh and management  |
+--------------------+----------------+-----------------+------+----------------------------------+
| collectdgw         | YES            | NO              | NO   | Consume data from collectd       |
+--------------------+----------------+-----------------+------+----------------------------------+
| consolidation      | NO             | YES             | NO   | Calculate series from metrics    |
+--------------------+----------------+-----------------+------+----------------------------------+
| crecord_dispatcher | NO             | NO              | YES  | Send new records to engines      |
+--------------------+----------------+-----------------+------+----------------------------------+
| eventduration      | YES            | NO              | YES  | Calculate time passed in engines |
+--------------------+----------------+-----------------+------+----------------------------------+
| scheduler          | YES            | NO              | YES  | Send job to task handlers        |
+--------------------+----------------+-----------------+------+----------------------------------+
| task_mail          | YES            | NO              | NO   | Task handler to send mail        |
+--------------------+----------------+-----------------+------+----------------------------------+
| task_linklist      | YES            | NO              | NO   | Task handler to list context URL |
+--------------------+----------------+-----------------+------+----------------------------------+
| task_dataclean     | YES            | NO              | NO   | Task handler to remove old data  |
+--------------------+----------------+-----------------+------+----------------------------------+
| amqp2tty           | YES            | NO              | NO   | Print events to console          |
+--------------------+----------------+-----------------+------+----------------------------------+
