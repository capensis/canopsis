==========================
Backend (re-)configuration
==========================

.. contents:
	maxdepth: 2

.. module:: canopsis.organisation
    :synopsis: organisation library for managing backend configuration

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Objective
=========

This document aims to specify Backend configuration functionalities which may be agnostic from configuration languages and specific to user requirements.

Such (re-)configuration is done by a canopsis.configuration.Configurable object and its apply_configuration method parameterized by a canopsis.configuration.Configuration object.

Configuration file management
=============================

In an evolutive and adaptive system, configuration files may be :

1. language and technology agnostic,
2. easy to read and to modify by a human,
3. without ambiguities,
4. user requirement specific,
5. optionnaly composed of other configuration resources for modular reasons.

Therefore, configuration is done with cascading parameters and accessed through an UI.

Package contents
================

.. toctree::

   parameters
   configurable
   watcher
   manager/index
   ui/index

Technical description
=====================

.. data:: __version__ = '0.1'
