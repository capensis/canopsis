=========================================
Topology: library for managing topologies
=========================================

.. module:: canopsis.topology

A topology aims to resolve root cause analysis thanks to a graph where nodes and edges are respectively system infrastructure components or status changing rules with technical and functional dependencies among them.

Topology processing is as follow:

1. Configure a topology in binding nodes to system components and in adding status dependencies and status changing rules.
2. Start topology engine(s) which get topology configuration and monitore the IT system in order to identify which status are changing.
3. topology execution can be done in two ways, synchronously and asynchronously to event consumption:
    - Synchronously: When a status changing rule is triggered, a topology engine updates the topology configuration and execute the rule.
    - Periodically: a topology engine gets a topology configuration and resolves all node status.

Specification
=============

A topology is a directed graph, where edges represent status changing dependencies and where nodes represent logical and technical status system component information.

----
Node
----

Description
-----------

There are three types of nodes:

- technical node: bound to a technical system component such as a (virtual) machine, service, DB, etc.
- logical node: bound to a logical system component such as cluster of machines, etc.
- root: First topology component which is also a logical node. Used to export the topology in other topologies.

Rule
####

A node can depend on child nodes. This dependency specify status changing rule propagation from child nodes to parent nodes, and is coordinated by topology rules.

Here are examples of rules:

- Operator.
- Worst_state.
- Best_state.
- And.
- Or.

Operator
>>>>>>>>

Every operation is a specialization of this one. Therefore, defined parameters will take different values depending on operation type.

Here are list of parameters:

- check state: state able to trigger the operation.
- at least: maximal number of status nodes which can impact the node state change.
- at most: minimal number of status nodes which can impact the node state change.
- default state: state when operation is not triggered.
- in state: state when operation is triggered.

If a parameter is not given, its value is setted by the user.

TODO: in a future version, ``in state`` should be a generic action, and other parameters except ``default state`` will be part of a generic condition.

Worst state
~~~~~~~~~~~

- check state: different state than actual.
- at least: 1 (very effective in synchronous mode if new state is worst than the old one).
- at most: 1
- default state: {bound entity state, ok} if no child node exist.
- in state: max(child node states).

Best state
~~~~~~~~~~

- check state: different state than actual.
- at least: 1 (very effective in synchronous mode if new state is better than the old one).
- at most: 1.
- default state: {bound entity state, ok} if no child node exist.
- in state: minimal(child node states).

All
~~~

- at least: -1 (all nodes)
- at most: -1 (all nodes)
- default state: {bound entity state, ok}.

Any
~~~

- at least: 1
- at most: 1
- default state: {bound entity state, ok}.

Structure
---------

It is important to separate topology and node structures in order to modify a node without busying the topology and other nodes (in a future situation, it could be interesting to propose a collaborative edition/visualization of a topology).

Topology
########

- ``id``: unique id among topologies
- ``nodes``: node ids
- ``root``: root node id


Topology-node
#############

- ``id``: unique id among one topology.
- ``rules``: rules.
- ``next``: nodes which are concerned by a change of state.
- ``state``: topology node state.

Rule
####

- ``name``: rule name (``'worst state'``, etc.).
- ``parameters``: couple of (name, value) where name is in {``at_least``, ``default state``, ...}.

---------------------------
Asynchronous vs Synchronous
---------------------------

When a node status change, impacted nodes can be among nodes pointed by edges from potentially same nodes, among several topologies.

That mean time and space complexity depends on number of nodes and edges.

Let E the set of entities, N the set of nodes and M the set of edges.

