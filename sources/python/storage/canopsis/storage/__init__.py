# -*- coding: utf-8 -*-
# --------------------------------
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

from functools import reduce

from time import sleep

try:
    from threading import Thread
except ImportError:
    from dummy_threading import Thread

from canopsis.common.init import basestring
from canopsis.common.utils import isiterable
from canopsis.configuration.parameters import Parameter
from canopsis.middleware import Middleware


class DataBase(Middleware):
    """
    Abstract class which aims to manage access to a data base.

    Related to a configuration file, it can connects to a database
    depending on several parameters like.

    :param host: db host name
    :type host: basestring
    """

    CATEGORY = 'DATABASE'

    DB = 'db'
    JOURNALING = 'journaling'

    SHARDING = 'sharding'

    CONF_RESOURCE = 'storage/storage.conf'

    class DataBaseError(Exception):
        pass

    def __init__(
        self, db='canopsis', journaling=False, sharding=False, *args, **kwargs
    ):
        """
        :param str db: db name
        :param bool journaling: journaling mode enabling.
        :param bool sharding: db sharding mode enabling.
        """

        super(DataBase, self).__init__(*args, **kwargs)

        # initialize instance properties with default values
        self._db = db
        self._journaling = journaling
        self._sharding = sharding

    @property
    def db(self):
        return self._db

    @db.setter
    def db(self, value):
        self._db = value
        self.reconnect()

    @property
    def journaling(self):
        return self._journaling

    @journaling.setter
    def journaling(self, value):
        self._journaling = value
        self.reconnect()

    @property
    def sharding(self):
        return self._sharding

    @sharding.setter
    def sharding(self, value):
        self._sharding = value
        self.reconnect()

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

    def _get_conf_paths(self, *args, **kwargs):

        result = super(DataBase, self)._get_conf_paths(*args, **kwargs)

        result.append(DataBase.CONF_RESOURCE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(DataBase, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=DataBase.CATEGORY,
            new_content=(
                Parameter(DataBase.DB, critical=True),
                Parameter(
                    DataBase.JOURNALING, parser=Parameter.bool, critical=True),
                Parameter(
                    DataBase.SHARDING, critical=True, parser=Parameter.bool)
            )
        )

        return result

    def restart(self, criticals, to_configure=None, *args, **kwargs):

        super(DataBase, self).restart(
            to_configure=to_configure, criticals=criticals, *args, **kwargs)

        if self._is_critical_category(DataBase.CATEGORY, criticals):
            if self.auto_connect:
                self.reconnect()


class Storage(DataBase):
    """
    Manage different kind of storages by data_scope.

    For example, perfdata and context are two data types.

    Related to such data types, it is possible to specialize the storage
        related to such data type structure thanks to the data attribute.
    And for better improvements, the indexes attribute permits to specify kind
        of indexes to use even if storages are data oriented.

    For technical improvements, a storage manages a query cache for processing
        multi queries at a time (reduce use of the network). Such feature is
        enabled by the cache_size which specified the size of the cache. If 0,
        cache is disabled.
    """

    __protocol__ = 'storage'
    """register itself and all subclasses to storage protocol"""

    DATA_ID = 'id'  #: db data id

    DATA = 'data'  #: collection/table data struct

    INDEXES = 'indexes'  #: storage indexes
    CACHE_SIZE = 'cache.size'  #: query cache size to send to the server
    CACHE_ORDERED = 'cache.ordered'  #: order query if cache is used
    CACHE_TIMEOUT = 'cache.timeout'  #: timeout cache before auto executing it

    DEFAULT_CACHE_SIZE = 1000  #: default cache size
    DEFAULT_CACHE_TIMEOUT = 1  #: default cache timeout

    CATEGORY = 'STORAGE'  #: storage category

    KEY = 'key'  #: data field key name
    TYPE = 'type'  #: data field type name
    DEFAULT = 'default'  #: data field default name
    NULL = 'null'  #: data field NULL name
    FOREIGN = 'foreign'  #: data field FOREIGN name

    ASC = 1  #: ASC order
    DESC = -1  #: DESC order

    class StorageError(Exception):
        """
        Handle Storage errors
        """
        pass

    def __init__(
        self,
        indexes=None, data=None,
        cache_size=0, cache_ordered=True,
        cache_timeout=DEFAULT_CACHE_TIMEOUT,
        *args, **kwargs
    ):
        """
        :param indexes: indexes to use.
        :type indexes: list or str
        :param dict data: data structure with expected fields, keys, etc.
        :param int cache_size: query cache size.
        :param bool cache_ordered: query cache order
        :param float cache_timeout: cache timeout before automatically execute
            queries.
        """

        super(Storage, self).__init__(*args, **kwargs)

        self._indexes = [] if indexes is None else indexes

        self._data = data

        self._cache_size = cache_size
        self._cache_ordered = cache_ordered
        self._cache_timeout = cache_timeout

    @property
    def indexes(self):

        return self._indexes

    def all_indexes(self):
        """
        :return: all self indexes (concatenation of id and additional indexes),
            such as a list of list of tuple(s).
        :rtype: list
        """

        result = [[(Storage.DATA_ID, 1)]]
        if self._indexes:
            result.append(self._indexes[:])

        return result

    @indexes.setter
    def indexes(self, value):
        """
        Indexes setter

        :param value: indexes such as::
            - one name
            - one tuple of kind (name, ASC/DESC)
            - a list of tuple or name [{(name, ASC/DESC), name}* ]
        :type value: str, tuple ot list
        """

        indexes = []

        # if value is a name, transform it into a list
        if isinstance(value, basestring):
            indexes = [[(value, Storage.ASC)]]
        elif isinstance(value, tuple):  # if value is a tuple
            indexes = [[value]]
        elif isinstance(value, list):  # if value is a list
            for index in value:
                index = self._ensure_index(index)
                indexes.append(index)
        else:  # error in other cases
            raise Storage.StorageError(
                "wrong indexes value %s. str, tuple or list accepted" % value)

        self._indexes = indexes
        self.reconnect()

    @property
    def cache_size(self):
        return self._cache_size

    @cache_size.setter
    def cache_size(self, value):
        self._cache_size = value
        self._init_cache()
        self.connect()

    @property
    def cache_ordered(self):
        return self._cache_ordered

    @cache_ordered.setter
    def cache_ordered(self, value):
        self._cache_ordered = value
        self._init_cache()
        self.connect()

    @property
    def cache_timeout(self):
        return self._cache_timeout

    @cache_timeout.setter
    def cache_timeout(self, value):
        self._cache_timeout = value

    @property
    def data(self):
        return self._data

    @data.setter
    def data(self, value):
        self._data = value

    def _init_cache(self):
        """
        Initialize cache processing.
        """

        # if cache size exists
        if self._cache_size > 0:
            # initialize all cache variables in order to process it
            self._cache_count = 0  # cache count equals 0
            self._cache = self._new_cache()  # (re)new cache
            self._updated_cache = False  # set false to updated cache
            # kill previous thread if it's alive
            if hasattr(self, '_cache_thread') and self._cache_thread.isAlive():
                # set self cache size to 0 in order to stop the thread
                cache_size, self._cache_size = self._cache_size, 0
                try:
                    self._cache_thread.join()
                except RuntimeError:
                    pass
                # set the right value of cache size
                self._cache_size = cache_size
            # start a new thread
            self._cache_thread = Thread(target=self._cache_async_execution)
            self._cache_thread.start()
        else:  # nullify _cache if it exists
            if hasattr(self, '_cache'):
                del self._cache
            self._cache = None

    def _new_cache(self):
        """
        Get self cache for query.
        """

        raise NotImplementedError()

    def _process_query(
        self,
        query_op, cache_op, query_kwargs={}, cache_kwargs={}, cache=True,
        **kwargs
    ):
        """
        Execute a query or the query cache depending on values of _cache_size
            and input cache parameter.

        :param function query_op: query operation.
        :param function cache_op: query operation to cache.
        :param dict query_kwargs: query operation kwargs.
        :param dict cache_kwargs: query operation kwargs to cache.
        :param bool cache: avoid cache operation if False (True by default).

        :return: query/cache operation result.
        """

        result = None

        if cache and self._cache_size > 0:
            if cache_op is not None:
                cache_op(**{} if cache_kwargs is None else cache_kwargs)
                # increment the counter
                self._cache_count += 1
                # check for updating cache
                self._updated_cache = True
                # if cache count is greater than cache size
                if self._cache_count >= self._cache_size:
                    # execute the cache
                    result = self._execute_cache()
                    # renew the cache
                    self._cache = self._new_cache()
                    # renew cache count
                    self._cache_count = 0

        else:  # process the query operation
            if query_kwargs is not None:
                kwargs.update(query_kwargs)
            result = query_op(**kwargs)

        return result

    def _cache_async_execution(self):
        """
        Threaded method which execute the cache.
        """

        while self._cache_size > 0:
            # wait cache timeout before trying to executing it
            sleep(self._cache_timeout)
            # if cache has not been updated
            if not self._updated_cache:
                try:  # try to execute self cache
                    if self._cache_count > 0:
                        self._execute_cache()
                except Exception as e:
                    self.logger.error(
                        'Interuption of cache execution: {}'.format(e)
                    )
            else:  # mark the cache such as not updated
                self._updated_cache = False

    def __del__(self):

        self._cache_size = 0

        if hasattr(self, '_cache_thread') and self._cache_thread.isAlive():
            self._cache_thread.join()

    def _execute_cache(self):
        """
        Execute the query cache.
        """

        raise NotImplementedError()

    def _ensure_index(self, index):
        """
        Get a right index structure related to input index.

        :return: depending on index:
            - str: [(index, Storage.ASC)]
            - tuple: (index, order): [(index, order)]
            - list: [{(index, order), (index)}+]: [(index, order)+]
        """

        result = []

        if isinstance(index, basestring):  # one value
            result = [(index, Storage.ASC)]
        elif isinstance(index, tuple):  # one value with order
            result = [index]
        elif isinstance(index, list) and index:  # not empty list of indexes
            for idx in index:  # convert
                if isinstance(idx, basestring):
                    idx = (idx, Storage.ASC)
                result.append(idx)

        return result

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
        self,
        ids=None, query=None, limit=0, skip=0, sort=None, with_count=False
    ):
        """
        Get a list of elements where id are input ids

        :param ids: element ids or an element id to get if not None
        :type ids: list of str

        :param limit: max number of elements to get
        :type limit: int

        :param skip: first element index among searched list
        :type skip: int

        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order
        :type sort: list of {(str, {ASC, DESC}}), or str}

        :param bool with_count: If True (False by default), add count to the
            result

        :return: input id elements, or one element if ids is an element
            (None if this element does not exist).
        :rtype: iterable of dict or dict or NoneType
        """

        raise NotImplementedError()

    def __getitem__(self, ids):
        """
        Python shortcut to the get_elements(ids) method.
        """

        result = self.get_elements(ids=ids)

        if result is None or ids and not result:
            raise KeyError('%s not in self' % ids)

        return result

    def __contains__(self, ids):
        """
        Python shortcut to the get_elements(ids) method.
        """

        result = True

        # self does not contain ids only if a KeyError is raised
        try:
            self[ids]

        except KeyError:
            result = False

        return result

    def find_elements(
        self,
        request, limit=0, skip=0, sort=None, with_count=False, cache=True,
    ):
        """
        Find elements corresponding to input request and in taking care of
        limit, skip and sort find parameters.

        :param dict request: set of couple of (field name, field value)
        :param int limit: max number of elements to get
        :param int skip: first element index among searched list
        :param list sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order
        :param bool with_count: If True (False by default), add count to the
            result
        :param bool cache: cache query.

        :return: input request elements
        :rtype: list of objects
        """

        raise NotImplementedError()

    def remove_elements(self, ids=None, _filter=None, cache=True):
        """
        Remove elements identified by the unique input ids

        :param ids: ids of elements to delete.
        :type ids: list of str
        :param dict _filter: removing filter.
        :param Filter _filter: additional filter to use if not None.
        :param bool cache: cache query.
        """

        raise NotImplementedError()

    def __delitem__(self, ids):
        """
        Python shortcut to the remove_elements method.
        """

        return self.remove_elements(ids=ids)

    def __isub__(self, ids):
        """
        Python shortcut to the remove_elements method.
        """

        self.remove_elements(ids=ids)

    def put_element(self, _id, element, cache=True):
        """
        Put an element identified by input id

        :param str _id: element id to update
        :param dict element: element to put (couples of field (name,value))
        :param bool cache: cache query.

        :return: True if updated
        :rtype: bool
        """

        raise NotImplementedError()

    def __setitem__(self, _id, element):
        """
        Python shortcut for the put_element method.
        """

        self.put_element(_id=_id, element=element)

    def __iadd__(self, element):
        """
        Python shortcut for the put_element method.
        """

        self.put_element(element=element)

    def count_elements(self, request=None, cache=True):
        """
        Count elements corresponding to the input request

        :param dict request: request which contain set of couples (key, value)
        :param bool cache: cache query.

        :return: Number of elements corresponding to the input request
        :rtype: int
        """

        raise NotImplementedError()

    def __len__(self):
        """
        Python shortcut to the count_elements method.
        """

        return self.count_elements()

    def _find(self, cache=True, *args, **kwargs):
        """
        Find operation dedicated to technology implementation.
        :param bool cache: cache query.
        """

        raise NotImplementedError()

    def _update(self, cache=True, *args, **kwargs):
        """
        Update operation dedicated to technology implementation.
        :param bool cache: cache query.
        """

        raise NotImplementedError()

    def _remove(self, cache=True, *args, **kwargs):
        """
        Remove operation dedicated to technology implementation.
        :param bool cache: cache query.
        """

        raise NotImplementedError()

    def _insert(self, cache=True, *args, **kwargs):
        """
        Insert operation dedicated to technology implementation.
        :param bool cache: cache query.
        """

        raise NotImplementedError()

    def _count(self, cache=True, *args, **kwargs):
        """
        Count operation dedicated to technology implementation.
        :param bool cache: cache query.
        """

        raise NotImplementedError()

    def get_table(self):
        """
        Table name related to self type and data_scope.

        :return: table name
        :rtype: str
        """

        prefix = self.data_type

        if isiterable(prefix, is_str=False):
            prefix = reduce(lambda x, y: '%s_%s' % (x, y), prefix)

        result = "{0}_{1}".format(prefix, self.data_scope).lower()

        return result

    def copy(self, target):
        """
        Copy self content into target storage.
        target type must implement the same class in cstorage packege as self.
        If self implements directly cstorage.Storage, we don't care about
        target type

        :param target: target storage where copy content
        :type target: same as self or any storage if type(self) is Storage
        """

        result = 0

        from canopsis.storage import Storage
        from canopsis.storage.periodic import PeriodicStorage
        from canopsis.storage.timed import TimedStorage
        from canopsis.storage.timedtyped import TimedTypedStorage
        from canopsis.storage.typed import TypedStorage

        storage_types = [
            Storage,
            PeriodicStorage,
            TimedStorage,
            TimedTypedStorage,
            TypedStorage]

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

    def _copy(self, target):
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

        prefix = self.data_type

        if isiterable(prefix):
            prefix = reduce(lambda x, y: '%s_%s' % (x, y), prefix)

        result = '{0}_{1}'.format(prefix, self.data_scope).lower()

        return result

    def _conf(self, *args, **kwargs):

        result = super(Storage, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Storage.CATEGORY,
            new_content=(
                Parameter(Storage.INDEXES, parser=eval),
                Parameter(Storage.CACHE_SIZE, parser=int),
                Parameter(Storage.CACHE_ORDERED, parser=Parameter.bool),
                Parameter(Storage.CACHE_TIMEOUT, parser=int)
            )
        )

        return result

    @staticmethod
    def _update_sort(sort):
        """
        Add ASC values by default if not specified in input sort.

        :param sort: sort configuration
        :type sort: list of {tuple(str, int), str}
        """

        sort[:] = [item if isinstance(item, tuple) else (item, Storage.ASC)
            for item in sort]
