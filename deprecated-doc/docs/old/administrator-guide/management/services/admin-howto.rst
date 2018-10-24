.. _admin-manage-services-howto:

How to find common informations
===============================

How to find Canopsis commit ?
-----------------------------

Once canopsis installed, it is possible to get informations about current build. When the webserver is started, go to the link

``http://your.canopsis.dns:8082/en/static/canopsis/canopsis-meta.json``

.. code-block:: javascript

   {
      "build-date": "Thu Feb 19 08:28:48 UTC 2015",
      "build-timestamp": 1424334528,
      "git-commit": "202396bd8b2d200938cc353dccf590f3d6c2422f"
   }

This is only available with the Ansible installation way.  

If you want to know which release you are executing, please have a look at 

.. code-block:: javascript

   cat ~/etc/canopsis-version 


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

	[http://IP_CANOPSIS_SERVER:15672](http://IP_CANOPSIS_SERVER:15672)

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
