.. _dev-backend-webserver:

Canopsis Web Server
===================

The Canopsis web server is a simple **WSGI** application.

Architecture
------------

The webserver is launched with *gunicorn*, a simple HTTP server for *WSGI* applications.
It is started on port 8082, and it is recommended to use a front-end webserver like *apache2*.

The *WSGI* application is provided by the module ``canopsis.webcore.wsgi`` and gets
its configuration from ``etc/webserver.conf``.

The package ``canopsis.webcore`` is built upon the **Bottle** micro-web-framework.

Web service
-----------

A *web service* is a set of algorithms and routes to provide them via HTTP.
It will be loaded by the web server before running the *WSGI* application.

They are provided through the package ``canopsis.webcore.services`` :

.. code-block:: python

   from canopsis.common.ws import route
   from bottle import HTTPError


   def fetch_items_from_db(db, collection, conditions):
       if conditions:
           return db[collection].find()

       else:
           raise HTTPError(404, 'No items found')

   def exports(ws):
       """
       This will generate the following route for POST requests :

            /find-items/:collection

       And will except a POST parameter : conditions
       """
       @route(ws.application.post, payload=['conditions'])
       def find_items(collection, conditions=None):
           return fetch_items_from_db(ws.db, collection, conditions)

When the web service will be loaded, the ``exports()`` function will be called with
the web server instance as parameter.

See the API documentation for the package ``canopsis.common.ws`` for more informations
about the ``route`` decorator.

Authentication
--------------

Authenticating the user against the web server is a 3 stages process :

 1. submit credentials to ``/auth``
 2. check if the submitted credentials :
     * exists ;
     * are not referencing an external account ;
     * are valid.
 3. if the submitted credentials are not existent, or the account is external, redirect to ``/auth/external``

The route ``/auth/external`` is never executed. The authentication back-ends act
before the route is called, and will try to log in the user using an alternative
authentication method :

 * using the user's *authkey* ;
 * or via an LDAP directory ;
 * or via a CAS server...

Authentication back-ends are provided through the package ``canopsis.auth``.
Here is the skeleton of a simple authentication plugin :

.. code-block:: python

   from canopsis.auth.base import BaseBackend
   from bottle import request, HTTPError


   class MyBackend(BaseBackend):
       name = 'MyBackend'
       handle_logout = False

       def __init__(self, *args, **kwargs):
           super(MyBackend, self).__init__(*args, **kwargs)

       def setup_config(self, context):
           # Content of ``wsgi_params`` option from the ``route`` decorator
           self.config = context['config']

       def apply(self, callback, context):
           self.setup_config(context)

           def decorated(*args, **kwargs):
               s = self.session.get()

               if not s.get('auth_on', False):
                   username, userrecord = self.do_auth()

                   # Create session
                   if not self.install_account(username, userrecord):
                       return HTTPError(403, 'Forbidden')

               return callback(*args, **kwargs)

           return decorated

       def do_auth(self):
           username = request.params.get('username')
           password = request.params.get('password')

           userrecord = self.ws.db.user.find({'_id': username})

           if password == userrecord['password']:
               return username, userrecord

           else:
               return False, None
  
   def get_backend(ws):
       return MyBackend(ws)

The plugin **must** inherits from ``BaseBackend`` directly or indirectly.
The ``apply()`` method will apply the decorator to the route handler.

The ``install_account()`` method will create the session for the user, if he have
enough permissions.

The plugin must be named, it is a standard for *Bottle* plugins, here we set the
plugin's name to ``MyBackend``.

If your plugin handles logout (like the CAS plugin for example), you must set the
property ``handle_logout`` to ``True``.

Then, in the ``apply()`` decorator, you **must** treat the route ``/logout`` :

.. code-block:: python

       def apply(self, callback, context):
           self.setup_config(context)

           def decorated(*args, **kwargs):
               s = self.session.get()

               if request.path == '/logout':
                   self.undo_auth()

               elif not s.get('auth_on', False):
                   username, userrecord = self.do_auth()

                   # Create session
                   if not self.install_account(username, userrecord):
                       return HTTPError(403, 'Forbidden')

               return callback(*args, **kwargs)

           return decorated

Finally, the ``get_backend()`` function will be called just after the module has
been loaded by the web server, to instantiate the plugin and apply it to the *WSGI*
application. Its first argument is the web server instance, needed by the authentication
back-end.

Configuration
-------------

Let's see the default configuration file, and explain what it means :

.. code-block:: ini

    [server]
    debug=False
    enable_crossdomain_send_events=True
    root_directory=~/var/www/

    [auth]

    #providers=authkey,ldap,cas
    providers=authkey

    [session]
    cookie_expires=300
    secret=canopsis
    data_dir=~/tmp/webcore_cache

    [webservices]

    auth=1
    calendar=1
    context=1
    entities=1
    event=1
    gui=1
    i18n=1
    perfdata=1
    rest=1
    rights=1
    session=1
    topology=1

Section: server
+++++++++++++++

Here is the list of accepted options :

 * ``debug`` : a boolean value, if True, all logs will be open with a debug level
 * ``enable_crossdomain_send_events`` : if enabled, will allow the ``/event`` route to act as an event relay to another Canopsis
 * ``root_directory`` : is the absolute path to static files

Section: auth
+++++++++++++

This section is intended for authentication back-ends loading, you will list the
Python modules from ``canopsis.auth`` to load.

Section: session
++++++++++++++++

This section will configure how ``beaker`` (the *Bottle* middleware for session
handling) works :

 * ``cookie_expires`` : duration in seconds of a user's session
 * ``secret`` : key used to encrypt the session
 * ``data_dir`` : folder containing locks for ``mongodb_beaker``

Section: webservices
++++++++++++++++++++

Every options of this section is considered as a boolean. If it evaluates to ``True``,
then the Python module named after the key, will be loaded from ``canopsis.webcore.services``.

