.. tr-topology-backend:

========
Topology
========

Project ``canopsis.topology`` description.

.. sectnum::

.. contents::
   :depth: 2

----------
References
----------

- fr-topology_
- `fr-topology`_
- :ref: `fr-topology`
- :ref: `graph <./graph>`

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jonathan Labéjof", "27/08/2015", "0.1", "Creation", ""

--------
Contents
--------

canopsis.topology.elements
==========================

Module which contains all topology elements

.. _topovertices:

TopoVertice
-----------

Base object for all other topology objects.

info contains:

- state: int.
- operation: str.
- entity: str.

.. _toponodes:

TopoNode
--------

Inherits from the TopoVertice_.

.. _topoedges:

TopoEdge
--------

Inherits from the TopoNode_ and the Edge_.

.. _graphs:

Graph
-----

Inherits from the TopoNode_ and the Graph_.

canopsis.topology.manager
=========================

TopologyManager
---------------

Inherits from the GraphManager_.

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

The webservice provides all graphmanager_ methods through the route 'topology/' plus topologymanager_ methods below:

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
