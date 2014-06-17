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


class LastValueStorage(Storage):
    """
    Storage dedicated to manage last value data.
    """

    VALUE = 'value'

    class LastValueStorageError(Exception):
        pass

    def get(
        self, data_ids=None, limit=0, skip=0, sort=None,
        *args, **kwargs
    ):
        """
        Get a list of data in limiting number of document, skipping and sorting
            the result

        :param data_ids: data ids to get
        :type data_id: list of str

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

    def put(self, values_by_id, *args, **kwargs):
        """
        Put a set of data by id

        :param values_by_id: couples of data by id
        :type data_id: dict
        """

        raise NotImplementedError()

    def remove(self, data_ids=None, *args, **kwargs):
        """
        Remove data_ids

        :param data_ids: list of data id
        :type data_ids: list
        """

        raise NotImplementedError()

    def _get_storage_type(self, *args, **kwargs):

        return 'last_value'
