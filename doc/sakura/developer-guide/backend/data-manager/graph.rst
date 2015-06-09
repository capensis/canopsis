.. _dev-backend-mgr-graph:

GRAPH: library for managing graph structure
===========================================

Graph elements
--------------

Functional
~~~~~~~~~~

A graph is a mean to construct information based on logical understanding of
information.

A graph is composed of vertices such as elementary information and edges which
are logical relationships between information.

With the previous definition, a graph is a complex vertice which results in
providing information.

In such way, an information can be elementary or composed of graph information
such as hypergraphs.

An hypergraph permits to add a context dimension over graph elements. In such
structure, vertices, edges and graphs exist in multiple graphs, and their
surround depends on the graph they are associated at a "time".

Technical
~~~~~~~~~

For simplification reasons, a graph is technically solved by such concepts.

Graph Element
*************

The graph element is the base concept of elements decribed here. It has a
unique identifier among all graph elements and a type for graph specialisation
reasons.

A graph element contains::

    - id: unique identifier among all graph elements.
    - type: type of graph element. A graph could be a topology or something
        else.
    - _cls: python class path.
    - _type: base type which permits to recognize the type of element.

Graph vertice
*************

A graph vertice inherits from the graph element and can contain data
information.

From a graph vertice, it is possible to resolve neighbour vertices thanks to
edges.

A graph vertice contains::

    - data: vertice information.
    - _type: equals vertice.

Graph edge
**********

Technically, a graph edge is more rich than its representation in the
functional definition in order to ease its manipulation in a richer context
instead of keeping only a logical use of edges. It becomes possible to describe
logical information between two edges.

A graph edge inherits from the graph vertice in order to transport information
and can bind several source with several targets, directly or not.

    - sources: source vertices.
    - targets: target vertices.
    - directed: directed orientation. If False, source and target vertices are
        directly connected, otherwise, only sources are directly connected to
        targets.
    - weight: edge weight. Default 1.
    - _type: equals edge.

A graph inherits from vertice and contains::

    - elts: elements existing in this graph.
    - _type: graph.

Graph Manager
-------------

.. module:: canopsis.graph.manager

This module defines the GraphManager which interacts between graph elements and
the DB.

Functional
~~~~~~~~~~

The role of the GraphManager is to ease graph element CRUD operations and
to retrieve graphs, vertices and edges thanks to methods with all element
parameters useful to find them.

Technical
~~~~~~~~~

The graph manager permits to get graph elements with any context information.

One, generic methods permit to get/put/delete elements in understanding such
elements such as dictionaries or GraphElement depending on serialize parameter
value.

Two, it is possible to find graphs, vertices and edges thanks to parameters
which correspond to their properties.

Tutorial
--------

Create a Vertice, an edge and a graph and save them
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: python

   from canopsis.graph.elements import Vertice, Edge, Graph
   from canopsis.graph.manager import GraphManager

   graphManager ~ GraphManager()
   # create a vertice task
   task ~ {'task': 'canopsis.task.task'}
   # create a vertice
   vertice ~ Vertice(data~task)
   # save it in DB
   vertice.save(graphManager)
   # create an edge which bind the vertice to itself
   edge ~ Edge(sources~vertice, targets~vertice)
   # save it in DB
   edge.save(graphManager)
   # create a graph wich contains vertice and edge
   graph ~ Graph(elts~[vertice, edge])
   # save it in DB
   graph.save(graphManager)

Find graph elements such as dictionaries
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: python

   from canopsis.graph.manager import GraphManager

   # graph id
   graph_ids ~ ['graph_id0', 'graph_id1']

   graphManager ~ GraphManager()

   # find elements from graph ids where types are dictionaries
   elts ~ graphManager.get_edges(graph_ids~graph_ids, serialize~False)

Delete graph elements
~~~~~~~~~~~~~~~~~~~~~

.. code-block:: python

   from canopsis.graph.manager import GraphManager

   # graph id
   graph_ids ~ ['graph_id0', 'graph_id1']

   graphManager ~ GraphManager()

   elts ~ []

   # find elements from graph ids where types are dictionaries
   elts +~ graphManager.get_edges(graph_ids~graph_ids, serialize~False)
   # and get graphs
   elts +~ graphManager.get_graphs(ids~graph_ids)
   # delete them from DB
   for elt in elts:
      elt.delete(graphManager)
