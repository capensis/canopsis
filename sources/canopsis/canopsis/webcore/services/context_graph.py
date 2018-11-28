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

from bottle import request
import json as j
import os
from uuid import uuid4

from canopsis.common import root_path
from canopsis.common.ws import route
from canopsis.alerts.manager import Alerts
from canopsis.common.converters import id_filter
from canopsis.confng import Configuration, Ini
from canopsis.context_graph.import_ctx import ImportKey, Manager
from canopsis.context_graph.manager import ContextGraph
from canopsis.webcore.utils import gen_json, gen_json_error, HTTP_ERROR

import_col_man = Manager(Configuration.load(Manager.CONF_FILE, Ini))
alerts_manager = Alerts(*Alerts.provide_default_basics())


__IMPORT_ID = "import_id"
__ERROR = "error"
__OTHER_ERROR = "An error occured : {0}."
__EVT_ERROR = "error while sending a event to the task : {0}."
__STORE_ERROR = "Impossible to store the import: {0}."


event_body = {ImportKey.EVT_IMPORT_UUID: None,
              ImportKey.EVT_JOBID: None}

RK = "task_importctx"


def get_uuid():
    """Return an UUID never used for an import. If the generated UUID is already
    used, try again until an UUID not used is created"""

    uuid = uuid4()
    while import_col_man.check_id(uuid):
        uuid = uuid4()

    return str(uuid)


