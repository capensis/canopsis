* **GET** ``/account/:_id`` and ``/account/`` parameters

.. code-block:: javascript

	{
		'limit':       // limit the number of items returned (does not impact 'total' value), default: 20
		'start':       // do no returns the #start firsts items (does not impact 'total' value), default: 0
	}


* **POST** ``/account/`` request body

.. code-block:: javascript

	{
		'user':        // account username
	}


* **PUT** ``/account/:_id`` and ``/account/`` request body

.. code-block:: javascript

	[
		{
			'_id':         // account id (override URL's _id)
			'id':          // account id (override _id field)
			'passwd':      // new password
			'aaa_group':   // new group owner
			'groups':      // new groups (do not erase previous ones)

			// every other fields overwrite the record's ones
		},
		// ...
	]


* **DELETE** ``/account/:_id`` and ``/account/`` request body

 .. code-block:: javascript

 	// this list override the URL's _id
 	[
 		{
 			'_id':         // account id
 			'id':          // account id (override _id field)
 		},
 		// ...
 	]