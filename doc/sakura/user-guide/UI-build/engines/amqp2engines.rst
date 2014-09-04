Configure amqp2engines
======================

*amqp2engines* is a group of process, where each process is a single engine.

The engines are launched by the script ``engine-launcher``.

Engine description
------------------

An engine has the following properties :

 * a type (the python module to load)
 * a name (must be unique)
 * an instance identifier (0, 1, 2, 3, ..., must be unique)
 * a logging level (debug, info, warning, or error)

You can run an engine with the following command line :

.. code-block:: bash

    $ engine-launcher -e <type> -n <name> -w <process id> -l info

Each engines is started automatically by *supervisord*, you can tweak the configuration
in ``etc/engines``.

amqp2engines.conf
-----------------

Every engines have a dedicated section named ``engine:<engine's name>`` in the file
``etc/amqp2engines.conf`` with the following attributes :

 * ``next`` : list of engines to pass the event after processing
 * ``next_balanced`` : if True, it will use round-robin algorithm on the ``next`` param, otherwise it will send the event to the whole list
 * ``beat_interval`` : interval in seconds between two beats
 * ``exchange_name`` : exchange on which the engine is listening (default: amq.direct)
 * ``routing_keys`` : list of routing key to listen for (default: #)

Every extra parameter will be passed as a string.

Example
-------

I want two cleaner engine on the exchange ``canopsis.events``.

Create the file ``etc/engines/cleaner-events.ini`` :

.. code-block:: ini

    [program:engine-cleaner-events]
    
    autostart=false
    
    directory=%(ENV_HOME)s
    numprocs=2
    process_name=%(program_name)s-%(process_num)d
    
    command=engine-launcher -e cleaner -n cleaner_events -w %(process_num)d -l info
    
    stdout_logfile=%(ENV_HOME)s/var/log/engines/cleaner_events.log
    stderr_logfile=%(ENV_HOME)s/var/log/engines/cleaner_events.log

Then, edit the file ``etc/amqp2engines.conf`` :

.. code-block:: ini

    [engine:cleaner_events]

    exchange=canospsis.events
    routing_keys=#
    next=event_filter

Now make sure the name ``engine-cleaner-events`` appears in the list ``programs`` in
the file ``etc/supervisord.d/amqp2engines.conf``, so the engine will be added to the
group.
