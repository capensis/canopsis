* **POST** ``/perfstore/perftop`` and ``/perfstore/perftop/:start/:stop`` parameters

.. code-block:: javascript

	{
		'limit':             // number of records to return (do not impact 'total'), default is 10
		'sort':              // sort records by their metric's value, default is True
		'mfilter':           // MongoDB filter for selection
		'output':            // True if we want to calculate the most recurrent output and return it
		'time_window':       // If start and stop aren't defined, start = now - time_window, and stop = now, default is 1h converted to seconds
	}
