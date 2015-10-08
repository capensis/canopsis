.. _FR__Storage:

=================
Canopsis Storages
=================

This document describes the concept of storage in Canopsis.

.. contents::
   :depth: 2

----------
References
----------

 - :ref:`FR::Middleware <FR__Middleware>`

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jonathan Labéjof", "2015/10/06", "0.1", "Document creation", ""

--------
Contents
--------

Description
===========

A storage is a common way to access to data, with contextual concerns.

.. _FR__Storage_Type:

Types of Storage
================

According to ORM (Object Relational Mapping) solutions, canopsis storages are an abstract layer to access to data, whatever data base implementation. Except than there are several types of canopsis storages which focuses on the contextual data properties, such as time, composition, etc. Instead of focusing only on data structure like ORM.

Here are types of Storages:

- **Storage**: default period, same behaviors than ORMs.
- **PeriodicStorage**: dedicated to data which have instances which exists on a period of time (for example, Obama is an instance of president which occured from past 7 years to next... let's "humanity" decide).
- **TimeStorage**: useful to data which have an instance at a specific moment. For example, a perfdata exists at a moment.
- **RelationalStorage**: specific to typed relational data. For example, a canopsis resource is embedded by a component.

According to storage types, the implementation must provide the best data model, which is independent from data structures.
