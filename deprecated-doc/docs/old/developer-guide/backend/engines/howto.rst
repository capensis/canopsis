.. _dev-backend-engines-howto:

How to write an Engine
======================

This document describes how the **dynamic** engine is working and how to implement new tasks to process.

Introduction
------------

The **dynamic** engine is a generic engine which loads a specific algorithm
to process events and executes the beat (periodic task).

Those algorithms are specified in the configuration file ``amqp2engines.conf`` :


Configuration
-------------


.. code-block:: ini

   [engine:myengine]

   event_processing=canopsis.project.module.task_name
   beat_processing=canopsis.project.module.task_name
   next=perfdata

   # ...

These event_processing and beat_processing values are path to engines methods stored accordingly in canopsis python files.
As example the dynamic engine topology is configured with ``event_processing=canopsis.topology.process.event_processing``
that will trigger the method event_processing in the following path
``canopsis_source_folder/sources/python/topology/canopsis/topology/process.py``


The queue order for engines is described in the amqp2engines.conf file using **next** field that waits for the name of any engine. The engine order is not fixed, however some event processing in some engine may depend on what happened in a previous engine. Engine processing dependency have to be avoided and eliminated in future developements in order to make each engine entierely independant.

The engine reference have to be written to the ``etc/supervisord/amqp2engines.conf`` file by adding the reference at the end of the **programs** line. Just add engine-myengine ad the end separated by a comma with no space.

Finally, The supervisord configuration for an engine is described in an ini file that looks like the code below. This sample is the configuration for the topology engine.


.. code-block:: ini

  [program:engine-topology]

  autostart=false

  directory=%(ENV_HOME)s
  numprocs=1
  process_name=%(program_name)s-%(process_num)d

  command=engine-launcher -e dynamic -n topology -w %(process_num)d -l info

  stdout_logfile=%(ENV_HOME)s/var/log/engines/topology.log
  stderr_logfile=%(ENV_HOME)s/var/log/engines/topology.log

This is a supervisord notation where for instance it is possible to change in order to tell canopsis to start **one or many instances** of the same engine with `numproc = 1 or 2 or more`.


Develop
-------

Edit python function that have to be triggered by the engine, then process a build install again, this is almost done. Your engine should be now registered, set in the engine pipeline and available for supervisord. (The build install process can be skipped using slinks)

Now restart all engines :

.. code-block:: bash

   service amqp2engines restart

check what happens:

.. code-block:: bash

   tail -F /opt/canopsis/var/log/engines/myengine.log

For debuging purposes, it is also possible to start the engine with direct ouput and in debug mode:

(first avoid two same engines run at the same time)

.. code-block:: bash

   supervisorctl stop amqp2engines:engine-myengine-0

Then start the engine directly

.. code-block:: bash

   $ engine-launcher -e dynamic -n myengine -w 0 -l debug


Writing a task
--------------

A task is just a decorated function. The decorator ``register_task`` will
register the function in a *TaskManager*.

With the configuration, the **dynamic** engine will ask the *TaskManager* for
the registered function, in order to execute it.

For example :

.. code-block:: python

   from canopsis.task import register_task


   @register_task
   def my_awesome_work(engine, event, logger=None, **kwargs):
       """
       This prototype is the one required to process an event in
       the work() method of the dynamic engine.

       :param Engine engine: engine consuming the event.
       :param dict event: event to process.
       :param Logger logger: engine's logger
       :return: event (modified or not)
       :rtype: dict
       """

       return event

   @register_task
   def my_awesome_beat(engine, logger=None, **kwargs):
       """
       This prototype is the one required to execute a periodic task
       in the beat() method of the dynamic engine.

       :param Engine engine: engine consuming the event.
       :param Logger logger: engine's logger
       """

       pass
