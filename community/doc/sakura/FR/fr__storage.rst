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
 - :ref:`FR::Webservice <FR__Webservice>`

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "David Delassus", "2015/10/20", "0.3", "Update storage definition"
   "David Delassus", "2015/10/20", "0.2", "Add webservice"
   "Jonathan Labéjof", "2015/10/06", "0.1", "Document creation", ""

--------
Contents
--------

.. _FR__Storage__

Description
===========

A storage is a common way to access to data, with contextual concerns.

.. _FR__Storage__Type:

Types of Storage
================

According to ORM (Object Relational Mapping) solutions, Canopsis storages are an abstract layer to access data, whatever database implementation. Except that there are several types of Canopsis storages which focus on the contextual data properties, such as time, composition, etc. Instead of focusing only on data structure like ORM.

Here are types of Storages:

- **Storage**: default period, same behaviors than ORMs.
- **PeriodicStorage**: dedicated to data which have instances which exists on a period of time (for example, Obama is an instance of president which occurred from past 7 years to next... let's "humanity" decide).
- **TimeStorage**: useful to data which have an instance at a specific moment. For example, a perfdata exists at a moment.
- **CompositeStorage**: specific to typed relational data. For example, a Canopsis resource is embedded by a component.
- **FileStorage**: useful to store/access files like a file system.
- **GraphStorage**: used to store complex relationships among data.

According to storage types, the implementation must provide the best data model, which is independent from data structures.

.. _FR__Storage__Definition:

Storage Definition
==================

A storage is defined by:

 - a protocol: specifying the database implementation
 - a data type: see :ref:`the above section <FR__Storage__Type>`
 - a data scope: used to separate data according to a scope (example: ``configuration``, ``context``, ``perfdata``, ...)

.. _FR__Storage__Webservice:

Serving storages
================

A :ref:`webservice <FR__Webservice>` **SHOULD** provide CRUD access to each storage.

The URL **SHOULD** be ``/storage/<protocol>/<data_type>/<data_scope>`` and each method
**SHOULD** resolve to the corresponding storage and CRUD operation.
