Manage connectors module
========================

Connectors module enable on a canopsis install to manage remote connectors from a single view. It is built on top of canopsis tools such as list widget or crud mixin. It also provide many action button allowing remote interaction on the remote connector host. Possibles actions are described below.

General principles of remote connectors management
--------------------------------------------------

Connectors module once installed in canopsis enable a custom view in the web user interface called ``/userview/view.connectors``. There it is possible to create records that tells Canopsis how to reach the remote connector host with infomration such as host adress and port.

Many other information are availabe for fullfill depending on connector type.

Remote connector management in Canopsis requires the canopsis host server to be able to reach a remote connector host thanks to an ssh shared key. It is possible to acheive this by registering the canopsis public key ; once logged as canopsis with command ``sudo su - cannopsis`` copy ``cat ~/.ssh_id_dsa.pub`` into the remote host file ``~/.ssh/authorized_keys`` for a valid unix user that is allowed to manage nagios service command (root for example).


Nagios connectors module setup
------------------------------

For this tutorial, the example of the Nagios will be explained. Nagios remote connector is more precisely a broker program that is integrated to nagios workflow and it allows producing canopsis events from Nagios informations that are sent to the Canopsis server.

Let's begin by setup the remote connector control script. This script is written in bash and allow interaction with the nagios system on the remote server. It allows several actions on the nagios environnment that are :

- getState : retrieve the connector and the service (Nagios) state
- getConf : fetch the canopsis broker module configuration as json document
- setConf : transform a json document into a broker module configuration
- enable : enable the broker (uncomment if needed)
- disable : disable the broker (add comment if needed)

These commands are triggered from the canopsis user interface buttons thanks to a fabric program that runs remote commands thanks to the ssh key exchange.

These actions are available from the command script in the cat package, and it's location is the following: ``[CAT_SOURCES_LOCATION]/connector-interfaces/cat-nagios`` This script have to be copied where the program is in the path of the user on the remote machine ``/usr/bin/cat-nagios for example``. This have to be done in order to allow a fabric remote commands to find the script and run it properly.

Connector module setup
----------------------

Now it is time to setup the connector module itself. in the Canopsis CAT source folder run the follwing command : ``sudo ./build-install.sh -o connectors -d /path_to_canopsis_sources``

This command will install the connectors module that installs an engine, some webservice routes, a fabfile that enable remote commands and the webcore module enabling the connector list and action management display.

Once this setup done, it is necessary to add the engine to the canopsis engine configuration files. edit ``/opt/canopsis/etc/amqp2engines.conf`` and add the connector engine as following:

.. code-block:: bash

   next=topology

   [engine:topology]
   event_processing=canopsis.topology.process.event_processing

   #Add the following lines at the end of the last serialised engine
   next=connectors

   [engine:connectors]

refer to the `administrator guide <../administration/amqp2engines.html>`_ for more information about this kind of configuration.

Then add the engine to the engine list in the file ``/opt/canopsis/etc/supervisord.d/amqp2engines.conf`` and add ``engine-connectors`` at the end of the programs line.

At the end, just restart canopsis engines and webserver by either :

.. code-block:: bash

   service webserver restart
   service amqp2engines restart

or

.. code-block:: bash

   hypcontrol restart
