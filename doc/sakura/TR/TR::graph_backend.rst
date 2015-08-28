=====
Graph
=====

Project ``canopsis.graph`` description.

.. sectnum::

.. contents::
   :depth: 2

----------
References
----------

- FR::graph

.. _graph: ./FR::graph/graph_

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jonathan Labéjof", "27/08/2015", "0.1", "Creation", ""

--------
Contents
--------

canopsis.graph.elements
=======================

Module which contains all graph elements

.. _graphelements:

Graph Element
-------------

Base object for all other graph objects.

Contains:

- uid: str.
- types: str(s).
- info: dict.

.. _vertices:

Vertice
-------

Inherits from the Vertice_.

.. _edges:

Edge
----

Inherits from the Vertice_.

Contains:

- sources: str(s).
- targets: str(s).
- oriented: bool.

.. _graphs:

Graph
-----

Inherits from the Vertice_.

Contains:

- elts: str(s)

canopsis.graph.manager
======================

GraphManager
------------

Manages graph storing and analyzing fonctions.

get_elts
>>>>>>>>

UTs
<<<

put_elt
>>>>>>>

UTs
<<<

put_elts
>>>>>>>>

UTs
<<<

remove_elts
>>>>>>>>>>>

UTs
<<<

del_edge_refs
>>>>>>>>>>>>>

UTs
<<<

update_elts
>>>>>>>>>>>

UTs
<<<

get_vertices
>>>>>>>>>>>>

UTs
<<<

get_edges
>>>>>>>>>

UTs
<<<

get_neighbourhood
>>>>>>>>>>>>>>>>>

UTs
<<<

get_sources
>>>>>>>>>>>

UTs
<<<

get_targets
>>>>>>>>>>>

UTs
<<<

get_orphans
>>>>>>>>>>>

UTs
<<<

get_graphs
>>>>>>>>>>

UTs
<<<

canopsis.graph.factory
======================

GraphFactory
------------

Instantiate a graph from a simple serialized format.

load
>>>>

UTs
<<<
