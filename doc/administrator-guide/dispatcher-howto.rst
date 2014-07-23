Dispatcher Howto
================

Overview
--------

The canopsis dispatcher (or Y) allows to broadcast rabbitmq input event to many targets. It is a combination of rabbitmq used as an event buffer with a python script that fetch events from input queue and then post them to every valid target. A valid target is defined in a configuration file and broadcast works when target connection is established. A failover system make failure recover possible.

How to install ?
----------------

1. Get dispatcher zip or tar file and extract archive or github checkout from canopsis_hac repository.

2. Install command will start dispatch process, you may setup now configuration file to complete install in one step (see below). Run command ``sudo ./install.sh`` within dispatcher directory. This will create a dispatcher user and install the dispatcher system to ``/opt/dispatcher``.

3. Edit /opt/dispatcher/etc/targets.json . This is a json file that describes where event will be routed once dispatcher daemon is started.
This file contains a 'target' key witch is a list of canopsis's rabbitmq connection information. You must specify at least the dns name of target host. optional fields are : ``port, user and vhost``.

.. code-block:: javascript

	{
		"targets": [
			{
				//simple target definition
				"host"	: "my.target.server.net"
			},
			{
				//optional/custom parameters
				"host"	: "192.168.0.246"
				"user" 	: "guest"
				"port"	: "5672"
				"vhost"	: "canopsis"
			}
		],
		"params": {
			//optional/custom parameters
			"max_failover_duration" : 30
		}
	}


In configuration file, another field allow you to customize dispatcher behavior. In **'param'** field, it is possible to set **'max_failover_duration'** value in order to tell dispatcher max time that dispatcher waits when a target is not reachable before trying again to resume event dispatch.


How to manage ?
---------------

* Once dispatcher setup, log as dispatcher user with ``sudo su - dispatcher`` command. It is possible to stop both rabbitmq and dispatcher thanks to the command:
``./configure start|stop``

* Is is also possible to only manage dispatcher activity letting rabbitmq act as an event buffer while dispatcher is out of service. This is useful when configuration file changes. Updating the ``targets.json`` configuration file should be done with following steps:
``./service stop``

* Then it is possible to edit the configuration file and restart the dispatcher as below.
``./service start``

* Ensure the service started by having a look at the log file:
``tail -F var/log/dispatcher.log``
Soon should appear lines saying **amqp Ready** and that events are sent to each target.

If last message is :  **Error while parsing configuration file for dispatcher targets**, then you may have a typo error in your json configuration file or malformed configuration information.


Behaviors
---------

* Log file is located in ``/opt/dispatcher/var/log/dispatcher.log`` and shows sent event count every 1000 event by target.

* When dispatcher is ready, it starts fetching events from it's local rabbitmq instance. Events are routed to python threads, one by target. events are then stacked in memory until they can be sent.

* Sometimes when canopsis rabbitmq instance on target host restarts for maintenance, the dispatcher loose connection to remote and starts waiting more and more time until the max_failover_duration seconds defined in the parameter file. When this limit is reached, the dispatcher tries to connect again to send in memory waiting events.

* If max_failover_duration is reached, in mempry events are dumped into a temporary file into ``/tmp`` with ``target-name_timestamp.json`` name. Every minutes, the dispatcher search for such files and if connection is restored for target, events are loaded from dump files and then sent, dump files are removed when they are restored.

