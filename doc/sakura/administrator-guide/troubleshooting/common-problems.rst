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


CentOS 7 python libs install
----------------------------

At build install a bug may appear

.. code-block:: bash

   /opt/canopsis/include/python2.7/modsupport.h:27:1: erreur: ‘PyArg_ParseTuple’ is an unrecognized format function type [-Werror=format=]
   PyAPI_FUNC(int) PyArg_ParseTuple(PyObject *, const char *, ...) Py_FORMAT_PARSETUPLE(PyArg_ParseTuple, 2, 3);

The way to solve it:

.. code-block:: bash

   sed -i "/define HAVE_ATTRIBUTE_FORMAT_PARSETUPLE/ s#1#0#g" /opt/canopsis/include/python2.7/pyconfig.h

This fix could have unknown side effect that we have not met until now.

