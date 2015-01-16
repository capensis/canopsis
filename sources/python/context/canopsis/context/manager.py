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

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.storage import Storage

CONF_RESOURCE = 'context/context.conf'  #: last context conf resource
CATEGORY = 'CONTEXT'  #: context category


@add_category(CATEGORY)
@conf_paths(CONF_RESOURCE)
class Context(MiddlewareRegistry):
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
        - connector: name is connector and type is connector
    """

    DATA_SCOPE = 'context'  #: default data scope

    CTX_STORAGE = 'ctx_storage'  #: ctx storage name
    CONTEXT = 'context'

    TYPE = 'type'  #: entity type field name
    NAME = Storage.DATA_ID  #: entity name field name
    EXTENDED = 'extended'  #: extended field name

    DEFAULT_CONTEXT = [
        TYPE, 'connector', 'connector_name', 'component', 'resource'
    ]

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

    def get_entities(self, ids):
        """
        Get entities by id

        :param ids: one id or a set of ids.
        """

        return self[Context.CTX_STORAGE].get_elements(ids=ids)

    def get_entity(
        self, event, from_db=False, create_if_not_exists=False, cache=False
    ):
        """
        Get event entity.

        :param bool from_base: If True (False by default), check return entity
            from base, otherwise, return entity information from the event.
        :param bool create_if_not_exists: Create the event entity if it does
            not exists (False by default).
        :param bool cache: use query cache if True (False by default).
        """

        result = {}

        _type = event['source_type']

        if Context.NAME in event:
            name = event[Context.NAME]
        else:
            name = event[_type]

        # get the right context
        context = {Context.TYPE: _type}
        for ctx in self.context:
            if ctx in event:
                context[ctx] = event[ctx]
        # remove field which is the name
        if _type in context:
            del context[_type]

        if from_db:
            result = self.get(_type=_type, names=name, context=context)

        else:
            result = context.copy()
            result[Context.NAME] = name

        # if entity does not exists, create it if specified
        if result is None and create_if_not_exists:
            result = {Context.NAME: name}
            self.put(_type=_type, entity=result, context=context, cache=cache)

        return result

    def get_by_id(
        self,
        ids=None, limit=0, skip=0, sort=None, with_count=False
    ):
        """
        Get a list of entities where id are input ids.

        :param ids: element ids or an element id to get if is a string.
        :type ids: list of str

        :param int limit: max number of elements to get.
        :param int skip: first element index among searched list.
        :param sort: contains a list of couples of field (name, ASC/DESC)
            or field name which denots an implicitelly ASC order.
        :type sort: list of {(str, {ASC, DESC}}), or str}
        :param bool with_count: If True (False by default), add count to the
            result.

        :return: a Cursor of input id elements, or one element if ids is a
            string (None if this element does not exist).
        :rtype: Cursor of dict elements or dict or NoneType
        """

        result = self[Context.CTX_STORAGE].get_elements(
            ids=ids,
            limit=limit,
            skip=skip,
            sort=sort,
            with_count=with_count
        )

        return result

    def get(self, _type, names, context=None, extended=False):
        """
        Get entities by name.

        :param str _type: entity type (connector, component, etc.)
        :param str names: entity names.
        :param dict context: entity context such as couples of name, value.
        :param bool extended: get extended entities if entity is shared.

        :return: one element (dict or None) if names is a str and not extended.
            List of elements if names is a list or extended and names is a str.
            List of list of elements if names is a list and extended.
        :rtype: dict, list or None
        """

        path = {Context.TYPE: _type}

        if context is not None:
            path = context.copy()
            path[Context.TYPE] = _type

        result = self[Context.CTX_STORAGE].get(
            path=path, data_ids=names, shared=extended
        )

        return result

    def find(
        self, _type=None, context=None, _filter=None, extended=False,
        limit=0, skip=0, sort=None, with_count=False
    ):
        """
        Find all entities which of input _type and context with an additional
        filter.

        :param extended: get extended entities if they are shared.
        """

        path = {}

        if context is not None:
            path.update(context)

        if _type is not None:
            path[Context.TYPE] = _type

        result = self[Context.CTX_STORAGE].get(
            path=path, _filter=_filter, shared=extended,
            limit=limit, skip=skip, sort=sort, with_count=with_count
        )

        return result

    def put(
        self,
        _type, entity, context=None, extended_id=None, add_parents=True,
        cache=False
    ):
        """
        Put an entity designated by the _type, context and entity.

        If parent entities do not exist, create them automatically.

        :param bool add_parents: ensure to add parents if child do not exists.
        :param bool cache: use query cache if True (False by default).
        """

        path = {Context.TYPE: _type}

        if context is not None:
            path.update(context)
            path[Context.TYPE] = _type

        name = entity[Context.NAME]

        entity_db = self.get(
            _type=_type, names=name, context=context
        )
        # check if entity exists in db
        if add_parents and context is not None and entity_db is None:
            # if entity does not exist in db

            # check if parent entities exists
            # path is a copy of context
            parent_path = path.copy()
            # ensure all parent context exist, or create them if necessary
            # get key context without type
            keys = self._context[1:]
            keys.reverse()
            for key in keys:
                if key in context:
                    # update path type with input key
                    parent_path['type'] = key
                    # parent name is path[key]
                    parent_name = parent_path[key]
                    # del path[key] in order to avoid wrong path resolution
                    del parent_path[key]
                    # get entity
                    parent_entity = self[Context.CTX_STORAGE].get(
                        path=parent_path, data_ids=parent_name
                    )
                    # if entity does not exist
                    if parent_entity is None:
                        # put a new entity in DB
                        parent_entity = {Context.NAME: parent_name}
                        self[Context.CTX_STORAGE].put(
                            path=parent_path,
                            data_id=parent_name,
                            data=parent_entity,
                            cache=cache
                        )
                    else:
                        break

        # initialize entity db for future update
        if entity_db is None:
            entity_db = {}

        # check if an update is necessary in comparing entity fields with
        to_update = False
        # if entity different than entity_db
        for field in entity:
            value = entity[field]
            to_update = (field not in entity_db) or entity_db[field] != value
            if to_update:
                break

        if to_update:
            # finally, put the entity if necessary
            self[Context.CTX_STORAGE].put(
                path=path,
                data_id=name,
                data=entity,
                shared_id=extended_id,
                cache=cache
            )

    def remove(
        self, ids=None, _type=None, context=None, extended=False, cache=False
    ):
        """
        Remove a set of elements identified by element_ids, an element type or
        a timewindow. If ids, _type and context are all None, remove all
        elements.

        :param bool cache: use query cache if True (False by default).
        """

        path = {}

        if context is not None:
            path.update(context)

        if _type is not None:
            path[Context.TYPE] = _type

        if path:
            self[Context.CTX_STORAGE].remove(
                path=path, shared=extended, cache=cache
            )

        if ids is not None:
            self[Context.CTX_STORAGE].remove_elements(ids=ids, cache=cache)
        # if all parameters are None, delete all elements
        elif (_type, context) == (None, None):
            self[Context.CTX_STORAGE].remove_elements()

    def get_entity_context_and_name(self, entity):
        """
        Get the right context related to input entity
        """

        result = self[Context.CTX_STORAGE].get_path_with_id(entity)

        return result

    def get_entity_id(self, entity):
        """
        Get unique entity id related to its context and name.
        """

        path, data_id = self.get_entity_context_and_name(entity=entity)

        result = self[Context.CTX_STORAGE].get_absolute_path(
            path=path, data_id=data_id
        )

        return result

    def unify_entities(self, entities, extended=False, cache=False):
        """
        Unify input entities as the same entity.

        :param bool cache: use query cache if True (False by default).
        """

        self[Context.CTX_STORAGE].share_data(
            data=entities, shared=extended, cache=cache
        )

    def _configure(self, unified_conf, *args, **kwargs):

        super(Context, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        if Context.CTX_STORAGE in self:
            self[Context.CTX_STORAGE].path = self.context
