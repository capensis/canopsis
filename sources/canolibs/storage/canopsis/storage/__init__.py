#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

__version__ = "0.1"

__all__ = ('DataBase', 'Storage')

from canopsis.configuration import Configurable, Parameter, MetaConfigurable
from canopsis.common.utils import resolve_element


class MetaDataBase(MetaConfigurable):
    """
    Meta class for DataBase classes.
    """

    def __call__(cls, *args, **kwargs):
        """
        call super class __call__ method and check if auto_connect is True.
        In this case, call result.auto_connect
        """
        result = super(MetaDataBase, cls).__call__(*args, **kwargs)

        if result.auto_connect:
            result.connect()

        return result


class DataBase(Configurable):
    """
    Abstract class which aims to manage access to a data base.

    Related to a configuration file, it can connects to a database
    depending on several parameters like.

    :param host: db host name
    :type host: basestring
    :param port: db port
    :type port: int
    :param db: db name
    :type db: basestring
    :param auto_connect: auto connect to database when initialised.
    :type auto_connect: bool
    :param backend: default backend to use.

    It provides a DataBaseError for internal errors
    """

    __metaclass__ = MetaDataBase

    class DataBaseError(Exception):
        """
        Errors raised by the DataBase class.
        """

        pass

    CATEGORY = 'DATABASE'

    URI = 'uri'
    HOST = 'host'
    PORT = 'port'
    DB = 'db'
    AUTO_CONNECT = 'auto_connect'
    JOURNALING = 'journaling'
    SAFE = 'safe'
    WTIMEOUT = 'wtimeout'
    SSL = 'ssl'
    SSL_KEY = 'ssl_key'
    SSL_CERT = 'ssl_cert'
    USER = 'user'
    PWD = 'pwd'

    CONF_FILE = '~/etc/database.conf'

    def __init__(
        self,
        uri=None,
        host='localhost', port=0, db='canopsis', auto_connect=True,
        journaling=False, safe=False, wtimeout=100,
        ssl=False, ssl_key=None, ssl_cert=None, user=None, pwd=None,
        *args, **kwargs
    ):
        """
        :param uri: db uri
        :param host: db host name
        :param port: db port
        :param db: db name
        :param auto_connect: auto connect to database when initialised.
        :param backend: default backend to use.
        :param journaling: journaling mode enabling.
        :param safe: ensure writing data.
        :param wtimeout: writing time out in milliseconds.
        :param ssl: ssl mode
        :param ssl_key: ssl keys file.
        :param ssl_cert: ssl certification file.
        :param user: user
        :param pwd: password

        :type uri: str
        :type host: str
        :type port: int
        :type db: str
        :type auto_connect: bool
        :type backend: str
        :type journaling: bool
        :type safe: bool
        :param wtimeout: int
        :type ssl: bool
        :type ssl_key: str
        :type ssl_cert: str
        :type user: str
        :type pwd: str
        """

        super(DataBase, self).__init__(*args, **kwargs)

        # initialize instance properties with default values
        self._uri = uri
        self._host = host
        self._port = port
        self._db = db
        self._auto_connect = auto_connect
        self._journaling = journaling
        self._safe = safe
        self._wtimeout = wtimeout
        self._ssl = ssl
        self._ssl_key = ssl_key
        self._ssl_cert = ssl_cert
        self._user = user
        self._pwd = pwd

    @property
    def uri(self):
        return self._uri

    @uri.setter
    def uri(self, value):
        self._uri = value
        self.reconnect()

    @property
    def host(self):
        return self._host

    @host.setter
    def host(self, value):
        self._host = value
        self.reconnect()

    @property
    def port(self):
        return self._port

    @port.setter
    def port(self, value):
        self._port = value
        self.reconnect()

    @property
    def db(self):
        return self._db

    @db.setter
    def db(self, value):
        self._db = value
        self.reconnect()

    @property
    def auto_connect(self):
        return self._auto_connect

    @auto_connect.setter
    def auto_connect(self, value):
        self._auto_connect = value
        self.reconnect()

    @property
    def journaling(self):
        return self._journaling

    @journaling.setter
    def journaling(self, value):
        self._journaling = value
        self.reconnect()

    @property
    def safe(self):
        return self._safe

    @safe.setter
    def safe(self, value):
        self._safe = value
        self.reconnect()

    @property
    def wtimeout(self):
        return self._wtimeout

    @wtimeout.setter
    def wtimeout(self, value):
        self._wtimeout = value
        self.reconnect()

    @property
    def ssl(self):
        return self._ssl

    @ssl.setter
    def ssl(self, value):
        self._ssl = value
        self.reconnect()

    @property
    def ssl_key(self):
        return self._ssl_key

    @ssl_key.setter
    def ssl_key(self, value):
        self._ssl_key = value
        self.reconnect()

    @property
    def ssl_cert(self):
        return self._ssl_cert

    @ssl_cert.setter
    def ssl_cert(self, value):
        self._ssl_cert = value
        self.reconnect()

    @property
    def user(self):
        return self._user

    @user.setter
    def user(self, value):
        self._user = value
        self.reconnect()

    @property
    def pwd(self):
        return self._pwd

    @pwd.setter
    def pwd(self, value):
        self._pwd = value
        self.reconnect()

    def connect(self, *args, **kwargs):
        """
        Connect this database.

        .. seealso:: disconnect(self), connected(self), reconnect(self)
        """

        raise NotImplementedError()

    def disconnect(self, *args, **kwargs):
        """
        Disconnect this database.

        .. seealso:: connect(self), connected(self), reconnect(self)
        """

        raise NotImplementedError()

    def connected(self, *args, **kwargs):
        """
        :returns: True if this is connected.
        """

        raise NotImplementedError()

    def drop(self, table=None, *args, **kwargs):
        """
        Drop related all tables or one table if given.

        :param table: table to drop
        :type table: str

        :return: True if dropped
        :rtype: bool
        """

        raise NotImplementedError()

    def size(self, table=None, criteria=None, *args, **kwargs):
        """
        Get database size in Bytes

        :param table: table from where get data size
        :type table: str

        :param criteria: dictionary of field/value which correspond to
            elements to get size.
        :type criteria: dict

        :return: database size in Bytes of elements if criteria is not None,
            else all storage size.
        :rtype: number
        """

        raise NotImplementedError()

    def reconnect(self, *args, **kwargs):
        """
        Try to reconnect and returns connection result

        :return: True if connected
        :rtype: bool
        """

        result = False

        try:
            self.disconnect()
        except Exception:
            pass
        else:
            try:
                result = self.connect()
            except Exception:
                pass

        return result

    def _configure(self, unified_conf, *args, **kwargs):

        super(DataBase, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        reconnect = False

        db_properties = (parameter.name for parameter
            in self.conf[DataBase.CATEGORY])

        for db_property in db_properties:
            updated_property = self._update_property(
                unified_conf=unified_conf,
                param_name=db_property,
                public_property=False)
            if updated_property:
                reconnect = True

        if reconnect and self.auto_connect:
            self.reconnect()

    def _get_conf_files(self, *args, **kwargs):

        result = super(DataBase, self)._get_conf_files(*args, **kwargs)

        result.append(DataBase.CONF_FILE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(DataBase, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=DataBase.CATEGORY,
            new_content=(
                Parameter(DataBase.URI, self.uri),
                Parameter(DataBase.HOST, self.host),
                Parameter(DataBase.PORT, self.port, int),
                Parameter(DataBase.DB, self.db),
                Parameter(
                    DataBase.AUTO_CONNECT, self.auto_connect, Parameter.bool),
                Parameter(DataBase.SAFE, self.safe, Parameter.bool),
                Parameter(DataBase.WTIMEOUT, self.wtimeout, int),
                Parameter(DataBase.SSL, self.ssl, Parameter.bool),
                Parameter(DataBase.SSL_KEY, self.ssl_key),
                Parameter(DataBase.SSL_CERT, self.ssl_cert),
                Parameter(DataBase.USER, self.user),
                Parameter(DataBase.PWD, self.pwd)))

        return result


class Storage(DataBase):
    """
    Manage different kind of storages by data_type.

    For example, perfdata and context are two data types.
    """

    DATA_TYPE = 'data_type'

    ASC = 1  # ASC order
    DESC = -1  # DESC order

    class StorageError(Exception):
        """
        Handle Storage errors
        """
        pass

    def __init__(self, data_type, *args, **kwargs):

        super(Storage, self).__init__(*args, **kwargs)

        self.data_type = data_type

    def bool_compare_and_swap(self, _id, oldvalue, newvalue):
        """
        Performs an atomic compare_and_swap operation on database related to \
        input _id.

        :remarks: this method is not atomic

        :returns: True if the swamp succeed
        """
        raise NotImplementedError()

    def val_compare_and_swap(self, _id, oldvalue, newvalue):
        """
        Performs an atomic val_compare_and_swap operation on database related \
        to input _id, oldvalue and newvalue.

        :remarks: this method is not atomic

        :returns: True if the comparison succeed
        """
        raise NotImplementedError()

    def get_elements(
        self, ids=None, limit=0, skip=0, sort=None, *args, **kwargs
    ):
        """
        Get a list of elements where id are input ids

        :param ids: element ids to get if not None
        :type ids: list of str

        :param limit: max number of elements to get
        :type limit: int

        :param skip: first element index among searched list
        :type skip: int

        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order
        :type sort: list of {(str, {ASC, DESC}}), or str}

        :return: input id elements.
        :rtype: list of dict
        """

        raise NotImplementedError()

    def remove_elements(self, ids, *args, **kwargs):
        """
        Remove elements identified by the unique input ids

        :param ids: ids of elements to delete
        :type ids: list of str
        """

        raise NotImplementedError()

    def put_element(self, _id, element, *args, **kwargs):
        """
        Put an element identified by input id

        :param id: element id to update
        :type id: str

        :param element: element to put (couples of field (name,value))
        :type element: dict

        :return: True if updated
        :rtype: bool
        """

        raise NotImplementedError()

    def get_table(self, *args, **kwargs):
        """
        Table name related to elf type and data_type.

        :return: table name
        :rtype: str
        """

        result = "{0}_{1}".format(
            self._get_storage_type(), self.data_type).upper()

        return result

    def copy(self, target, *args, **kwargs):
        """
        Copy self content into target storage.
        target type must implement the same class in cstorage packege as self.
        If self implements directly cstorage.Storage, we don't care about
        target type

        :param target: target storage where copy content
        :type target: same as self or any storage if type(self) is Storage
        """

        result = 0

        from cstorage.periodic import PeriodicStorage
        from cstorage.timed import TimedStorage
        from cstorage.timedtyped import TimedTypedStorage
        from cstorage.typed import TypedStorage

        storage_types = [
            PeriodicStorage, TimedStorage, TimedTypedStorage, TypedStorage]

        if not isinstance(self, storage_types):
            pass
        else:
            for storage_type in storage_types:
                if isinstance(self, storage_types):
                    if not isinstance(target, storage_types):
                        raise Storage.StorageError(
                            'Impossible to copy {0} content into {1}. \
Storage types must be of the same type.'.format(self, target))
                    else:
                        self._copy(target)

            result = -1

        return result

    def _copy(self, target, *args, **kwargs):
        """
        Called by Storage.copy(self, target) in order to ensure than target
        type is the same as self
        """

        for element in self.get_elements():
            _id = self._element_id(element)
            target.put_element(_id=_id, element=element)

        raise NotImplementedError()

    def _element_id(self, element):
        """
        Get element id related to self behavior
        """

        raise NotImplementedError()

    def _get_category(self, *args, **kwargs):
        """
        Get configuration category for self storage
        """

        result = '{0}_{1}'.format(
            self._get_storage_type().upper(), self.data_type.upper())

        return result

    def _get_storage_type(self, *args, **kwargs):
        """
        Get storage type (last_value and timed are two storage types)

        :return: storage type name
        :rtype: str
        """

        return 'storage'

    @staticmethod
    def _update_sort(sort):
        """
        Add ASC values by default if not specified in input sort.

        :param sort: sort configuration
        :type sort: list of {tuple(str, int), str}
        """

        sort[:] = [item if isinstance(item, tuple) else (item, Storage.ASC)
            for item in sort]

    @staticmethod
    def get_storage(storage_type, data_type, *args, **kwargs):

        if isinstance(storage_type, str):
            storage_type = resolve_element(storage_type)

        result = storage_type(data_type=data_type, *args, **kwargs)

        return result
