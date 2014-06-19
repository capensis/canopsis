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

from cstorage import Storage


class TypedStorage(Storage):
    """
    Storage dedicated to manage typed data.
    """

    VALUE = 'value'
    TYPE = 'type'

    class TypedStorageError(Exception):
        pass

    def get(
        self, _ids=None, data_type=None, limit=0, skip=0, sort=None,
        *args, **kwargs
    ):
        """
        Get a list of data identified among data_ids or a type

        :param data_ids: data ids to get
        :type data_id: list of str

        :param data_type: data_id type to get if not None
        :type data_type: str

        :param limit: max number of data to get
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

    def put(self, _id, data, data_type=None, *args, **kwargs):
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

    def remove(self, _ids=None, data_type=None, *args, **kwargs):
        """
        Remove data from ids or type

        :param _ids: list of data id
        :type _ids: list

        :param data_type: data type to remove if not None
        :type data_type: str
        """

        raise NotImplementedError()

    def _get_storage_type(self, *args, **kwargs):

        return 'typed'
