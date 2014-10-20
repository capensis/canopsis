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

from functools import reduce

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

    CONF_RESOURCE = 'storage/storage.conf'

    class DataBaseError(Exception):
        pass

    def __init__(self, db='canopsis', journaling=False, *args, **kwargs):
        """
        :param db: db name
        :param journaling: journaling mode enabling.

        :type db: str
        :type journaling: bool
        """

        super(DataBase, self).__init__(*args, **kwargs)

        # initialize instance properties with default values
        self._db = db
        self._journaling = journaling

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
                    DataBase.JOURNALING, parser=Parameter.bool, critical=True)
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
    """

    __protocol__ = 'storage'
    """register itself and all subclasses to storage protocol"""

    DATA_ID = 'id'  #: db data id

    INDEXES = 'indexes'

    CATEGORY = 'STORAGE'

    ASC = 1  #: ASC order
    DESC = -1  #: DESC order

    class StorageError(Exception):
        """
        Handle Storage errors
        """
        pass

    def __init__(self, indexes=None, *args, **kwargs):

        super(Storage, self).__init__(*args, **kwargs)

        self._indexes = [] if indexes is None else indexes

    @property
    def indexes(self):

        return self._indexes

    def all_indexes(self):
        """
        :return: all self indexes.
        """

        result = [[(Storage.DATA_ID, 1)]]
        _indexes = self._indexes
        if _indexes:
            result.append(self._indexes[:])

        return result

    @indexes.setter
    def indexes(self, value):
        """
        Indexes setter

        :param value: indexes such as::
            - one name
            - one tuple of kind (name, ASC/DESC)
            - a list of tuple or name [((name, ASC/DESC) | name)* ]
        :type value: str, tuple ot list
        """

        # if value is a name, transform it into a list
        if isinstance(value, (str, tuple)):
            value = [value]
        # if value is iterable
        elif not isinstance(value, list):
            raise Storage.StorageError(
                "wrong indexes value %s. str, tuple or list accepted" % value)

        self._indexes = value
        self.reconnect()

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
            (None if this element does not exist)
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
        self, request, limit=0, skip=0, sort=None, with_count=False
    ):
        """
        Find elements corresponding to input request and in taking care of
        limit, skip and sort find parameters.

        :param request: set of couple of (field name, field value)
        :type request: dict(str, object)

        :param limit: max number of elements to get
        :type limit: int

        :param skip: first element index among searched list
        :type skip: int

        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order
        :type sort: list of {(str, {ASC, DESC}}), or str}

        :param bool with_count: If True (False by default), add count to the
            result

        :return: input request elements
        :rtype: list of objects
        """

        raise NotImplementedError()

    def remove_elements(self, ids):
        """
        Remove elements identified by the unique input ids

        :param ids: ids of elements to delete
        :type ids: list of str
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

    def put_element(self, _id, element):
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

    def count_elements(self, request=None):
        """
        Count elements corresponding to the input request

        :param id: request which contain set of couples (key, value)
        :type id: dict

        :return: Number of elements corresponding to the input request
        :rtype: int
        """

        raise NotImplementedError()

    def __len__(self):
        """
        Python shortcut to the count_elements method.
        """

        return self.count_elements()

    def _find(self, *args, **kwargs):
        """
        Find operation dedicated to technology implementation.
        """

        raise NotImplementedError()

    def _update(self, *args, **kwargs):
        """
        Update operation dedicated to technology implementation.
        """

        raise NotImplementedError()

    def _remove(self, *args, **kwargs):
        """
        Remove operation dedicated to technology implementation.
        """

        raise NotImplementedError()

    def _insert(self, *args, **kwargs):
        """
        Insert operation dedicated to technology implementation.
        """

        raise NotImplementedError()

    def _count(self, *args, **kwargs):
        """
        Count operation dedicated to technology implementation.
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
            new_content=(Parameter(Storage.INDEXES, parser=eval)))

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
