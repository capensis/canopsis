.. _FR__Configurable:

=====================
Canopsis configurable
=====================

This document describes the concept of configurable in Canopsis.

.. contents::
   :depth: 2

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

A configurable is the base class for all python classes in Canopsis.

A configurable uses a configuration model, which is given to configuration drivers in order to:

- retrieve a specific configuration from at least one configuration resources.
- write a specific configuration to configuration resources.

The driver permits to be agnostic from the type of configuration resource. Therefore a configuration resource can be a file, or a database document, etc.

The configuration model is the common language for configuring canopsis python objects, which is agnostic from configuration resource.

.. _FR__Configurable__configuration:

Configuration
=============

The configuration uses a unified langage which is agnostic from configuration resources, and several configuration resources in order to allow configuration overriding.

.. _FR__Configurable__Configuration__Overriding:

Configuration overriding/inheritance
------------------------------------

According to class inheritance, a Configurable can use self configuration resources, and ones defined at a base class level. In such way, the final configuration is based from base class configuration, and it is possible to apply a fined grained configuration with a specific Configuration instance.

.. _FR__Configurable__Registry:

Configurable Registry
---------------------

A configurable registry is a configurable which is composed of configurables.

.. _FR__Configurable__Model:

Configuration Model
===================

The configuration model is composed of three concepts:

- **Configuration**: set of categories by name.
- **Category**: set of parameters by name, and is identified by a name.
- **Parameter**: identified by a name. It uses also a default value and uses a parser in order to get a value from a serialized string value.

.. _FR__Configurable__Driver:

Driver
======

A configuration driver permits to access to a configurable configuration, whatever the nature of the existing configuration resource.

Example of drivers are:

- FileDriver: dedicated to configuration files.
- IniDriver: dedicated to ini files.
- JSonDriver: dedicated to json files.
- DBDriver: dedicated to configuration database.
