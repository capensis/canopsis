.. _dev-backend-engines-scheduler:

Canopsis Task Scheduling
========================

Task scheduling in Canopsis is based on a set of engines :

 * the schedulers
 * the task handlers (one set for each type of job)

Each engine scheduler is configured to load pending jobs from database, and route
them to the correct task handler.

After receiving the job, the task handler executes it and emits a check event to
notify about the job's execution.

Nomenclature
------------

 * **task** : describes what to do
 * **job** : a scheduled task
 * **scheduler** : engine which is in charge of routing pending jobs
 * **task handler** : engine which execute the task

Pending Jobs
------------

At each *beat* in one of the scheduler engine, jobs are loaded from the database.
For each job, we decide which type of task handler will have to execute the task.

If you send a job record directly to the scheduler queue, one of the engine will
route this job to the correct task handler. This is used to ask a job's execution
without caring about the schedule.

Task Handler
------------

A task handler is a function, decorated with ``task_handler`` (returning a
function that can be registered as a task).

This function returns a tuple composed of the execution state, and the output.

For example :

.. code-block:: python

   from canopsis.task import register_task
   from canopsis.engines.core import task_handler


   @register_task
   @task_handler
   def task_processing(engine, logger, job, **params):
       return (0, 'OK')

Here is, basically, what is done by the decorator ``task_handler`` :

.. code-block:: python

   state, output = task_processing(engine, logger, job)

   if state == 0:
       print 'OK:', output

   elif state == 1:
       print 'Warning:', output

   elif state == 2:
       print 'Critical:', output

   elif state == 3:
       print 'Unknown Error:', output

The state and the output is used to generate a *check* event, in order to notify
Canopsis about the job's execution :

.. code-block:: python

   start = time.time()
   state, output = task_processing(engine, logger, job)
   end = time.time()

   event = {
       'timestamp': end,
       'connector': 'taskhandler',
       'connector_name': self.name,
       'event_type': 'check',
       'source_type': 'resource',
       'component': 'job',
       'resource': job['jobid'],
       'state': state,
       'state_type': 1,
       'output': output,
       'execution_time': end - start
   }

Each task handler is associated to a schema, describing the expected content of
``job``. This schema is identified by ``'task.' + self.name``.

If the job doesn't validate the schema, then it won't be executed.
