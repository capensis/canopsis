===================================
Manager: configuration file manager
===================================

.. contents:
    maxdepth: 2

.. module:: canopsis.configuration.manager
    :synopsis: configuration manager library

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

Manage reading or updating configuration files thanks to classes which inherit from the base class ConfigurationManager.

This class uses a meta-class in order to register every defined manager into a global manager list.

All managers to register must have the class attribute ``__register__`` to True (False by default).

.. class:: ConfigurationManager

.. data:: __register__

    If True, this class is automatically registered among a global list of ConfigurationManagers

.. method:: get_parameters(self, configuration_file, parsing_rules, logger)

    Same logic than for a class Configurable.

.. method:: handle(self, conf_file, logger)

.. method:: set_configuration(self, configuration_file, parameter_by_categories, logger)

    Same logic than for a class Configurable.

.. method:: get_configuration(self, conf_file, logger, conf=None, fill=False)

.. method:: get_managers() [static]

.. method:: get_manager(path) [static]

.. method:: _get_categories(self, conf_resource, logger)

.. method:: _get_parameters(self, conf_resource, category, logger)

.. method:: _has_category(self, conf_resource, category, logger)

.. method:: _has_parameter(self, conf_resource, category, param, logger)

.. method:: _get_conf_resource(self, logger, conf_file=None)

.. method:: _get_value(self, conf_resource, category, param)

.. method:: _set_category(self, conf_resource, category, logger)

.. method:: _set_parameter(self, conf_resource, category, param, logger)

.. method:: _update_conf_file(self, conf_resource, conf_file)
