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

from uuid import uuid4 as uuid

from canopsis.configuration.model import Parameter
from canopsis.common.utils import ensure_iterable, isiterable, get_first
from canopsis.storage.core import Storage

from urllib import quote_plus, unquote_plus


class CompositeStorage(Storage):
    """Storage dedicated to manage composite data identified by a name in a path
    of ordered fields.

    For example, a metric is identified by a unique name in the path
    (type=metric, connector, component, resource) or
    (type=metric, connector, component).

    In addition to such composity, data of the same name and type can be the
    same data with different path. In such case, they are called shared and
    share the same value which is unique among all composite data.
    """

    __datatype__ = 'composite'  #: registered such as a composite storage.

    PATH_SEPARATOR = '/'  #: char separator between path values.

    SHARED = 'shared'  #: shared field name.
    VALUE = 'value'  #: data value.
    PATH = 'path'  #: path value.

    NAME = 'name'  #: name field name.

    EXTENDED = 'extended'  #: additional fields over self.path.

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
        :return: tuple of ordered field names.
        """

        return self._path

    @path.setter
    def path(self, value):

        self._path = value
        self.reconnect()

    def all_indexes(self, *args, **kwargs):

        result = super(CompositeStorage, self).all_indexes(*args, **kwargs)

        # add index to shared property
        result.append([(CompositeStorage.SHARED, Storage.ASC)])

        return result

    def get_shared_data(self, shared_ids):
        """Get all shared data related to input shared ids.

        :param shared_ids: one or more data.
        :type shared_ids: list or str

        :return: depending on input shared_ids::

            - one shared id: one list of shared data
            - list of shared ids: list of list of shared data
        """

        result = []

        sids = ensure_iterable(shared_ids, iterable=set)

        for shared_id in sids:
            query = {CompositeStorage.SHARED: shared_id}
            shared_data = self.get_elements(query=query)
            result.append(shared_data)

        # return first or data if data is not an iterable
        if not isiterable(shared_ids, is_str=False):
            result = get_first(result)

        return result

    def share_data(
            self, data, shared_id=None, share_extended=False, cache=False
    ):
        """Set input data as a shared data with input shared id.

        :param data: one data
        :param str shared_id: unique shared id. If None, the id is generated.
        :param bool share_extended: if True (False by default), set shared
            value to all shared data with input data
        :param bool cache: use query cache if True (False by default).

        :return: shared_id value (generated if None)
        """

        result = str(uuid()) if shared_id is None else shared_id

        # get an iterable version of input data
        if isinstance(data, dict):
            data_to_share = [data]
        else:
            data_to_share = data

        for dts in data_to_share:
            # update extended data if necessary
            if share_extended:
                path, name = self.get_path_with_name(dts)
                extended_data = self.get(
                    path=path, names=name, shared=True
                )
                # decompose extended data into a list
                dts = []
                for ed in extended_data:
                    dts.append(ed)
            else:
                dts = [dts]

            for dt in dts:
                path, name = self.get_path_with_name(dt)
                dt[CompositeStorage.SHARED] = result
                self.put(path=path, name=name, data=dt, cache=cache)

        return result

    def unshare_data(self, data, cache=False):
        """Remove share property from input data.

        :param data: one or more data to unshare.
        :param bool cache: use query cache if True (False by default).
        """
        data = ensure_iterable(data)

        for d in data:
            if CompositeStorage.SHARED in d:
                d[CompositeStorage.SHARED] = str(uuid())
                path, name = self.get_path_with_name(d)
                self.put(path=path, name=name, data=d, cache=cache)

    def get(
            self,
            path, names=None, _filter=None, shared=False,
            limit=0, skip=0, sort=None, with_count=False
    ):
        """Get data related to input names, input path and input filter.

        :param dict path: dictionnary of path valut by path name
        :param names: data names in the input path.
        :type names: str or iterable of str
        :param _filter: additional filter condition to input path
        :type _filter: storage filter
        :param bool shared: if True, convert result to list of list of data
            where list of data are list of shared data.
        :param int limit: max number of data to get. Useless if name exists.
        :param int skip: starting index of research if multi data to get
        :param dict sort: couples of field (name, value) to sort with ASC/DESC
            Storage fields
        :param bool with_count: If True (False by default), add count to the
            result

        :return: a data or a list of data respectively to ids such as a str or
            an iterable of str
        :rtype: dict or list of dict
        """

        raise NotImplementedError()

    def find(
            self, path, _filter, shared=False, limit=0, skip=0, sort=None,
            with_count=False
    ):
        """Get a list of data identified among a dictionary of composite values
        by name.

        :param path: path
        :type path: storage filter
        :param _filter: additional filter condition to input path
        :type _filter: storage filter
        :param bool shared: if True, convert result to list of list of data
            where list of data are list of shared data.
        :param int limit: max number of data to get. Useless if name exists.
        :param int skip: starting index of research if multi data to get
        :param dict sort: couples of field (name, value) to sort with ASC/DESC
            Storage fields
        :param bool with_count: If True (False by default), add count to the
            result

        :return: a list of data.
        :rtype: list of dict
        """

        raise NotImplementedError()

    def put(self, path, name, data, shared_id=None, cache=False):
        """Put a data related to an id and a path.

        :param path: path
        :type path: storage filter
        :param str name: data id
        :param dict data: data to update
        :param str shared_id: shared_id id not None
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def remove(self, path, names=None, shared=False, cache=False):
        """Remove data from ids or type.

        :param path: path to remove
        :type path: storage filter
        :param names: data id or list of data id
        :type names: list or str

        :param bool shared: remove shared data if data ids are related to
            shared data.
        :param bool cache: use query cache if True (False by default).
        """

        raise NotImplementedError()

    def get_path_with_name(self, data):
        """Get input data path and id.

        :type data: dict

        :return: data path, data id
        :rtype: tuple
        """

        path = {
            field: data[field]
            for field in data
            if field in self.path and data.get(field) is not None
        }

        result = path, data[CompositeStorage.NAME]

        return result

    def get_data_from_id(self, data_id):

        result = {}

        names = data_id.split(CompositeStorage.PATH_SEPARATOR)[1:]

        for path, name in zip(self.path, names):
            result[path] = unquote_plus(name)

        if len(names) > len(self.path):
            result['extended'] = names[-1]

        return result

    def get_absolute_path(self, path, name=None):
        """Get input data absolute path.

        :param dict path: path from where get absolute path.
        :param str name: data id
        """

        result = ''

        for field in self.path:
            if path.get(field) is not None:
                result = '{0}{1}{2}'.format(
                    result,
                    CompositeStorage.PATH_SEPARATOR,
                    quote_plus(path[field])
                )
            else:
                break

        if name is not None and result:
            result = '{0}{1}{2}'.format(
                result,
                CompositeStorage.PATH_SEPARATOR,
                quote_plus(name)
            )

        return result

    def _conf(self, *args, **kwargs):

        result = super(CompositeStorage, self)._conf(*args, **kwargs)

        category = None

        for _category in result:
            category = _category

        # add path property to the last category
        if category is not None:
            category += Parameter(CompositeStorage.PATH, parser=eval)

        return result
