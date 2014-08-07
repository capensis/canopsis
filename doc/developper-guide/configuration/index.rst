=====================
Backend configuration
=====================

.. contents:
	maxdepth: 2

.. module:: canopsis.organisation
    :synopsis: organisation library for managing human resources

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
5. optionnaly composed of other configuration files for modular reasons.

Therefore, configuration is done with cascading parameters and accessed through an UI.

Package contents
================

.. toctree::

   watcher
   manager/index
   ui/index

   :ref: _configuration
   :ref: _category
   :ref: _parameter
   :ref: _configurable

Technical description
=====================

.. class:: Configuration

    Dictionary of (Category name, Category) with implementation of methods __iter__, __delitem__, __getitem__, __len__, __iadd__ and __contains methods related to Category

.. data:: categories

.. function:: get

.. function:: put

.. function:: unify

.. function:: get_unified_category

.. function:: add_unified_category

.. function:: clean

.. function:: copy

.. function:: update

.. class:: Category

    Dictionary of (Property name, Property) with implementation of methods __iter__, __delitem__, __getitem__, __len__, __iadd__ and __contains methods related to Property

.. data:: name

.. data:: params

.. function:: setdefault

.. function:: get

.. function:: put

.. function:: clean

.. function:: copy

.. class:: Parameter

.. data:: name

.. data:: value

.. function:: copy

.. function:: clean

.. class:: Configurable

    In charge of persisting a Configuration coming from runtime or configuration files.

    Properties defined in configuration files can be overriden depending on the order of reading (from the first to the last).

    In such case, the last read parameter value is the parameter value to use.

    configuration files can be read in respect of configuration managers. Such configuration managers are given by the runtime or by global registration (see _configurationmanager).

    If you want to specialise your Configurable, at most two methods should be overriden: _configure and _get_conf_files

.. data:: auto_conf

.. data:: once

.. data:: conf_files

.. data:: managers

.. data:: conf

.. data:: logger

.. data:: log_lvl

.. data:: log_name

.. data:: log_debug_format

.. data:: log_info_format

.. data:: log_warning_format

.. data:: log_error_format

.. data:: log_critical_format

.. function:: apply_configuration(self, conf=None, conf_files=None, managers=None)

.. function:: get_configuration(self, conf=None, conf_files=None, logger=None, managers=None, fill=False)

.. function:: set_configuration(self, conf_file, conf, manager=None, logger=None)

.. function:: configure(self, conf)

    apply input conf to self

.. function:: _configure(self, unified_conf)

    protected method to override in order to do a local configuration.

    unified_conf is a Configuration which contains respectively categories VALUES and ERRORS

.. function:: _update_property(self, unified_conf, param_name, public_property)

    protected method which update an attribute of self related to an unified_conf, a param_name and public_property boolean. If public_property is True, the attribute is the param_name, else it's prefixed by '_'

.. function:: _get_conf_files(self)

    protected method to override in order to get the list of conf files.
