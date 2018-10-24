.. dev-spec-event:

Event specification
===================

AMQP Informations
-----------------

    * Vhost:              canopsis
    * Routing key:        <connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]
    * Exchange:           canopsis.events
    * Exchange Options:   type: "topic", durable: true, auto_delete: false
    * Content Type:       "application/json"

Basic Event Structure
---------------------

Here is the basic event structure, common to all event-types :

.. code-block:: javascript

    {
        'connector':        // Connector Type (gelf, nagios, snmp, ...)
        'connector_name':   // Connector Identifier (nagios1, nagios2, ...)
        'event_type':       // Event type (see below)
        'source_type':      // Source of event ('component', or 'resource')
        'component':        // Component's name
        'resource':         // Resource's name (only if source_type is 'resource')

        // The following is optional
        'hostgroups':       // Nagios hostgroups for component, default []
        'servicegroups':    // Nagios servicegroups for resource, default []
        'timestamp':        // UNIX timestamp for when the event  was emitted (optional: set by the server to now)

        'output':           // Message
        'long_output':      // Description
        'display_name':     // Name to display in Canopsis
        'tags':             // Tags for the event (optional, the server adds connector, connector_name, event_type, source_type, component and resource if present)

        'perf_data':        // Nagios formatted perfdata string
        'perf_data_array':  // array of metrics (see below)
    }

Event Check Structure
---------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'check',

        'state':                // Check state (0 - OK, 1 - WARNING, 2 - CRITICAL, 3 - UNKNOWN), default is 0
        'state_type':           // Check state type (0 - SOFT, 1 - HARD), default is 1
        'status':               // 0 == Ok | 1 == En cours | 2 == Furtif | 3 == Bagot | 4 == Annule
        // The following is optional
        'scheduled':            // True if the check was scheduled, False otherwise

        'check_type':           // Nagios Check Type (host or service)
        'current_attempt':      // Attempt ID for the check
        'max_attempts':         // Max attempts before sending HARD state
        'execution_time':       // Check duration
        'latency':              // Check latency (time between schedule and execution)
        'command_name':         // Check command
    }

Event Log Structure
-------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'log',

        'output':           // Becomes mandatory
        'long_output':      // Remains optional
        'display_name':     // Remains optional

        'level':            // Optional log level
        'facility':         // Optional log facility
    }

Event Acknowledgment Structure
------------------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'ack',

        'ref_rk':               // Routing Key of acknowledged event
        'author':               // Acknowledgment author
        'output':               // Acknowledgment comment
    }

Event Cancel Structure
----------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'cancel',

        'ref_rk':               // Routing Key of event
        'author':               // author
        'output':               // comment
    }

Event Undo Cancel Structure
---------------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'uncancel',

        'ref_rk':               // Routing Key of event
        'author':               // author
        'output':               // comment
    }

Event Ackremove Structure
-------------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'ackremove',

        'ref_rk':               // Routing Key of event
        'author':               // author
        'output':               // comment
    }

Event Downtime Structure
------------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'downtime',

        'author':               // Downtime author
        'output':               // Downtime comment
        'start':                // UNIX timestamp for downtime's start
        'end':                  // UNIX timestamp for downtime's end
        'duration':             // Downtime's duration
        'entry':                // Downtime's schedule date/time (as a UNIX timestamp)
        'fixed':                // Does the downtime starts at 'start' or at next check after 'start' ?
        'downtime_id':          // Downtime's identifier
    }

Event SNMP Structure
--------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'trap',
        'snmp_severity':        // SNMP severity
        'snmp_state':           // SNMP state
        'snmp_oid':             // SNMP oid
    }

Event Calendar Structure
------------------------

After defining the basic event structure, add the following fields :

.. code-block:: javascript

    {
        'event_type': 'calendar',
        'resource':                 // iCal event UID
        'start':                    // iCal event start UNIX timestamp
        'end':                      // iCal event end UNIX timestamp
        'all_day':                  // True or False
        'output':                   // iCal event title
    }

Event Perf Structure
--------------------

An event of type 'perf' will never be saved in database, it is used to send only
perfdata :

.. code-block:: javascript

    {
        'event_type': 'perf',

        'perf_data':
        'perf_data_array':
    }

See bellow for more informations about those fields.

Metrology
^^^^^^^^^

To send perfdata to Canopsis, you just need to specify one of the following fields :

