New Engine Format
=================

This document describes how the **dynamic** engine is working and how to
implement new tasks to process.

Introduction
------------

The **dynamic** engine is a generic engine which loads a specific algorithm
to process events and executes the beat (periodic task).

Those algorithms are specified in the configuration file ``amqp2engines.conf`` :


.. code-block:: ini

   [engine:myengine]

   event_processing=canopsis.project.module.task_name
   beat_processing=canopsis.project.module.task_name

   # ...

And then, you can start the engine :

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
