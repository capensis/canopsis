Metric specification
====================

Once stored in the ``perfdata2`` MongoDB collection, the metric's informations have
the following specification :

.. code-block:: javascript

	{
		'lts':    // last update as a UNIX timestamp
		'lv':     // last metric's value
		'me':     // metric's name
		'co':     // component who owns the metric
		're':     // resource who owns the metric (optional)
		't':      // metric's type
		'tg':     // event's tags
		'_id':    // node ID (hash of the concatenation of 'co', 're' and 'me')
	}

