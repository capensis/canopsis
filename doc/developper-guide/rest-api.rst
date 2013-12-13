Rest API
========

.. danger:: *TODO* : improve this page ergonomy

General points
--------------

* In the following document when there is an ':' in url it means that this value is a variable that you must provide .
* Each webservice anwser with a json object with the following scheme :

.. code-block:: javascript

	{
	    'total': int,        // The total of data returned (ex: if 5 accounts were found, the total is 5).
	    'success': bool,     // If the request is a success or not. In case of failure the output message is given in data field.
	    'data': array        // Usually an array.
	}

* If an error occurred the webservice return the appropriate HTTP error with the python exception as comment.

Rest
----

*GET : /rest/:namespace/:ctype/:_id*
*GET : /rest/:namespace/:ctype*
*GET : /rest/:namespace*
Main service to get any canopsis database record.

* arguments :

.. code-block:: plain

	namespace - mongo collection used
	ctype (optional) - Canopsis type record ('account','group','view',etc ...).
	_id (optional) - Requested item id.

* GET arguments :

.. code-block:: plain

	limit (default 20) - Number of item in answer.
	start (default 0) - Start index of the answer (ex: in a list of 100 you can start from the index 40 to the specified limit).
	search (optional) - String that the object id must contain.
	filter (optional) - Mongo filter to add to general filter.
	sort (optional) - Sorting option DESC/ASC.
	query (optional) - String that the object name must contain.
	onlyWritable (optional) - Only return object writable by the account.
	ids (optional) - List of the requested items id.
	_id (optional) - The id of the requested item.


*POST : /rest/:namespace/:ctype/:_id*
*POST : /rest/:namespace/:ctype*
Main service to create any database record.

* arguments :

.. code-block:: plain

	namespace - mongo collection used
	ctype (optional) - Canopsis type record ('account','group','view',etc ...).
	_id (optional) - Item id.

* POST arguements :

.. code-block:: plain

	The full object in json form to insert in database.


*PUT : /rest/:namespace/:ctype/:_id*
*PUT : /rest/:namespace/:ctype*
Main service to update any database record.

* arguments :

.. code-block:: plain

	namespace - mongo collection used
	ctype - Canopsis type record ('account','group','view',etc ...).
	_id (optional) - Item id.

* PUT arguments :

.. code-block:: plain

	the object with its id (if not given in url) and all the field to modify


*DELETE : /rest/:namespace/:ctype/:_id*
*DELETE : /rest/:namespace/:ctype*
Main service to delete any database record.

* arguments :

.. code-block:: plain

	namespace - mongo collection used
	ctype - Canopsis type record ('account','group','view',etc ...).
	_id (optional) - Item id.

* DELETE arguments :

.. code-block:: plain

	the id of the item to removed in a json object

Account webservice
------------------

*GET : /account/me*
Use to know if you're logged and to retrieve your personnal informations (they're returned in the data field)

* arguments : None
* result : The user account in json object form


_recheck this one_
*POST : /account/setConfig/:_id*
Use to change one value of the current user record.

* arguments :
* result : The status of the request


*GET : /account/getAuthKey/:account_name*
Use to retrieve your unique authentication key

* arguments :

.. code-block:: plain

	account_name - the name of the account that you want the key

* result : the authkey in a json object


*GET : /account/getNewAuthKey/:account_name*
Use to get a new authentication key for account in the url

* arguments :

.. code-block:: plain

	account_name - the name of the account that you want a new key

* result : the new authkey in a json object

*GET : /account/:account_id*
Use to retrieve an account or the full list of account

* arguments :

.. code-block:: plain

	account_id(optionnal) - use to request only one account

* GET arguments : those arguments are optionnal

.. code-block:: plain

	limit - the maximal length of returned list (default: 20)
	start - the request start index (ex: the second page of 60 item long list is between the index 20 and 40)

* result : list of account


