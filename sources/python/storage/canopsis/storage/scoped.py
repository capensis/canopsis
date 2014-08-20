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

    __datatype__ = 'scoped'  #: registered such as a scoped storage

    SCOPE_SEPARATOR = '/'  #: char separator between scope values

    VALUE = 'value'  #: data value
    SCOPE = 'scope'
    ID = 'id'  #: data id

    def __init__(self, scope=None, *args, **kwargs):
        """
        :param scope: iterable of ordered lists of scope names
        :type scope: Iterable
        """

        super(ScopedStorage, self).__init__(*args, **kwargs)

        self._scope = scope

    @property
    def scope(self):
        """
        tuple of ordered field names.
        """
        return self._scope

    @scope.setter
    def scope(self, value):
        self._scope = value
        self.reconnect()

    def get(
        self, scope, ids=None, _filter=None, limit=0, skip=0, sort=None
    ):
        """
        Get data related to input ids, input scope and input filter.

        :param scope: dictionnary of scope valut by scope name
        :type scope: dict

        :param ids: data ids in the input scope.
        :type ids: str or iterable of str

        :param filter: additional filter condition to input scope
        :type filter: storage filter

        :param limit: max number of data to get. Useless if data_id is given.
        :type limit: int

        :param skip: starting index of research if multi data to get
        :type skip: int

        :param sort: couples of field (name, value) to sort with ASC/DESC
            Storage fields
        :type sort: dict

        :return: a data or a list of data respectively to ids such as a str or
            an iterable of str
        :rtype: dict or list of dict
        """

        raise NotImplementedError()

    def find(self, scope, _filter, limit=0, skip=0, sort=None):
        """
        Get a list of data identified among a dictionary of scoped values by
        name.

        :param scope: scope
        :type scope: storage filter

        :param _filter: additional filter condition to input scope
        :type _filter: storage filter

        :param limit: max number of data to get. Useless if data_id is given.
        :type limit: int

        :param skip: starting index of research if multi data to get
        :type skip: int

        :param sort: couples of field (name, value) to sort with ASC/DESC
            Storage fields
        :type sort: dict

        :return: a list of data.
        :rtype: list of dict
        """

        raise NotImplementedError()

    def put(self, scope, _id, data):
        """
        Put a data related to an id

        :param scope: scope
        :type scope: storage filter

        :param _id: data id
        :type _id: str

        :param data: data to update
        :type data: dict
        """

        raise NotImplementedError()

    def remove(self, scope, ids=None):
        """
        Remove data from ids or type

        :param scope: scope to remove
        :type scope: storage filter

        :param _ids: data id or list of data id
        :type _ids: list or str
        """

        raise NotImplementedError()

    def get_absolute_path(self, scope, _id):
        """
        Get data absolute path among scope for input data_id in input scope.

        :param scope: dictionary of scope value by name
        :type scope: dict

        :param data_id: data id in input scope
        :type data_id: str
        """

        result = None

        result = ''
        for scope_name in scope:
            result = '%s%s%s' % (
                result, ScopedStorage.SCOPE_SEPARATOR,
                scope[scope_name])

        if result is not None:
            result = '%s%s%s' % (
                result, ScopedStorage.SCOPE_SEPARATOR, _id)

        return result
