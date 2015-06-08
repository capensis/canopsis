Canopsis on boot
================

For start Canopsis on boot you can use ``/opt/canopsis/etc/init.d/canopsis`` as init script.

Create symlink

.. code-block:: bash

    ln -s /opt/canopsis/etc/init.d/canopsis /etc/init.d/

RHEL and CentOS
---------------

Add script to startup process

.. code-block:: bash

    chkconfig --add canopsis

Check configuration

.. code-block:: bash

    chkconfig | grep canopsis

Debian and Ubuntu
-----------------

Add script to startup process

.. code-block:: bash

    update-rc.d canopsis defaults

Control daemons
---------------

Now you can use init script for ``start``, ``stop``, ``restart`` and ``status``

.. code-block:: text

    Usage: /etc/init.d/canopsis {start|stop|restart|status}

    # service canopsis status
    State of canopsis:
    amqp2engines                     RUNNING    pid 21469, uptime 0:23:02
    apsd                             RUNNING    pid 21489, uptime 0:22:55
    celeryd                          RUNNING    pid 21479, uptime 0:22:58
    collectd                         RUNNING    pid 21559, uptime 0:22:43
    rabbitmq-server                  RUNNING    pid 20918, uptime 0:24:12
    webserver                        RUNNING    pid 21530, uptime 0:22:48
    websocket                        RUNNING    pid 21500, uptime 0:22:52

