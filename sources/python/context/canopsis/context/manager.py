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

from canopsis.configuration import conf_paths, add_category
from canopsis.middleware.manager import Manager


CONF_RESOURCE = 'context/context.conf'
CATEGORY = 'CONTEXT'


@add_category(CATEGORY)
@conf_paths(CONF_RESOURCE)
class Context(Manager):
    """
    Manage access to a context (connector, component, resource) elements
    and context data (metric, downtime, etc.)

    It uses a composite storage in order to modelise composite data.

    For example, let a resource ``R`` in the component ``C`` and connector
    ``K`` is identified through the context [``K``, ``C``], the name ``R`` and
    the type ``resource``.

    In addition to those composable data, it is possible to extend two entities
    which have the same name and type but different context.

    For example, in following entities:
        - component: name is component_id and type is component
        - connector: data_id
    """

    DATA_SCOPE = 'context'

    CTX_STORAGE = 'ctx_storage'
    CONTEXT = 'context'

    TYPE = 'type'  #: entity type field name
    NAME = 'name'  #: entity name field name
    EXTENDED = 'extended'  #: extended field name

    DEFAULT_CONTEXT = [
        'type', 'connector', 'connector_name', 'component', 'resource']

    def __init__(
        self, context=DEFAULT_CONTEXT, ctx_storage=None, *args, **kwargs
    ):

        super(Context, self).__init__(self, *args, **kwargs)

        self._context = context
        if ctx_storage is not None:
            self[Context.CTX_STORAGE] = ctx_storage

    @property
    def context(self):
        """
        List of context element name.
        """
        return self._context

    @context.setter
    def context(self, value):
        self._context = value

    def get_entity_context(self, entity, context=None):
        """
        Get the right context related to input entity
        """

        result = {}

        if context is None:
            context = self.context

        for value in context:
            if value in entity:
                result[value] = entity[value]

        return result

    def get_entities(self, ids):
        """
        Get entities by id

        :param ids: one id or a set of ids
        """

        return self[Context.CTX_STORAGE].get_elements(ids=ids)

    def get(self, _type, name, context=None, extended=False):
        """
        Get one entity related to:
            - an entity type.
            - a context (connector, ..., other).
            - a name.

        :param extended: get extended entities if they are shared.

        :return: one element or None
        :rtype: dict
        """

        path = {
            Context.TYPE: _type,
        }

        if context is not None:
            path.update(context)

        result = self[Context.CTX_STORAGE].get(path=path, ids=name)

        if extended:
            result = self[Context.CTX_STORAGE].get_shared_data(data=result)

        return result

    def find(
        self, context=None, _filter=None, extended=False,
        limit=0, skip=0, sort=None
    ):
        """
        Find all entities which of input _type and context with an additional
        filter.

        :param extended: get extended entities if they are shared
        """

        path = {} if context is None else context

        result = self[Context.CTX_STORAGE].get(
            path=path, _filter=_filter, limit=limit, skip=skip, sort=sort)

        if extended:
            result = self[Context.CTX_STORAGE].get_shared_data(data=result)

        return result

    def put(self, _type, entity, context=None):
        """
        Put an element designated by the element_id, element_type and element.
        If timestamp is None, time.now is used.
        """

        path = {
            Context.TYPE: _type
        }

        if context is not None:
            path.update(context)

        name = entity[Context.NAME]

        self[Context.CTX_STORAGE].put(path=path, _id=name, data=entity)

    def remove(self, ids=None, _type=None, context=None):
        """
        Remove a set of elements identified by element_ids, an element type or
        a timewindow
        """

        path = {}

        if _type is not None:
            path[Context.TYPE] = _type

        if context is not None:
            path.update(context)

        self[Context.CTX_STORAGE].remove(path=path, ids=ids)

    def get_entity_id(self, entity, context=None):
        """
        Get unique entity id related to its context and name.
        """

        if context is None:
            context = self.context

        path = self.get_entity_context(entity=entity, context=context)

        name = entity['name']

        result = self[Context.CTX_STORAGE].get_absolute_path(
            path=path, _id=name)

        return result

    def unify_entities(self, entities):
        """
        Unify input entities as the same entity
        """

        # get unique and shared id
        pass

    def _configure(self, unified_conf, *args, **kwargs):

        super(Context, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        if Context.CTX_STORAGE in self:
            self[Context.CTX_STORAGE].path = self.context
