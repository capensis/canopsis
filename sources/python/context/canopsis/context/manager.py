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

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.configuration.model import Parameter

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.event import Event, forger
from canopsis.storage.composite import CompositeStorage

from urllib import unquote_plus


CONF_RESOURCE = 'context/context.conf'  #: last context conf resource
CATEGORY = 'CONTEXT'  #: context category
CONTENT = [Parameter('accept_event_types', Parameter.array())]


@add_category(CATEGORY, content=CONTENT)
@conf_paths(CONF_RESOURCE)
class Context(MiddlewareRegistry):
    """
    Manage access to a context (connector, component, resource) elements
    and context data (metric, downtime, etc.)

    It uses a composite storage in order to modelise composite data.

    For example, let a resource ``R`` in the component ``C`` and connector
    ``K``. ``R`` is identified through the context [``K``, ``C``],
    the name ``R`` and the type ``resource``.

    In addition to those composable data, it is possible to extend two entities
    which have the same name and type but different context.

    For example, in following entities:
        - component: name is component_id and type is component.
        - connector: name is connector and type is connector.
    """

    DATA_SCOPE = 'context'  #: default data scope

    CTX_STORAGE = 'ctx_storage'  #: ctx storage name
    CONTEXT = 'context'

    DATA_ID = '_id'  #: temporary id field
    TYPE = 'type'  #: entity type field name
    NAME = CompositeStorage.NAME  #: entity name field name
    EXTENDED = 'extended'  #: extended field name

    DEFAULT_CONTEXT = [
        TYPE, 'connector', 'connector_name', 'component', 'resource'
    ]

    ENTITY = Event.ENTITY  #: entity id in event

    EID = 'eid'  #: entity id.

    def __init__(
            self, context=DEFAULT_CONTEXT, ctx_storage=None, *args, **kwargs
    ):

        super(Context, self).__init__(self, *args, **kwargs)

        self._context = context

        if ctx_storage is not None:
            self[Context.CTX_STORAGE] = ctx_storage

    @property
    def context(self):
        """List of context element name."""

        return self._context

    @context.setter
    def context(self, value):

        self._context = value

    @property
    def accept_event_types(self):
        if not hasattr(self, '_accept_event_types'):
            self.accept_event_types = None

        return self._accept_event_types

    @accept_event_types.setter
    def accept_event_types(self, value):
        if value is None:
            value = [
                'perf',
                'check',
                'ack',
                'ackremove',
                'declareticket',
                'assocticket',
                'cancel',
                'uncancel',
                'changestate',
                'snooze'
                'downtime',
                'snooze',
                'comment'
            ]

        self._accept_event_types = value

    def get_entities(self, ids):
        """Get entities by id.

        :param ids: one id or a set of ids.
        :type ids: list or str
        """

        return self[Context.CTX_STORAGE].get_elements(ids=ids)

    def iter_ids(self):
        """Returns a cursor on all context ids.
        """

        cursor = self[Context.CTX_STORAGE].get_elements(projection={
            Context.DATA_ID: True
        })

        for doc in cursor:
            yield doc[Context.DATA_ID]

    def clean(self, entity):
        """Remove entity properties which are not in self.context."""

        result = {}

        for ctx in self.context:
            if ctx in entity:
                result[ctx] = entity[ctx]
            else:
                break

        result[Context.NAME] = entity[Context.NAME]

        return result

    def get_children(self, entity):
        """Get children entities of input entity.

        For example, if entity is a component, you can retrieve component
        resources and metrics.

        :param dict entity: parent entity.
        :return: list of children entity.
        :rtype: list
        """

        # construct query with parent fields which can be found among children
        query = self.clean(entity)
        del query[Context.TYPE]

        # iterate on context information without the type
        for ctx in self.context[1:]:
            if ctx not in entity:
                query[ctx] = query.pop(Context.NAME)
                break

        # execute query in order to get children
        children = self[Context.CTX_STORAGE].find_elements(query=query)

        return children

    # FIXME find a better replacement name.
    def get_entity_old(
            self, event, from_db=False, create_if_not_exists=False, cache=False
    ):
        """Get event entity.

        :param bool from_base: If True (False by default), check return entity
            from base, otherwise, return entity information from the event.
        :param bool create_if_not_exists: Create the event entity if it does
            not exists (False by default).
        :param bool cache: use query cache if True (False by default).
        """

        result = {}

        _event = event.copy()

        # get the right type which is type, or event_type or component/resource
        # if event_type is not an entity
        _type = _event['source_type']
        # try to get the right type if the event corresponds to the old system
        if _type in self.context:
            event_type = _event['event_type']
            if event_type not in self.accept_event_types:
                _type = event_type

        # set type in event
        _event[Context.TYPE] = _type

        # set name if not given
        if Context.ENTITY not in _event:
            for ctx in reversed(self.context):
                if ctx in _event and _event[ctx]:
                    _event[Context.NAME] = _event[ctx]
                    del _event[ctx]
                    break
        else:  # delete ctx field which matches with the ctx
            for ctx in reversed(self.context):
                if ctx in _event:
                    if _event[ctx] == _event[Context.ENTITY]:
                        del _event[ctx]
                    break

        ctx, name = self.get_entity_context_and_name(_event)

        # remove type from ctx
        _type = ctx[Context.TYPE]
        del ctx[Context.TYPE]

        if from_db:
            result = self.get(_type=_type, names=name, context=ctx)
            # if entity does not exists, create it if specified
            if result is None and create_if_not_exists:
                result = {Context.NAME: name}
                self.put(_type=_type, entity=result, context=ctx, cache=cache)
                result.update(ctx)
                result[Context.TYPE] = _type
        else:
            result = ctx.copy()
            result[Context.NAME] = name
            result[Context.TYPE] = _type

        return result

    def get_name(self, entity_id, _type=None):
        """Get entity name related to an entity_id.

        :param str entity_id: entity id (in the form '/a/b/c/...').
        :param str _type: context type such as component, resource, other, etc.
            If None, get last entity name in entity_id.
        """

        result = None

        names = entity_id.split('/')[1:]  # get names after the first '/'
        try:
            index = self.context.index(_type)  # get name index
        except ValueError:
            result = names[-1]  # in case of index error, get the last name
        else:
            if index < len(names):  # else get names[index]
                result = names[index]

        result = unquote_plus(result)

        return result

    def get_entity_by_id(self, _id, _type=None):
        """Generate an entity related to input id.

        :param str _id: entity id from where get entity properties.
        """

        result = {}
        # get ctx values from _id

        result = self[Context.CTX_STORAGE].get_data_from_id(data_id=_id)

        if _type is None:
            _type = result[self.TYPE]

        else:
            result[self.TYPE] = _type

        if _type in result:
            result[Context.NAME] = result.pop(_type)

        elif 'extended' in result:
            result[Context.NAME] = result.pop('extended')

        else:
            result[Context.NAME] = result.pop(self.context[len(result) - 1])

        return result

    def get_event(self, entity, event_type='check', **kwargs):
        """Get an event from an entity.

        :param dict entity: entity to convert to an event.
        :param str event_type: specific event_type. Default check.
        :param dict kwargs: additional fields to put in the event.
        :rtype: dict
        """

        kwargs['event_type'] = event_type

        # In some cases, name is present but is component in fact
        if 'name' in entity:
            if 'component' not in entity:
                entity['component'] = entity['name']
            entity.pop('name')

        # fill kwargs with entity values
        for field in entity:
            kwargs[field] = entity[field]

        # forge the event
        result = forger(**kwargs)

        return result

    def get_by_id(
            self,
            ids=None, limit=0, skip=0, sort=None, with_count=False
    ):
        """Get a list of entities where id are input ids.

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
        """Get entities by name.

        :param str _type: entity type (connector, component, etc.)
        :param names: entity name(s).
        :type names: list or str
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
            path=path, names=names, shared=extended
        )

        return result

    def find(
            self, _type=None, context=None, _filter=None, extended=False,
            limit=0, skip=0, sort=None, with_count=False
    ):
        """Find all entities which of input _type and context with an additional
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

        path = {}

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
            keys = [k for k in reversed(self._context[1:]) if k in context]
            for key in keys:
                # update path type with input key
                parent_path['type'] = key
                # parent name is path[key]
                parent_name = parent_path.pop(key)
                # get entity
                parent_entity = self[Context.CTX_STORAGE].get(
                    path=parent_path, names=parent_name
                )
                # if entity does not exist
                if parent_entity is None:
                    # put a new entity in DB
                    parent_entity = {Context.NAME: parent_name}
                    self[Context.CTX_STORAGE].put(
                        path=parent_path,
                        name=parent_name,
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
                name=name,
                data=entity,
                shared_id=extended_id,
                cache=cache
            )

    def remove(
            self, ids=None, _type=None, context=None, extended=False,
            cache=False
    ):
        """Remove a set of elements identified by element_ids, an element type
        or a timewindow. If ids, _type and context are all None, remove all
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
        """Get the right context related to input entity.

        :param dict entity: must contain at least name and type fields.
        """

        _entity = entity.copy()

        # ensure name exists once
        # _name = entity[Context.NAME]
        # remove useless fields where ctx field is the name
        """for ctx in reversed(self.context):
            if ctx in _entity and _entity[ctx] == _name:
                del _entity[ctx]
                break
        """

        result = self[Context.CTX_STORAGE].get_path_with_name(_entity)

        return result

    def get_entity_id(self, entity):
        """Get unique entity id related to its context and name.

        :param dict entity: must contain at least name and type fields.
        """

        path, name = self.get_entity_context_and_name(entity=entity)

        result = self[Context.CTX_STORAGE].get_absolute_path(
            path=path, name=name
        )

        return result

    def get_entity_id_context_name(self, entity):
        """Get the right id, context and name of input entity."""

        path, name = self.get_entity_context_and_name(entity=entity)

        result = self[Context.CTX_STORAGE].get_absolute_path(
            path=path, name=name
        ), path, name

        return result

    def unify_entities(self, entities, extended=False, cache=False):
        """Unify input entities as the same entity.

        :param bool cache: use query cache if True (False by default).
        """

        self[Context.CTX_STORAGE].share_data(
            data=entities, shared=extended, cache=cache
        )

    def _configure(self, unified_conf, *args, **kwargs):

        super(Context, self)._configure(
            unified_conf=unified_conf, *args, **kwargs
        )

        if Context.CTX_STORAGE in self:
            self[Context.CTX_STORAGE].path = self.context
