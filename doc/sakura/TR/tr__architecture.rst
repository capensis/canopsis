.. _TR__Architecture:

===============================
Canopsis Technical Architecture
===============================

This document specifies the Canopsis Architecture, describing the toolchain and
deployment.

References
==========

 - :ref:`FR__Architecture <FR__Architecture>`
 - :ref:`FR__Engine <FR__Engine>`
 - :ref:`TR__Package <TR__Package>`

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/09/01", "0.3", "Update references", ""
   "David Delassus", "2015/09/01", "0.2", "Rename document", ""
   "David Delassus", "2015/08/02", "0.1", "Document creation", ""

Contents
========

Choosing the toolchain
----------------------

Data source
~~~~~~~~~~~

Data sources will be developed on our own, see the :ref:`connector <FR__Connector>`
reference for more informations.

Canopsis have two data sources deployed by default :

 - CollectD : gather informations about the host (CPU/MEM usage)
 - :ref:`engines <FR__Engine>` : gather informations about their performance

Messaging Queue System
~~~~~~~~~~~~~~~~~~~~~~

RabbitMQ is a software implementing the **A**dvanced **M**essaging **Q**ueue **P**rotocol.
Thus, it is providing the following features :

 - exchange and types of exchange (answering the needs for consumers subscriptions)
 - queues and message acknowledgment
 - authentication via credentials or SSL
 - support for high availability and load-balancing

Data Storage
~~~~~~~~~~~~

One of the only database system that was supporting JSON data format was MongoDB
back then.

It also provides hierarchical indexing and in-memory cache for your working dataset.

Scaling the storage is made easy by the MongoDB command-line and provide some useful
ways to tweak and enlarge the working dataset.

Data Exposure
~~~~~~~~~~~~~

Data will be exposed to the user via a REST API. Since Canopsis is written in
Python, it is common sense to use the WSGI standard.

gunicorn is a simple webserver transmitting HTTP requests to the WSGI application.

The support for authentication and external authentication MUST be handled by the
WSGI application, since gunicorn just transmit requests without touching them.

Building and Installing
-----------------------

Build System
~~~~~~~~~~~~

_ package:

The build-system is actually a bunch of scripts-shell describing :

 - how to build a Canopsis package
 - how to install/update/remove a Canopsis package

The entry point of the whole process is the ``build-install.sh`` script :

.. figure:: ../_static/images/architecture/buildinstall.png

Running services
~~~~~~~~~~~~~~~~

.. _service:

In Canopsis, the software stack we rely on is distributed as a set of services :

 - mongodb
 - rabbitmq-server
 - collectd
 - amqp2engines
 - webserver

All of those services are managed by supervisord, who's in charge of loading them,
logging their output, restarting them if they suddenly stop, ...

supervisord configuration is located at ``~canopsis/supervisord.conf`` ans launchers
configuration are stored in ``~canopsis/etc/supervisord.d``.

There is a special case for engines, which have their configuration in ``~canopsis/etc/engines``.

.. figure:: ../_static/images/architecture/supervisord.png

**NB:** A command ``service`` is provided which is used to start/stop/... services.
It ensures that supervisord is started.

**NB2:** A command ``hypcontrol`` is also provided which is used to start/stop the
whole Canopsis system. The services which are run are read from ``~canopsis/etc/hypcontrol.conf``, they are categorized in 3 sections and launched in parallel.

Deploying data sources
~~~~~~~~~~~~~~~~~~~~~~

CollectD is built inside the Canopsis environment with the AMQP plugin enabled.
Its configuration is located at ``~canopsis/etc/collectd.conf``.

The Canopsis package in charge of this is ``collectd`` and ``collectd-libs``.

Engines as data sources are not configurable, there is no more requirements for
deployment.

All other data sources MUST be distributed with their own deployment process.

Deploying messaging queue system
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

RabbitMQ, depending on Erlang, is also built inside the Canopsis environment.
Its configuration is located in the folder ``~canopsis/etc/rabbitmq`` and its
logs are written in ``~canopsis/var/log/rabbitmq/``.

The file ``~canopsis/var/log/rabbitmq-server.log`` is the logfile for our custom
launcher ``rabbitmq-server-wrapper``, which writes nothing interesting for debug.

For the client part, the file ``etc/amqp.conf`` is used to configure the old messaging
implementation.

Deploying database
~~~~~~~~~~~~~~~~~~

MongoDB binaries are distributed with the Canopsis environment.
The configuration is found at ``~canopsis/etc/mongodb.conf`` and logs are written
to ``~canopsis/var/log/mongodb.log``.

It needs at least 20GB of free disk space to preallocate database files, otherwise
it won't start.

MongoDB tries to fit the working set into RAM. If the whole data occupies 10GB and
only 1GB of data is accessed regularly and its index is also sized at 1GB, then
the working set is 2GB and will be the RAM requirement for MongoDB.

Deploying data exposure
~~~~~~~~~~~~~~~~~~~~~~~

The webserver is configured in two files :

 * ``~canopsis/etc/webserver.conf`` : configures the WSGI application ran with gunicorn
 * ``~canopsis/etc/supervisord.d/webserver.conf`` : contains the gunicorn command ran by supervisord

In order to change the listened port, you'll have to modify the call to gunicorn.
In order to change available webservices, you'll have to modify the general configuration.
