.. _admin-troubleshooting-common:

Common problems
===============

MongoDB won't start
-------------------

You need to verify the following:

 * at least 4GB free on your *Canopsis* partition
 * no other service are listening on port 27017

RabbitMQ won't start
--------------------

Sometimes, *RabbitMQ* fails to stop properly, a child process remains : ``beam.smp``.
Check if it exists by typing:

.. code-block:: bash

   # su - canopsis
   $ ps ux | grep beam.smp
   if a process is found
   $ kill -9 <found pid>

Then try to restart *RabbitMQ*.

Another problem that can occur is wrong permissions on the *Erlang* cookie.
Check that you have the following:

.. code-block:: bash

   # su - canopsis
   $ ls -l .erlang.cookie
   -r-------- 1 canopsis canopsis 20 june  10 00:00 .erlang.cookie

If not, correct the permissions, and try to restart *RabbitMQ*.

No events found
---------------

When no events are found by *Canopsis UI*, the root cause can be:

 * an exception occurred in the engines chain:
    * check the logs in ``var/log/engines`` to verify
    * if an error is found, report it with the traceback
 * an engine failed to start, and events are stacking up in the queue:
    * check if the engine is listed in ``etc/supervisord.d/amqp2engines.conf``
    * check the engines status with: ``service amqp2engines status``
    * look into the logs to determine why it didn't start
    * report the error if found
 * the engines chain is mis-configured:
    * verify that the file ``etc/amqp2engines.conf`` chain the event to the ``eventstore``
 * the event isn't received at all:
    * using ``amqp2tty``, verify that the event is listed after its emission

Python project won't install
----------------------------

This means there is an error in the *Canopsis* Python package.
The ``setup.py`` will import ``canopsis.<project>`` in order to fetch the package's
version.

If there is an error in this module, the ``setup.py`` will fail with an ambiguous error.

Since the Python projects are installed **after** the Python libs, you can try this:

.. code-block:: bash

   $ cd sources/python
   $ PYTHONPATH="$(pwd):$PYTHONPATH" /opt/canopsis/bin/python
   >>> from canopsis import <project>
   the real error should be printed here

