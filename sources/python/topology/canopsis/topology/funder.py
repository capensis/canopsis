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

from canopsis.topology.elements import TopoVertice, Vertice
from canopsis.topology.manager import TopologyManager
from canopsis.ctxinfo.funder import CTXInfoFunder


class TopologyFunder(CTXInfoFunder):
    """In charge of binding a topology information to context entities.
    """

    __datatype__ = 'topology'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(TopologyFunder, self).__init__(*args, **kwargs)

        self.manager = TopologyManager()

    def _get_documents(self, entity_ids, query, *args, **kwargs):

        result = []

        entity_id_field = self._entity_id_field()

        ENTITY = TopoVertice.ENTITY
        INFO = Vertice.INFO

        entity = (
            {'$exists': True} if entity_ids is None else {'$in': entity_ids}
        )
        info = {ENTITY: entity}

        docs = self.manager.get_elts(
            info=info, serialize=False, query=query
        )

        for doc in docs:
            doc[entity_id_field] = doc[INFO][ENTITY]
            result.append(doc)

        return result

    def _get(self, entity_ids, query, *args, **kwargs):

        return self._get_documents(entity_ids=entity_ids, query=query)

    def _delete(self, entity_ids, query, *args, **kwargs):

        result = self._get_documents(entity_ids=entity_ids, query=query)

        ids = [doc['_id'] for doc in result]
        self.manager.del_elts(ids=ids)

        return result

    def entity_ids(self, query=None):

        result = set()

        elts = self.manager.get_elts(query=query)

        for elt in elts:
            result.add(elt.entity)

        return result
