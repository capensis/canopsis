.. _dev-backend-mgr-topology:

Topology: library for managing topologies
=========================================

.. module:: canopsis.topology

Description
===========

Functional
----------

A topology is a graph dedicated to enriches status model of entities with state
dependency between entities.

Topological tasks consist to update status vertice information and to propagate
the change of state in sending check events.

vertices could be finally connected to the topology in order to propagate all
 change of state to the topology itelf.

An example of application is root cause analysis where a topology may react
when an entity change of state and can propagate over topology nodes the change
of state in some propagation conditions such as operations like ``worst state``

Topology tasks are commonly rules of (condition, actions). A condition takes in
parameter the execution context of an event, the engine and the vertice which
embeds the rule.

Technical
---------

Three types of vertices exist in topology::

- cluster: operation between vertice states.
- node: vertice bound to an entity, like components, resources, etc.

A topology vertice contains::

- state: vertice state which change at runtime depending on bound entity state and event propagation.

A topology edge contains::

- weight: vertice weight in the graph.

Rule
====

.. module:: canopsis.topology.rule

Condition
---------

.. module:: canopsis.topology.rule.condition

Module in charge of defining main topological rule conditions.

A topological node doesn't require a condition.

In addition to this condition, it is possible to test another condition which
refers to source nodes if they exist.

Such conditions are::
    - ``new_state``: test if state (or event state) is not equal to node state
    - ``at_least``: test if source node states match with an input state.
    - ``_all``: test if all source node states match with an input state.
    - ``nok``: test if source node states are not OK.

The ``new_state`` condition may be used by nodes bound to entities in order to
update such nodes when the entity change of state.

Related rule actions are defined in ``canopsis.topology.rule.action`` module.

new_state
---------

Condition triggered when state is different than node state.

   - dict event: event from where get state if input state is None.
   - Node node: node from where check if node state != input state.
   - int state: state to compare with input node state.

at_least
--------

Generic condition applied on sources of node which check if at least source nodes check a condition.

   - dict event: processed event.
   - dict ctx: rule context which must contain rule node.
   - Node node: node to check.
   - int state: state to check among sources nodes.
   - float min_weight: minimal weight (default 1) to reach in order to
        validate this condition. If None, condition results in checking all
            sources.
   - rrule rrule: rrule to consider in order to check condition in time.
   - f: function to apply on source node state. If None, use equality
        between input state and source node state.

   - return True if condition is checked among source nodes.

_all
----

Check if all source nodes match with input check_state.

   - dict event: processed event.
   - dict ctx: rule context which must contain rule node.
   - Node node: node to check.
   - int min_weight: minimal node weight to check.
   - int state: state to check among sources nodes.
   - rrule rrule: rrule to consider in order to check condition in time.
   - f: function to apply on source node state. If None, use equality
        between input state and source node state.

   - True if condition is checked among source nodes.

nok
---

Condition which check if source nodes are not ok.

   - dict event: processed event.
   - dict ctx: rule context which must contain rule node.
   - Node node: node to check.
   - int min_weight: minimal node weight to check.
   - int state: state to check among sources nodes.
   - rrule rrule: rrule to consider in order to check condition in time.
   - f: function to apply on source node state. If None, use equality
        between input state and source node state.

   - True if condition is checked among source nodes.

Action
######

.. module:: canopsis.topology.rule.action


A topological node has at least one of four actions in charge of changing
of state::

    - ``change_state``: change of state related to an input or event state.
    - ``state_from_sources``: change of state related to source nodes.
    - ``best_state``: change of state related to the best source node state.
    - ``worst_state``: change of state related to the worst source node state

change_state
------------

Change of state on node and propagate the change of state on bound entity if necessary.

   - event: event to process in order to change of state.
   - node: node to change of state.
   - state: new state to apply on input node. If None, get state from
        input event.
   - bool update_entity: update entity state if True (False by default).
        The topology graph may have this flag to True.
   - int criticity: criticity level. Default HARD.


state_from_sources
------------------

Change ctx node state which equals to f result on source nodes.

worst_state
-----------

Check the worst state among source nodes.

best_state
----------

Get the best state among source nodes.

Tutorial
========

Create a node with the change_state task and save it in DB
----------------------------------------------------------

.. code-block:: python

   from canopsis.topology.elements import Node
   from canopsis.topology.manager import TopologyManager
   from canopsis.check import Check

   topologyManager = TopologyManager()

   # create a parameterized task
   task = {
      'id': 'canopsis.topology.rule.action.change_state',
      'params': {'update_entity': True}
   }
   entity_id = '/component/connector_name/connector/component'
   # create a node with previous task, default state to (WARNING) and bound to an entity
   node = Node(task=task, state=Check.WARNING, entity=entity_id)
   node.save(topologyManager)
