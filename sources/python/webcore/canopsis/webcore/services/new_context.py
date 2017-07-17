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

    @ws.application.route(
        '/api/v2/context/_filter:=_filter&start:=start&sort:=sort&limit:=limit'
    )
    def context(
        _filter=None,
        limit=0,
        start=0,
        sort=None
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
