* **GET** ``/files`` and ``/files/:metaId`` parameters

.. code-block:: javascript

	{
		'as_attachment':  // True to download the file if 'metaId' is defined
	}

* **PUT** ``/files`` and ``/files/:metaId`` request body

.. code-block:: javascript

	{
		'id':         // file id (override URL's 'metaId')
		'file_name':  // new file's name
	}

* **DELETE** ``/files`` and ``/files/:metaId`` request body

.. code-block:: javascript

	// override URL's metaId
	[
		{
			'id':    // file id to delete
		}
	]