.. _admin-manage-services-init:

Canopsis Initialization Process
===============================

supervisord
-----------

Most of the Canopsis services are handled by *supervisord*, which allows to
daemonize a process and log its output.

 * the main configuration file is located at **~/etc/supervisord.conf**
 * each service has its configuration file in **~/etc/supervisord.d**

See the documentation of *supervisord* for more informations about those files.

There is a special case : **amqp2engines**.
This service is in fact a group of programs where each program is defined in
**~/etc/engines**.

LSB initscript
--------------

Some services can't be run with *supervisord*, so we have a LSB compatibility
layer. Put your initscripts in **~/etc/init.d** and let the magic happens.

Running services
----------------

The script **service** check if there is a LSB initscript present for the service
to start, if not it will check if there is a *supervisord* program configured.

The available commands are :

.. code-block:: bash

   service <service name> start
   service <service name> stop
   service <service name> restart
   service <service name> status


**Special tip for `amqp2engines` groups of processes**

You can use that command

.. code-block:: bash

   service amqp2engines* mstart|mstop|mrestart

to play with engines in parallel instead of in serie. (`m` stands for `massive`).


Configuring and running the hypervisor
--------------------------------------

The hypervisor is handled by the script **hypcontrol**.

There is 3 stages in the initialization process :

 * backend : launch the backend services (RabbitMQ, MongoDB, ...)
 * middleware : launch the connectors and engines for data processing
 * frontend : launch the webserver and other services with direct user access

Here is the default **~/etc/hypcontrol.conf** file :

.. code-block:: bash

   #!/bin/bash

   BACKEND=(rabbitmq-server mongodb redis-server)
   MIDDLEWARE=(collectd amqp2engines)
   FRONTEND=(webserver)

The content of those arrays are services names known by the **service** script.

Available commands for **hypcontrol** are the same as the ones for **service** :

.. code-block:: bash

   hypcontrol start
   hypcontrol stop
   hypcontrol restart
   hypcontrol status
