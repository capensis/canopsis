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
from uuid import uuid4
import json as j
import os

manager = ContextGraph()
import_manager = ContextGraphImport()

__FILE = "~/tmp/import-{0}.json"
__IMPORT_ID = "import_id"
__ERROR = "error"
__STORE_ERROR = "Impossible to store the file on the disk : {0}."
__OTHER_ERROR = "An error occured : {0}."
__CANNOT_EXEC_IMPORT = "Error while calling the process responsible"\
                       " of the import"

def get_uuid():
    """Return an UUID never used for an import. If the generated UUID is already
    used, try again until an UUID not used is created"""

    uuid = uuid4()
    while not import_manager.check_id(uuid):
        uuid = uuid4()

    return str(uuid)

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
        uuid = get_uuid()
        try:
            file_ = __FILE.format(uuid)

            if os.path.exists(file_):
                return {__ERROR: __STORE_ERROR.format("A file with the same "\
                                                      "name already exists")}

            with open(file_, 'x') as fd:
                j.dump(json, fd)

            status = os.spawnl(os.P_NOWAIT, "import.py", file_)

            if status == 127:
                return {__ERROR: __CANNOT_EXEC_IMPORT}

            return {__IMPORT_ID : str(uuid)}

        except IOError as ioerror:
            return {__ERROR: __STORE_ERROR.format(str(ioerror))}

        except:
            {__ERROR: __OTHER_ERROR.format(str(ioerror))}

    @route(
        ws.application.get,
        name='truc/machin'
    )
    def get_graph():
        entities_list = manager.get_entities()

        ret_json = {
            'links':[],
            'nodes':[]
        }

        for i in entities_list:
            ret_json['nodes'].append({'group':1, 'id': i['_id']})

        for i in entities_list:
            source = i['_id']
            for target in i['depends']:
                ret_json['links'].append({'value': 1, 'source': source, 'target': target})

        f = open('/opt/canopsis/tmp/graph.json', 'w')
        f.write(j.dumps(ret_json))
        f.close()

        return ret_json

    @route(
        ws.application.get,
        name='getgraphimpact',
        payload=['_id', 'deepness']
    )
    def get_graph_impact(_id, deepness=None):
        return manager.get_graph_impact(_id, deepness)

    @route(
        ws.application.get,
        name='getgraphdepends',
        payload=['_id', 'deepness']
    )
    def get_graph_depends(_id, deepness=None):
        return manager.get_graph_depends(_id, deepness)

    @route(
        ws.application.get,
        name='getleavesdepends',
        payload=['_id', 'deepness']
    )
    def get_leaves_depends(_id, deepness=None):
        return manager.get_leaves_depends(_id, deepness)

    @route(
        ws.application.get,
        name='getleavesimpact',
        payload=['_id', 'deepness']
    )
    def get_leaves_impact(_id, deepness=None):
        return manager.get_leaves_impact(_id, deepness)
