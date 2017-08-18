.. _FR__Middleware:

===================
Canopsis Middleware
===================

This document describes the concept of middleware in Canopsis.

.. contents::
   :depth: 2

----------
References
----------

 - :ref:`FR::Configurable <FR__Configurable>`

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jonathan Labéjof", "2015/10/07", "0.1", "Document creation", ""

--------
Contents
--------

Description
===========

A middleware is a common way to access to data, with contextual concerns.

Example of middleware are database, MOM (Middleware Oriented Message), data sharing or RPC (Remote-Procedure Call) technologies.

.. _FR__Middleware__Configuration:

Configuration
=============

Cause a middleware inherits from the Configurable, it is possible to configurate it in the same way than Configurables.

A more consistent way to configure it is to use an URI which contains all middleware information such as:

- protocol : middleware type to use. For example, 'storage' could designates the storage type of middleware.
- data_scope : middleware accessed data scope. For example, 'canopsis' could concern accessible data from canopsis.
- data_type : middleware accessed data type. For example, 'perfdata' could concern performance data.

Final URI example might respect this form:

{protocol}-{data_scope}-{data_type}://({user}(:{password})?@)?{host}(:{port})?/({path})?('?'{params})?

.. note::

   In a future release, the URI will look like this:

   {protocol}://({user}(:{password})?@)?{host}(:{port})?/({data_scope})?/({data_type})?/({path})?('?'{params})?

Where '{X}' designates X properties from this documentation or from expected URI parts.

For example, a 'mom' using:

- a broker named 'broker'.
- a scope 'canopsis' with data of type 'perfdata'.
- respective user and password 'b3j0f', 'f0j3b'.
- a sin of name 'sin'.
- a ttl = 60 s quality of service for message lifespan.

``mom://b3j0f:f0j3b@broker/canopsis/perfdata/sin?ttl=60``

And with the data_scope by default:

``mom://b3j0f:f0j3b@broker//perfdata/sin?ttl=60``
