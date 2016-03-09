# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from urlparse import urlparse

from canopsis.middleware.loader import Loader
from canopsis.common.utils import ensure_iterable
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.model import (
    Parameter, Configuration, Category
)
from canopsis.configuration.configurable import MetaConfigurable


SCHEME_SEPARATOR = '-'  #: char separator in uri to proto/data_type/data_scope
PROTOCOL_INDEX = 0  #: protocol name index in an uri scheme
DATA_TYPE_INDEX = 1  #: data type name index in an uri scheme
DATA_SCOPE_INDEX = 2  #: data scope name index in an uri scheme
DEFAULT_DATA_TYPE = 'default'  #: default data type
DEFAULT_DATA_SCOPE = 'canopsis'  #: default data_scope


def parse_scheme(uri):
    """Get a tuple of protocol, data_type names from input uri.

    :return: (protocol, data_type, data_scope) from uri scheme
    :rtype: tuple
    """

    result = None, None

    protocol = None
    data_type = None
    data_scope = DEFAULT_DATA_SCOPE

    parsed_url = urlparse(uri)
    protocol = parsed_url.scheme

    if protocol and SCHEME_SEPARATOR in parsed_url.scheme:
        splitted_scheme = protocol.split(SCHEME_SEPARATOR)

        if len(splitted_scheme) > PROTOCOL_INDEX:
            protocol = splitted_scheme[PROTOCOL_INDEX]

        if len(splitted_scheme) > DATA_TYPE_INDEX:
            data_type = splitted_scheme[DATA_TYPE_INDEX]

        if len(splitted_scheme) > DATA_SCOPE_INDEX:
            data_scope = splitted_scheme[DATA_SCOPE_INDEX]

    result = protocol, data_type, data_scope

    return result


def get_uri(
        protocol, data_type=None, data_scope=None, host=None, port=None,
        user=None, pwd=None, path=None, parameters=None
):
    """Get a scheme related to input protocol, data_type and data_scope.

    :return: {protocol[-{data_type}[-{data_scope}]]}://[{user}[:{pwd}]]@{host}?
        localhost[:port][/{path}][?{parameters}]
    :rtype: str
    """

    result = ''

    if not protocol:
        raise Middleware.Error('protocol must be set')

    result = protocol

    if data_type:
        result = '{0}{1}{2}{1}'.format(result, SCHEME_SEPARATOR, data_type)

    if data_scope:
        result = '%s%s%s' % (result, SCHEME_SEPARATOR, data_scope)

    result += '://'

    if user:
        result += user

        if pwd:
            result += ':%s' % pwd

        result += '@'

    if not host:
        host = 'localhost'
    result += host

    if port:
        result += ':%s' % port

    if path:
        result += '/%s' % path

    if parameters:
        result += '?%s' % parameters

    return result


class MetaMiddleware(MetaConfigurable):
    """Middleware meta class which register all middleware in a global
    set of middlewares.
    """

    def __init__(self, name, bases, attrs):

        super(MetaMiddleware, self).__init__(name, bases, attrs)

        # if the class claims to be registered
        if self.__register__:
            self.register_middleware()