In a connected graph of n nodes, there is at least n-1 edges. That means \|N\| in [\|M\|-1; inf[

And when a node has to be calculated, it consists to get all child node status which could come from bound entities or sub-child nodes.

Asynchronous mode
-----------------

This mode consist to get periodically all topologies and to resolve node status on the entire topology.

Each step, such mode consists to get all nodes and all bound entites in order to set the rights states on each nodes.

Even if there is a root node in a topology, the graph is not a tree, and so, it is important to ensure an order of status change from parent to children. In such situation, it is required to resolve such mode in a width path instead of a depth path.

In time complexity, this mode depends on number of edges where for each edge, we require to get a node information and the bound entity status if exist.

Therefore the time complexity is O(N) = 2|M|.

In space complexity, this mode is hard because dependency between topologies implies to load all topology each time we need to solve topologies.

Advantages
##########

The topology consistency.

Weaknesses
##########

- Space complexity: require to load all topologies at a time. Impossible to use distributed calculus to solve topologies.
- During topology solving, it is possible most node rule checking are useless.
- The period of resolution is a time gap between real-time system and close-real-time system. This delta has an impact in monitoring and root cause analysis time (and so in precision). Several status changing propagations will be fired during the period of topology resolution.

Synchronous mode
----------------

This mode aims to focus on a topology nodes which corresponds to an event processed by an engine.

A topology engine in this mode gets only events of type check. When such event arrives, it tries to get all nodes related to such event. If nodes are found, the engine runs all node rules. Some rule require to get child nodes, in this case, the engine get all child nodes. For every nodes which have been modified after the rule execution, it sends an event of type ``topology-propagation`` in order to ensure than only one engine will do this operation in a distributed calculus.

Advantages
##########

- Very close-real-time solution.
- Limit number of rule execution to required ones.
- Allows distributed calculus.

Weaknesses
##########

TODO:
    ... to determinate

-------------------
Root cause analysis
-------------------

The root cause analysis is of two types:

- dynamically (close-)real-time.
- static.

Dynamically close-real-time
---------------------------

In such situation, it is useful to be notified as soon as possible about propagation of a change of status.

A dedicated widget in the UI could permit to see how a topology change in close-real-time.

TODO:
    - engines and managers

Static
------

This mode should permit to solve root cause analysis from historical data, and to do analysis about the system consistency/state.

TODO:
    - functionalities

CRUD on topology
================

Here are CRUD operations useful to monitor topologies.

--------
Topology
--------

.. module:: canopsis.topology.ws

.. function:: get(ids=None, add_nodes=False)

    .. data:: REST route = 'GET:/topology[/ids[/add_nodes[/limit[/skip]]]]'

    :param ids: None, one id or a list of ids
    :type ids: NoneType, list(str) or str

    :param add_nodes: if True (default False), replace nodes references in the result by topology node data.
    :type add_nodes: bool

    :return: depending on ids

        - one id: one topology or None if topology does not exist.
        - list of ids: list of topologies.
        - None: all topologies.

.. function:: find(regex, add_nodes=False)

    .. data:: REST route = 'GET:/topology-regex[/regex[/add_nodes]]'

    :param regex: regex expression which could match with topology ids.
    :type regex: str

    :param add_nodes: if True (default False), replace nodes references in the result by topology node data.
    :type add_nodes: bool

    :return: find all topologies where ids correspond to input regex

.. function:: put(topology)

    .. data:: REST route = 'PUT:/topology[/topology]'

    put input topology in DB. If topology already exist, update existing information with input information, or add them if they does not exist.

    :param topology: topology
    :type topology: dict

.. function:: remove(ids)

    .. data:: REST route = 'DELETE:/topology[/ids]'

    remove one or more topologies depending on ids:

    - None: remove all topologies
    - id: remove one topology where id is input ids
    - list of id: remove topologies where id is in input ids

-------------
Topology node
-------------

.. module:: canopsis.topology.node.ws

.. function:: get(ids=None)

    .. data:: route = 'GET:/rest/topology-node[/ids]'

    get topology-nodes depending on a ``topology`` and ``ids``:

    The related topology is the one where the id equals ``topology_id``

    - None: get all topology-nodes which are in the related topology.
    - id: get one topology-node which is in the related topology.
    - lit of id: get one or more topology-nodes which are in the related topology and where id are in input ids.

.. function:: find_nodes_by_entity(entity_id)

    .. data:: route = 'GET:/rest/topology-node-entity[/id]'

    Find all topology-nodes which are bound to entity_id.

.. function:: find_nodes_by_next(next_id)

    .. data:: route = 'GET:/rest/topology-node-next[/id]'

    Find all topology-nodes which are sources of a topology-node id given in parameter.

.. function:: remove_nodes(ids=None)

    .. data:: route = 'DELETE:/rest/topology-node[/ids]'

    Remove all topology-nodes which are in a topology and identified among one id or a list of ids.


User Interface
==============

TODO: ...

