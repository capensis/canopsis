Common problems
===============

No datas retreived
------------------

Look at rabbitmq's logs in var/log/rabbitmq/rabbit@<hostname>.log for any connection problems.

Did you changed machine's name after Canopsis installation?

If so, you can stop all canopsis services:

.. code-block:: bash
    hypcontrol stop
    pkill -9 -u canopsis

And then reinstall Canopsis the way you want (only tested with source install).
