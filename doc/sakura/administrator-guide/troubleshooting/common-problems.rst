Common problems
===============

No datas retreived
------------------

Look at rabbitmq's logs in ``var/log/rabbitmq/rabbit@<hostname>.log`` for any connection problems.

* Did you changed machine's name after Canopsis installation?

If yes, then modify ``etc/rabbitmq/rabbitmq-env.conf`` this way:

.. code-block:: bash

    NODENAME=rabbit@<old hostname>

The hostname is used by RabbitMQ to store datas (mnesia database, logs…), if you change it, be sure to modify this variable.

Now you can restart canopsis:

.. code-block:: bash

    su - canopsis
    hypstart canopsis restart

MongoDB won't start
-------------------

MongoDB needs the locales to be configured in order to run properly. If you meet this line :

    [initandlisten] exception in initAndListen std::exception: locale::facet::_S_create_c_locale name not valid, terminating

Then you must install locales, example with Debian :

    # apt-get install locales
