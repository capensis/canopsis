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

from canopsis.context.manager import Context
from canopsis.topology.format import formatter
from canopsis.topology.elements import Topology, TopoEdge, TopoNode
from canopsis.topology.manager import TopologyManager
from canopsis.task import new_conf
from canopsis.task.condition import condition as cond
from canopsis.topology.rule import action, condition


class Factory(object):
    """docstring for Factory"""
    def __init__(self, arg):
        self.arg = arg

    def create_topology(self, topo_name, topoEdge, topoNode):
        '''
            TODO
        '''
        # Initialize the Toplogy Manager
        manager = TopologyManager()
        # Create the topology name
        topo = Topology(_id=topo_name)
        # Add the topology Edge
        topo.add_elts(topoEdge)
        # Add the topology nodes
        topo.add_elts(topoNode)
        # Save the topology
        topo.save(manager=manager)

    def create_component(self, id_component, top_ctx, dict_op=None):
        '''
            Create a component
        '''
        id = self.get_topo_id(top_ctx)
        topoNode = TopoNode(_id=id_component, entity=id, task=dict_op)
        return topoNode

    def create_connections(self, source, target):
        '''
            Create a connection between components
            :param source: the source
            :param target: the target

            :return: a TopoEdge
            :rtype: TopoEdge
        '''
        topoEdge = TopoEdge(sources=source, targets=target)
        return topoEdge

    def get_topo_id(self, top_ctx):
        '''
            Get the context ID
        '''
        # Initialize the context
        ctx = Context()
        return ctx.get_entity_id(top_ctx)

    def matcher(self):
        top = self.create_topology()
        top.add_elts('topoNode')
        top.add_elts('topoEdge')

    def at_least(self, dict_data):
        return new_conf(condition.at_least, **dict_data)

    # condition(condition=None, statement=None, _else=None, **kwargs):

    # def at_least(state=Check.OK, min_weight=1)

    def cluster(self, condition, statement, _else):
        '''
        condition = at_least
        statement = then
        _else = else

        value = state dans at_least du V2
        min_weight = min (v1 vers v2)
        '''
        return new_conf(cond, condition=condition, statement=statement, _else=_else)

    def build(self):
        '''
        '''
        f = formatter()
        # kind=1 means to get formatted components
        components = f.get_event_type(kind=1)
        # Operator components
        opcomps = f.match_operator(components)
        # List of Nodes (TopoNodes)
        node_list = []
        # List of connections (topoEdge)
        conn_list = []
        # Create components
        for c in components.get(f.EVENT_TYPE[1]):
            node_list.append(self.create_component(c.keys()[0], c.values[0]))
        for c in components.get(f.EVENT_TYPE[2]):
            node_list.append(self.create_component(c.keys()[0], c.values[0]))
        # OPERATOR_ID[0] --> Cluster
        for cmps in opcomps.get(f.OPERATOR_ID[0]):
            for c in cmps:
                value = c.values()[0].get('options')
                least_value = value.get('least')
                cond_value = value.get('state')
                stat_value = value.get('then')
                else_value = value.get('else')

                dict_cluster = {}
                dict_cluster['state'] = int(cond_value)

                #dict_least = {}
                #dict_least['state'] = int(least_value)
                # Create at_least (Voir avec Jonathan)
                least_conf = new_conf(condition.at_least, min_weight=int(least_value))

                # Create Condition
                conf = new_conf(condition.at_least, **dict_cluster)
                condition = self.create_component(c.keys()[0], c.values[0], conf)

                # Create statement/action
                if stat_value != '-1':
                    statement = new_conf(action.change_state, state=int(stat_value))
                else:
                    statement = new_conf(action.worst_state)

                # Create the else
                if else_value != '-1':
                    _else = new_conf(action.change_state, state=int(else_value))
                else:
                    _else = new_conf(action.worst_state)
                node_list.append(self.create_component(self.cluster(condition, statement, _else)))
        for cmps in opcomps.get(f.OPERATOR_ID[1]):
            for c in cmps:
                node_list.append(new_conf(action.worst_state))
        # OPERATOR_ID[2] --> And
        for cmps in opcomps.get(f.OPERATOR_ID[2]):
            for c in cmps:
                value = c.values()[0].get('form').get('items')
                cond_value = value[0].get('value')
                stat_value = value[1].get('value')
                else_value = value[2].get('value')

                dict_and = {}
                dict_and['state'] = int(cond_value)

                # Create the condition
                conf = new_conf(condition.at_least, **dict_and)
                condition = self.create_component(c.keys()[0], c.values[0], conf)

                # Create the statement/action
                if stat_value != '-1':
                    statement = new_conf(action.change_state, state=int(stat_value))
                else:
                    statement = new_conf(action.worst_state)

                # Create the else
                if else_value != '-1':
                    _else = new_conf(action.change_state, state=int(else_value))
                else:
                    _else = new_conf(action.worst_state)
                node_list.append(self.create_component(self.cluster(condition, statement, _else)))
        # OPERATOR_ID[3] --> Or
        for cmps in opcomps.get(f.OPERATOR_ID[3]):
            for c in cmps:
                value = c.values()[0].get('form').get('items')
                cond_value = value[0].get('value')
                stat_value = value[1].get('value')
                else_value = value[2].get('value')
                condition = ""
                statement = ""
                _else = ""

                dict_or = {}
                dict_or['state'] = int(cond_value)

                # Create the condition
                conf = new_conf(condition.at_least, **dict_or)
                condition = self.create_component(c.keys()[0], c.values[0], conf)

                # Create the statement/action
                if stat_value != '-1':
                    statement = new_conf(action.change_state, state=int(stat_value))
                else:
                    statement = new_conf(action.worst_state)

                # Create the _else
                if else_value != '-1':
                    _else = new_conf(action.change_state, state=int(else_value))
                else:
                    _else = new_conf(action.worst_state)
                #statement = new_conf(action.change_state, state='value')
                #_else = new_conf(action.change_state, state='value')
                #worst = new_conf(action.worst_state) # if value = -1 Worst state
                node_list.append(self.create_component(self.cluster(condition, statement, _else)))

        # OPERATOR_ID[4] --> Best State
        for cmps in opcomps.get(f.OPERATOR_ID[4]):
            for c in cmps:
                node_list.append(new_conf(action.best_state))

        # Create connections between components
        for tween in f.get_comp_graph():
            for val in tween:
                conn_list.append(self.create_connections(val[0], val[1]))

        # Create the Topology
        root_id = f.get_root_id()
        self.create_topology(root_id, conn_list, node_list)
