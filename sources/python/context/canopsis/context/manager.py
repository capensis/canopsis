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
from md5 import new as md5
from canopsis.configuration import Parameter


class Context(Manager):
    """
    Manage access to a context (connector, component, resource) elements
    and context data (metric, downtime, etc.)

    For example, in following entities:
        - component: data_id is component_id and data_type is component
        - connector: data_id
    """

    CONF_RESOURCE = 'context/context.conf'

    DATA_TYPE = 'context'

    CATEGORY = 'CONTEXT'

    CTX_STORAGE = 'ctx_storage'
    CONTEXT = 'context'

    def __init__(self, ctx_storage=None, *args, **kwargs):

        super(Context, self).__init__(self, *args, **kwargs)

        self.ctx_storage = ctx_storage

    @property
    def ctx_storage(self):
        return self._ctx_storage

    @ctx_storage.setter
    def ctx_storage(self, value):
        self._ctx_storage = value
        if value is not None:
            self._ctx = self.get_timed_typed_storage(
                data_type=Context.DATA_TYPE)

    def get(
        self,
        element_type,
        connector=None, connector_type=None, component=None, resource=None,
        other=None,
        timewindow=None,
        *args, **kwargs
    ):
        """
        Get one element related to:
            - an element_type,
            - a path (connector, ..., other),
            - a timewindow

        :return: array of couple(timestamp, element) in the ASC order
        :rtype: list(tuple(float, dict))
        """

        element_id = Context.get_element_id(
            connector, connector_type, component, resource, other)

        return self._ctx.get(
            data_ids=[element_id], data_type=element_type,
            timewindow=timewindow, *args, **kwargs)

    def find(
        self,
        element_type,
        connector=None, connector_type=None, component=None, resource=None,
        other=None,
        element_ids=None, timewindow=None,
        *args, **kwargs
    ):
        """
        Find all elements which have a element_id among input element_ids, or
        type equals to element_type or inside timewindow if specified
        """

        if element_ids is None:
            element_id = Context.get_element_id(
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

    def _get_conf_files(self, *args, **kwargs):

        result = super(Context, self)._get_conf_files(*args, **kwargs)

        result.append(Context.CONF_RESOURCE)

        return result

    def _configure(self, unified_conf, *args, **kwargs):

        super(Context, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        self._update_property(
            unified_conf=unified_conf, param_name=Context.CTX_STORAGE,
            public=True)

    def _conf(self, *args, **kwargs):

        result = super(Context, self)._conf(*args, **kwargs)

        result.add_unified_category(
            name=Context.CATEGORY,
            new_content=Parameter(Context.CTX_STORAGE, self.ctx_storage))

        return result

    @staticmethod
    def get_element_id(*context_elements):
        """
        Get element id from information context
        """

        md5_result = md5()

        # remove None values from context_elements
        context_elements = [ce for ce in context_elements if ce is not None]

        for context_element in context_elements:
            if context_element is None:
                break
            md5_result.update(context_element.encode('ascii', 'ignore'))

        # resolve md5
        result = md5_result.hexdigest()

        return result
