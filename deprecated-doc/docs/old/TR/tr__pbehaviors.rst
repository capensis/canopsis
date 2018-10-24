TR massive pbehaviors
---------------------

.. contents:: Table of contents


Context
=======

Introduction
^^^^^^^^^^^^

Entities can have specific behavior during a period of time:

 * the entity is in downtime (checks and metrics doesn't interfere with SLA)
 * the entity is disabled (no checks or metrics are saved)
 * the entity is in benchmark (no alerts generated when threshold are exceeded)
 * ...

A periodic behavior is:

 * a period (date-time of start, date-time of end, recursion rule, ...) in iCal format
 * a behavior (downtime, disable, benchmark, ...)
 * the entity it acts on


Defects
^^^^^^^

* A pbehavior can be attached to only an entity at a time
* The pbehavior UI editor is not perfect and permit only to create
* User cannot add comments on pbehaviors
* pbehaviors are accessed by an entity only (not by period)


Work
^^^^

* Rewrite pbehaviors model to fit new needs.
* We won't work on current pbehavior code, new work from scratch.


Backend
=======


Data model
^^^^^^^^^^

Here is the wanted definition of a pbehavior



This need to be adapted like :

.. code-block::

 {
 	"_id" : string,
 	"name" : string,
 	"filter" : string,
 	"comments" : [ {
 	    "_id": string,
 	    "author": string,
 	    "ts": timestamp,
 	    "message": string
 	} ],
 	"tstart" : timestamp,
 	"tstop" : timestamp,
 	"rrule" : string,
 	"enabled" : boolean,
 	"eids" : [ ],
 	"connector" : string,
 	"connector_name" : string,
 	"author" : string
 }

* _id : the unique id of the pbehavior, handled by mongodb itself
* name : the name of the pbehavior
* filter : a mongo filter that match entities from canopsis context
* comments : a list of comments made by users
* tstart : timestamp that correspond to the start of the pbehavior
* tstop : timestamp that correspond to the end of the pbehavior
* rrule : reccurent rule that is compliant with rrule spec
* enabled : boolean to know if pbhevior is enabled or disabled
* eids : list of entity id that match previous filter
* connector : a string representing the type of connector that has generated the pbehavior (**canopsis** if not specified)
* connector_name : a string representing the name of connector that has generated the pbehavior (**canopsis** if not specified)
* author : the name of the user/app that has generated the pbehavior (Logged_in user)


Manager
^^^^^^^

.. csv-table:: Methods
   :header: "#", "Name", "Args", "Comments"
   :widths: 5, 40, 80, 80

   "1", "create_pbehavior", "name, filter, tstart, tstop, rrule, enabled, comments, connector, connector_name, author", "Except comments, rrule, enabled, all args are mandatory"
   "2", "read_pbehavior", "_id", "_id arg is optional, if not defined, get all pbehaviors"
   "3", "update_pbehavior", "_id, name, filter, tstart, tstop, rrule, enabled, connector, connector_name, author", "_id arg is mandatory"
   "4", "delete_pbehavior", "_id", "_id arg is mandatory"
   "5", "create_pbehavior_comment", "pbehavior_id, author, message", "Note that a comment must have an _id"
   "6", "update_pbehavior_comment", "pbehavior_id, _id, author, message", "both pbehavior_id and _id are mandatory"
   "7", "delete_pbehavior_comment", "pbehavior_id, _id", "both pbehavior_id and _id are mandatory"
   "8", "get_pbehaviors", "entity_id", "Return a list of pbehaviors (name, tstart, tstop, rrule, enabled, comments, connector, connector_name, author) that match entity_id in eids. entity_id is mandatory"
   "9", "compute_pbehaviors_filters", "", "Compute all filters and update eids attributes. * Explained in engine section"



Engines
^^^^^^^

In canopsis, events are processed by engines.
Engines generaly receive events by consuming AMQP Queue.

An engine consists in 2 methods :

1. Work : executed when an event is consumed
2. Beat : every beat interval




pbehavior
^^^^^^^^^

This new engine **pbehavior** has 2 goals :

1. Compute pbehaviors_filters in order to build **eids** list in a pbehavior record
2. Consume event_type = 'pbehavior'


**1 Compute pbehaviors filter**

At every engine beat, the engine must

* Iterate on pbehaviors filters
* For each filter, ask the canopsis context to get a list of entity ids
* Insert the **eids** attribute in a pbehavior record


**2 Consume events**

When there is a message (event_type: pbehavior) in a AMQP queue (corresponding to the engine queue), the engine must

* Read the message
* Understand which kind of operation the message deals with
* Create or Delete the corresponding pbehavior

Here is the structure of a compliant message

.. code-block::

    {
      "event_type": "pbehavior",
      "pbehavior_name": "downtime",
      "start": ts,
      "end": ts,
      "comment": "a comment about the pbehaviour",
      "author": "the author",
      "action": "create" or "delete" only
      "rrule": "an rrule",
      + classic event fields
    }

for creating actions, create an entry with this fields:

.. code-block::

    "name": pbehavior_name,
    "filter": {
      "entity_id": /resource/connector/connector_name/component/resource (for a resource, test type field value (component or resource))
    },
    "comments": no comments,
    "tstart": start,
    "tstop": end,
    "rrule": "",
    "enabled": True,
    "eids": [],
    "connector": comes from the event,
    "connector_name": comes from the event
    "author": comes from the event


for deleting actions, delete an entry with exactly the same fields as it is given for these ones:
remove an entry in mongo with this filter:

.. code-block::

    "filter": {
      "entity_id": /resource/connector/connector_name/component/resource (for a resource, test type field value (component or resource))
    },
    "name": pbehavior_name,
    "tstart": start,
    "tstop": end,
    "rrule": "",
    connector": comes from the event,
    "connector_name": comes from the event

And whatever for the others.


Event Filter
^^^^^^^^^^^^

The event filter handle one or several events thanks to a filter and apply action(s) to these events
The current event filter engine must be compliant with Pbehaviors.

Actually, an event filter is composed of

* A filter
* A list of actions

The event filter must now provide a way to match pbehaviors you implement, so it has to become

* A filter
* A list of within/without pbehaviors
* A list of actions

This new point let the user include pbehaviors to handle or not an event.
For example, if an event match the given filter AND this event is included in pbehavior1 OR this event is not included in pbehavior2, then the event filter handle this event and apply given actions.

In order to add this new point in event filter, you have to add a new method in the manager, **check_pbehaviors** with following args :

* entity_id
* in=[ pbehavior_name, ]
* out=[ pbehavior_name, ]

This method return a boolean if the entity_id is currently in **in** arg and out **out** arg.
**in** and **out** are evaluated with **tstart**, **tstop**, and **rrule** timestamps compared to **now**

To fit your changes, you also have to change the schemas in the database related to the event_filter:

* Connect to mongo `mongo canopsis -u cpsmongo -p canopsis`
* Find the event filter schema in object collection
* Modify this schema: add the pbehaviors settings.

Selector / SLA
^^^^^^^^^^^^^^

Selector
~~~~~~~~

The selector engine aggregates entity states to build new entities.

* First, selectors are defined by the user with a filter. This filter aggregates several entities. To do that, you have to see in the context all entities that match the filter.
* Then the engine selector will find all alarms related to these entities. (The old selector system uses 'events' collection but we want to migrate to 'alerts' collections)
* Finally, the engine will compute the entity state with the rule given by the user (worst state, best state, mean, etc...)

This new selector engine has to be adaptable, in order to let us use it with our new context, actually in development that let this engine do only one request (context and alerts at the same time).

You can retrieve the old engine in develop branch of this repo at this place: canopsis/sources/python/engines/canopsis/engines/selector.py

SLA
~~~

SLA ARE NOT READY FOR THE MOMENT - WE HAVE TO SPECIFY THEM BEFORE IMPLEMENT THEM. PLEASE DO NOT PAY ATTENTION TO ANY METHOD RELATED TO THE SLA IF YOU FIND THEM.

It can be combined to the sla lib that will calcultate availability rates.
Some behaviors affects SLA rates.
For example, if a entity is in downtime, a selector which uses that entity must not be affected.


Connectors
^^^^^^^^^^

**neb2amqp**


As the `neb2amqp` connector is handling `downtime` events, it needs to fit the new pvehavior schema.


Getting started
^^^^^^^^^^^^^^^

Sources
~~~~~~~

A new branch has been created on canopsis open core project :
``feature-pbehaviors``. Changes must be commited in this branch. This branch
contains :

* pbehavior manager in
  ``sources/python/pbehavior/canopsis/pbehavior/manager.py``.  Business logic
  must be written in this file (methods described in the table above).
* pbehavior webservice in
  ``sources/python/webcore/canopsis/webcore/services/pbehavior.py``. This file
  contains proxy functions that must rely on the manager. Routes should be
  requestable : <ip>:<port>/pbehavior/create
* pbehavior engine in
  ``sources/python/pbehavior/canopsis/pbehavior/process.py``.  An engine is
  composed of 2 functions :

  - event_processing : called each time an event is received
  - beat_processing : called once a minute

Examples showing how to achieve main operations are provided in source files.

Environment
~~~~~~~~~~~

A development environment is available at ... . It has been deployed from the
``feature-pbehaviors`` branch.

You should work with local sources and push your modifications on the
environment to test. Here at capensis we tend to use ``rsync``.

Once you changed some code, you can reload it with :

  * ``service amqp2engines* mrestart`` for the engine
  * or ``service webserver restart``

Logs
~~~~

Log files that should be used are :

  * /opt/canopsis/var/log/pbehaviormanager.log
  * /opt/canopsis/var/log/engines/pbehavior.log

Frontend
========

Pbehavior editor
^^^^^^^^^^^^^^^^

By :

 * alarms view
 * context view
 * dedicated view with filter possibilities


Action buttons
^^^^^^^^^^^^^^

As ack, ticket, and other actions, `pbehaviors` must have it's own action buttons.

It has to be available on the following views :

* alarms
* context
* other widgets : have a look at https://git.canopsis.net/canopsis-ui-bricks/brick-service-weather/blob/master/doc/TR/TR_service_weather.rst
Work has maybe be already done



Calendar widget
^^^^^^^^^^^^^^^

Be able to set pbehavior to entities using the calendar widget
