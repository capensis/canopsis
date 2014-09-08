=================================
Manager: manager of configurables
=================================

.. contents:
   maxdepth: 2

.. module:: canopsis.configuration.configurable.manager

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Objective
=========

This document aims to specify manager of configurables which ease the configuration of a set of configurables.

Functional description
======================

A configurable manager is a configurable which is composed of configurables (getter property ``configurables``) and configurable type restrictions (getter property ``configurable_types``).

Let's call those configurables ``sub-configurable``.

It allows to easily specify how a configurable is (re-)configured, from its simple object path (example: ``canopsis.configuration.configurable.Configurable``) to fine grained (re-)configuration thanks to dedicated categories created from the same politic.

When a sub configurable is specified, it becomes accessible from ``__getitem__``, ``__setitem__`` and ``__delitem__`` accessors (same as dictionary elements).

In a configuration, a sub-configurable is specified in the manager category in all parameters which ends by ``{name}_value``. And all sub-configurable uses a generated category which unifies their default configuration and which is named {NAME}_CONF.

Finally, it is possible to specifying types of sub-configurable because such configurables may have dedicated API.

Example:
Let a ini configuration file like::

   [MANAGER]
   auto_conf=false
   test_value=canopsis.configuration.configurable.Configurable
   test_type=canopsis.configuration.configurable.Configurable
   test2_type=canopsis.configuration.configurable.manager.Manager

   [TEST_CONF]
   auto_conf=false

This configuration file specify a manager with the configurable ``test`` and two configurable type restrictions on configurables ``test`` and ``test2``.

After configuring the manager, it is possible to get the sub-configurable ``test`` in using item accessors manager['test']. And related type restriction from manager['test_type'].

The configurable ``test`` is an instance of `̀Configurable``.

Doing the affectation manager['test'] = Configurable() will be ok only if manager['test'] is an instance of manager['test_type'], Configurable otherwise.

Finally, it is possible to delete a configurable in doing del manager['test'].

In above configuration, ``manager['test2'] = Configurable`` will raise a ``Manager.Error`` exception because the type restriction ``test2_type`` is not a base class of ``Configurable``.
And test.auto_conf is false as specified in the configuration in the category ``TEST_CONF``.

Technical description
=====================

.. class:: Configurables(dict)

   With a ConfigurableTypes, it is in charge of a Manager sub-configurables.
   When a configurable is trying to be setted, the type is checked related
   to its Manager ConfigurableTypes.

   .. method:: __init__(manager, values=None, *args, **kwargs)

      :param manager: related manager
      :type manager: Manager

      :param values: default values if not None (default)
      :type values: dict

   .. method:: __setitem__(name, value)

      Set a new configurable value only if value inherits from manager.configurable_types[name].

      :param name: new configurable name
      :type name: str

      :param value: new configurable value
      :type value: str (path) or class or instance

.. class:: ConfigurableTypes(dict)

   With a Configurables, it is in charge of a set of configurable type.
   When a new type is setted but the old configurable value does not inherits
   from it, then the old value is removed automatically.

   .. method:: __init__(manager, values=None, *args, **kwargs)

      :param manager: related manager
      :type manager: Manager

      :param values: default values if not None (default)
      :type values: dict

   .. method:: __setitem__(name, value)

      Set a new configurable type.

      :param name: new configurable name.
      :type name: str

      :param value: new type value.
      :type value: str (path) or class

.. class:: Manager(canopsis.configuration.configurable.Configurable):

   Manage a set of configurables which are accessibles
   from self.configurables.

   Each configurable can be defined in conf parameters where names are like
   {name}_configurable={configurable_path, configurable_class, configurable}

   And a configurable configuration are in categories {NAME}_CONF.

   .. class:: Error(Exception):
      """handle manager errors"""
        pass

   CONF_PATH = 'configuration/manager.conf'

   CATEGORY = 'MANAGER'

   CONFIGURABLE_SUFFIX = '_value'
   CONFIGURABLE_TYPE_SUFFIX = '_type'

   .. method:: __init__(configurables=None, configurable_types=None, *args, **kwargs)

      :param configurables: dictionary of configurables by name.
      :type configurables: dict

      :param configurable_types: dictionary of configurable types by name
      :type configurable_types: dict

   .. method:: _get_category()

      Get category.

   .. property:: configurables

      Configurables which manages sub-configurables

   .. property:: configurable_types

      ConfigurableTypes which manages restriction of sub-configurable types

   .. method:: __contains__(name)

      Redirection to self.configurables.__contains__

   .. method:: __getitem__(name)

      Redirection to self.configurables.__getitem__

   .. method:: __setitem__(name, value)

      Redirection to self.configurables.__setitem__

   .. method:: __delitem__(name)

      Redirection to self.configurables.__delitem__

   .. staticmethod:: get_configurable_category(name):

      Get generated sub-configurable category name

   .. staticmethod:: get_configurable(configurable, *args, **kwargs)

      Get a configurable instance from a configurable class/path/instance and args, kwargs, None otherwise.

      :param configurable: configurable path, class or instance
      :type configurable: str, class or Configurable

      :return: configurable instance or None if input configurable can not be solved such as a configurable.
