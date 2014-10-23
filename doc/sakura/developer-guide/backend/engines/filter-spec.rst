Engine Filter specification
===========================


Adding functionalities
=======================

Adding an action
-----------------

To add a new action, simply follow this prototype :

.. code-block:: python

	def a_ACTIONAME(self, event, derogation, action, name):
	    """                                                                                                                                                                                   
            Args:                                                                                                                                                                                 
                event map of the event to be modified                                                                                                                                             
                derogation map of the filter that matched the event                                                                                                                               
                action map of type action                                                                                                                                                         
                name of the rule                                                                                                                                                                 
            Returns:                                                                                                                                                                              
                ``None``                                                                                                                                                                          
            """

Only the ``pass`` and ``drop`` actions return something meaningful.
Returning something else than ``None`` will prevent any further actions and changes to be applied to the event.

The action's name and function also needs to be added  to the ``actionMap`` in ``work()``

Adding a modification action
----------------------------

To add a new mofication action, simply follow this prototype :

.. code-block:: python

    def a_ACTIONANEM(self, event, action):
        """                                                                                                                                                                                   
	Args:
	    event map of the event to be modified
	    action map of type action
	Returns:
	    ``True`` if the event was modified
	    ``False`` otherwise
        """
       
The action's name and function also needs to be added  to the ``actionMap`` in ``a_modify()``
        

Rule specification
===================

Basic Rule Structure
---------------------

Here is the basic rule structure :

.. code-block:: javascript

    {
        'name':		    // string - Rule name
        'crecord_name':     // string - cRecord name
        'description':      // string - Short description of the rule
	'mfilter':	    // dictionary - Filter to match
        'actions':          // list - Actions to apply
        'time_conditions':  // list - Optional - specific to downtime events
        'priority':	    // integer - Priority of the rule
    }

mFilter Structure
---------------------

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
---------------------

.. code-block:: javascript

    'actions': [{
        'type': 'pass|drop|override|remove|route',

	// Specific to override action
	'field':            // Field to override
	'value':	    // Value to override with

	// Specific to remove action
	// Field 'key' must be a dict or list if element is specified
	'key':		    // Field to remove
	'element':          // Element from field 'key' to remove - optional
	'met':		    // Should be specified if a metric is to be removed

	// Specific to route action
	// Field 'route' must be a string
	'route':	    // Engine to send event to

    },...]

Time Structure
---------------------

.. code-block:: javascript

	'time_conditions': [{
		'type': 'time_interval',
		'always': True|False,
		'startTs':		//Timestamp of start time
		'stopTs':		//Timestamp of stop time
		},...]

See `event_filter-Myunittest <https://github.com/capensis/canopsis/blob/NRPUIV2/sources/python/engines/test/event_filter.py>`_ for examples

Below is a simplified example on how the rules work

.. image:: ../../../_static/images/dev_engines/schema_event_filter_rule.png
