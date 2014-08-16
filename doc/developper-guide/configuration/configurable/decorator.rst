==================================
Decorator: configurable decorators
==================================

.. contents:
   maxdepth: 2

.. module:: canopsis.configuration.configurable.decorator

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Objective
=========

This document aims to specify decorators useful to configurable objets. Those decorators are provided in order to ease the use of the complex configurable library.

Functional description
======================

This library aims to simplify the use of the configurable library.

The conf_paths decorator permits to specify which conf_path a Configurable object should use.

Technical description
=====================

.. decorator:: conf_paths(*conf_paths):

   Specify which conf_paths to use by a configurable class.

   :param conf_paths: configuration paths to apply on configurable class.
