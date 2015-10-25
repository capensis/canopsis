.. _TR__Configurable:

=====================
Canopsis configurable
=====================

This document describes the technical solution of configurable in Canopsis.

.. contents::
   :depth: 2

-------
Updates
-------

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Jonathan Labéjof", "2015/10/22", "0.1", "Document creation"

----------
References
----------

- :ref:`FR::Configurable <fr__configurable>`

--------
Contents
--------

.. _TR__Configurable__Model:

Architecture
============

This project contains those modules:

- model: configuration model classes and tools.
- configurable: configurable base class and tools.
- drivers: configuration driver classes and tools.

.. _TR__Configurable__Model:

.. include:: ../developer-guide/api/configuration/canopsis.rst

TOTO

.. include:: ../developer-guide/api/configuration/canopsis.configuration.rst

Model
=====

Four classes are provided in order to modelize a configuration.

.. _FR__Configurable__configuration:

Configuration
=============

The configuration uses a unified langage which is agnostic from configuration resources, and several configuration resources in order to allow configuration overriding.

.. _FR__Configurable__Configuration__Overriding:

Configuration overriding/inheritance
------------------------------------

According to class inheritance, a Configurable can use self configuration resources, and ones defined at a base class level. In such way, the final configuration is based from base class configuration, and it is possible to apply a fined grained configuration with a specific Configuration instance.

.. _FR__Configurable__Configuration__Language:

Configuration Language
----------------------

The configuration langage is composed of three concepts:

- **Configuration**: set of categories by name.
- **Category**: set of parameters by name, and is identified by a name.
- **Parameter**: identified by a name. It uses also a default value and uses a parser in order to get a value from a serialized string value.

.. _FR__Configurable__Registry:

Configurable Registry
=====================

A configurable registry is a configurable which is composed of configurables.

.. _FR__Configurable__Driver:

Driver
======

A configuration driver permits to access to a configurable configuration, whatever the nature of the existing configuration resource.

Example of drivers are:

- FileDriver : dedicated to condfiguration files.
- DBDriver : dedicated to configuration database.
