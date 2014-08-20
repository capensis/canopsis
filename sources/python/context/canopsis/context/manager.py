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

from canopsis.configuration import conf_paths
from canopsis.middleware.manager import Manager
from md5 import new as md5

CONF_RESOURCE = 'context/context.conf'


@conf_paths(CONF_RESOURCE)
class Context(Manager):
    """
    Manage access to a context (connector, component, resource) elements
    and context data (metric, downtime, etc.)

    For example, in following entities:
        - component: data_id is component_id and data_type is component
        - connector: data_id
    """

    DATA_SCOPE = 'context'

    CATEGORY = 'CONTEXT'

    CTX_STORAGE = 'ctx_storage'
    CONTEXT = 'context'

    DEFAULT_CONTEXT = [
        'type', 'connector', 'connector_name', 'component', 'resource']

    def __init__(self, scope=DEFAULT_CONTEXT, ctx_storage=None, *args, **kwargs):

        super(Context, self).__init__(self, *args, **kwargs)

        self.scope = scope
        self['ctx_storage'] = ctx_storage

    @property
    def scope(self):
        return self.scope

    @scope.setter
    def scope(self, value):
        self.scope = value

    def get(self, _type, name, context=None, *args, **kwargs):
        """
        Get one element related to:
            - an element_type,
            - a path (connector, ..., other),
            - a timewindow

        :return: one element or None
        :rtype: dict
        """

        scope = {
            'type': _type,
        }

        if context is not None:
            scope.update(context)

        result = self['ctx_storage'].get(scope=scope, ids=name)

        return result

    def find(self, _type, context, _filter, *args, **kwargs):
        """
        Find all elements which have an element_id among input element_ids, or
        type equals to element_type or inside timewindow if specified
        """

        if element_ids is None:
            element_id = Context.get_element_id(
                connector, connector_type, component, resource, other)
            element_ids = [element_id]

        return self._ctx.get(
            data_ids=element_ids, data_type=_type,
            *args, **kwargs)

    def get_by_name(
        self,
        element_type, name=None, *args, **kwargs
    ):

        return self.get(element_id=name, element_type=element_type)

    def put(
        self, element_id, element_type, element,
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

    def _configure(self, unified_conf, *args, **kwargs):

        super(Context, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        self['ctx_storage'].scope = self.scope

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
