* **GET** ``/rest/:namespace`` and ``/rest/:namespace/:ctype`` and ``/rest/:namespace/:ctype/:_id`` parameters

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

* **POST** ``/rest/:namespace/:ctype`` and ``/rest/:namespace/:ctype/:_id`` request body

.. code-block:: javascript

	[
		// list of records
		// the specification depends on namespace and ctype
	]

* **PUT** ``/rest/:namespace/:ctype`` and ``/rest/:namespace/:ctype/:_id`` request body

.. code-block:: javascript

	{
		'_id':   // record id to update (ignored if URL's _id is defined)
		'id':    // override field _id (ignored if URL's _id is defined)

		// signle record
		// the specification depends on namespace and ctype
	}

* **DELETE** ``/rest/:namespace/:ctype`` and ``/rest/:namespace/:ctype/:_id`` request body

.. code-block:: javascript

	// optional if URL's _id is defined
	[
		// list of records to delete
	]
