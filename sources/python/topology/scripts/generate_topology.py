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
from canopsis.topology.manager import TopologyManager
from canopsis.topology.elements import TopoEdge, TopoNode, Topology
from canopsis.topology.rule.action import change_state, worst_state, best_state
from canopsis.topology.rule.condition import at_least, _all, nok
from canopsis.task.core import new_conf, TASK
from canopsis.task.condition import condition

from argparse import ArgumentParser

manager = TopologyManager()


def generate_test(topo, name='test'):
    """Generate a test topology.
    """

    # create compositions of all actions and conditions
    actions = [
        new_conf(action) for action in [change_state, worst_state, best_state]
    ]
    # generate all conditions with all states
    conditions = []
    #for state in range(4):
    for cond in (at_least, _all, nok):
        conditions.append(new_conf(cond))
    # list of all possible operators
    operators = actions[:]
    for cond in conditions:
        for statement in actions:
            for _else in actions:
                if statement != _else:
                    operator = new_conf(
                        condition,
                        condition=cond,
                        statement=statement,
                        _else=_else
                    )
                    operators.append(operator)
    # for all operatorscreate all node operators and bind them to the topology
    for operator in operators:
        # create an operator node
        op = TopoNode(info={TASK: operator})
        topo.add_elts(op)
        # bind it to the topo
        edge = TopoEdge(sources=op, targets=topo)
        topo.add_elts(edge)
        # and add three mode nodes bound to the operator node
        nodes = [TopoNode() for i in range(3)]
        for node in nodes:
            topo.add_elts(node)
            edge = TopoEdge(sources=node, targets=op)
            topo.add_elts(edge)


def generate_topology(name, _type):
    """Generate a topology related to a name and type.
    """

    # delete old topology
    topo = manager.get_graphs(ids=name)
    if topo is not None:
        topo.delete(manager=manager)
    else:
        topo = Topology(id=name)

    if _type == 'context':
        generate_context_topology(topo=topo, name=name)
    elif _type == 'rules':
        generate_rules_topology(topo=topo, name=name)
    elif _type == 'test':
        generate_test(topo=topo, name=name)
    # save topology
    topo.save(manager=manager)


def generate_context_topology(topo, name='context'):
    """Generate a context topology where nodes are components and resources,
    and edges are dependencies from components to resources, or from resources
    to the topology.

    :param str name: topology name.
    """

    # initialize context and topology
    context = Context()

    def addElt(elt):
        """
        Add input elt in topology.

        :param GraphElement elt: elt to add to topology.
        """

        topo.add_elts(elt.id)
        elt.save(manager)

    components = context.find({'$in': ['component', 'topo', 'selector']})
    for component in components:
        component_id = context.get_entity_id(component)
        component_node = TopoNode(entity=component_id)
        addElt(component_node)

        ctx, name = context.get_entity_context_and_name(component)
        ctx['component'] = name
        resources = context.find('resource', context=ctx)
        if resources:  # link component to all its resources with the same edge
            edge = TopoEdge(sources=component_node.id, targets=[])
            addElt(edge)  # add edge in topology
            for resource in resources:
                resource_id = context.get_entity_id(resource)
                resource_node = TopoNode(entity=resource_id)
                addElt(resource_node)  # save resource node
                # add resource from component
                edge.targets.append(resource_node.id)
                res2topo = TopoEdge(
                    sources=resource_node.id, targets=topo.id
                )
                addElt(res2topo)
            if not edge.targets:  # bind topology from component if not sources
                edge.targets.append(topo.id)
            addElt(edge)  # save edge in all cases
        else:  # if no resources, link the component to the topology
            edge = TopoEdge(sources=component_node.id, targets=topo.id)
            addElt(edge)  # add edge in topology


def generate_rules_topology(topo, name):
    """Generate a topology with rules.
    """

    elts = []
    # create a simple rule
    tn1 = TopoNode()
    # and bind it to the topo
    tn1topo = TopoEdge(sources=tn1, targets=topo)
    elts += [tn1, tn1topo]
    tn0 = TopoNode()
    tn0tn1 = TopoEdge(sources=tn0, targets=tn1)
    elts += [tn0, tn0tn1]
    # add rules in the topo
    topo.add_elts(elts)


if __name__ == '__main__':

    parser = ArgumentParser(description='Generate a topology')
    parser.add_argument(
        dest='name',
        help='topology name to generate (default: context)',
        default='context'
    )
    parser.add_argument(
        dest='type',
        help='topology type among context, rules, test (default: context)',
        default='context'
    )
    args = parser.parse_args()

    topology_name = args.name
    topology_type = args.type
    generate_topology(name=topology_name, _type=topology_type)