def exports(ws):

    ws.application.router.add_filter('id_filter', id_filter)

    manager = alerts_manager.context_manager

    @route(ws.application.get, name='contextgraph/all')
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
            projection=None,
            limit=0,
            sort=False,
            with_count=False
    ):
        return manager.get_entities(
            query=query,
            projection=projection,
            limit=limit,
            sort=sort,
            with_count=with_count
        )

    @route(
        ws.application.put,
        name='api/contextgraph/import',
        payload=['json']
    )
    def put_graph(json='{}'):
        uuid = get_uuid()
        # FIXME: A race condition may occur here
        import_col_man.create_import_status(uuid)

        file_ = ImportKey.IMPORT_FILE.format(uuid)

        if os.path.exists(file_):
            return {__ERROR: __STORE_ERROR.format(
                "an import already exist with the same id on the disk")}

        try:
            with open(file_, 'w') as fd:
                j.dump(json, fd)
        except IOError as ioerror:
            return {__ERROR: __STORE_ERROR.format(str(ioerror))}

        try:
            event = event_body.copy()
            event[ImportKey.EVT_IMPORT_UUID] = uuid
            event[ImportKey.EVT_JOBID] = ImportKey.JOB_ID.format(uuid)
            ws.amqp_pub.json_document(event, 'amq.direct', RK)
        except Exception as e:
            ws.logger.error(e)
            return {__ERROR: __EVT_ERROR.format(repr(e))}

        return {__IMPORT_ID: str(uuid)}

    def get_state(_id):
        """
            va chercher si il y a une alarme ouverte d'une entitée
            et si oui choppe l'etat si non return 0
        """
        al = alerts_manager.get_alarm_with_eid(_id, resolved=False)
        if al == []:
            return 0
        return al[0]['v']['state']['val']

    @route(
        ws.application.get,
        name='api/contextgraph/d3graph'
    )
    def get_graph():
        entities_list = manager.get_entities()

        entities_dico = {}
        for i in entities_list:
            entities_dico[i['_id']] = i

        ret_json = {
            'links': [],
            'nodes': []
        }

        for i in entities_list:
            ret_json['nodes'].append({'group': 1,
                                      'id': i['_id'],
                                      'name': i['name'],
                                      'state': get_state(i['_id'])})

        for i in entities_list:
            source = i['_id']
            for target in i['impact']:
                if entities_dico[source]['type'] == 'resource' and \
                   entities_dico[target]['type'] == 'connector':
                    pass
                else:
                    ret_json['links'].append({
                        'value': 1,
                        'source': source,
                        'target': target
                    })

        directory = os.path.join(root_path, 'var/www/src/canopsis/d3graph')
        if not os.path.exists(directory):
            os.makedirs(directory)
            htmldoc = """<!DOCTYPE html>
<meta charset="utf-8">
<style>

.links line {
  stroke: #999;
  stroke-opacity: 0.6;
}

.nodes circle {
  stroke: #fff;
  stroke-width: 0.5px;
}

</style>
<svg width="1000" height="900"></svg>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script>

var svg = d3.select("svg"),
    width = +svg.attr("width"),
    height = +svg.attr("height");

var color = d3.scaleOrdinal(d3.schemeCategory20);

var manybody = d3.forceManyBody()
    .strength(-500)

var simulation = d3.forceSimulation()
    .force("link", d3.forceLink().id(function(d) { return d.id; }))
    .force("charge", manybody)
    .force("center", d3.forceCenter(width / 2, height / 2));

d3.json("graph.json", function(error, graph) {
  if (error) throw error;

  var link = svg.append("g")
      .attr("class", "links")
    .selectAll("line")
    .data(graph.links)
    .enter().append("line")
      .attr("stroke-width", function(d) { return Math.sqrt(d.value); });

  var node = svg.append("g")
      .attr("class", "nodes")
    .selectAll("circle")
    .data(graph.nodes)
    .enter().append("circle")
      .attr("r", 5)
      .attr("fill", function(d) {if(d.state == 0){return "#A1D490"}if(d.state == 1){return "#ffff1a"}if(d.state == 2){return "#ff9900"}if(d.state == 3){return "#E30B1A"}})
      .call(d3.drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended));

  var text = svg.append("g")
      .attr("class", "text")
    .selectAll("text")
    .data(graph.nodes)
    .enter().append("text")
    .text(function(d) {return d.name});

  node.append("title")
      .text(function(d) { return d.id; });

  simulation
      .nodes(graph.nodes)
      .on("tick", ticked);

  simulation.force("link")
      .links(graph.links);

  function ticked() {
    link
        .attr("x1", function(d) { return d.source.x; })
        .attr("y1", function(d) { return d.source.y; })
        .attr("x2", function(d) { return d.target.x; })
        .attr("y2", function(d) { return d.target.y; });

    node
        .attr("cx", function(d) { return d.x; })
        .attr("cy", function(d) { return d.y; });
    text
        .attr("x", function(d) { return d.x + 5})
        .attr("y", function(d) { return d.y})
  }
});

function dragstarted(d) {
  if (!d3.event.active) simulation.alphaTarget(0.3).restart();
  d.fx = d.x;
  d.fy = d.y;
}

function dragged(d) {
  d.fx = d3.event.x;
  d.fy = d3.event.y;
}

function dragended(d) {
  if (!d3.event.active) simulation.alphaTarget(0);
  d.fx = null;
  d.fy = null;
}

</script>
"""
            a = open(os.path.join(directory, 'index.html', 'a'))
            a.write(htmldoc)
            a.close()

        f = open(os.path.join(directory, 'graph.json', 'w'))
        f.write(j.dumps(ret_json))
        f.close()

        return ret_json

    @route(
        ws.application.get,
        name='api/contextgraph/graphimpact',
        payload=['_id', 'deepness']
    )
    def get_graph_impact(_id, deepness=None):
        return manager.get_graph_impact(_id, deepness)

    @route(
        ws.application.get,
        name='api/contextgraph/graphdepends',
        payload=['_id', 'deepness']
    )
    def get_graph_depends(_id, deepness=None):
        return manager.get_graph_depends(_id, deepness)

    @route(
        ws.application.get,
        name='api/contextgraph/leavesdepends',
        payload=['_id', 'deepness']
    )
    def get_leaves_depends(_id, deepness=None):
        return manager.get_leaves_depends(_id, deepness)

    @route(
        ws.application.get,
        name='api/contextgraph/leavesimpact',
        payload=['_id', 'deepness']
    )
    def get_leaves_impact(_id, deepness=None):
        return manager.get_leaves_impact(_id, deepness)

    @ws.application.get('/api/contextgraph/import/status/<cid>')
    def get_status(cid):
        return import_col_man.get_import_status(cid)

    @ws.application.delete('/api/v2/context/<entity_id:id_filter>')
    def delete_entity_v2(entity_id):
        """
        Remove entity from context
        """
        try:
            res = manager.delete_entity(entity_id)
        except ValueError as vale:
            return gen_json_error({'description': str(vale)}, HTTP_ERROR)

        return gen_json(res)

    @ws.application.post('/api/v2/context_graph/get_id/')
    def get_entity_id():
        """
        Get the generated id tfrom an event.
        """
        event = request.json

        if event is None:
            return gen_json_error({'description': 'no event givent'},
                                  HTTP_ERROR)

        return gen_json(ContextGraph.get_id(event))

    @route(
        ws.application.get,
        name='api/v2/entities',
        payload=['query', 'limit', 'offset'],
        response=lambda x, **kwargs: x
    )
    def get_entities_with_open_alarms(query={}, limit=0, offset=0):
        """
        Return the entities filtered with a mongo filter.
        Each entity contain a list of the currently open alarms.
        """
        try:
            res = manager.get_entities_with_open_alarms(query, limit, offset)
            return res
        except Exception as err:
            return gen_json_error({'description': str(err)}, HTTP_ERROR)
