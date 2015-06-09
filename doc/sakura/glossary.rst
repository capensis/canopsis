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


.. NOTE::

   TODO: Add link

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

.. NOTE::

   TODO: Add link

Manager
-------

Using a storage, allows parts of *Canopsis* to interact with any kind of database
transparently.

.. NOTE::

   TODO: Add link

Web Service
-----------

Set of WSGI routes, using one or more manager to provide data to the client.

.. NOTE::

   TODO: Add link

Frontend
~~~~~~~~

UI Brick
--------

Set of UI Adapters/Components/Editors/Widgets.

.. NOTE::

   TODO: Add link

UI Adapter
----------

Used to interact with the *WSGI API*.

.. NOTE::

   TODO: Add link

UI Component
------------

Used to display interactive data to the user.

.. NOTE::

   TODO: Add link

UI Editor
---------

Using a component, provides a way to edit a model, before persisting it to the
*WSGI API*.

.. NOTE::

   TODO: Add link

UI Widget
---------

Improved component, based on a **MVC** design (unlike the UI components).
They can be directly added to a view, and can have mixins.

.. NOTE::

   TODO: Add link

UI Mixin
--------

Set of business code that can be applied to any widget.

.. NOTE::

   TODO: Add link

UI Container
------------

Component containing widgets, used to dispose them in a specific layout.
There is only one widget container, which can have different layout mixins.

.. NOTE::

   TODO: Add link

UI View
-------

Editable view which contains by default a single widget container.

.. NOTE::

   TODO: Add link

Miscellaneous
~~~~~~~~~~~~~

Context
-------

Contextual informations about an event, organized in graph. All other stored data
are referencing the associated context, for example :

 * a perfdata document reference the metric context
 * a periodic behavior reference the component or resource context
 * ...

.. NOTE::

   TODO: Add link
