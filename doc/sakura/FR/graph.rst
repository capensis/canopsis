.. fr-graph:

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

Data structure
==============

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

I contains:

- a ``weight`` which designates the edge weight.
- an ``oriented`` flag which distinguishes or not targets from sources.

.. _graphs:

Graph
-----

A graph is a vertice used such as a set of vertices.

Behavior
========

In addition to data information embedded by those elements, it might be important to use a specific element lifecycle behavior (creation, updating, deletion).

In order to ensure this idea, all elements can be used such as objects from the OOP paradigm

This related object is referenced by its class full name in the element.

Factory
=======

The factory is dedicated to instanciate a graph from a simple data format, which respects both `data structure`_ and behavior_.

Graph parsers might help any data format translation to the expected graph factory data format.
