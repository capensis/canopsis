============
Configurable
============

.. contents:
   maxdepth: 2

.. module:: canopsis.configuration.configurable

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Objective
=========

This document aims to specify Backend configuration functionalities which may be agnostic from configuration languages and specific to user requirements.

Package contents
================

.. toctree::

   manager
   decorator

Functional description
======================

Such (re-)configuration is done by a ``canopsis.configuration.configurable.Configurable`` object and its ``apply_configuration`` method parameterized by a ``canopsis.configuration.parameters.Configuration`` object.

A configurable object uses configuration managers which are able to parse/write configuration resources (files, db, etc.) identified by a name. When a configurable is instantiated, it can watch its configuration resources in order to reconfigurate itself automatically if configuration content has been modified.

A configurable object has its own logger.

Technical description
=====================

.. class:: Configurable(object)

   In charge of persisting a Configuration coming from runtime or configuration files.

   Properties defined in configuration files can be overriden depending on the order of reading (from the first to the last).

   In such case, the last read parameter value is the parameter value to use.

   Configuration resources can be read in respect of configuration managers. Such configuration managers are given by the runtime or by global registration (see _configurationmanager).

   If you want to specialise your Configurable, at most three methods should be overriden: ``_conf``, ``_configure`` and ``_get_conf_files``.

   .. data:: CONF_RESOURCE = 'configuration/configurable.conf'

      Configuration resource where self properties can be found.

   .. data:: CONF = 'CONFIGURATION'

      Configuration category.

   .. data:: LOG = 'LOG'

      logging configuration category.

   .. data:: AUTO_CONF = 'auto_conf'

      If true, configurate this after initialization or when an automatic reconfiguration is required.

   .. data:: RECONF_ONCE = 'reconf_once'

      Ensure than the next automatic reconfiguration task will be done unless ``auto_conf`` is True. After doing the reconfiguration, this property will be setted to false.

   .. data:: CONF_FILES = 'conf_files'

      Files from where get configuration.

   .. data:: MANAGERS = 'managers'

      Managers able to read configuration files.

   .. data:: LOG_LVL = 'log_lvl'

      Logging level

   .. data:: LOG_NAME = 'log_name'

      Logging name

   .. data:: LOG_DEBUG_FORMAT = 'log_debug_format'

      Logging debug format

   .. data:: LOG_INFO_FORMAT = 'log_info_format'

      Logging info format

   .. data:: LOG_WARNING_FORMAT = 'log_warning_format'

      Logging warning format

   .. data:: LOG_ERROR_FORMAT = 'log_error_format'

      Logging error format

   .. data:: LOG_CRITICAL_FORMAT = 'log_critical_format'

      Logging critical format

   .. data:: DEBUG_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] [%(process)d] [%(thread)d] [%(pathname)s] [%(lineno)d] %(message)s"

      Default logging debug format

   .. data:: INFO_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s"

      Default logging debug format

   .. data:: WARNING_FORMAT = INFO_FORMAT

      Default logging debug format

   .. data:: ERROR_FORMAT = WARNING_FORMAT

      Default logging debug format

   .. data:: CRITICAL_FORMAT = ERROR_FORMAT

      Default logging debug format

   .. method:: apply_configuration(conf=None, conf_files=None, managers=None, logger=None)

      Apply conf on a destination in 5 phases:

      1. identify the right manager to use with conf_files to parse.
      2. for all conf_files, get conf which match with input conf.
      3. apply parsing rules on conf_file params.
      4. put values and parsing errors in two different dictionaries.
      5. returns both dictionaries of param values and errors.

      :param conf: conf from where get conf
      :type conf: Configuration

      :param conf_files: conf files to parse. If conf_files is a str, it is automatically putted into a list
      :type conf_files: list of str

   .. method:: get_configuration(conf=None, conf_files=None, logger=None, managers=None, fill=False)

      Get a dictionary of params by name from conf, conf_files and conf_managers

      :param conf: conf to update. If None, use self.conf.
      :type conf: Configuration

      :param conf_files: list of conf files. If None, use self.conf_files.
      :type conf_files: list of str

      :param logger: logger to use for logging info/error messages. If None, use self.logger
      :type logger: logging.Logger

      :param managers: conf managers to use. If None, use self.managers
      :type managers: list of ConfigurationManager

      :param fill: if True (False by default) load in conf all conf_files content.
      :type fill: bool

   .. method:: set_configuration(conf_file, conf, manager=None, logger=None)

      Set params on input conf_file.

      :param conf_files: conf_file to udate with params.
      :type conf_files: str

      :param conf: configuration to set.
      :type conf: (dict(str: dict(str: object))
      :param logger: logger to use to set params.
      :type logger: logging.Logger

   .. method:: configure(conf, logger=None)

      Update self properties with input params only if:
      - self.configure is True
      - self.auto_conf is True
      - param conf 'configure' is True
      - param conf 'auto_conf' is True

      This method may not be overriden. see _configure instead

      :param conf: object from where get paramters
      :type conf: Configuration

   .. method:: configure(unified_conf, logger=None)

      protected method to override in order to do a local configuration.

      unified_conf is a Configuration which contains respectively categories VALUES and ERRORS

   .. method:: _update_property(unified_conf, param_name, public_property)

      protected method which update an attribute of self related to an unified_conf, a param_name and public_property boolean. If public_property is True, the attribute is the param_name, else it's prefixed by '_'.

   .. method:: _get_conf_files()

      protected method to override in order to get the list of conf files.

   .. method:: _update_property(unified_conf, param_name, public=False)

      True if a property update is required and do it.

      Check if a param exist in paramters where name is param_name.
      Then update self property depending on input public:

      - True => name is param_name
      - False => name is '_{param_name}'

      The idea of the public argument permits to avoid to run an auto_conf in changing a private attribute in using its setter method.

      :param unified_conf: unified conf
      :type params: Configuration

      :param param_name: param name to find in params
      :type param_name: str

      :param public: If False (default), update directly private property, else update public property in using the property.setter
      :type property_name: bool

   .. method:: _get_conf_files()

      Get the first manager able to handle input conf_file. None if no manager is able to handle input conf_file.

      :return: first ConfigurationManager able to handle conf_file.
      :rtype: ConfigurationManager
