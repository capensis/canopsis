* **GET** ``/ui/view`` parameters

.. code-block:: javascript

	{
		'node': // we don't care about its value, but it must be present (yeah... i know...)
	}

* **DELETE** ``/ui/view`` request body

.. code-block:: javascript

	[
		{
			'_id': // view id to delete
			// other fields are the view's record fields, but are unused here
		},
		// ...
	]

* **POST** ``/ui/view`` request body

.. code-block:: javascript

	{
		// view record
	}

* **PUT** ``/ui/view`` request body

.. code-block:: javascript

	{
		// view record
	}

* **GET** ``/ui/view/export/:_id`` and ``/ui/export/object/:_id`` parameters

.. code-block:: javascript

	{
		'_id':   // id of the record to export (ignored if URL's id is defined)
	}

* **GET** ``/ui/export/objects`` parameters

.. code-block:: javascript

	{
		'ids':  // list of records id to export
	}