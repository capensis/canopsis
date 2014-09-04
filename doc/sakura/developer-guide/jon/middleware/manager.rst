========================================
Manager: Library for managing middleware
========================================

.. module:: canopsis.middleware.manager
    :synopsis: middleware library for managing middleware.

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

A Manager has the same logic than the configurable manager, but is dedicated to middleware in resolving them thanks to URI.

The configuration of sub-middleware uses the parameter name of shape '{name}_uri'.

For example, the folowing ini file content::

   [MANAGER]
   test_uri=prototest://
   test_type=test.mom.Sender

Is dedicated to a manager which manages a middleware named ``test``, and of type ``test.mom.Sender``

In a python code, if ``manager`` is an instance of such Manager, the accessor ``manager['test']`` equals an instance of ``test.mom.Sender`` with the protocol ``prototest``.

Package contents
================

.. class:: Manager(ConfigurableManager)

   Manages middlewares like a Manager manages sub-configurables.

   Attributes are related to middlewares where the data_scope corresponds to attribute names.

   Middleware instances can be shared through a sharing_scope in the same processus. By default, this sharing_scope is the same for all Managers.

   .. data:: CONF_RESOURCE = 'middleware/manager.conf'

      conf path

   .. data:: SHARED = 'shared'

      conf shared name

   .. data:: SHARING_SCOPE = 'sharing_scope'

      conf sharing scope name
   .. data:: AUTO_CONNECT = 'auto_connect'

      conf auto connect name
   .. data:: DATA_SCOPE = 'data_scope'

      configuration data scope name

   .. data:: CATEGORY = 'MANAGER'

      middleware manager

   .. data:: MIDDLEWARE_SUFFIX = '_uri'

      middleware attribute suffix

   .. data:: __MIDDLEWARES__ = {}

      Global dict of {sharing_scope: {protocol: {data_type: {data_scope: middleware}}}}


   .. class:: Error(Exception)

      Handle Manager errors.

   .. method:: __init__(shared=True, sharing_scope=None, auto_connect=True,data_scope=None, *args, **kwargs)

      :param shared: sub-middleware shared usage (default:True)
      :type shared: bool

      :param sharing_scope: sub-middleware sharing scope usage (default:None)
      :type sharing_scope: object

      :param auto_connect: sub-middleware auto connect (default:True)
      :type auto_connect: bool

      :param data_scope: sub-middleware data_scope property (default:None)
      :type data_scope: str

   .. property:: shared

   .. property:: sharing_scope

   .. property:: auto_connect

   .. property:: data_scope

   .. method:: get_middleware(protocol, data_type=None, data_scope=None, auto_connect=None, shared=None, sharing_scope=None, *args, **kwargs)

      Load a middleware related to input uri.

      If shared, the result instance is shared among sharing_scope, protocol, data_type and data_scope.

      :param protocol: protocol to use
      :type protocol: str

      :param data_type: data type to use
      :type data_type: str

      :param data_scope: data scope to use
      :type data_scope: str

      :param auto_connect: middleware auto_connect parameter
      :type auto_connect: bool

      :param shared: if True, the result is a shared middleware instance among managers of the same class. If None, use self.shared.
      :type shared: bool

      :param sharing_scope: scope sharing
      :type sharing_scope: bool

      :return: middleware instance corresponding to input uri and data_scope.
      :rtype: Middleware

   .. method:: get_middleware_by_uri(uri, auto_connect=None, shared=None, sharing_scope=None, *args, **kwargs)

      Load a middleware related to input uri.

      If shared, the result instance is shared among same middleware type and self class type.

      :param uri: middleware uri
      :type uri: str

      :param auto_connect: middleware auto_connect parameter
      :type auto_connect: bool

      :param shared: if True, the result is a shared middleware instance among managers of the same class. If None, use self.shared.
      :type shared: bool

      :param sharing_scope: scope sharing
      :type sharing_scope: bool

      :return: middleware instance corresponding to the input uri.
      :rtype: Middleware
