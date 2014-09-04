Entities
========

Entities references to all objects known by Canopsis :

 * Component
 * Resource
 * HostGroup
 * ServiceGroup
 * Acknowledgment
 * Downtime
 * Metric
 * Criticality / Macro

All those entities are transmitted via informations stored in events.

The engine ``entities`` extract those informations to a collection ``entities`` in MongoDB.
Each document of this collection represent a single entity.

Component
---------

.. code-block:: javascript

     {
          'type': 'component',
          'name':                  // component field in event
          'hostgroups':            // hostgroups field in event
          'mCrit':                 // Criticality when in state CRITICAL
          'mWarn':                 // Criticality when in state WARNING
     }

Resource
--------

.. code-block:: javascript

     {
          'type': 'resource',
          'name':                  // resource field in event
          'component':             // component field in event
          'hostgroups':            // hostgroups field in event
          'servicegroups':         // servicegroups field in event
          'mCrit':                 // Criticality when in state CRITICAL
          'mWarn':                 // Criticality when in state WARNING
     }

HostGroup
---------

.. code-block:: javascript

     {
          'type': 'hostgroup',
          'name':                  // hostgroup name
     }

ServiceGroup
------------

.. code-block:: javascript

     {
          'type': 'servicegroup',
          'name':                  // servicegroup name
     }

Acknowledgment
--------------

.. code-block:: javascript

     {
          'type': 'ack',
          'component':        // component field in event
          'resource':         // resource field in event
          'timestamp':        // timestamp field in event

          'author':           // author field in event
          'comment':          // output field in event
     }

Downtime
--------

.. code-block:: javascript

     {
          'type': 'downtime',
          'component':        // component field in event
          'resource':         // resource field in event
          'id':               // downtime_id field in event

          'author':           // author field in event
          'comment':          // output field in event

          'start':            // start field in event
          'end':              // end field in event
          'duration':         // duration field in event
          'fixed':            // fixed field in event
          'entry':            // entry field in event
     }

Metric
------

.. code-block:: javascript

     {
          'type': 'metric',
          'component':        // component field in event
          'resource':         // resource field in event
          'name':             // metric field in perf_data_array item
          'nodeid':           // node ID used to fetch data from REST API
     }

Criticality / Macro
-------------------

Those are the only entities who are not defined in this collection.
This kind of document is stored in the collection ``object``, and is used for the SLA view :

.. code-block:: javascript

     {
          'crecord_type': 'sla',
          'objclass': 'crit',

          'crit':                  // Criticality name
          'delay':                 // Number of seconds before SLA time break
     }

     // Must have only one document with objclass=macro
     {
          'crecord_type': 'sla',
          'objclass': 'macro',

          'mCrit':                 // Criticality for CRITICAL state field's name
          'mWarn':                 // Criticality for WARNING state field's name
     }