*POST : /account/*
Use to create or update a new account

* POST arguments :

.. code-block:: plain

	data - json object with information to build the account, if the id already exist the old account is updated with the new provided options.

* result : no result (implanted soon)


*DELETE : /account/:account_id*
Delete and account

* arguments :

.. code-block:: plain

	account_id - the id of the account to remove

* result : return an httpError 404 if not found


*POST : /account/addToGroup/:group_id/:account_id*
Add an account to a group

* arguments :

.. code-block:: plain

	The id of group/account

* result : standart json output or http error with corresponding error in output


*POST : /account/removeFromGroup/:group_id/:account_id*
Remove an account from a group

* arguments :

.. code-block:: plain

	The id of group/account

* result : standart json output or http error with corresponding error in output

Auth webservice
---------------

*GET : /auth/:login/:password*
Used to log into canopsis

* GET arguments :

.. code-block:: plain

	password - the password
	cryptedKey - set true if the given password is hashed with sha1 + hexdigest + timestamp (strongest security)
	shadow - set true if the given password is hashed with sha1 + hexdigest

* result : return http error if auth failed or log user in


*GET : /autoLogin/:key*
Log user with his authkey

* arguments :

.. code-block:: plain

	the personal authkey

* result : json object with account or http error

*GET /logout*
*GET /disconnect*
Logout the user, clean connection cookie and close session

* arguements : none
* result : json object with success

Event webservice
----------------

*POST : /event/*
*POST : /event/:routing_key*
Used to post an event, then the event is process by engines like standard supervision event.

* arguments :

.. code-block:: plain

	routing_key (optional) - you can provide the routing key in the url, or put the elements in POST form

* POST arguements :

.. code-block:: plain

	connector
	connector_name
	event_type
	source_type
	component
	resource
	state
	state_type  (default : 1)
	perf_data  (optional)
	perf_data_array  (optional)
	output  (optional)
	long_output  (optional)

File webservice
---------------

*GET : /files/:file_metaId*
*GET : /files*
Use to retrieve file. Give the whole list of files if metaId not given.

* arguments:

.. code-block:: plain

	file_metaId - This id is returned by all webservices dealing with file

* result : file, or list of file in json format

*POST : /files*
Update the name of a file.

* POST arguments :

.. code-block:: plain

	metaId - The file to update metaId.
	file_name - The new file name.

*DELETE /files/:file_metaId*
Delete a file.

* arguments :

.. code-block:: plain

	file_metaId - the file to remove metaId.

Perfstore webservice
--------------------

*POST : /perfstore/values*
*POST : /perfstore/values/:start/:stop*
Use to get metrics on the specified time

* arguments :

.. code-block:: plain

	start - the timestamp used like a beginning index
	stop - the timestamp used like an ending index

* POST arguments :

.. code-block:: plain

	nodes - list of metric nodes requested
	interval (optional) - the interval time between two points
	aggregate_method (optional) - the used method for aggregate points, mean/max/min/last/first/delta (Mean by default)
	use_window_ts (optional) - use the timestamp of the time window

*GET : /perfstore/get_all_metrics*
Used to get the full list of metrics

* GET arguments :

.. code-block:: plain

	limit - the max number of returned result (default : 20)
	start - the start index (default : 0)
	search - string that metrics must contain

Reporting webservice
--------------------

*Get : /reporting/:startTime/:stopTime/:view_name/:mail*
*Get : /reporting/:startTime/:stopTime/:view_name*
Used to launch a view export and optionally send it by mail.


* arguments :

.. code-block:: plain

	startTime - timestamp used as beggening time index.
	stopTime - timestamp used as ending time index.
	view_name - the name of the view to export.
	mail (optional) - the recipient email.

* result : json file with action success state


*POST : /sendreport*
Used to send a report file to an email recipient.

* arguments :

.. code-block:: plain

	recipients - List of recipients
	_id - Id of the file to send
	body - body of the email
	subject - subject of the email

* result : json file with action success state



*POST : /export_svg*
Webservice used by highchart in order to export graph to svg file.

* POST arguments:

.. code-block:: plain

	filename - The svg file name
	svg - SVG file

result : return svg file

Right webservice
----------------

*PUT : /rights/:namespace/:_id*
Change object owner/rights.

* arguments : 

.. code-block:: plain

	namespace - Mongo collection used.
	_id - Id of item to modify access.

* PUT arguements:

.. code-block:: plain

	aaa_owner (optional) - New item owner.
	aaa_group (optional) - New item group.
	aaa_access_owner (optional) - New w/r owner rights.
	aaa_access_group (optional) - New w/r group rights.
	aaa_access_other (optional) - New w/r other rights.


View webservice
---------------

*GET : /ui/view*
Get view tree.

* arguments : None
* result : View tree as json object.

*DELETE : /ui/view/:name*
Delete view.

* arguments :
>name - The name of the view to delete.


*POST : /ui/view*
*POST : /ui/view/:name*
Used to create/update a view.

* arguments :
>name (optional) - View name.
* POST arguments :

.. code-block:: plain

	The view as json object (with field to update in case of update).


*GET : /ui/view/export/:_id*
Used to get a json file with the specified view

* arguments :

.. code-block:: plain

	_id - Id of the view to export.

* result : The json file corresponding to the view.

Widgets webservice
------------------

*GET : /ui/widgets*
Get widget list.

* arguments : None
* result : List of widgets with their configuration


*GET : /ui/widgets.css*
Compile all widget css.

* arguments : None
* result : A compiled css file of all widget css.