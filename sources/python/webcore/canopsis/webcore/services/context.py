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

from canopsis.common.ws import route
from canopsis.context_graph.manager import ContextGraph

manager = ContextGraph()


def exports(ws):

    @route(ws.application.get)
    def context(_type, names=None, context=None, extended=None):
        if names:
            names = [n.strip() for n in names.split(',')]
        """
        result = manager.get(
            _type=_type, names=names, context=context, extended=extended)
        """
        # this is a test before adapter refactoring
        result = manager.get_entities_by_id(names)
        return result

    @route(ws.application.get, name='context/ids')
    @route(
        ws.application.post,
        payload=['ids', 'limit', 'start', 'sort', 'with_count'],
        name='context/ids'
    )
    def context_by_id(
        ids=None, limit=0, start=0, sort=None, with_count=False
    ):
        """
        result = manager.get(
            ids=ids,
            limit=limit,
            skip=start,
            sort=sort,
            with_count=with_count
        )
        """
        result = manager.get_entities_by_id(ids)
        return result

    @route(ws.application.post, payload=['limit', 'start', 'sort', '_filter'])
    def context(
        context=None, _filter=None, extended=False,
        limit=0, start=0, sort=None
    ):

        query = {}
        if _filter is not None:
            query.update(_filter)

        cursor, count = manager.get_entities(
            query=query,
            limit=limit,
            start=start,
            sort=sort,
            with_count=True
        )

        data = []
        for ent in cursor:
            data.append(ent)

        return data, count

    @route(ws.application.put, payload=[
        '_type', 'entity', 'context', 'extended_id'
    ])
    def context(_type, entity, context=None, extended_id=None):
        """
        manager.put(
            _type=_type,
            entity=entity,
            context=context,
            extended_id=extended_id
        )
        """
        manager.create_entity(entity=entity)
        return entity

    @route(ws.application.delete, payload=[
        'context', 'ids', '_type', 'extended'
    ])
    def context(ids=None, _type=None, context=None, extended=False):
        """
        manager.remove(
            ids=ids,
            _type=_type,
            context=context,
            extended=extended
        )"""
        manager.delete_entity(
            ids
        )
