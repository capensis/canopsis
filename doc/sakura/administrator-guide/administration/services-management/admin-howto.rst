Admin Howto
===========

How to find BDD Size ?
----------------------

Open shell with ``Canopsis`` user and type:

.. code-block:: bash

    cps_mongostat

you can also type:

.. code-block:: bash

    du -hcs var/lib/mongodb/*


How to access to RabbitMQ UI ?
------------------------------

you could access to RabbitMQ UI with this url :

.. code-block:: text

	[http://IP_CANOPSIS_SERVER:55672](http://IP_CANOPSIS_SERVER:55672)

	* login : admin
	* mdp : admin_password

With that, you could see messages in AMQP queues.


How to list queues from RabbitMQ ?
----------------------------------

.. code-block:: bash

    rabbitmqctl  list_queues -p canopsis name messages messages_ready messages_unacknowledged

How to force db repair on MongoDB ?
-----------------------------------

.. code-block:: bash

	hypcontrol stop
	service rabbitmq-server start
	service mongodb start

	mongo -u mongo_user -p mongo_password canopsis
	> db.repairDatabase();
	> exit

	python opt/mongodb/filldb.py update
    schema2db

	hypcontrol start
