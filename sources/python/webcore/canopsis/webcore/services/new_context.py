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
from bottle import request
from json import loads
from canopsis.webcore.utils import gen_json_error, HTTP_ERROR, gen_json
manager = ContextGraph()


def exports(ws):

    @ws.application.route(
        '/api/v2/context/<_filter>'
    )
    def context(
        _filter=None,
    ):
        ws.logger.critical(_filter)
        limit = 0
        sort = None
        start = 0
        payload = {}
        if request.json is not None:
            payload = request.json
        if 'limit' in payload.keys():
            limit = payload['limit']
        if 'start' in payload.keys():
            start = payload['skip']
        if 'sort' in payload.keys():
            sort = payload['sort']


        query = {}
        if _filter is not None:
            try:
                query = loads(_filter)
            except ValueError:
                return gen_json_error({'description': 'can t load filter'}, HTTP_ERROR)


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

        return gen_json(data)
