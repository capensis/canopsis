#!/usr/bin/env python
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

from bottle import get, delete, put, post

from canopsis.common.ws import response, route
from canopsis.context.manager import Context

manager = Context()


@route(get)
def context(_type, name=None, context=None, extended=None):

    entities = manager.get(
        _type=_type, name=name, context=context, extended=extended)

    result = response(entities)

    return result


@route(get)
def context_find(
    _type=None, context=None, _filter=None, extended=False,
    limit=0, skip=0, sort=None
):

    entities = manager.find(
        _type=_type, context=context, _filter=_filter, extended=extended,
        limit=limit, skip=skip, sort=sort)

    result = response(entities)

    return result


@route(put)
def context(_type, entity, context=None, extended_id=None):

    manager.put(
        _type=_type, entity=entity, context=context, extended_id=extended_id)

    result = response(entity)

    return result


@route(delete)
def delete(ids=None, _type=None, context=None, extended=False):

    manager.remove(ids=ids, _type=_type, context=context, extended=extended)

    result = response(None)

    return result


@route(post)
def unify(entities, extended=False):

    nodes = manager.unify_entities(entities=entities, extended=extended)

    result = response(nodes)

    return result
