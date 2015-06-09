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

- state information (OK, WARNING, etc.) which can get value from entity states or not (entity_binding__).
- change of state propagation rules (state_propagation_rule__).
- weight which permits to calculate state propagation from a set of nodes.

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

   + state: state to check among source nodes.
   + min_weight (default 1): minimal weight of source nodes where state equals check_state.
   + rrule (default None): before and after rrule conditions (see RFC2445 http://tools.ietf.org/html/rfc2445).

- all_nodes: check if all source node states respect input check_state.

   + state: state to check in all source nodes.

Additional conditions permit to check multiple conditions with logical ``and`` and ``or`` conditions and ``time`` condition.

- ``and``: check if all input conditions are true.
- ``or``: check if one input condition is True
- ``time``: check if current time is in the interval described by a rrule (RFC2445) and a duration.

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
- weight: topology node weight. Default 1.

Rule
####

- condition:

   + none
   + string value among:

      * topology.new_state
      * topology.condition
      * topology.all_nodes
      * any
      * all
      * during

   + dictionary:

      * task_path: string value such as previously
      * params: dictionary of condition parameter values

         - state: state parameter for change_state, condition, all and any functions.
         - at_least: number of source nodes to check for condition function.
         - required rrule and optional timestamp and duration (``during`` function)

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

A mock-up is available to this url:

https://cacoo.com/diagrams/BzcENww2MapkAhqx

Visualisation
#############

This widget is composed of two parts.

- A graph view which displays nodes and edges.
- An array view which display information related to graph selected nodes.

Both views can be hidden/shown at any time.

Graph
#####

Nodes and edges are displayed related to node distinguishable properties.

For example, node size are related to their weight (more a node is a weighted source, more it should be visible). Node colors are related to their states.

Edge size is related to node weight compared to other node weights. Line can be doted/normal related to respectives node condition ``default``/``all_nodes``.

OK is green, KO is red, UNKNOWN is white, etc.
Related to such colors, a downed node is transparent (ok and down). And if a related event is in an ACK status, the contour circle color is blue.

A filter permits to temporarely display nodes which match an input filter. This filter applies a regex on any (entity) name/rule names or weight values. In the mock-up, the filter is at the top right. It contains the value ``service`` which allows to display all nodes which has properties values matching with ``service``.

Number of neighboor nodes can be displayed at the top/bottom for ``next``/``source`` nodes. In addition to those information, displaying node buttons are available if the mouse is over the node or if the node is selected.

- ``-``: hide nodes. Available only if nodes are displayed.
- ``+``: show neighbor nodes. Available only if nodes are hidden.
- ``++``: show all nodes and not only neighbor. Always available.

Description in the mock-up:
###########################

https://cacoo.com/diagrams/BzcENww2MapkAhqx

- serviceA is OK and selected (double contour lines). At the bottom, 4 source nodes over 18, and at the top, 2 next nodes over 10. The edge from serviceA to serviceL describes an ``default`` operation (doted line). The edge from serviceA to serviceR specifies an ``all_nodes`` condition.
- serviceL is KO (red color) but in ACK status (blue contour). Source nodes can be totally expanded (``++`` button) or hidden (``-``).
- serviceR is OK and downed (green transparent). Its total weight is the lower because its size is the smaller. A red flashlight shows that this node use the ``worst_state`` (green for ``best_state``).

Array
#####

In order to improve visibility, an array is dynamically linked to selected nodes. In this example, selected nodes are ``serviceA`` and ``serviceL``. Related properties are displayed in the different columns.

Interaction
-----------

Graph
#####

Zoom in/out is possible with the wheel or in moving two with fingers.

An auto-layout is possible and deactivable if necessary (check box under the filter inside the mock-up).

This auto-layout permits to choose node disposition and to reorganize nodes when focusing/selecting on them.

All those buttons are displayed only if neighbor nodes exist.

Double clicking on a node permits to go to the event box and to focus on selected events with optionally an historic of change of state propagation.

Adding a node/edge is possible from a dedicated entity/rules set where such elements are dragable to the graph view (and rules can be added).

Array
#####

Array items can be sorted by columns.
Clicking on the entity name force the view to go to the event box view with cliqued item filter.
