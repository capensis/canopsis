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

Technical description
=====================

.. class:: Configuration

   Dictionary of (Category name, Category) with implementation of methods __iter__, __delitem__, __getitem__, __len__, __iadd__ and __contains methods related to Category

   .. data:: categories

   .. method:: get(category_name, default=None)

   .. method:: setdefault(category_name, category)

   .. method:: put(category)

   .. method:: unify(copy=False)

   .. method:: get_unified_category(name, copy=False)

   .. method:: add_unified_category(name, copy=False, new_content=None)

   .. method:: clean

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

   .. method:: copy(name=None)

   .. method:: clean

   .. staticmethod:: array(item_type=str)

   .. staticmethod:: bool(value)
