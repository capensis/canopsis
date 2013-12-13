Event specification
===================

AMQP Informations
-----------------

    * Vhost:              canopsis
    * Routing key:        <connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]
    * Exchange:           canopsis.events
    * Exchange Options:   type: "topic", durable: true, auto_delete: false
    * Content Type:       "application/json"

Event structure
---------------

.. code-block:: javascript

    {
      '_id':               // (reserved) MongoDB
      'event_id':          // (reserved) Original in alerts exchange
      'connector':         // Connector type (gelf, nagios, snmp, ...),
      'connector_name':    // Connector name (nagios1, nagios2 ...),
      'event_type':        // Event type (check, log, trap, ...),
      'source_type':       // Source type ('component' or 'resource'),
      'component':         // Component name,
      'resource':          // Ressource name,
      'timestamp':         // UNIX seconds timestamp (UTC),
      'state':             // State (0 (Ok), 1 (Warning), 2 (Critical), 3 (Unknown)),
      'state_type':        // State type (O (Soft), 1 (Hard)),
      'scheduled':         // (optional) True if this is a scheduled event
      'last_state_change': // (reserved) Last timestamp after state change,
      'previous_state':    // (reserved) Previous state (after change),
      'output':            // Event message,
      'long_output':       // Event long message,
      'tags':              // Event Tags (default: []),
      'display_name':      // The name to display (customization purpose)
    }

Metrology
^^^^^^^^^
.. code-block:: javascript

    'perf_data':      Performance data ("Nagios format":http://nagiosplug.sourceforge.net/developer-guidelines.html#AEN201)

or

.. code-block:: javascript

    'perf_data_array': Array of performance data with metric's type ('GAUGE', 'DERIVE', 'COUNTER', 'ABSOLUTE'), Ex:
                            [
                              {'metric': 'shortterm', 'value': 0.25, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
                              {'metric': 'midterm',   'value': 0.16, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
                              {'metric': 'longterm',  'value': 0.12, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' }
                            ]


Event Type
^^^^^^^^^^

* **check**: Result of control
* **comment**: Result's Comment
* **trap**: SNMP Trap
* **log**: Log
* **user**: User input
* **selector**:  Selector result
* **sla**:  Sla result
* **eue**:  EUE result

Nagios/Icinga
-------------

.. code-block:: javascript

    'connector':  'nagios'
    'event_type': 'check', 'ack', 'notification', 'downtime'

Extra - check
^^^^^^^^^^^^^

.. code-block:: javascript

    'check_type':       ,
    'current_attempt':  ,
    'max_attempts':     ,
    'execution_time':   ,
    'latency':          ,
    'command_name':     ,
    'address':          ,

Gelf (Graylog)
--------------

.. code-block:: javascript

    'connector':  'gelf'
    'event_type': 'log'

Extra
^^^^^
.. code-block:: javascript

    'level':    Log level
    'facility': Log facility

SNMP Traps
----------

.. code-block:: javascript

    'connector':  'snmp'
    'event_type': 'trap'

Extra
^^^^^
.. code-block:: javascript

    'snmp_severity':
    'snmp_state':
    'snmp_oid':
    'address':

Shinken
-------

.. code-block:: javascript

    'connector':  'shinken'
    'event_type': 'check', 'ack', 'notification', 'downtime'

Cucumber (EUE)
--------------

.. code-block:: javascript

    'connector': 'cucumber'
    'component': name_of_the_application
    'event_type': eue
    'connector_name': name_of_the_bot
    'source_type': resource

For the eue stack, three types of messages will be published:

* Concerning the feature
* Concerning the scenario
* Concerning the step

En fonction du type message le nom de la ressource sera:

* For the feature : ```'resource': feature_name```
* For the scenario: ```'resource': feature_name.scenario_name.localization.OS.browser```
* For the step :    ```'resource': feature_name.scenario_name.step_name.localization.OS.browser```

Global
^^^^^^

.. code-block:: javascript

    'type_message' : ( feature, scenario, step)

Pour la feature
^^^^^^^^^^^^^^^

.. code-block:: javascript

    'description': global description (propre a la fonctionnalité)
    'media_bin' : binary content of the file, encoded in base64
    'media_type': mime type
    'media_name': name of the media

Pour le scenario
^^^^^^^^^^^^^^^^

.. code-block:: javascript

    'child':               routing_key_feature
    'cntxt_env' :          (prod / test)
    'cntxt_os':            OS type
    'cntxt_browser':       browser name,
    'cntxt_localization':  localisation of the bot who plays the scenario
    'media_bin':           binary content of the file, encoded in base64
    'media_type':          mime type
    'media_name':          name of the media

Pour la step
^^^^^^^^^^^^

.. code-block:: javascript

    'child': 'routing_key_scenario'
    'media_bin': binary content of the file, encoded in base64
    'media_type': mime type
    'media_name': media_name

Calendar
--------

.. code-block:: javascript

    'event_type': 'calendar'
    'connector':  Source name
    'ressource':  event UID
    'start':  event start timestamp
    'end':  event end timestamp
    'all_day' : (True/False)
    'output' : event title
