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

The add_category decorator permits to specify a new category to add to previous configuration.

Technical description
=====================

.. decorator:: conf_paths(*conf_paths):

   Specify which conf_paths to use by a configurable class.

   :param conf_paths: configuration paths to apply on configurable class.

.. decorator:: add_category(name, unified=True, content=None)

   Add a category to a configurable configuration.

   :param name: category name
   :type name: str

   :param unified: if True (by default), the new category is unified from previous conf
   :type unified: bool

   :param content: category or list of parameters to add to the new category
   :type content: Category or list(Parameter)
