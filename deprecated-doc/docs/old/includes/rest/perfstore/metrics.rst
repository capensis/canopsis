* **GET** ``/perfstore/get_all_metrics`` parameters

.. code-block:: javascript

	{
		'limit':           // number of records to return (do not impact 'total'), default is 20
		'start':           // number of records to skip (do not impact 'total'), default is 0
		'filter':          // MongoDB filter for metrics selection, default is {}
		'sort':            // True to sort metrics by value, default is False
		'show_internals':  // True to show internal metrics, default is False
	}

