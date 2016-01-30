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

__version__ = '0.1'

__all__ = ['DataBase', 'Storage']

from functools import reduce

from time import sleep

try:
    from threading import Thread, current_thread, Lock
except ImportError:
    from dummy_threading import Thread, current_thread

from collections import Iterable

from canopsis.common.init import basestring
from canopsis.common.utils import isiterable
from canopsis.configuration.model import Parameter
from canopsis.middleware.core import Middleware


class DataBase(Middleware):
    """Abstract class which aims to manage access to a data base.

    Related to a configuration file, it can connects to a database
    depending on several parameters like.

    :param host: db host name
    :type host: basestring
    """

    CATEGORY = 'DATABASE'

    DB = 'db'  #: database name.
    JOURNALING = 'journaling'  #: journaling flag.
    SHARDING = 'sharding'  #: sharding name.
    REPLICASET = 'replicaset'  #: replication set name.
    #: retention rule : INF|{number}[ywdhm]. INF by default.
    RETENTION = 'retention'

    CONF_RESOURCE = 'storage/storage.conf'

    class DataBaseError(Exception):
        """Handle DataBase errors."""

    def __init__(
            self,
            db='canopsis', journaling=False, sharding=False, replicaset=None,
            retention=None,
            *args, **kwargs
    ):
        """
        :param str db: db name
        :param bool journaling: journaling mode enabling.
        :param bool sharding: db sharding mode enabling.
        :param str replicaset: db replicaset.
        :param str retention: retention rule.
        """

        super(DataBase, self).__init__(*args, **kwargs)

        # initialize instance properties with default values
        self._db = db
        self._journaling = journaling
        self._sharding = sharding
        self._replicaset = replicaset
        self._retention = retention

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
    def retention(self):
        """Get retention rule.

        :rtye: str"""

        return self._retention

    @retention.setter
    def retention(self, value):
        """Change of retention rule.

        :param str value: new retention rule to apply."""

        self._retention = value
        self.reconnect()

    @property
    def sharding(self):
        return self._sharding

    @sharding.setter
    def sharding(self, value):
        self._sharding = value
        self.reconnect()

    @property
    def replicaset(self):
        return self._replicaset

    @replicaset.setter
    def replicaset(self, value):
        self._replicaset = value
        self.reconnect()

    def drop(self, table=None, *args, **kwargs):
        """Drop related all tables or one table if given.

        :param table: table to drop
        :type table: str

        :return: True if dropped
        :rtype: bool
        """

        raise NotImplementedError()

    def size(self, table=None, criteria=None, *args, **kwargs):
        """Get database size in Bytes.

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
                    DataBase.JOURNALING, parser=Parameter.bool, critical=True
                ),
                Parameter(
                    DataBase.SHARDING, critical=True, parser=Parameter.bool
                ),
                Parameter(DataBase.REPLICASET, critical=True)
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
    """Manage different kind of storages by data_scope.

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
    TABLE = 'table'  #: table field name

    INDEXES = 'indexes'  #: storage indexes
    CACHE_SIZE = 'cache_size'  #: query cache size to send to the server
    CACHE_ORDERED = 'cache_ordered'  #: order query if cache is used
    CACHE_AUTOCOMMIT = 'cache_autocommit'  #: duration before auto-commit cache

    DEFAULT_CACHE_SIZE = 1000  #: default cache size
    DEFAULT_CACHE_AUTOCOMMIT = 1  #: default cache auto-commit in seconds
    DEFAULT_CACHE_ORDERED = True  #: default cache ordered

    CATEGORY = 'STORAGE'  #: storage category

    KEY = 'key'  #: data field key name
    TYPE = 'type'  #: data field type name
    DEFAULT = 'default'  #: data field default name
    NULL = 'null'  #: data field NULL name
    FOREIGN = 'foreign'  #: data field FOREIGN name

    ASC = 1  #: ASC order
    DESC = -1  #: DESC order

    class StorageError(Exception):
        """Handle Storage errors"""

    def __init__(
            self,
            indexes=None, data=None,
            cache_size=DEFAULT_CACHE_SIZE, cache_ordered=DEFAULT_CACHE_ORDERED,
            cache_autocommit=DEFAULT_CACHE_AUTOCOMMIT, table=None,
            *args, **kwargs
    ):
        """
        :param str table: default table name.
        :param indexes: indexes to use.
        :type indexes: list or str
        :param dict data: data structure with expected fields, keys, etc.
        :param int cache_size: (default 1000) query cache size.
        :param bool cache_ordered: (default True) query cache order
        :param float cache_autocommit: (default 1) duration in seconds before
            auto-commit cache if no activity. If not greater than 0, auto
            commit is deactivated.
        """

        super(Storage, self).__init__(*args, **kwargs)

        self._indexes = [] if indexes is None else indexes

        self._data = data
        self._table = table

        self._updated_cache = False
        self._cache = None
        self._cache_size = cache_size
        self._cache_count = 0
        self._cache_ordered = cache_ordered
        self._cache_autocommit = cache_autocommit
        self._cached_thread = self._parent_thread = None
        self._lock = Lock()  # lock for asynchronous autocommit

    @property
    def indexes(self):
        """Get storage indexes.

        :return: storage indexes such as a list of list of (name, direction).
        :rtype: list
        """
        return self._indexes

    def all_indexes(self):
        """
        :return: all self indexes (concatenation of id and additional indexes),
            such as a list of list of tuple(s).
        :rtype: list
        """

        result = [[(Storage.DATA_ID, 1)]]
        # add indexes from self indexes
        if self._indexes:
            result += self._indexes[:]
        # add indexes from self data
        if self._data:
            data = self._data
            # search key among data fields
            for field in data:
                value = data[field]
                if isinstance(value, dict):
                    if Storage.KEY in value:
                        index = [(field, value[Storage.KEY])]
                        result.append(index)

        return result

    @indexes.setter
    def indexes(self, value):
        """Indexes setter.

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
    def table(self):
        return self._table

    @table.setter
    def table(self, value):

        self._table = value
        self.reconnect()

    @property
    def data(self):
        return self._data

    @data.setter
    def data(self, value):
        self._data = value
        self.reconnect()

    @property
    def cache_size(self):
        return self._cache_size

    @cache_size.setter
    def cache_size(self, value):
        self._lock.acquire()
        try:
            self.execute_cache()
            self._cache_size = value
        finally:
            self._lock.release()

    @property
    def cache_ordered(self):
        return self._cache_ordered

    @cache_ordered.setter
    def cache_ordered(self, value):
        self._lock.acquire()
        try:
            self.execute_cache()
            self._cache_ordered = value
        finally:
            self._lock.release()

    @property
    def cache_autocommit(self):
        return self._cache_autocommit

    @cache_autocommit.setter
    def cache_autocommit(self, value):
        self._lock.acquire()
        try:
            self.execute_cache()
            self._cache_autocommit = value
        finally:
            self._lock.release()

    def queries_in_cache(self):
        """
        :return: number of queries in cache to commit.
        :rtype: int
        """

        return self._cache_count

    def _init_cache(self):
        """Initialize cache processing."""

        # if cache size exists
        if self._cache_size > 0:
            self._parent_thread = current_thread()
            # initialize all cache variables in order to process it
            self._cache_count = 0  # cache count equals 0
            self._cache = self._new_cache()  # (re)new cache
            self._updated_cache = False  # set false to updated cache
            # kill previous thread if it's alive
            self.halt_cache_thread()
            # start a new thread if self cache auto commit greater than 0
            if self._cache_autocommit > 0:
                self._cached_thread = Thread(
                    target=self._cache_async_execution
                )
                self._cached_thread.start()
        else:  # nullify _cache if it exists
            if hasattr(self, '_cache'):
                del self._cache
            self._cache = None

    def _new_cache(self):
        """Get self cache for query."""

        raise NotImplementedError()

    def _process_query(
            self,
            query_op, cache_op,
            query_kwargs=None, cache_kwargs=None, cache=False,
            **kwargs
    ):
        """Execute a query or the query cache depending on values of _cache_size
        and input cache parameter.

        :param function query_op: query operation.
        :param function cache_op: query operation to cache.
        :param dict query_kwargs: query operation kwargs.
        :param dict cache_kwargs: query operation kwargs to cache.
        :param bool cache: avoid cache operation if False (True by default).

        :return: query/cache operation result.
        """

        result = None

        if query_kwargs is None:
            query_kwargs = {}
        if cache_kwargs is None:
            cache_kwargs = {}

        if cache and self._cache_size > 0:
            # if self cache is None, that means thisd is the first use to cache
            if (
                    self._cached_thread is None or
                    not self._cached_thread.isAlive()
            ):
                # init cache
                self._init_cache()
            self._lock.acquire()  # avoid concurrent calls to cache execution
            try:
                if cache_op is not None:
                    cache_op(**cache_kwargs)
                    # check for updating cache
                    self._updated_cache = True
                    # increment the counter
                    self._cache_count += 1
                    # if cache count is greater than cache size
                    if self._cache_count >= self._cache_size:
                        # execute the cache
                        result = self.execute_cache()
            finally:
                self._lock.release()
        else:  # process the query operation
            if query_kwargs is not None:
                kwargs.update(query_kwargs)
            result = query_op(**kwargs)

        return result

    def _cache_async_execution(self):
        """Threaded method which execute the cache."""

        # while parent thread is alive and cache size is greater than 0
        while (
                self._parent_thread.isAlive()
                and self._cache_autocommit > 0
                and self._cache_size > 0
        ):
            # wait cache timeout before trying to executing it
            sleep(self._cache_autocommit)
            self._lock.acquire()  # avoid concurrent calls to cache execution
            try:
                # if cache has not been updated
                if not self._updated_cache:
                    # execute cache
                    self.execute_cache()
                else:  # mark the cache such as not updated
                    self._updated_cache = False
            finally:
                self._lock.release()

    def halt_cache_thread(self, timeout=None):
        """Halt cache auto_commit. This method aims to wait cache at most
        ``cache_autocommit`` or input timeout seconds before finishing.

        :param float timeout: max time to wait before waiting for this halting
            cache thread. Default value is self cache autocommit.
        """

        # change value of cache auto commit in order to stop thread
        cache_autocommit, self._cache_autocommit = self._cache_autocommit, 0

        if self._cached_thread is not None and self._cached_thread.isAlive():
            try:  # wait for cache thread end
                self._cached_thread.join(timeout)
            except RuntimeError:
                pass

        # recover cache auto commit
        self._cache_autocommit = cache_autocommit

    def execute_cache(self):
        """Execute the query cache and return execution processing."""

        result = None
        # do something only if there are cached query to execute
        if self._cache_count > 0:
            try:
                result = self._execute_cache()
            except Exception as ex:
                self.logger.error(
                    'Interruption of cache execution: {}'.format(ex)
                )
            else:  # if no error, renew the cache
                self._cache = self._new_cache()
            # initialize cache count
            self._cache_count = 0

        return result

    def _execute_cache(self):
        """Private cache execution. May be overriden."""

        raise NotImplementedError()

    def _ensure_index(self, index):
        """Get a right index structure related to input index.

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
        """Performs an atomic compare_and_swap operation on database related to
        input _id.

        :remarks: this method is not atomic

        :returns: True if the swamp succeed
        """
        raise NotImplementedError()

    def val_compare_and_swap(self, _id, oldvalue, newvalue):
        """Performs an atomic val_compare_and_swap operation on database related
        to input _id, oldvalue and newvalue.

        :remarks: this method is not atomic

        :returns: True if the comparison succeed
        """

        raise NotImplementedError()

    def get_elements(
            self,
            ids=None, query=None, limit=0, skip=0, sort=None, projection=None,
            tags=None, with_count=False
    ):
        """Get a list of elements where id are input ids.

        :param ids: element ids or an element id to get if is a string.
        :type ids: list of str
        :param dict query: set of couple of (field name, field value).
        :param int limit: max number of elements to get.
        :param int skip: first element index among searched list.
        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order.
        :type sort: list of {(str, {ASC, DESC}}), or str}
        :param dict projection: key names to keep from elements.
        :param list tags: search tags.
        :param bool with_count: If True (False by default), add count to the
            result.

        :return: a Cursor of input id elements, or one element if ids is a
            string (None if this element does not exist).
        :rtype: Cursor of dict elements or dict or NoneType
        """

        raise NotImplementedError()

    def __getitem__(self, ids):
        """Python shortcut to the get_elements(ids) method."""

        result = self.get_elements(ids=ids)

        if result is None or ids and not result:
            raise KeyError('%s not in self' % ids)

        return result

    def __contains__(self, ids):
        """Python shortcut to the get_elements(ids) method."""

        result = True

        # self does not contain ids only if a KeyError is raised
        try:
            self[ids]

        except KeyError:
            result = False

        return result

    def distinct(self, field, query):
        """Find distinct elements from a query into the given storage.

        :param string field: The distinct field to projection
        :param dict query: set of couple of (field name, field value).
        """

        raise NotImplementedError()

    def find_elements(
            self,
            query=None, limit=0, skip=0, sort=None, projection=None,
            tags=None, with_count=False
    ):
        """Find elements corresponding to input request and in taking care of
        limit, skip and sort find parameters.

        :param dict query: set of couple of (field name, field value).
        :param int limit: max number of elements to get.
        :param int skip: first element index among searched list.
        :param list sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order.
        :param dict projection: key names to keep from elements.
        :param list tags: search tags.
        :param bool with_count: If True (False by default), add count to the
            result.

        :return: a cursor of input request elements.
        :rtype: Cursor
        """

        raise NotImplementedError()

    def remove_elements(self, ids=None, _filter=None, tags=None, cache=False):
        """Remove elements identified by the unique input ids.

        :param ids: ids of elements to delete.
        :type ids: list of str
        :param dict _filter: removing filter.
        :param Filter _filter: additional filter to use if not None.
        :param list tags: element tags to remove.
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def __delitem__(self, ids):
        """Python shortcut to the remove_elements method."""

        return self.remove_elements(ids=ids)

    def __isub__(self, ids):
        """Python shortcut to the remove_elements method."""

        self.remove_elements(ids=ids)

    def put_element(
        self, element, _id=None, tags=None, version=None, cache=False
    ):
        """Put an element identified by input id.

        :param str _id: element id to update.
        :param dict element: element to put (couples of field (name,value)).
        :param list tags: element indexed tags.
        :param bool cache: use query cache if True (False by default).

        :return: True if updated.
        :rtype: bool
        """

        raise NotImplementedError()

    def __setitem__(self, _id, element):
        """Python shortcut for the put_element method."""

        self.put_element(_id=_id, element=element)

    def __iadd__(self, element):
        """Python shortcut for the put_element method."""

        if isinstance(element, list):
            self.put_elements(elements=element)

        else:
            self.put_element(element=element)

    def put_elements(self, elements, tags=None, cache=False):
        """Put list of elements at a time.

        :param list elements: elements to put.
        :param list tags: element indexed tags.
        :param bool cache: use query cache if True (False by default)."""

        for element in elements:
            self.put_element(element)

    def count_elements(self, query=None, tags=None):
        """Count elements corresponding to the input query.

        :param dict query: query which contain set of couples (key, value)

        :return: Number of elements corresponding to the input query
        :rtype: int
        """

        cursor = self.find_elements(query=query, tags=None)

        result = len(cursor)

        return result

    def __len__(self):
        """Python shortcut to the count_elements method."""

        return self.count_elements()

    def _find(self, *args, **kwargs):
        """Find operation dedicated to technology implementation."""

        raise NotImplementedError()

    def _update(self, cache=False, *args, **kwargs):
        """Update operation dedicated to technology implementation.

        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def _remove(self, cache=False, *args, **kwargs):
        """Remove operation dedicated to technology implementation.

        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def _insert(self, cache=False, *args, **kwargs):
        """Insert operation dedicated to technology implementation.

        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def _count(self, *args, **kwargs):
        """Count operation dedicated to technology implementation."""

        raise NotImplementedError()

    def get_table(self):
        """Table name related to self table or type and data_scope.

        :return: table name.
        :rtype: str
        """

        # try to use local table
        result = self.table

        if not result:

            prefix = self.data_type

            if isiterable(prefix, is_str=False):
                prefix = reduce(lambda x, y: '%s_%s' % (x, y), prefix)

            result = "{0}_{1}".format(prefix, self.data_scope).lower()

        return result

    def copy(self, target):
        """Copy self content into target storage.

        target type must implement the same class in cstorage packege as self.
        If self implements directly cstorage.Storage, we don't care about
        target type

        :param target: target storage where copy content
        :type target: same as self or any storage if type(self) is Storage
        """

        result = 0

        from canopsis.storage.core import Storage
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

        if isinstance(self, storage_types):
            for _ in storage_types:
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
        """Called by Storage.copy(self, target) in order to ensure than target
        type is the same as self.
        """

        for element in self.get_elements():
            target.put_element(element=element)

        raise NotImplementedError()

    def _element_id(self, element):
        """Get element id related to self behavior."""

        raise NotImplementedError()

    def _get_category(self, *args, **kwargs):
        """Get configuration category for self storage."""

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
                Parameter(Storage.DATA, parser=eval),
                Parameter(Storage.TABLE),
                Parameter(Storage.CACHE_SIZE, parser=int, critical=True),
                Parameter(
                    Storage.CACHE_ORDERED, parser=Parameter.bool, critical=True
                ),
                Parameter(
                    Storage.CACHE_AUTOCOMMIT, parser=float, critical=True
                )
            )
        )

        return result

    @staticmethod
    def _resolve_sort(sort):
        """Resolve input sort in transforming it to a list of tuple of (name,
        direction).

        :param sort: sort configuration. Can be a string, or
        :type sort: list of {tuple(str, int), str}
        :return: depending on type of sort:
            - str: [(sort, Storage.ASC)]
            - dict: [(sort['property'], sort.get('direction', Storage.ASC))]
            - tuple: [(sort[0], sort[1])]
            - list:
                - str
        :rtype: str, dict, tuple or list
        """

        result = []
        if isinstance(sort, basestring):

            result.append((sort, Storage.ASC))

        elif isinstance(sort, dict):

            sort_tuple = None
            field = sort.get('property', None)

            if field is not None:
                direction = sort.get('direction', Storage.ASC)
                if isinstance(direction, basestring):
                    direction = getattr(Storage, direction.upper())
                    # Need field property filled in the sort document
                    sort_tuple = (field, direction)

            if sort_tuple is not None:
                result.append(sort_tuple)

        elif isinstance(sort, tuple):

            direction = sort[1]
            if isinstance(direction, basestring):
                direction = getattr(Storage, direction.upper())
            result.append((sort[0], direction))

        elif isinstance(sort, Iterable):

            for item in sort:
                result += Storage._resolve_sort(item)

        return result


class Cursor(object):
    """Query cursor object.

    An iterable object in order to retrieve data from a Storage.
    A reference to the technology cursor is provided by the cursor getter.
    """

    __slots__ = ('_cursor', )

    def __init__(self, cursor):
        """
        :param cursor: Technology implementation cursor.
        """

        super(Cursor, self).__init__()

        self._cursor = cursor

    @property
    def cursor(self):
        """Get technology implementation cursor."""

        return self._cursor

    def __len__(self):
        """Get number of cursor items."""

        raise NotImplementedError()

    def __iter__(self):
        """Iterate on cursor items."""

        raise NotImplementedError()

    def __getitem__(self, index):
        """Get a single document or a slice of documents from this cursor.

        :param index: An integer or slice index to be applied to this cursor.
        :type index: int or slice
        """

        raise NotImplementedError()
