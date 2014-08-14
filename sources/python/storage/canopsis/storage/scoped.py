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

from canopsis.storage import Storage


class ScopedStorage(Storage):
    """
    Storage dedicated to manage scoped data identified by a data id in a scope
    of ordered fields.

    For example, a metric is identified by a unique name in the scope
    (type=metric, connector, component, resource) or
    (type=metric, connector, component).
    """

    __storage_type__ = 'scoped'

    SCOPE_SEPARATOR = '/'  #: char separator between scope values

    VALUE = 'value'
    SCOPE = 'scope'
    ID = 'id'

    def __init__(self, scope, *args, **kwargs):
        """
        :param scope: iterable of ordered lists of scope names
        :type scope: Iterable
        """

        super(ScopedStorage, self).__init__(*args, **kwargs)

        self.scope = scope

    @property
    def scope(self):
        return self._scope

    @scope.setter
    def scope(self, value):
        self._set_scope(scope=value)

    def _set_scope(self, scope):
        """
        Self scope setter.
        """
        self._scope = scope

    def get(self, scope, data_id):
        """
        Get a data related to a dictionnary of scope value by name and a
        data_id in the input scope.

        :param scope: dictionnary of scope valut by scope name.
        :type scope: dict

        :param data_id: data id in the input scope.
        :type data_id: str

        :return: a data
        :rtype: dict
        """

        raise NotImplementedError()

    def find(self, scope, filter, limit=0, skip=0, sort=None):
        """
        Get a list of data identified among a dictionary of scoped values by
        name and a data_id.

        :param scope: scope
        :type scope: storage filter

        :param filter: additional filter condition to input scope
        :type filter: storage filter

        :param data_id: if not None, data id in the input scope.
        :type data_id: storage filter

        :param limit: max number of data to get. Useless if data_id is given.
        :type limit: int

        :param skip: starting index of research if multi data to get
        :type skip: int

        :param sort: couples of field (name, value) to sort with ASC/DESC
            Storage fields
        :type sort: dict

        :return: a list of couples of field (name, value) or None respectivelly
            if such data exist or not
        :rtype: list of dict of field (name, value)
        """

        raise NotImplementedError()

    def put(self, scope, data_id, data):
        """
        Put a data related to an id

        :param _id: data id
        :type _id: str

        :param data_type: data type to update
        :type data_type: str

        :param data: data to update
        :type data: dict
        """

        raise NotImplementedError()

    def remove(self, _ids=None, data_type=None):
        """
        Remove data from ids or type

        :param _ids: list of data id
        :type _ids: list

        :param data_type: data type to remove if not None
        :type data_type: str
        """

        raise NotImplementedError()

    def get_absolute_path(self, scope, data_id):
        """
        Get data absolute path among scope for input data_id in input scope.

        :param scope: dictionary of scope value by name
        :type scope: dict

        :param data_id: data id in input scope
        :type data_id: str
        """

        result = None

        appropriated_scope = self.get_appropariate_scope(scope)

        if appropriated_scope is not None:
            result = ScopedStorage.SCOPE_SEPARATOR
            for scope_name in appropriated_scope:
                if scope_name in scope:
                    result = '%s%s%s' % (
                        result, ScopedStorage.SCOPE_SEPARATOR,
                        scope[scope_name])

        if result is not None:
            result = '%s%s%s' % (
                result, ScopedStorage.SCOPE_SEPARATOR, data_id)

        return result
