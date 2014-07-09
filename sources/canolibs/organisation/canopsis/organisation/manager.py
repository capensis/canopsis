#!/usr/bin/env python
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

from canopsis.storage.manager import Manager


class Organisation(Manager):
    """
    Manage access to a context (connector, component, resource) elements
    and context data (metric, downtime, etc.)

    For example, in following entities:
        - component: data_id is component_id and data_type is component
        - connector: data_id
    """

    CONF_FILE = '~/etc/organisation.conf'

    DATA_TYPE = 'organisation'

    def get_users(self):

        raise NotImplementedError()

    def get(
        self,
        element_type,
        connector=None, connector_type=None, component=None, resource=None,
        other=None,
        element_ids=None, timewindow=None,
        *args, **kwargs
    ):
        """
        Get all elements which have a element_id among input element_ids, or
        type equals to element_type or inside timewindow if specified
        """

        if element_ids is None:
            element_id = Organisation.get_element_id(
                connector, connector_type, component, resource, other)
            element_ids = [element_id]

        return self._ctx.get(
            data_ids=element_ids, data_type=element_type,
            timewindow=timewindow, *args, **kwargs)

    def get_by_name(
        self,
        element_type, name=None, *args, **kwargs
    ):

        return self.get(element_id=name, element_type=element_type)

    def put(
        self, element_id, element_type, element, timestamp=None,
        *args, **kwargs
    ):
        """
        Put an element designated by the element_id, element_type and element.
        If timestamp is None, time.now is used.
        """

        return self._ctx.put(
            data_id=element_id, data_type=element_type, data=element,
            timestamp=timestamp, *args, **kwargs)

    def remove(
        self, element_ids=None, element_type=None, timewindow=None,
        *args, **kwargs
    ):
        """
        Remove a set of elements identified by element_ids, an element type or
        a timewindow
        """

        return self._ctx.remove(
            data_ids=element_ids, data_type=element_type,
            timewindow=timewindow, *args, **kwargs)

    def _get_conf_files(self, conf_files, *args, **kwargs):

        result = super(Organisation, self)._get_conf_files(
            conf_files=conf_files, *args, **kwargs)

        result.append(Organisation.CONF_FILE)

        return result
