* **GET** ``/auth`` parameters

.. code-block:: javascript

	{
		'login':     // username
		'password':  // password
		'crypted':   // if True, password is encrypted using 'CRYPT' method
		'shadow':    // if True, password is encrypted using 'SHA1' method
		// if both crypted and shadow are False, password isn't encrypted
	}