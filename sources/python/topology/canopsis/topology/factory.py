#!/usr/bin/env python
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
from canopsis.topology.rule import action
from canopsis.topology.rule.condition import at_least


class Factory(object):
    """docstring for Factory"""
    def __init__(self):
        pass

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

    def create_component(self, id_component, top_ctx=None, dict_op=None, data=None):
        '''
            Create a component

            :param id_component: the component ID.
            :param top_ctx: the topology context.
            :param dict_op: the operator dictionnary.
            :param data: the label of each component.

            :return: a TopoNode.
            :rtype: TopoNode.
        '''
        #data = {'label':'value'}
        if top_ctx is not None:
            id = self.get_topo_id(top_ctx)
        else:
            id = None
        topoNode = TopoNode(_id=id_component, entity=id, operator=dict_op, data=data)
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
            :param top_ctx: the context.
            :return: Context ID.
            :rtype: Context.
        '''
        # Initialize the context
        ctx = Context()
        return ctx.get_entity_id(top_ctx)

    def at_least(self, dict_data):
        '''
        Create the condition statement.
        '''
        return new_conf(condition.at_least, **dict_data)

    def cluster(self, condition, statement, _else):
        '''
        condition = at_least
        statement = then
        _else = else

        value = state dans at_least du V2
        min_weight = min (v1 vers v2)
        '''
        return new_conf(cond, condition=condition, statement=statement, _else=_else)

    def build(self, topo_id):
        '''
        Create the topology V2 using the topology V1 data.
        '''
        f = formatter.Formatter(topo_id)
        # kind=1 means to get formatted components
        components = f.get_event_type(kind=1)
        # Operator components
        opcomps = f.match_operator(components)
        # List of Nodes (TopoNodes)
        node_list = []
        # List of connections between components (topoEdge)
        conn_list = []
        # Label dictionnary
        data = {}
        # Create components (type=check)
        comp_check  = components.get(f.EVENT_TYPE[1])
        if comp_check is not None:
            for c in comp_check:
                entity = {'component': unicode(c.values()[0].get('component')),'resource': unicode(c.values()[0].get('resource')),'connector': unicode(c.values()[0].get('connector')),'connector_name':unicode(c.values()[0].get('connector_name')),'type':unicode(c.values()[0].get('type'))}
                entity['id'] = c.values()[0].get(c.values()[0].get('source_type'))
                if c.values()[0].get('label') is None:
                    data['label'] = c.values()[0].get('component')
                else:
                    data['label'] = c.values()[0].get('label')
                node_list.append(self.create_component(c.keys()[0], entity, data=data))
                data = {}
        # Create components (type=selector)
        comp_selct = components.get(f.EVENT_TYPE[2])
        if comp_selct is not None:
            for c in comp_selct:
                entity = {'component': unicode(c.values()[0].get('component')),'resource': unicode(c.values()[0].get('resource')),'connector': unicode(c.values()[0].get('connector')),'connector_name':unicode(c.values()[0].get('connector_name')),'type':unicode(c.values()[0].get('type'))}
                entity['id'] = c.values()[0].get(c.values()[0].get('source_type'))
                if c.values()[0].get('label') is None:
                    data['label'] = c.values()[0].get('component')
                else:
                    data['label'] = c.values()[0].get('label')
                node_list.append(self.create_component(c.keys()[0], entity, data=data))
                data = {}
        # Create components (type=topology)
        comp_topo = components.get(f.EVENT_TYPE[3])
        if comp_topo is not None:
            for c in comp_topo:
                entity = {'component': unicode(c.values()[0].get('component')),'resource': unicode(c.values()[0].get('resource')),'connector': unicode(c.values()[0].get('connector')),'connector_name':unicode(c.values()[0].get('connector_name')),'type':unicode(c.values()[0].get('type'))}
                entity['id'] = c.values()[0].get(c.values()[0].get('source_type'))
                if c.values()[0].get('label') is None:
                    data['label'] = c.values()[0].get('component')
                else:
                    data['label'] = c.values()[0].get('label')
                node_list.append(self.create_component(c.keys()[0], entity, data=data))
                data = {}
        # Create components (type=operator)
        # OPERATOR_ID[0] --> Cluster
        comp_opera = opcomps.get(f.OPERATOR_ID[0])
        if comp_opera is not None:
            for cmps in comp_opera:
                tmpdict = cmps.values()[0]
                value = tmpdict.get('options')
                clt_list = tmpdict.get('form').get('items')
                if value is None:
                    least_value = clt_list[0].get('value')
                    cond_value = clt_list[1].get('value')
                    stat_value = clt_list[2].get('value')
                    else_value = clt_list[3].get('value')
                else:
                    least_value = value.get('least')
                    cond_value = value.get('state')
                    stat_value = value.get('then')
                    else_value = value.get('else')

                least_conf = new_conf(at_least, min_weight=int(least_value), state=int(cond_value))
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
                if tmpdict.get('label') is None:
                    data['label'] = tmpdict.get('component')
                else:
                    data['label'] = tmpdict.get('label')
                node_list.append(self.create_component(cmps.keys()[0], dict_op=self.cluster(least_conf, statement, _else), data=data))
                data = {}
        # OPERATOR_ID[1] --> Worst Sate
        comp_worst = opcomps.get(f.OPERATOR_ID[1])
        if comp_worst is not None:
            for cmps in comp_worst:
                tmpdict = cmps.values()[0]
                if tmpdict.get('label') is None:
                    data['label'] = tmpdict.get('component')
                else:
                    data['label'] = tmpdict.get('label')
                node_list.append(self.create_component(cmps.keys()[0], dict_op=new_conf(action.worst_state), data=data))
                data = {}
        # OPERATOR_ID[2] --> And
        comp_and = opcomps.get(f.OPERATOR_ID[2])
        if comp_and is not None:
            for cmps in comp_and:
                mydict = cmps.values()[0]
                value = mydict.get('form').get('items')
                cond_value = value[0].get('value')
                stat_value = value[1].get('value')
                else_value = value[2].get('value')

                entity = {'component': unicode(mydict.get('component')),'resource': unicode(mydict.get('resource')),'connector': unicode(mydict.get('connector')),'connector_name':unicode(mydict.get('connector_name')),'type':unicode(mydict.get('type'))}
                entity['id'] = mydict.get(mydict.get('source_type'))

                dict_and = {}
                dict_and['state'] = int(cond_value)
                # Create the condition
                condition = new_conf(at_least, **dict_and)
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
                if tmpdict.get('label') is None:
                    data['label'] = tmpdict.get('component')
                else:
                    data['label'] = tmpdict.get('label')
                node_list.append(self.create_component(cmps.keys()[0], dict_op=self.cluster(condition, statement, _else), data=data))
                data = {}

        # OPERATOR_ID[3] --> Or
        comp_or = opcomps.get(f.OPERATOR_ID[3])
        if comp_or is not None:
            for cmps in comp_or:
                mydict = cmps.values()[0]
                value = mydict.get('form').get('items')
                cond_value = value[0].get('value')
                stat_value = value[1].get('value')
                else_value = value[2].get('value')
                condition = ""
                statement = ""
                _else = ""

                dict_or = {}
                dict_or['state'] = int(cond_value)

                entity = {'component': unicode(mydict.get('component')),'resource': unicode(mydict.get('resource')),'connector': unicode(mydict.get('connector')),'connector_name':unicode(mydict.get('connector_name')),'type':unicode(mydict.get('type'))}
                entity['id'] = mydict.get(mydict.get('source_type'))

                # Create the condition
                condition = new_conf(at_least, **dict_or)

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
                if tmpdict.get('label') is None:
                    data['label'] = tmpdict.get('component')
                else:
                    data['label'] = tmpdict.get('label')
                node_list.append(self.create_component(cmps.keys()[0], dict_op=self.cluster(condition, statement, _else), data=data))
                data = {}
        # OPERATOR_ID[4] --> Best State
        comp_best = opcomps.get(f.OPERATOR_ID[4])
        if comp_best is not None:
            for cmps in comp_best:
                mydict = cmps.values()[0]
                if tmpdict.get('label') is None:
                    data['label'] = tmpdict.get('component')
                else:
                    data['label'] = tmpdict.get('label')
                node_list.append(self.create_component(cmps.keys()[0], dict_op=new_conf(action.best_state), data=data))
                data = {}

        # Create connections between components
        for tween in f.get_comp_graph():
            conn_list.append(self.create_connections(tween[0], tween[1]))
        # Create the Topology
        root_id = f.get_root_id()
        self.create_topology(root_id, conn_list, node_list)
        print "Toplogy ", root_id, " is created successfully."

    def delete_topology(self, comp_ID):
        # Initialize the Toplogy Manager
        manager = TopologyManager()
        # Create the topology name
        top = manager.get_graphs(ids=comp_ID)
        if top is not None:
            top.delete(manager)
            print "component: ", comp_ID, " is deleted ..."
        else:
            print "component: ", comp_ID, " does no exist in the Database ..."


if __name__ == '__main__':
    fact = Factory()
    fact.build('canopsis_arbre') # create the topology 'canopsis_arbre'
    #fact.delete_topology('component-1852') # delete the topology 'canopsis_arbre'