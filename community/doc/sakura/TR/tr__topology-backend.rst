.. _TR__Topology_backend:

========
Topology
========

Project ``canopsis.topology`` description.

.. contents::
   :depth: 2

----------
References
----------

- :ref:`FR__Topology <FR__Topology>`
- :ref:`TR__graph <TR__Graph_backend>`

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "01/09/2015", "0.3", "Update references", ""
   "David Delassus", "01/09/2015", "0.2", "Rename document", ""
   "Jonathan Labéjof", "2015/08/27", "0.1", "Creation", ""

--------
Contents
--------

canopsis.topology.elements
==========================

Module which contains all topology elements

.. _TR__Topology_backend__TopoVertice:

TopoVertice
-----------

Base object for all other topology objects.

info contains:

- state: int.
- operation: str.
- entity: str.

.. _TR__Topology_backend__TopoNode:

TopoNode
--------

Inherits from the :ref:`TopoVertice <TR__Topology_backend__TopoVertice>`.

.. _TR__Topology_backend__TopoEdge:

TopoEdge
--------

Inherits from the :ref:`TopoNode <TR__Topology_backend__TopoNode>` and the
:ref:`TR__Graph_backend__Edge`.

.. _TR__Topology_backend__Graph:

Graph
-----

Inherits from the :ref:`TopoNode <TR__Topology_backend__TopoNode>` and the
:ref:`Graph <TR__Graph_backend__Graph>.

canopsis.topology.manager
=========================

.. _TR__Topology_backend__Manager:

TopologyManager
---------------

Inherits from the :ref:`GraphManager <TR__Graph_backend__Manager>`.

Manages topology storing and analyzing fonctions.

get_causals
>>>>>>>>>>>

UTs
<<<

get_consequences
>>>>>>>>>>>>>>>>

UTs
<<<

get_causalsandconsequencespertopo
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

UTs
<<<

get_consequencesandcausalspertopo
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

UTs
<<<

canopsis.topology.webservices
=============================

The webservice provides all :ref:`GraphManager <TR__Graph_backend__Manager>`
methods through the route 'topology/' plus
:ref:`GraphManager <TR__Topology_backend__Manager>` methods below:

get_causals
-----------

'topology/causals'

get_consequences
----------------

'topology/consequences'

get_causalsandconsequencespertopo
---------------------------------

'topology/causalsandconsequences'

get_consequencesandcausalspertopo
---------------------------------

'topology/consequencesandcausals'
