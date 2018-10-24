.. _dev-backend-engines-filter:

Engine Filter
=============

This document describes how to modify the filter engine.

This engine is implemented in ``canopsis.event.filter``.

Adding an action
----------------

Actions are method of the class ``Filter``, with the following prototype:

.. code-block:: python

   def a_ACTIONNAME(self, engine, logger, manager, event, action, name):
       """
       :param Engine engine: engine running the task
       :param Logger logger: logger to use
       :param DBConfiguration manager: configuration manager
       :param dict event: event to process
       :param dict action: action parameters
       :param str name: action's name
       :return: True if the event was modified, False otherwise
       """

After adding your method, you have to modify the method ``apply_actions()`` in
order to add your method to the dict ``actionMap``.

Adding a modification action
----------------------------

If your action modifies the event, it must be mapped to ``self.a_modify`` in the
method ``apply_actions()``, then add your method to the ``actionMap`` dictionary
of the method ``a_modify()``.
        
Rule specification
------------------

Basic Rule Structure
~~~~~~~~~~~~~~~~~~~~

Here is the basic rule structure :

.. code-block:: javascript

   {
       'name':             // string - Rule name
       'crecord_name':     // string - cRecord name
       'description':      // string - Short description of the rule
       'mfilter':          // dictionary - Filter to match
       'actions':          // list - Actions to apply
       'time_conditions':  // list - Optional - specific to downtime events
       'priority':         // integer - Priority of the rule
       'break':            // boolean - Allow or stop the processing of further filters
   }

mFilter Structure
~~~~~~~~~~~~~~~~~

.. code-block:: javascript

    'mfilter': {
        FIELD: VALUE | {OPERATOR: VALUE | [VALUE_LIST]}
    }

With :

* ``FIELD``: a valid field of event (see event-spec)
* ``VALUE``: a value to match
* ``OPERATOR``: ``['$eq', '$ne', '$gt', '$gte', '$lt', '$lte']``
* ``VALUE_LIST``: a list of ``VALUE``


Action Structure
~~~~~~~~~~~~~~~~

.. code-block:: javascript

    'actions': [{
        'type': 'pass|drop|override|remove|route',

         // Specific to override action
         'field':    // Field to override
         'value':    // Value to override with

         // Specific to remove action
         // Field 'key' must be a dict or list if element is specified
         'key':      // Field to remove
         'element':  // Element from field 'key' to remove - optional
         'met':      // Should be specified if a metric is to be removed

         // Specific to route action
         // Field 'route' must be a string
         'route':    // Engine to send event to
    }, ...]

Time Structure
~~~~~~~~~~~~~~

.. code-block:: javascript

   'time_conditions': [{
      'type': 'time_interval',
      'always': True|False,
      'startTs':              // Timestamp of start time
      'stopTs':               // Timestamp of stop time
   },...]

Workflow
--------

.. image:: ../../../_static/images/dev_engines/schema_event_filter_rule.png
