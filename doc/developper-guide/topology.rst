========
Topology
========

A topology aims to resolve root cause analysis thanks to a graph where nodes and edges are respectively system infrastructure components or status changing rules and technical and functional dependencies among them.

Topology processing is as follow:

1. Configure a topology in binding nodes to system components and in adding status dependencies and status changing rules.
2. Start topology engine(s) which get topology configuration and monitore the IT system in order to identify which status are changing.
3. When a status changing rule is triggered, a topology engine updates the topology configuration and execute the rule.

Configuration
=============

A topology is a directed graph, where edges represent status changing dependencies and where nodes represent logical and technical status system component information.

There are three types of nodes:

- technical node: bound to a monitored system component.
- logical node: contains status changing rules.
- root: First topology component which is also a logical node. Used to export the topology in other topologies.

All nodes are entities (entities_), where types equal `̀topology-root̀`` for a root node, and ̀̀topology-nodè` for other ones.

--------------
Technical node
--------------

A technical node is bound to an entity, such as hosts, services, DB, topology-node, etc. (entities_)

This binding is ensured by an "entity_id" field.

------------
Logical node
------------

A logical node contains status changing rules which are event_rules_ dedicated topology.

Those rules support are coupled to a maximal action execution count. This counter is called ``max_actions``.

For example, let 3 rules configured on a logical node, and three rules always triggered when an event ``E`` occured. If you want to limit the number of triggered rules to the first both ones, you have to set the ``max_actions`` to 2. Then when the second action will be done, the number of rules to test will be stopped.

A condition is a cfilter_ which permits to apply a filter on any event properties.

An action is a job_, such as task_mail, send_event, etc.

In the case of a topology, 

Here are simplified rules dedicated to the status changing rules in topologies.

Worst state
-----------

----
Rule
----

Such as an operational part, a node should contain aggregation operation information in order to apply rules when neighboor components change of status.

Operations are Operator, Worst state, Best state, And and Or.

Operator
--------

Every operation is a specialization of this one. Therefore, defined parameters will takes different values depending on the type of the operation.

Here are list of parameters:

- actions: action to do when the operation is triggered.
- check state: status able to trigger the operation.
- at least: nomber of nodes which can impact the node status change.
- in state: status when operation is not triggered.

Wors state
----------

Best state
----------

And
---

Or
--


Asynchronous vs Synchronous
===========================

When a node status change, impacted nodes can be among nodes pointed by edges, among several topologies.

Therefore, even the complexity of topologies, it is possible to reduce their path complexity in focusing recursively on nodes which notice the changement of status, or on nodes pointed by previous nodes.

Engine
======

Actions
=======
