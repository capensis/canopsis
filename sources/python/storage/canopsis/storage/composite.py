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

from uuid import uuid4 as uuid

from canopsis.storage import Storage


class CompositeStorage(Storage):
    """
    Storage dedicated to manage composite data identified by a data id in a
    path of ordered fields.

    For example, a metric is identified by a unique name in the path
    (type=metric, connector, component, resource) or
    (type=metric, connector, component).

    In addition to such composity, data of the same name and type can be the
    same data with different path. In such case, they are called shared and
    share the same value which is unique among all composite data.
    """

    __datatype__ = 'composite'  #: registered such as a composite storage

    PATH_SEPARATOR = '/'  #: char separator between path values

    SHARED = 'shared'  #: shared field name
    VALUE = 'value'  #: data value
    PATH = 'path'
    ID = 'id'  #: data id

    def __init__(self, path=None, *args, **kwargs):
        """
        :param path: iterable of ordered lists of path names
        :type path: Iterable
        """

        super(CompositeStorage, self).__init__(*args, **kwargs)

        self._path = path

    @property
    def path(self):
        """
        tuple of ordered field names.
        """
        return self._path

    @path.setter
    def path(self, value):
        self._path = value
        self.reconnect()

    def get_shared_data(self, data):
        """
        Get all shared data related to input data. If input data is not shared,
        returns a list containing only data.
        """

        result = [data]

        if CompositeStorage.SHARED in data:

            shared = data[CompositeStorage.SHARED]
            request = {CompositeStorage.SHARED: shared}
            result = self.find_elements(request=request)

        return result

    def set_shared_data(self, data, shared=None):
        """
        Set input data as a shared data with input shared id

        If input data is already shared, update all shared data with input
        shared id

        :param shared: unique id
        """
        if shared is None:
            shared = uuid()

        shared_data = [data]

        # if data is alraedy shared, update all shared data with input shared
        if CompositeStorage.SHARED in data \
                and data[CompositeStorage.SHARED] != shared:
            shared_data += self.get_shared_data(data)

        # for all shared data, update the shared property and put them
        for _shared_data in shared_data:
            _shared_data[CompositeStorage.SHARED] = shared
            _id = self.get_absolute_path(_shared_data)
            self.put(_id=_id, element=_shared_data)

    def get(
        self, path, ids=None, _filter=None, limit=0, skip=0, sort=None
    ):
        """
        Get data related to input ids, input path and input filter.

        :param path: dictionnary of path valut by path name
        :type path: dict

        :param ids: data ids in the input path.
        :type ids: str or iterable of str

        :param filter: additional filter condition to input path
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

    def find(self, path, _filter, limit=0, skip=0, sort=None):
        """
        Get a list of data identified among a dictionary of composite values by
        name.

        :param path: path
        :type path: storage filter

        :param _filter: additional filter condition to input path
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

    def put(self, path, _id, data):
        """
        Put a data related to an id

        :param path: path
        :type path: storage filter

        :param _id: data id
        :type _id: str

        :param data: data to update
        :type data: dict
        """

        raise NotImplementedError()

    def remove(self, path, ids=None):
        """
        Remove data from ids or type

        :param path: path to remove
        :type path: storage filter

        :param _ids: data id or list of data id
        :type _ids: list or str
        """

        raise NotImplementedError()

    def get_absolute_path(self, path, _id):
        """
        Get data absolute path among path for input data_id in input path.

        :param path: dictionary of path value by name
        :type path: dict

        :param data_id: data id in input path
        :type data_id: str
        """

        result = None

        result = ''
        for path_name in path:
            result = '%s%s%s' % (
                result, CompositeStorage.PATH_SEPARATOR,
                path[path_name])

        if result is not None:
            result = '%s%s%s' % (
                result, CompositeStorage.PATH_SEPARATOR, _id)

        return result
