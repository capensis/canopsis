* **GET** ``/cal/:source/:interval_start/:interval_end`` parameters

This request makes a call to the route ``rest_get()`` from the webservice REST.

.. code-block:: javascript

	{
		'limit':        // number of events returned (do not impact 'total')
		'start':        // number of events to skip (do not impact 'total')
		'filter':       // MongoDB filter for events selection
		'sort':         // sort items if True
		'query':        // selection of events matching {'crecord_name': {'$regex': '.*' + query + '.*', '$options': 'i'}}
		'onlyWritable': // returns only events that are writable by current logged in user
		'noInternal':   // exclude events with 'internal' field set to True
		'ids':          // list of records id to retrieve
		'_id':          // comma separated list of records id to retrieve (override 'ids' field)
		'fields':       // list of fields to return
	}