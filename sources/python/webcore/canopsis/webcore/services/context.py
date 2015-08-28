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
from canopsis.context.manager import Context

manager = Context()


def exports(ws):

    @route(ws.application.get)
    def context(_type, names=None, context=None, extended=None):
        if names:
            names = [n.strip() for n in names.split(',')]

        result = manager.get(
            _type=_type, names=names, context=context, extended=extended)

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

        result = manager.get_by_id(
            ids=ids,
            limit=limit,
            skip=start,
            sort=sort,
            with_count=with_count
        )

        return result

    @route(ws.application.post, payload=['limit', 'start', 'sort', '_filter'])
    def context(
        _type=None, context=None, _filter=None, extended=False,
        limit=0, start=0, sort=None
    ):

        result = manager.find(
            _type=_type,
            context=context,
            _filter=_filter,
            extended=extended,
            limit=limit,
            skip=start,
            sort=sort,
            with_count=True
        )

        return result

    @route(ws.application.put, payload=[
        '_type', 'entity', 'context', 'extended_id'
    ])
    def context(_type, entity, context=None, extended_id=None):
        manager.put(
            _type=_type,
            entity=entity,
            context=context,
            extended_id=extended_id
        )

        return entity

    @route(ws.application.delete, payload=[
        'context', 'ids', '_type', 'extended'
    ])
    def context(ids=None, _type=None, context=None, extended=False):
        manager.remove(
            ids=ids,
            _type=_type,
            context=context,
            extended=extended
        )

    @route(ws.application.post, payload=['entities', 'extended'])
    def unify(entities, extended=False):
        result = manager.unify_entities(entities=entities, extended=extended)

        return result
