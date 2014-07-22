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

from canopsis.configuration import Configurable, Parameter, Configuration
from canopsis.common.utils import resolve_element
from canopsis.storage import Storage


class Manager(Configurable):

    CONF_FILE = '~/etc/manager.conf'

    DATA_TYPE = 'data_type'

    TIMED_STORAGE = 'timed_storage'
    PERIODIC_STORAGE = 'periodic_storage'
    STORAGE = 'storage'
    TYPED_STORAGE = 'typed_storage'
    TIMED_TYPED_STORAGE = 'timed_typed_storage'
    SHARED = 'shared'

    CATEGORY = 'MANAGER'

    STORAGE_SUFFIX = 'storage'

    _STORAGE_BY_DATA_TYPE_BY_TYPE = dict()

    def __init__(
        self,
        data_type,
        shared=True,
        storage=None, timed_storage=None, periodic_storage=None,
        typed_storage=None, timed_typed_storage=None,
        *args, **kwargs
    ):

        super(Manager, self).__init__(*args, **kwargs)

        self.shared = shared

        self.data_type = data_type

        self.periodic_storage = periodic_storage
        self.timed_storage = timed_storage
        self.storage = storage
        self.typed_storage = typed_storage
        self.timed_typed_storage = timed_typed_storage

    @property
    def data_type(self):
        return self._data_type

    @data_type.setter
    def data_type(self, value):
        self._data_type = value

    @property
    def shared(self):
        return self._shared

    @shared.setter
    def shared(self, value):
        self._shared = value

    def _get_property_storage(self, value, *args, **kwargs):
        """
        Get property storage where value is given in calling property.setter
        """
        result = value

        if value is not None and not isinstance(value, Storage):
            result = self.get_storage(storage_type=value)

        return result

    @property
    def periodic_storage(self):
        return self._periodic_storage

    @periodic_storage.setter
    def periodic_storage(self, value):
        self._periodic_storage = self._get_property_storage(value)

    @property
    def timed_storage(self):
        return self._timed_storage

    @timed_storage.setter
    def timed_storage(self, value):
        self._timed_storage = self._get_property_storage(value)

    @property
    def storage(self):
        return self._storage

    @storage.setter
    def storage(self, value):
        self._storage = self._get_property_storage(value)

    @property
    def typed_storage(self):
        return self._typed_storage

    @typed_storage.setter
    def typed_storage(self, value):
        self._typed_storage = self._get_property_storage(value)

    @property
    def timed_typed_storage(self):
        return self._timed_typed_storage

    @timed_typed_storage.setter
    def timed_typed_storage(self, value):
        self._timed_typed_storage = self._get_property_storage(value)

    def get_storage(
        self, data_type=None, storage_type=None, shared=None, *args, **kwargs
    ):
        """
        Load a storage related to input data type and storage type.

        If shared, the result instance is shared among same storage type and
        data type.

        :param data_type: storage data type
        :type data_type: str

        :param storage_type: storage type (among timed, last_value ,etc.)
        :type storage_type: Storage or str

        :param shared: if True, the result is a shared storage instance among
            managers. If None, use self.shared
        :type shared: bool

        :return: storage instance corresponding to input storage_type
        :rtype: Storage
        """

        result = None

        if data_type is None:
            data_type = self.data_type

        if storage_type is None:
            storage_type = self.storage

        if shared is None:
            shared = self.shared

        if isinstance(storage_type, str):
            storage_type = resolve_element(storage_type)

        elif callable(storage_type):
            pass

        # if shared, try to find an instance with same storage and data types
        if shared:
            # search among isntances registred on storage_type
            storage_by_data_type = \
                Manager._STORAGE_BY_DATA_TYPE_BY_TYPE.setdefault(
                    storage_type, dict())

            if data_type not in storage_by_data_type:
                storage_by_data_type[data_type] = storage_type(
                    data_type=data_type, *args, **kwargs)

            result = storage_by_data_type[data_type]

        else:
            result = storage_type(data_type=data_type, *args, **kwargs)

        return result

    def get_timed_storage(
        self, data_type, timed_type=None, shared=None,
        *args, **kwargs
    ):

        if timed_type is None:
            timed_type = self.timed_storage

        result = self.get_storage(
            data_type=data_type, storage_type=timed_type, shared=shared,
            *args, **kwargs)

        return result

    def get_periodic_storage(
        self, data_type, periodic_type=None, shared=None,
        *args, **kwargs
    ):

        if periodic_type is None:
            periodic_type = self.periodic_storage

        result = self.get_storage(data_type=data_type, shared=shared,
            storage_type=periodic_type, *args, **kwargs)

        return result

    def get_typed_storage(
        self, data_type, typed_type=None, shared=None,
        *args, **kwargs
    ):

        if typed_type is None:
            typed_type = self.typed_storage

        result = self.get_storage(data_type=data_type, shared=shared,
            storage_type=typed_type, *args, **kwargs)

        return result

    def get_timed_typed_storage(
        self, data_type, timed_typed_type=None, shared=None,
        *args, **kwargs
    ):

        if timed_typed_type is None:
            timed_typed_type = self.timed_typed_storage

        result = self.get_storage(data_type=data_type, shared=shared,
            storage_type=timed_typed_type, *args, **kwargs)

        return result

    def _get_conf_files(self, *args, **kwargs):

        result = super(Manager, self)._get_conf_files(*args, **kwargs)

        result.append(Manager.CONF_FILE)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Manager, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Manager.CATEGORY,
            new_content=(
                Parameter(Manager.TIMED_STORAGE, self.timed_storage),
                Parameter(Manager.PERIODIC_STORAGE, self.periodic_storage),
                Parameter(Manager.STORAGE, self.storage),
                Parameter(Manager.TYPED_STORAGE, self.typed_storage),
                Parameter(
                    Manager.TIMED_TYPED_STORAGE, self.timed_typed_storage),
                Parameter(Manager.SHARED, parser=bool)))

        return result

    def _configure(self, unified_conf, *args, **kwargs):

        super(Manager, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        # set data_type
        self._update_property(
            unified_conf=unified_conf, param_name=Manager.DATA_TYPE)

        # set shared
        self._update_property(
            unified_conf=unified_conf, param_name=Manager.SHARED)

        values = unified_conf[Configuration.VALUES]

        # set all storages
        for parameter in values:
            if parameter.name.endswith(Manager.STORAGE_SUFFIX):
                self._update_property(
                    unified_conf=unified_conf, param_name=parameter.name,
                    public_property=True)
