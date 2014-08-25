Uninstall
=========

(execute as ``root``)

.. code-block:: bash

    pkill -9 -u canopsis
    userdel canopsis
    rm -Rf /opt/canopsis
    rm /etc/init.d/canopsis

