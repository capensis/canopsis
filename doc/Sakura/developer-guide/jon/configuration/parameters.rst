==========
Parameters
==========

.. contents:
    maxdepth: 2

.. module:: canopsis.configuration.parameters

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Objective
=========

This document aims to specify how to parameterize your backend configuration.

Functional description
======================

A Configuration is a structure composed of categories, and all categories are composed of parameters.

All three are item objects (close to ordered dictionary).

Categories and parameters have a name.

Parameters of different categories can have the same name in order to apply overriding configuration in processing such configuration.

Finally, a parameter has two more attributes. A default value (the attribute of a configurable for example) and a parser function (str by default).

Technical description
=====================

.. class:: Configuration

   Dictionary of (Category name, Category) with implementation of methods __iter__, __delitem__, __getitem__, __len__, __iadd__ and __contains methods related to Category

   .. data:: ERRORS = 'ERRORS'

      Used in unified configuration to get the category which contain all wrong parameters (where values are exception instances).

   .. data:: VALUES = 'VALUES'

      Used in unified configuration to get the category which contain all parameter values.

   .. data:: categories

      Ordered dict of categories.

   .. method:: get(category_name, default=None)

   .. method:: setdefault(category_name, category)

   .. method:: put(category)

   .. method:: unify(copy=False)

      Unifiy this content in copying or not parameters (input copy).

   .. method:: get_unified_category(name, copy=False)

      Get a new configuration with only two categories, the first one named 'VALUES' contain all parameters where values are not None and not Exception in respecting the order of parameters for overriding concerns.
      The second one is named 'ERRORS'.

   .. method:: add_unified_category(name, copy=False, new_content=None)

   .. method:: clean

      Removes parameters without values.

   .. method:: copy

   .. method:: update(conf)

.. class:: Category

   Dictionary of (Property name, Property) with implementation of methods __iter__, __delitem__, __getitem__, __len__, __iadd__ and __contains methods related to Property

   .. data:: name

   .. data:: params

   .. method:: setdefault(param_name, param)

   .. method:: get(param_name, default=None)

   .. method:: put(param)

   .. method:: clean

   .. method:: copy(name=None)

.. class:: Parameter

   .. data:: name

   .. data:: value

   .. data:: parser

   .. method:: copy(name=None)

   .. method:: clean

   .. staticmethod:: array(item_type=str)

      Parser to use in order to get an array from a str where items are separated by ','.

   .. staticmethod:: bool(value)

      Parser to use in order to get a boolean value from a str which can equal to ``True``, ``true`` or ``1``. Any other value will set the parameter to False.
