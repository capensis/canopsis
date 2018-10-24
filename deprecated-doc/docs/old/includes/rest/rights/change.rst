* **PUT** ``/rights/:namespace/:crecord_id`` parameters

.. code-block:: javascript

	{
		// all fields are optional

		'aaa_owner':             // new owner id
		'aaa_group':             // new group owner id
		'aaa_access_owner': [
			'r',                // if the owner can read
			'w'                 // if the owner can write
		],
		'aaa_access_group': [
			'r',                // if the group owner can read
			'w'                 // if the group owner can write
		],
		'aaa_access_other': [
			'r',                // if every user can read
			'w'                 // if every user can write
		]
	}