class Middleware(Configurable):
    """Abstract class which aims to manage middleware.

    A middleware is a resource which connects itself to a foreign resource such
    as a database, a mom broker, etc.
    Optionnaly, it is related to a data_type:

    - pubsub dedicated to events, status, etc.
    - database dedicated to relational data, no-sql data, timed data, etc.

    And a data scope ('canopsis' by default) which permits to define a domain.

    - perfdata for managing perfdata.
    - entities for managing entities.
    - etc.

    When a class is registered, it register itself with all protocols defined
    in bases classes.

    That means the last registered middleware will be registered to all
    protocol specified in its base class hierarchy.
    """

    __metaclass__ = MetaMiddleware

    __register__ = False  #: if True, automatically register this class
    __protocol__ = 'canopsis'  #: protocol registration name if __register__
    __datatype__ = DEFAULT_DATA_TYPE
    """ data_type registration name if __register__. Default None"""

    # private class attribute which manages middlewares classes per data_type
    __MIDDLEWARES__ = {}

    # dynamic loader which is initialized during 1st call to resolve_middleware
    __LOADER__ = None

    CATEGORY = 'MIDDLEWARE'

    URI = 'uri'
    PROTOCOL = 'protocol'
    DATA_TYPE = 'data_type'
    DATA_SCOPE = 'data_scope'
    HOST = 'host'
    PORT = 'port'
    PATH = 'path'
    AUTO_CONNECT = 'auto_connect'
    SAFE = 'safe'
    CONN_TIMEOUT = 'conn_timeout'
    IN_TIMEOUT = 'in_timeout'
    OUT_TIMEOUT = 'out_timeout'
    SSL = 'ssl'
    SSL_KEY = 'ssl_key'
    SSL_CERT = 'ssl_cert'
    USER = 'user'
    PWD = 'pwd'

    CONF_RESOURCE = 'middleware/middleware.conf'

    class Error(Exception):
        """Handle Middleware erros."""

    def __init__(
            self,
            uri=None, data_type=None, data_scope=None,
            protocol=None, host='localhost', port=0, path='canopsis',
            auto_connect=True, safe=False,
            conn_timeout=20000, in_timeout=100, out_timeout=2000,
            ssl=False, ssl_key=None, ssl_cert=None, user=None, pwd=None,
            proxies=None,
            *args, **kwargs
    ):
        """
        :param str uri: middleware uri to connect to. If other uri params are
            avoided (protocol, data_type, data_scope, host, port, path, user,
            pwd).
        :param str protocol: protocol managed if uri is not given.
        :param data_type: data type (typed, timed, event, etc.) managed if not
            uri.
        :param str data_scope: data scope (perfdata, entities, etc.).
        :param str host: host name.
        :param int port: port.
        :param str path: path name (end of the uri, could contains information
            such as db, virtual_host, etc.).
        :param bool auto_connect: auto connect when (re-)configured.
        :param bool safe: ensure output data.
        :param int conn_timeout: connection timeout in milliseconds.
        :param int in_timeout: input timeout in milliseconds.
        :param int out_timeout: output timeout in milliseconds.
        :param bool ssl: ssl mode
        :param str ssl_key: ssl keys file.
        :param str ssl_cert: ssl certification file.
        :param str user: user
        :param str pwd: password
        :param dict proxies: proxies to use.
        """

        super(Middleware, self).__init__(*args, **kwargs)

        self._conn = None

        # initialize instance properties with default values
        self._uri = uri
        self._protocol = type(self).__protocol__ if protocol is None \
            else protocol
        self._data_type = type(self).__datatype__ if data_type is None \
            else data_type
        self._data_scope = DEFAULT_DATA_SCOPE if data_scope is None \
            else data_scope
        self._host = host
        self._port = port
        self._path = path
        self._auto_connect = auto_connect
        self._safe = safe
        self._conn_timeout = conn_timeout
        self._in_timeout = in_timeout
        self._out_timeout = out_timeout
        self._ssl = ssl
        self._ssl_key = ssl_key
        self._ssl_cert = ssl_cert
        self._user = user
        self._pwd = pwd
        self._proxies = proxies

    @property
    def conn(self):
        return self._conn

    @property
    def uri(self):

        result = self._get_uri()

        return result

    def _get_uri(self):
        """Get uri in constructing it from individual uri parameters."""

        result = self._uri

        # if self._uri is not resolved, generate it related to other parameters
        if not self._uri:
            result = self.host

            if self.user:
                if self.pwd:
                    result = '%s:%s@%s' % (self.user, self.pwd, result)

                else:
                    result = '%s@%s' % (self.user, result)

            if self.path:
                result = '%s/%s' % (result, self.path)

            if self.protocol:
                data_type = self.data_type if self.data_type else ''
                data_scope = self.data_scope if self.data_scope else ''
                scheme = '{0}{1}{2}{1}{3}'.format(
                    self.protocol, SCHEME_SEPARATOR, data_type, data_scope
                )

                result = '%s://%s' % (scheme, result)

        return result

    @uri.setter
    def uri(self, value):

        self._set_uri(value)

    def _set_uri(self, value):
        """Set uri in getting values from uri parameters if value is None or
        empty.
        """

        self._uri = value

        # update other properties if value is not None
        if not value:
            self.protocol, self.data_type = parse_scheme(value)

            parsed_url = urlparse(value)

            self._host = parsed_url.hostname
            self._port = parsed_url.port
            self._user = parsed_url.username
            self._password = parsed_url.password

    @property
    def protocol(self):
        return self._protocol

    @protocol.setter
    def protocol(self, value):
        self._protocol = value

    @property
    def data_type(self):
        return self._data_type

    @data_type.setter
    def data_type(self, value):
        self._data_type = value

    @property
    def data_scope(self):
        return self._data_scope

    @data_scope.setter
    def data_scope(self, value):
        self._data_scope = value

    @property
    def host(self):
        return self._host

    @host.setter
    def host(self, value):
        self._host = value

    @property
    def port(self):
        return self._port

    @port.setter
    def port(self, value):
        self._port = value

    @property
    def path(self):
        return self._path

    @path.setter
    def path(self, value):
        self._path = value

    @property
    def auto_connect(self):
        return self._auto_connect

    @auto_connect.setter
    def auto_connect(self, value):
        self._auto_connect = value

    @property
    def safe(self):
        return self._safe

    @safe.setter
    def safe(self, value):
        self._safe = value

    @property
    def conn_timeout(self):
        return self._conn_timeout

    @conn_timeout.setter
    def conn_timeout(self, value):
        self._conn_timeout = value

    @property
    def in_timeout(self):
        return self._in_timeout

    @in_timeout.setter
    def in_timeout(self, value):
        self._in_timeout = value

    @property
    def out_timeout(self):
        return self._out_timeout

    @out_timeout.setter
    def out_timeout(self, value):
        self._out_timeout = value

    @property
    def ssl(self):
        return self._ssl

    @ssl.setter
    def ssl(self, value):
        self._ssl = value

    @property
    def ssl_key(self):
        return self._ssl_key

    @ssl_key.setter
    def ssl_key(self, value):
        self._ssl_key = value

    @property
    def ssl_cert(self):
        return self._ssl_cert

    @ssl_cert.setter
    def ssl_cert(self, value):
        self._ssl_cert = value

    @property
    def user(self):
        return self._user

    @user.setter
    def user(self, value):
        self._user = value

    @property
    def pwd(self):
        return self._pwd

    @pwd.setter
    def pwd(self, value):
        self._pwd = value

    @property
    def proxies(self):
        return self._proxies

    @proxies.setter
    def proxies(self, value):
        self._proxies = value

    def connect(self):
        """Connect this database.

        .. seealso:: disconnect(self), connected(self), reconnect(self)

        :return: True if connection has succeed
        """

        if not self.connected():

            self.logger.debug('Trying to connect to %s' % self.uri)

            self._conn = self._connect()

            # initialize the environment if connection is connected
            if self.connected():
                self.logger.debug('Initialize the environment')
                self._init_env(self._conn)

            else:
                self.logger.error("Connection failure to %s" % self.uri)

        else:
            self.logger.debug('Already connected to %s' % self.uri)

        return self.connected()

    def _connect(self):
        """Protected connection which has to be implemented by specialization
        classes.

        :return: new connection object, or None if it does not exist
        """

        raise NotImplementedError()

    def _init_env(self, conn):
        """Initialize the environement related to a newly connection.

        :param conn: newly created connection.
        """

    def disconnect(self):
        """Disconnect this database.

        .. seealso:: connect(self), connected(self), reconnect(self)
        """

        if self.connected():

            self.logger.debug("Disconnect %s from %s" % (self, self.uri))

            self._disconnect()

            self._conn = None

            self.logger.debug("Disconnected from %s" % (self.uri))

        else:
            self.logger.debug("%s is already disconnected" % self)

    def _disconnect(self):
        """Method to implement in order to disconnect this middleware."""

        raise NotImplementedError()

    def connected(self):
        """
        :returns: True if this is connected.
        """

        return False

    def reconnect(self):
        """Try to reconnect and returns connection result.

        :return: True if connected
        :rtype: bool
        """

        result = False

        try:
            self.disconnect()

        except Exception as err:
            self.logger.warning(
                'Disconnection problem while attempting to reconnect %s: %s' %
                (self, err)
            )

        else:
            try:
                result = self.connect()

            except Exception as err:
                self.logger.warning(
                    'Connection problem while attempting to reconnect %s: %s' %
                    (self, err)
                )

        return result

    def restart(self, criticals, to_configure=None, *args, **kwargs):

        super(Middleware, self).restart(
            to_configure=to_configure, criticals=criticals, *args, **kwargs
        )

        is_critical = self._is_critical_category(
            category=Middleware.CATEGORY,
            criticals=criticals
        )

        if is_critical:
            if self.auto_connect:
                self.reconnect()

    def _get_conf_paths(self, *args, **kwargs):

        result = super(Middleware, self)._get_conf_paths(*args, **kwargs)

        result.append(Middleware.CONF_RESOURCE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Middleware, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Middleware.CATEGORY,
            new_content=(
                Parameter(Middleware.URI, critical=True),
                Parameter(Middleware.PROTOCOL, critical=True),
                Parameter(Middleware.DATA_TYPE, critical=True),
                Parameter(Middleware.DATA_SCOPE, critical=True),
                Parameter(Middleware.HOST, critical=True),
                Parameter(Middleware.PORT, int, critical=True),
                Parameter(
                    Middleware.AUTO_CONNECT, Parameter.bool, critical=True
                ),
                Parameter(Middleware.SAFE, Parameter.bool, critical=True),
                Parameter(Middleware.CONN_TIMEOUT, int, critical=True),
                Parameter(Middleware.IN_TIMEOUT, int, critical=True),
                Parameter(Middleware.OUT_TIMEOUT, int, critical=True),
                Parameter(Middleware.SSL, Parameter.bool, critical=True),
                Parameter(Middleware.SSL_KEY, critical=True),
                Parameter(Middleware.SSL_CERT, critical=True),
                Parameter(Middleware.USER, critical=True),
                Parameter(Middleware.PWD, critical=True)
            )
        )

        return result

    @classmethod
    def register_middleware(cls, protocol=None, data_type=None):
        """Register a middleware class with input protocol name and data_type.

        :param protocol: one or many protocols to add to protocols of input cls
        :type protocol: str or Iterable(protocol)

        :param data_type: one or many data_types to add to data_types of input
            cls.
        :type data_type: str or Iterable(data_type)
        """

        protocols = cls.get_protocols()
        data_types = ensure_iterable(cls.__datatype__, iterable=set)

        if protocol is not None:
            protocol = ensure_iterable(protocol, iterable=set)
            protocols |= protocol

        if data_type is not None:
            data_type = ensure_iterable(data_type, iterable=set)
            data_types |= data_type

        for protocol in protocols:
            _data_types = Middleware.__MIDDLEWARES__.setdefault(protocol, {})

            for datatype in data_types:
                _data_types[datatype] = cls

    @classmethod
    def get_protocols(cls):
        """Get all protocols declared in the class hierarchy of cls.

        :return: set of protocols registered in the class tree
         of input cls
        :rtype: set([str])
        """

        protocols = ensure_iterable(cls.__protocol__, iterable=set)

        for base_cls in cls.__bases__:
            if issubclass(base_cls, Middleware):
                base_protocols = base_cls.get_protocols()

                protocols |= base_protocols

        result = protocols

        return result

    @staticmethod
    def resolve_middleware(protocol, data_type=None):
        """Get a reference to a middleware class registered by a protocol and a
        data_scope.

        :param str protocol: protocol name

        :param str data_type: data type name

        :return: Middleware type
        :rtype: type

        :raise: Middleware.Error if no middleware is registered related to
            input protocol and data_type.
        """

        # call dynamic middleware library loader
        if Middleware.__LOADER__ is None:
            Middleware.__LOADER__ = Loader()

        if data_type is None:
            data_type = DEFAULT_DATA_TYPE

        # try to get protocol
        if protocol not in Middleware.__MIDDLEWARES__:
            raise Middleware.Error(
                "No protocol %s found in registered middleware classes." %
                protocol
            )

        # try to get data_type
        data_types = Middleware.__MIDDLEWARES__[protocol]

        if data_type not in data_types:
            raise Middleware.Error(
                "No data type %s found in middleware protocol %s" %
                (data_type, protocol)
            )

        result = data_types[data_type]

        return result

    @staticmethod
    def resolve_middleware_by_uri(uri):
        """Get a reference to a middleware class corresponding to input uri.

        :param uri: the uri may contains a protocol of type 'protocol' or
            'protocol-data_type'.
        :type uri: str

        :return: Middleware type
        :rtype: type

        :raise: Middleware.Error if the uri is not reliable to a registered
            middleware.
        """

        result = None

        protocol, data_type, _ = parse_scheme(uri)

        result = Middleware.resolve_middleware(
            protocol=protocol, data_type=data_type
        )

        return result

    @staticmethod
    def get_middleware(
            protocol, data_type=None, data_scope=None, auto_connect=True,
            *args, **kwargs
    ):
        """Instantiate the right middleware related to input protocol,
        data_type and specific parameters (in args and kwargs).

        :param str protocol: protocol name.

        :param str data_type: data type name.

        :param args: list of args given to the middleware to instantiate.
        :param kwargs: kwargs given to the middleware to instantiate.

        :return: Middleware.
        :rtype: Middleware

        :raise: Middleware.Error if no middleware is registered related to
            input protocol and data_type.
        """

        result = None

        if data_type is None:
            data_type = DEFAULT_DATA_TYPE

        middleware_class = Middleware.resolve_middleware(
            protocol=protocol, data_type=data_type
        )

        if middleware_class is not None:
            # instantiate a new middleware
            result = middleware_class(auto_connect=False, *args, **kwargs)
            conf = Configuration(
                Category(
                    "get_middleware",
                    Parameter(Middleware.PROTOCOL, value=protocol),
                    Parameter(Middleware.DATA_TYPE, value=data_type),
                    Parameter(Middleware.DATA_SCOPE, value=data_scope)
                )
            )

            result._auto_connect = auto_connect
            result.configure(conf=conf)

            if auto_connect:
                result.reconnect()

        return result

    @staticmethod
    def get_middleware_by_uri(uri, auto_connect=True, *args, **kwargs):
        """Instantiate the right middleware related to input uri.

        :param str uri: the uri may contains a protocol of type 'protocol' or
            'protocol-data_type' or 'protocol-data_type-data_scope.

        :param args: list of args given to the middleware to instantiate.
        :param kwargs: kwargs given to the middleware to instantiate.

        :return: Middleware type.
        :rtype: type

        :raise: Middleware.Error if the uri is not reliable to a registered
            middleware.
        """

        result = None

        protocol, data_type, data_scope = parse_scheme(uri)

        if data_scope:
            kwargs["data_scope"] = data_scope

        middleware_class = Middleware.resolve_middleware_by_uri(uri)

        # if a new middleware class has been founded
        if middleware_class is not None:
            # instantiate it without connecting it automatically
            result = middleware_class(auto_connect=False, *args, **kwargs)

            # create a configuration with protocol, data_type and data_scope
            conf = Configuration(
                Category(
                    "get_middlewaqre_by_uri",
                    Parameter(Middleware.PROTOCOL, value=protocol),
                    Parameter(Middleware.DATA_TYPE, value=data_type),
                    Parameter(Middleware.DATA_SCOPE, value=data_scope)
                )
            )

            # set auto_connect to true
            result._auto_connect = auto_connect
            # and configure the result
            result.configure(conf=conf)

            if auto_connect:
                result.reconnect()

        return result
