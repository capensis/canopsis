* **POST** ``/sendreport`` parameters

.. code-block:: javascript

	{
		'recipients':    // list of mail address who will receive the generated PDF report
		'_id':           // id of the file to send via mail
		'body':          // body of the mail that will be sent
		'subject':       // mail subject
	}