.. code-block:: javascript

    {
        'perf_data':        // Performance data ("Nagios format":http://nagiosplug.sourceforge.net/developer-guidelines.html#AEN201)
        'perf_data_array':  // Array of performance data with metric's type ('GAUGE', 'DERIVE', 'COUNTER', 'ABSOLUTE'), Ex:
        [
            {'metric': 'shortterm', 'value': 0.25, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
            {'metric': 'midterm',   'value': 0.16, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' },
            {'metric': 'longterm',  'value': 0.12, 'unit': None, 'min': None, 'max': None, 'warn': None, 'crit': None, 'type': 'GAUGE' }
        ]
    }

Basic Alert Structure
---------------------

An alert is a notification of a state change after the event were saved in MongoDB,
it contains the following fields :

.. code-block:: javascript

    {
        '_id':          // MongoDB document ID
        'event_id':     // Event identifier (the routing key)
    }


Integration with Nagios/Icinga or Shinken
-----------------------------------------

The Nagios Event Broker module will send, to Canopsis, events with the following informations :

.. code-block:: javascript

    {
        'connector': 'nagios' or 'shinken'
        'event_type': 'check' or 'ack' or 'downtime'
    }

Integration with Graylog
------------------------

The GELF connector will send, to Canopsis, events with the following informations :

.. code-block:: javascript

    {
        'connector': 'gelf',
        'event_type': 'log'
    }


Integration with Cucumber (EUE)
-------------------------------

After defining the basic event structure, set the following fields as described :

.. code-block:: javascript

    {
        'event_type': 'eue',
        'connector': 'cucumber',
        'source_type': 'resource',

        'connector_name':           // Name of the bot
        'component':                // Name of the application

        'media_bin':                // Base64 encoded binary content of associated media
        'media_type':               // Media mime-type
        'media_name':               // Media name
    }

For the EUE stack, three types of messages will be published:

* Concerning the feature
* Concerning the scenario
* Concerning the step

According to the message's type, the resource's name will be :

* For the feature : ```'resource': feature_name```
* For the scenario: ```'resource': feature_name.scenario_name.localization.OS.browser```
* For the step :    ```'resource': feature_name.scenario_name.step_name.localization.OS.browser```

Message Feature structure
^^^^^^^^^^^^^^^^^^^^^^^^^

Add the following fields to your event :

.. code-block:: javascript

    {
        'type_message': 'feature',
        'description':              // Feature's description
    }

Message Scenario structure
^^^^^^^^^^^^^^^^^^^^^^^^^^

Add the following fields to your event :

.. code-block:: javascript

    {
        'type_message': 'scenario',
        'child':                    // Routing Key of feature event
        'cntxt_env':                // Environment identifier (prod, test, ...)
        'cntxt_os':                 // Environment OS
        'cntxt_browser':            // Browser type
        'cntxt_localization':       // Bot's localization
    }

Message Step structure
^^^^^^^^^^^^^^^^^^^^^^

Add the following fields to your event :

.. code-block:: javascript

    {
        'type_message': 'step',
        'child':                    // Routing Key of scenario event
    }


List of event types
-------------------

+---------------+---------------------------------------------------------------------------+
| calendar      | Used to send ICS events to Canopsis                                       |
+---------------+---------------------------------------------------------------------------+
| check         | Used to send the result of a check (from Nagios, Icinga, Shinken, ...)    |
+---------------+---------------------------------------------------------------------------+
| comment       | Used to send a comment                                                    |
+---------------+---------------------------------------------------------------------------+
| consolidation | Sent by the consolidation engine                                          |
+---------------+---------------------------------------------------------------------------+
| eue           | Used to send Cucumber informations                                        |
+---------------+---------------------------------------------------------------------------+
| log           | Used to log informations                                                  |
+---------------+---------------------------------------------------------------------------+
| perf          | Used to send perfdata only                                                |
+---------------+---------------------------------------------------------------------------+
| selector      | Sent by the selector engine                                               |
+---------------+---------------------------------------------------------------------------+
| sla           | Sent by the sla engine                                                    |
+---------------+---------------------------------------------------------------------------+
| topology      | Sent by the topology engine                                               |
+---------------+---------------------------------------------------------------------------+
| trap          | Used to send SNMP traps                                                   |
+---------------+---------------------------------------------------------------------------+
| user          | Used by user to send informations                                         |
+---------------+---------------------------------------------------------------------------+
| ack           | Used to acknowledge an alert                                              |
+---------------+---------------------------------------------------------------------------+
| downtime      | Used to schedule a downtime                                               |
+---------------+---------------------------------------------------------------------------+
| cancel        | Used to cancel an event and put it's status in cancel state.              |
|               | removes also referer event's ack if any.                                  |
+---------------+---------------------------------------------------------------------------+
| uncancel      | Used to uncancel an event. previous status is restored and ack too if any.|
+---------------+---------------------------------------------------------------------------+
| ackremove     | Used to remove an ack from an event.                                      |
|               | (ack field removed and ack collection updated)                            |
+---------------+---------------------------------------------------------------------------+

