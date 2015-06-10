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

See :ref:`dev-backend-engines` for more informations.

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

See :ref:`dev-backend-storage` for more informations.

Manager
-------

Using a storage, allows parts of *Canopsis* to interact with any kind of database
transparently.

See :ref:`dev-backend-mgr` for more informations.

Web Service
-----------

Set of WSGI routes, using one or more manager to provide data to the client.

See :ref:`dev-backend-webserver` for more informations.

Frontend
~~~~~~~~

UI Brick
--------

Set of UI Adapters/Components/Editors/Widgets.

See :ref:`dev-frontend-architecture` for more informations.

UI Adapter
----------

Used to interact with the *WSGI API*.

.. NOTE::

   TODO: Add link

UI Component
------------

Used to display interactive data to the user.

See :ref:`dev-frontend-cmp` for more informations.

UI Editor
---------

Using a component, provides a way to edit a model, before persisting it to the
*WSGI API*.

See :ref:`dev-frontend-cmp-editors` for more informations.

UI Widget
---------

Improved component, based on a **MVC** design (unlike the UI components).
They can be directly added to a view, and can have mixins.

See :ref:`dev-frontend-widgets` for more informations.

UI Mixin
--------

Set of business code that can be applied to any widget.

See :ref:`dev-frontend-widgets-mixins` for more informations.

UI Container
------------

Component containing widgets, used to dispose them in a specific layout.
There is only one widget container, which can have different layout mixins.

See :ref:`user-ui-widgets-containers` for more informations.

UI View
-------

Editable view which contains by default a single widget container.

See :ref:`user-ui-view` for more informations.

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

See :ref:`user-engines-selector` for more informations.

SLA
---

Feature providing availability informations.

See :ref:`dev-spec-sla` for more informations.

Periodic Behavior
-----------------

An entity of the context can be configured to have a specific behavior during a
specified period of time.

See :ref:`dev-backend-mgr-pbehavior` for more informations.

Downtime
++++++++

A downtime is configured when we must ignore eventual alerts on an entity.

See :ref:`dev-backend-mgr-pbehavior` for more informations.
