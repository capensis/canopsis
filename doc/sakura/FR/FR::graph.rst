=====
Graph
=====

This feature is a tool dedicated to analyze a system thanks to both information and relationships between information.

.. sectnum::

.. contents::
   :depth: 2

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jonathan Labéjof", "27/08/2015", "0.1", "Creation", ""

--------
Contents
--------

Description
===========

A graph is inherited from the mathematical structure of hypergraphs with some technical features such as information on edges.

Three type of elements describe a graph, the vertices_, the edges_ and the graphs_.

.. _vertices:

Vertice
-------

A vertice is an element which can embed an information.

It contains:

- unique id.
- several type names in order to ease classification of vertices.

Its lifecycle is independent from graphs, therefore, one vertice can be used by several graphs, but it does not directly graphs which use it.

.. _edges:

Edge
----

An edge inherits from vertice and links source vertice(s) with target vertice(s).

An edge contains the boolean property ``oriented`` which permits to distinguish or not target from sources.

.. _graphs:

Graph
-----

A graph is a vertice used such as a set of vertices.
