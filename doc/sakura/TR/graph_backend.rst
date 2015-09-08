.. tr-graph:

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

- :ref: `graph <../FR/graph>`

.. _graph: ../FR/graph_

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

.. _trgraphelements:
.. _trgraphelement:

Graph Element
-------------

Base object for all other graph objects.

Contains:

.. csv-table::
   :header: property, type, description

   uid, str, element id
   types, str(s), element types
   info, dict, element information
   cls, str, element class

A graph element exists in two formats, a serialized one, which is a dicitonary and an object. Both have same attributes.

.. _trvertices:
.. _trvertice:

Vertice
-------

Inherits from the Vertice_.

.. _tredges:
.. _tredge:

Edge
----

Inherits from the Vertice_.

Contains:

.. csv-table::
   :header: property, type, description

   sources, str(s), source element ids
   targets, str(s), target element ids
   weight, float, edge weight in [0; 1] (default 1)
   oriented, bool, oriented flag (True by default)

.. _trgraphs:
.. _trgraph:

Graph
-----

Inherits from the Vertice_.

Contains:

.. csv-table::
   :header: property, type, description

   elts, str(s), element ids used by the graph

canopsis.graph.manager
======================

.. _graphmanager:

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

load(elts)
>>>>>>>>>>

UTs
<<<

GraphParser
-----------

Translate a data format to the graph data format expected by the graph factory.

parse(data)
>>>>>>>>>>>

UTs
<<<
