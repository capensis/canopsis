==========================================
Topology: solution for root cause analysis
==========================================

This document specifies root cause analysis requirements.

Description
===========

A root cause analysis is a mean to identify origins of failures in an infrastructure.

In order to solve this problem, Canopsis provide topologies.

A topology is an entity composed of nodes partially connected.

Nodes use:

- state information (OK, WARNING, etc.) which can be synchronized to entity states or not (entity_binding__).
- change of state propagation rules (state_propagation_rule__).

In considering two directed connected nodes, connection source and target nodes are called respectivelly source node and next node.

For topology composition reasons, a node may be used by several topologies.

Entity binding
##############

A topology node can be bound to an entity in order to update its state if entity state change, whatever node state propagation rule (state_propagation_rule__). This behavior is not reciprocal.

State propagation Rule
######################

A propagation rule permits to define state propagation over node edges. Such rules use at most one condition and at least one action.

Condition
---------

Conditions are boolean functions used among source nodes (opposite of edge direction).

- default: base condition for all others.

   + check_state: state to check among source nodes
   + at_least: number of source nodes where state equals check_state

- All: check if all source node states respect input check_state.

   + state: state to check in all source nodes.

- Any: check if one source node state equals to input chec_state.

   + state: state to check in one source node.

Action
------

Actions are done when condition does not exist or is checked.

- change_state: default action.

   + state: new state to apply on node.
- best_state: update the state related to the best one among source nodes.
- worst_state: update the state related to the worst one among source nodes.

Data Model
==========

Topology
########

- type: 'topology'.
- id: topology name. Unique among all topologies.
- nodes: list of topology node ids.
- root: root topology node id.

Topology-node
#############

- type: 'topology-node'.
- id: topology-node id.
- entity_id: bound entity id. Equals to id if no bound entity.
- rule: topology node rule.
- nexts: list of next topology node ids.
- state: topology node state.

Rule
####

- condition:

   + none
   + string value among:

      * canopsis.topology.rule.condition.new_state
      * canopsis.topology.rule.condition.condition
      * canopsis.topology.rule.condition.all
      * canopsis.topology.rule.condition.any
   + dictionary:

      * task_path: string value such as previously
      * params: dictionary of condition parameter values

         - state: state parameter for change_state, condition, all and any functions.
         - at_least: number of source nodes to check for condition function.

Engine
======

The engine listen to events of type ``check``. When such event is received, the engine get all topology nodes bound to the event entity.

For all bound nodes, the engine execute their rules. And for all nodes which have change of state after applying their rules, the engine iterate on all next nodes.

For all next node, the engine create a new event of type check and send it to Canopsis::

   - 'type': 'check'
   - 'source_type': 'topology-node'
   - 'id': next node id
   - 'source': source node id

In order to be processed recursively by a topology engine and other canopsis engines such as a check event.

API REST
========

It is possible to interact with the topology model from an API REST.

In all routes::
   - REST operation prefix the route.
   - (optional) parameters are prefixed by the character ('::') ':'.

For example, the route::

   "GET:/topology/:mandatory/::optional"

Specifies the REST operation ``GET``, the required parameter ``mandatory`` and the optional parameter ``optional``.

Topology
########

Get topology
------------

Route
>>>>>

GET:/topology/::ids/::add_nodes

Parameters
>>>>>>>>>>

- ids (str or list of str): one topology id or a list of topology ids. If not specified, get all existing topologies.
- add_nodes (bool): add topology node values instead of keeping only topology node ids.

Find topology
-------------

Route
>>>>>

GET:/topology/::regex/::add_nodes

Parameters
>>>>>>>>>>

- regex (str): if given, find all topologies where the name matches with the regex. Otherwise, get all topologies.
- add_nodes (bool): add topology node values instead of keeping only topology node ids.

Put a topology
--------------

Route
>>>>>

PUT:/topology

Parameters
>>>>>>>>>>

- topology (dict): topology to put.

Delete a topology
-----------------

Route
>>>>>

DELETE:/topology/::ids

Parameters
>>>>>>>>>>

- ids (str or list of str): topology ids to delete. If not given, delete all topologies.

Topology node
#############

Get topology node(s)
--------------------

GET:/topology_nodes/::ids

Parameters
>>>>>>>>>>

- ids (str or list of str): topology node ids to get. If not, get all topology nodes.

Find topology nodes
-------------------

Route
>>>>>

GET:/topology_nodes_find/::entity_id

Parameters
>>>>>>>>>>

- entity_id (str): entity id

Put a topology node
-------------------

PUT:/topology_node

Parameters
>>>>>>>>>>

- topology_node (dict): topology node to put.

Delete a topology node
----------------------

Widget
======

The widget propose to visualize and edit topology nodes.
