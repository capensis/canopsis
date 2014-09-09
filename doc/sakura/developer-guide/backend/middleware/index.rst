===============================================
Middleware: Library for managing data exchanges
===============================================

.. module:: canopsis.middleware
    :synopsis: middleware library for exchanging data of different types and different scopes.

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Objective
=========

This library aims to provide to canopsis developers to focus even more on data to exchange instead of thinking about the way to do the exchange.

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

The API is generic in respecting data exchange paradigm common properties, and agnostic from technologies.

Its configuration can be done with fine grains or done in a common way to middleware concerns thanks to an uri.

In a data oriented approach, it is useful to take care about the data behavior which is changing related to quality of services (QoS). For example http://www.omg.org/spec/dds4ccm/1.1/PDF/ describes a data in a pubsub middleware which can become an event or a status information depending on QoS values.

Therefore, a middleware configuration depends on a ``data type``.

In addition to the `̀data typè`, a ̀̀data scope`̀ permits to specify scope of use of data.

URI
---

A Canopsis URI respects the convention in taking care only about the scheme and in using other parameters such as user name, password, host, etc. to configure middlewares.

A canopsis scheme is as follow ``protocol[-data_type[-data_scope]]`` where the ``data_type`̀ and `̀data_scope`` are optionals but the protocol is required and takes values of embedded technologies in middleware.

By convention, an URI scheme has forbidden characters such as '_', '$', etc. And in respect of this library, the character '-' is also forbidden to compose protocol, data_type and data_scope values because it is used by the library such as the separator char.

Paradigms
---------

This library is further declined into paradigms such as :

.. toctree::
   :maxdepth: 2
   :titlesonly:

   mom/index
   rpc/index
   storage/index

Package contents
================

.. toctree::
   :maxdepth: 2
   :titlesonly:

   manager
   sync

.. data:: __version__

    Current package version : 0.1

.. data:: SCHEME_SEPARATOR = '-'

    Char separation between protocol name and a data_type/data_scope in a Middleware URI.

.. data:: PROTOCOL_INDEX = 0

   protocol name index in an uri scheme

.. data:: DATA_TYPE_INDEX = 1

   data_type name index in an uri scheme

.. data:: DATA_SCOPE_INDEX = 2

   data scope name index in an uri scheme

.. data:: DEFAULT_DATA_SCOPE = 'canopsis'

   default data_scope

.. function:: parse_scheme(uri)

   Get a tuple of protocol and data_type names from input uri

   :return: (protocol, data_type, data_scope) from uri scheme
   :rtype: tuple

.. function:: get_uri(protocol, data_type=None, data_scope=None, host=None, port=None, user=None, pwd=None, path=None, parameters=None)

   Get a scheme related to input protocol, data_type and data_scope.

   :return: {protocol[-{data_type}[-{data_scope}]]}://[{user}[:{pwd}]]@{host}?localhost[:port][/{path}][?{parameters}]
   :rtype: str

.. class:: MetaMiddleware(canopsis.configuration.configurable.MetaConfigurable)

   Middleware meta class which register all middleware in a global set of middlewares, depending on their ``protocol`` name and ``data_type``.

   Middleware which want to be automatically registered may have set to True the class ``__register__`` attribute.
   Other class attributes permit to register ``protocol`` and ``data_type`` :

   - ``__protocol__`` : protocol name to register. ``canopsis`` by default.
   - ``__datatype__`` : data_type name to register. None by default (in case of data_type ommited in uri definition).

.. class:: Middleware(canopsis.configuration.Configurable)

   Multi middleware paradigm class.

   .. data:: __metaclass__ = MetaMiddleware

   .. data:: __register__ = False

      If True, automatically register this class

   .. data:: __protocol__ = 'canopsis'

      Protocol registration name if ``__register__``

   .. data:: __datatype__ = None

      Data type registration name if ``__register__``

   .. data:: CATEGORY = 'MIDDLEWARE'

      Configuration category name

   .. data:: CONF_RESOURCE = 'middleware/middleware.conf'

      Middleware conf resource (in addition to ones from the base class Configurable).

   .. data:: URI = 'uri'

      configuration uri. If not empty, then other uri parameters are avoided (protocol, data_type, host, port, path, user, pwd)

   .. data:: PROTOCOL = 'protocol'

      configuration protocol. Handled if not uri

   .. data:: DATA_TYPE = 'data_type'

      configuration data type. Handled if not uri

   .. data:: DATA_SCOPE = 'data_scope'

      configuration data scope.

   .. data:: HOST = 'host'

      configuration host. Handled if not uri

   .. data:: PORT = 'port'

      configuration port. Handled if not uri

   .. data:: PATH = 'path'

      configuration path. Handled if not uri

   .. data:: AUTO_CONNECT = 'auto_connect'

      configuration auto connect property. Tries to connect the middleware as soon as possible (after initialization or when a connection property is modified).

   .. data:: SAFE = 'safe'

      configuration safe output data property. If true, ensure than an output data operation succeed.

   .. data:: CONN_TIMEOUT = 'conn_timeout'

      configuration connection timeout property in milliseconds.

   .. data:: INPUT_TIMEOUT = 'in_timeout'

      configuration output data timeout property in milliseconds.

   .. data:: OUTPUT_TIMEOUT = 'out_timeout'

      configuration input data timeout property in milliseconds.

   .. data:: SSL = 'ssl'

      configuration ssl handling. If true, ssl_key and ssl_cert must be not None.

   .. data:: SSL_KEY = 'ssl_key'

      configuration ssl key.

   .. data:: SSL_CERT = 'ssl_cert'

      configuration ssl certificat.

   .. data:: USER = 'user'

      configuration user name. Handled if not uri

   .. data:: PWD = 'pwd'

      configuration password. Handled if not uri

   .. attribute:: conn

      Connection object

   .. method:: connect()

      Connect this middleware and return connected status.

      :return: True iif this is connected

   .. method:: _connect()

      get a new connection object.

   .. method:: _init_env(conn)

      Initialize the environment. Called if a new connection is successful.

      :param conn: newly created connection.

   .. method:: disconnect()

      Disconnect this middleware.

   .. method:: _disconnect()

      Method to implement in order to disconnect this middleware.

   .. method:: reconnect()

      Disconnect, then connect this middleware.

   .. method:: connected()

      True iif the middleware is connected

      :return: True iif self is connected
      :rtype: bool

   .. classmethod:: register_middleware(cls, protocol=None, data_type=None)

      Register the middleware class ``cls`` with input ``protocol`` and ``data_type``.

   .. staticmethod:: resolve_middleware(protocol, data_type=None)

      Get a reference to a middleware class registered by a protocol and a data_scope.

      :param protocol: protocol name
      :type protocol: str

      :param data_type: data type name
      :type data_type: str

      :return: Middleware type
      :rtype: type

      :raise: Middleware.Error if no middleware is registered related to input protocol and data_type.

   .. staticmethod:: resolve_middleware_by_uri(uri)

      Get a reference to a middleware class corresponding to input uri.

      :param uri: the uri may contains a protocol of type 'protocol' or 'protocol-data_type'.
      :type uri: str

      :return: Middleware type
      :rtype: type

      :raise: Middleware.Error if the uri is not reliable to a registered middleware.

   .. staticmethod:: get_middleware(protocol, data_type=None, *args, **kwargs)

      Instantiate the right middleware related to input protocol, data_type and specific parameters (in args and kwargs).

      :param protocol: protocol name
      :type protocol: str

      :param data_type: data type name
      :type data_type: str

      :param args: list of args given to the middleware to instantiate.
      :param kwargs: kwargs given to the middleware to instantiate.

      :return: Middleware
      :rtype: Middleware

      :raise: Middleware.Error if no middleware is registered related to input protocol and data_type.

   .. staticmethod:: get_middleware_by_uri(uri, *args, **kwargs)

      Instantiate the right middleware related to input uri.

      :param uri: the uri may contains a protocol of type 'protocol' or 'protocol-data_type'.
      :type uri: str

      :param args: list of args given to the middleware to instantiate.
      :param kwargs: kwargs given to the middleware to instantiate.

      :return: Middleware type
      :rtype: type

      :raise: Middleware.Error if the uri is not reliable to a registered middleware.
