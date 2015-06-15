Glossary
========

In this glossary, you will find a short definition of all concepts in *Canopsis*
with a link to the complete documentation.

Backend
~~~~~~~

Connector
---------

Daemon fetching data from a source (local system performance data, service, ...)
and publishing it to *Canopsis* via the *AMQP Bus* or the *WSGI API*.

.. NOTE::

   TODO: Add link

AMQP Bus
--------

Messaging system based on exchanges and queues, used by *Canopsis* to process
events asynchronously.

.. NOTE::

   TODO: Add link

Engine
------

Daemon with an associated AMQP queue. This daemon can consume the queue and/or
publish messages to another queue.

For more informations, see:

 * :ref:`dev-backend-engines`
 * :ref:`admin-manage-engines`
 * :ref:`user-engines`

Middleware
----------

*Canopsis* object which allows two parts of *Canopsis* to communicate (via AMQP,
HTTP, ...).

.. NOTE::

   TODO: Add link

Storage
-------

*Canopsis* object used to interact with a database (*MongoDB*, *PostgreSQL*,
*ElasticSearch*, ...)

For more informations, see:

 * :ref:`dev-backend-storage`

Manager
-------

Using a storage, allows parts of *Canopsis* to interact with any kind of database
transparently.

For more informations, see:

 * :ref:`dev-backend-mgr`

Web Service
-----------

Set of WSGI routes, using one or more manager to provide data to the client.

For more informations, see:

 * :ref:`dev-backend-webserver`

Frontend
~~~~~~~~

UI Brick
--------

Set of UI Adapters/Components/Editors/Widgets.

For more informations, see:

 * :ref:`dev-frontend-architecture`

UI Adapter
----------

Used to interact with the *WSGI API*.

.. NOTE::

   TODO: Add link

UI Component
------------

Used to display interactive data to the user.

For more informations, see:

 * :ref:`dev-frontend-cmp`

UI Editor
---------

Using a component, provides a way to edit a model, before persisting it to the
*WSGI API*.

For more informations, see:

 * :ref:`dev-frontend-cmp-editors`

UI Widget
---------

Improved component, based on a **MVC** design (unlike the UI components).
They can be directly added to a view, and can have mixins.

For more informations, see:

 * :ref:`dev-frontend-widgets`

UI Mixin
--------

Set of business code that can be applied to any widget.

For more informations, see:

 * :ref:`dev-frontend-widgets-mixins`

UI Container
------------

Component containing widgets, used to dispose them in a specific layout.
There is only one widget container, which can have different layout mixins.

For more informations, see:

 * :ref:`user-ui-widgets-containers`

UI View
-------

Editable view which contains by default a single widget container.

For more informations, see:

 * :ref:`user-ui-view`

Miscellaneous
~~~~~~~~~~~~~

Event
-----

JSON object containing specific informations for *Canopsis*, must be emitted on
the *AMQP Bus*.

For more informations, see:

 * :ref:`dev-spec-event`
 * :ref:`dev-backend-event`
 * :ref:`user-events`

Metric
------

Measurable information, associated to a component, or a resource. Can be used to
render in a widget graph, progress-bar, and/or text.
It is a contextual information referenced by each new inserted value.

For more informations, see:

 * :ref:`dev-frontend-widgets-perfdata`
 * :ref:`user-ui-view-perfdata`

Context
-------

Contextual informations about an event, organized in graph. All other stored data
are referencing the associated context, for example:

 * a perfdata document reference the metric context
 * a periodic behavior reference the component or resource context
 * ...

A view is available in order to manipulate the context.

For more informations, see:

 * :ref:`user-ui-view-context`
 * :ref:`dev-backend-mgr-vevent`
 * :ref:`dev-backend-mgr-pbehavior`

Selector
--------

.. NOTE::

   TODO: add short description

For more informations, see:

 * :ref:`user-engines-selector`

SLA
---

Feature providing availability informations.

For more informations, see:

 * :ref:`dev-spec-sla`

Periodic Behavior
-----------------

An entity of the context can be configured to have a specific behavior during a
specified period of time.

For more informations, see:

 * :ref:`dev-backend-mgr-pbehavior`

Downtime
++++++++

A downtime is configured when we must ignore eventual alerts on an entity.

For more informations, see:

 * :ref:`dev-backend-mgr-pbehavior`
