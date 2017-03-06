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
    def all():
        """
            :return all json for d3 representation
        """
    @route(ws.application.get)
    def get_entities(eids, is_active=True, **kwargs):
        """
            :return get_entities
        """

    @route(ws.application.put)
    def put_entities(entities):
        """
            put entities in db
        """
    @route(ws.application.post)
    def update_entity(entity):
        """
            update entity in db
        """
    
    @route(ws.application.delete)
    def delete_entity(eid):
        """
            remove  etity
        """
