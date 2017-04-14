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
from canopsis.context_graph.import_ctx import ContextGraphImport
import json as j

manager = ContextGraph()
import_manager = ContextGraphImport()

def exports(ws):

    @route(ws.application.get, name='context_graph/all')
    def all():
        """
            :return all json for d3 representation
        """
        return manager.get_entities()

    @route(
        ws.application.put,
        payload=['entity']
    )
    def put_entities(entity):
        """
            put entities in db
        """
        return manager.create_entity(entity)

    @route(
        ws.application.post, 
        payload=['entity']    
    )
    def update_entity(id_, entity):
        """
            update entity in db
        """
        return manager.update_entity(id_, entity)

    
    @route(
        ws.application.delete,
        payload=['id_']
    )
    def delete_entity(id_):
        """
            remove  etity
        """
        return manager.delete_entity(id_)


    @route(
        ws.application.get,
        payload=['query', 'projection', 'limit', 'sort', 'with_count']
    )
    def get_entities(
            query={},
            projection={},
            limit=0, 
            sort=False,
            with_count=False
    ):
        return get_entities(
            query=query,
            projection=projection,
            limit=limit,
            sort=sort, 
            with_count=with_count
        )
    
    @route(
        ws.application.put,
        name='coucou/bouh',
        payload=['json']
    )
    def put_graph(json='{}'):
        if isinstance(json, dict):
            import_manager.import_context(json)
        elif isinstance(json, str):
            js = j.loads(json)
            import_manager.import_context(js)

