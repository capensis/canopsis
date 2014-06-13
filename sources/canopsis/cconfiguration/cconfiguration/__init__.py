#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

__version__ = "0.1"

__all__ = ('Configurable')

from logging import Formatter, getLogger, FileHandler, Filter

from stat import ST_SIZE

from os.path import expanduser, exists, sep
from os import stat

from collections import OrderedDict

from inspect import isclass


class Configuration(object):
    """
    Manage configuration such as a list of Categories.

    The order of categories permit to ensure parameter overriding.
    """

    def __init__(self, *categories, **kwargs):
        """
        :param categories: categories to configure.
        :type categories: list of Category.
        """

        super(Configuration, self).__init__(**kwargs)

        # set categories
        self.categories = OrderedDict()
        for category in categories:
            self.categories[category.name] = category

    def __iter__(self, *args, **kwargs):

        return iter(self.categories.values())

    def __delitem__(self, category_name):

        del self.categories[category_name]

    def __getitem__(self, category_name, *args, **kwargs):

        return self.categories[category_name]

    def __contains__(self, category_name, *args, **kwargs):

        return category_name in self.categories

    def __len__(self):

        return len(self.categories)

    def get(self, category_name, default=None, *args, **kwargs):

        return self.categories.get(category_name, default)

    def setdefault(self, category_name, category, *args, **kwargs):

        return self.categories.setdefault(category_name, category)

    def put(self, category, *args, **kwargs):
        """
        Put a category and return the previous one if exist
        """

        result = self.get(category.name)
        self.categories[category.name] = category
        return result

    def get_parameters(self, *args, **kwargs):
        """
        Get values and errors of parameters in respecting parameter overriding.

        :return: two dictionaries respectively with values and errors by name.
        :rtype: dict, dict
        """

        values = dict()
        errors = dict()

        result = values, errors

        for category in self:

            for parameter in category:

                if parameter.value is not None:

                    to_update, to_delete = (errors, values) if \
                        isinstance(parameter.value, Exception) \
                        else (values, errors)

                    to_update[parameter.name] = parameter.value

                    if parameter.name in to_delete:
                        del to_delete[parameter.name]

        return result

    def clean(self, *args, **kwargs):
        """
        Clean this parameters in setting value to None.
        """

        for category in self:

            category.clean()

    def copy(self, *args, **kwargs):
        """
        Copy this Configuration
        """

        result = Configuration()

        for category in self:
            result.put(category.copy())

        return result

    def update(self, configuration, *args, **kwargs):
        """
        Update this content with input configuration
        """

        for category in configuration:
            category = self.setdefault(
                category.name, category.copy())

            for parameter in category:
                parameter = category.setdefault(
                    parameter.name, parameter.copy())


class Category(object):
    """
    Parameter category which contains a dictionary of parameters.
    """

    def __init__(self, name, *parameters, **kwargs):
        """
        :param name: unique in a configuration.
        :type name: str

        :param parameters: Parameters
        :type parameters: list of Parameter
        """
        super(Category, self).__init__(**kwargs)
        self.name = name
        # set parameter by names.
        self.parameters = {
            parameter.name: parameter for parameter in parameters}

    def __iter__(self, *args, **kwargs):

        return iter(self.parameters.values())

    def __delitem__(self, parameter_name, *args, **kwargs):

        del self.parameters[parameter_name]

    def __getitem__(self, parameter_name, *args, **kwargs):

        return self.parameters[parameter_name]

    def __contains__(self, parameter_name, *args, **kwargs):

        return parameter_name in self.parameters

    def __len__(self):

        return len(self.parameters)

    def setdefault(self, parameter_name, parameter, *args, **kwargs):

        return self.parameters.setdefault(parameter_name, parameter)

    def get(self, parameter_name, default=None, *args, **kwargs):

        return self.parameters.get(parameter_name, default)

    def put(self, parameter, *args, **kwargs):
        """
        Put a parameter and return the previous one if exist
        """

        result = self.get(parameter.name)
        self.parameters[parameter.name] = parameter
        return result

    def clean(self, *args, **kwargs):
        """
        Clean this parameters in setting value to None.
        """

        for parameter in self.parameters.values():

            parameter.clean()

    def copy(self, name=None, *args, **kwargs):

        if name is None:
            name = self.name

        result = Category(name)

        for parameter in self:
            result.put(parameter.copy())

        return result


class Parameter(object):
    """
    Parameter identified among a category by its name.
    Provide a value and a parser (str by default).
    """

    def __init__(self, name, value=None, parser=str, *args, **kwargs):
        """
        :param name: unique by category
        :type name: str

        :param value: parameter value. None if not given.
        :type value: object

        :param parser: parameter test deserializer which takes in parameter
            a str.
        :type parser: callable
        """

        super(Parameter, self).__init__(*args, **kwargs)
        self.name = name
        self._value = value
        self.parser = parser

    @property
    def value(self):
        return self._value

    @value.setter
    def value(self, value):
        if isinstance(value, str):
            # parse value if str
            try:
                self._value = self.parser(value)

            except Exception as e:
                self._value = e

        else:
            self._value = value

    def copy(self, name=None, *args, **kwargs):

        if name is None:
            name = self.name

        result = Parameter(name, value=self.value, parser=self.parser)

        return result


class Configurable(object):
    """
    Manages class configuration synchronisation with configuration files.
    """

    CONF_FILE = '~/etc/conf.conf'

    CONF = 'CONF'
    LOG = 'LOG'

    AUTO_CONF = 'auto_conf'
    CONFIGURE = 'configure'
    MANAGERS = 'conf_managers'
    LOG_NAME = 'log_name'
    LOG_LVL = 'log_lvl'
    LOG_DEBUG_FORMAT = 'log_debug_format'
    LOG_INFO_FORMAT = 'log_info_format'
    LOG_WARNING_FORMAT = 'log_warning_format'
    LOG_ERROR_FORMAT = 'log_error_format'
    LOG_CRITICAL_FORMAT = 'log_critical_format'

    DEFAULT_CONFIGURATION = Configuration(
        Category(CONF,
            Parameter(AUTO_CONF, parser=bool),
            Parameter(MANAGERS),
            Parameter(CONFIGURE, parser=bool)),
        Category(LOG,
            Parameter(LOG_NAME),
            Parameter(LOG_LVL),
            Parameter(LOG_DEBUG_FORMAT),
            Parameter(LOG_INFO_FORMAT),
            Parameter(LOG_WARNING_FORMAT),
            Parameter(LOG_ERROR_FORMAT),
            Parameter(LOG_CRITICAL_FORMAT)))

    DEBUG_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] \
[%(process)d] [%(thread)d] [%(pathname)s] [%(lineno)d] %(message)s"
    INFO_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s"
    WARNING_FORMAT = INFO_FORMAT
    ERROR_FORMAT = WARNING_FORMAT
    CRITICAL_FORMAT = ERROR_FORMAT

    def __init__(
        self,
        conf_files=None, auto_conf=True,
        managers=None, configuration=DEFAULT_CONFIGURATION.copy(),
        log_lvl='INFO', log_name=None, log_info_format=INFO_FORMAT,
        log_debug_format=DEBUG_FORMAT, log_warning_format=WARNING_FORMAT,
        log_error_format=ERROR_FORMAT, log_critical_format=CRITICAL_FORMAT,
        _ready_to_conf=True,
        *args, **kwargs
    ):
        """
        :param conf_files: conf_files to parse
        :type conf_files: Iterable or str

        :param auto_conf: true force auto conf as soon as parameter change
        :type auto_conf: bool

        :param configuration: configuration with parsing rules
        :type configuration: Configuration

        :param log_lvl: logging level
        :type log_lvl: str

        :param _ready_to_conf: protected parameter permetting to deactivate
            auto_conf processing in this call
        :type _ready_to_conf: bool
        """

        super(Configurable, self).__init__(*args, **kwargs)

        if conf_files is None:
            conf_files = [Configurable.CONF_FILE]

        self.conf_files = conf_files

        self.auto_conf = auto_conf

        # set logging properties
        self._log_lvl = log_lvl
        self._log_name = log_name if log_name is not None else \
            type(self).__name__
        self._log_debug_format = log_debug_format
        self._log_info_format = log_info_format
        self._log_warning_format = log_warning_format
        self._log_error_format = log_error_format
        self._log_critical_format = log_critical_format

        self._logger = self.newLogger()

        self.configuration = configuration

        # set managers
        if managers is None:
            from cconfiguration.manager import ConfigurationManager

            managers = ConfigurationManager.get_managers()

        self.managers = set(managers)

        if _ready_to_conf and self.auto_conf:
            self.apply_configuration()

    def newLogger(self):
        """
        Get a new logger related to self properties.
        """

        result = getLogger(self.log_name)
        result.setLevel(self.log_lvl)

        def setHandler(logger, lvl, path, format):

            class filter(Filter):
                def filter(self, record):
                    return record.levelno == lvl

            handler = FileHandler(path)
            handler.addFilter(Filter())
            handler.setLevel(lvl)
            formatter = Formatter(format)
            handler.setFormatter(formatter)
            logger.addHandler(handler)
            setattr(logger, lvl, handler)

        filename = self.log_name.replace('.', sep)
        path = expanduser('~/var/log/{0}.log'.format(filename))

        setHandler(result, 'DEBUG', path, self.log_debug_format)
        setHandler(result, 'INFO', path, self.log_info_format)
        setHandler(result, 'WARNING', path, self.log_warning_format)
        setHandler(result, 'ERROR', path, self.log_error_format)
        setHandler(result, 'CRITICAL', path, self.log_critical_format)

        return result

    @property
    def log_debug_format(self):
        return self._log_debug_format

    @log_debug_format.setter
    def log_debug_format(self, value):
        self._log_debug_format = value
        self._logger = self.newLogger()

    @property
    def log_info_format(self):
        return self._log_info_format

    @log_info_format.setter
    def log_info_format(self, value):
        self._log_info_format = value
        self._logger = self.newLogger()

    @property
    def log_warning_format(self):
        return self._log_warning_format

    @log_warning_format.setter
    def log_warning_format(self, value):
        self._log_warning_format = value
        self._logger = self.newLogger()

    @property
    def log_error_format(self):
        return self._log_error_format

    @log_error_format.setter
    def log_error_format(self, value):
        self._log_error_format = value
        self._logger = self.newLogger()

    @property
    def log_critical_format(self):
        return self._log_critical_format

    @log_critical_format.setter
    def log_critical_format(self, value):
        self._log_critical_format = value
        self._logger = self.newLogger()

    @property
    def log_name(self):
        return self._log_name

    @log_name.setter
    def log_name(self, value):
        self._log_name = value
        self._logger = self.newLogger()

    @property
    def log_lvl(self):
        """
        Get this logger lvl.

        :return: self logger lvl
        :rtype: str
        """
        return self._log_lvl

    @log_lvl.setter
    def log_lvl(self, value):
        """
        Change of logging level.

        :param value: new log_lvl to set up.
        :type value: str
        """

        self._logger.setLevel(value)

    @property
    def logger(self):
        return self._logger

    @property
    def conf_files(self):
        """
        :return: self configuration files
        :rtype: tuple
        """

        if not hasattr(self, '_conf_files'):
            self._conf_files = list()

        result = self._conf_files

        return result

    @conf_files.setter
    def conf_files(self, value):
        """
        Change of conf_files in adding it in watching list.
        """

        from cconfiguration.watcher import add_configurable,\
            remove_configurable

        # remove previous watching
        remove_configurable(self)
        self._conf_files = tuple(value)
        # add new watching
        add_configurable(self)

    def apply_configuration(
        self, configuration=None, conf_files=None,
        managers=None, *args, **kwargs
    ):
        """
        Apply configuration on a destination in 5 phases:

        1. identify the right manager to use with conf_files to parse.
        2. for all conf_files, get configuration which match
            with input configuration.
        3. apply parsing rules on conf_file parameters.
        4. put values and parsing errors in two different dictionaries.
        5. returns both dictionaries of parameter values and errors.

        :param configuration: configuration from where get parsers
        :type configuration: Configuration

        :param conf_files: configuration files to parse. If
            conf_files is a str, it is automatically putted into a list
        :type conf_files: list of str
        """

        if configuration is None:
            configuration = self.configuration.copy()

        configuration = self.get_configuration(
            configuration=configuration, conf_files=conf_files,
            managers=managers, *args, **kwargs)

        self.configure(configuration=configuration, *args, **kwargs)

    def get_configuration(
        self,
        configuration=None, conf_files=None, logger=None,
        managers=None, fill=False, *args, **kwargs
    ):
        """
        Get a dictionary of parameters by name from configuration,
        conf_files and conf_managers

        :param configuration: configuration to update. If None, use \
            self.configuration
        :type configuration: Configuration

        :param conf_files: list of configuration files. If None, use \
            self.conf_files
        :type conf_files: list of str

        :param logger: logger to use for logging info/error messages.
            If None, use self.logger
        :type logger: logging.Logger

        :param managers: conf managers to use. If None, use self.managers
        :type managers: list of ConfigurationManager

        :param fill: if True (False by default) load in configuration all \
            conf_files content
        :type fill: bool
        """

        # start to initialize input parameters
        if logger is None:
            logger = self._logger

        if configuration is None:
            configuration = self.configuration

        if conf_files is None:
            conf_files = self._conf_files

        if isinstance(conf_files, str):
            conf_files = [conf_files]

        # clean configuration file list
        conf_files = [
            expanduser(conf_file) for conf_file
            in conf_files]

        if managers is None:
            managers = self.managers

        # iterate on all conf_files
        for conf_file in conf_files:

            if not exists(conf_file) or stat(conf_file)[ST_SIZE] == 0:
                continue

            conf_manager = self._get_manager(
                conf_file=conf_file,
                logger=logger, managers=managers)

            # if a config_resource is not None
            if conf_manager is not None:

                configuration = conf_manager.get_configuration(
                    configuration=configuration, fill=fill,
                    conf_file=conf_file, logger=logger)

            else:
                # if no conf_manager, display a warning log message
                logger.warning('No manager found among {0} for {1}'.format(
                    conf_file))

        return configuration

    def set_configuration(
        self, conf_file, configuration, manager=None,
        logger=None, *args, **kwargs
    ):
        """
        Set parameters on input conf_file.

        Args:
            - conf_files (str): conf_file to udate with
                parameters
            - parameter_by_categories (dict(str: dict(str: object)):
            - logger (logging.Logger): logger to use to set parameters.
        """

        result = None

        if logger is None:
            logger = self._logger

        # first get content of input conf_file
        prev_manager = self._get_manager(
            conf_file=conf_file,
            logger=logger,
            managers=self.managers)

        if prev_manager is not None:
            prev_conf = prev_manager.get_configuration(
                conf_file=conf_file,
                logger=logger)

        # try to find a good manager if manager is None
        if manager is None:
            manager = self._get_manager(
                conf_file=conf_file,
                logger=logger,
                managers=self.managers)

        elif isclass(manager):
            manager = manager()

        # if prev manager is not the new manager
        if type(manager) is not type(prev_manager):
            # update prev_conf with input configuration
            prev_conf.update(configuration)
            configuration = prev_conf

        if manager is not None:
            manager.set_configuration(
                conf_file=conf_file,
                configuration=configuration,
                logger=logger)

        else:
            logger.error(
                'No ConfigurationManager found for \
                configuration file {0}'.format(
                    conf_file))

        return result

    def configure(self, configuration, *args, **kwargs):
        """
        Update self properties with input parameters only if:
        - self.configure is True
        - self.auto_conf is True
        - parameter configuration 'configure' is True
        - parameter configuration 'auto_conf' is True

        :param configuration: object from where get paramters
        :type configuration: Configuration
        """

        parameters, error_parameters = configuration.get_parameters()

        # set configure
        self.configure = parameters.get(Configurable.CONFIGURE, self.configure)

        # set auto_conf
        self.auto_conf = parameters.get(Configurable.AUTO_CONF, self.auto_conf)

        if self.configure or self.auto_conf:
            self._configure(parameters, error_parameters, *args, **kwargs)

    def _configure(self, parameters, error_parameters, *args, **kwargs):
        """
        Configure this class with input parameters only if auto_conf or
        configure is true.

        This method shouldn't be overriden.

        :param parameters: dictionary of parameter value by name
        :type parameters: dict

        :param error_parameters: dictionary of parameter parsing error by name
        :type error_parameters: dict

        :param configure: if True, force full self configuration
        :type configure: bool
        """

        # set log_lvl
        self.log_lvl = parameters.get(Configurable.LOG_LVL, self.log_lvl)

        # set log_debug_format
        self.log_debug_format = parameters.get(
            Configurable.LOG_DEBUG_FORMAT, self.log_debug_format)

        # set log_info_format
        self.log_info_format = parameters.get(
            Configurable.LOG_INFO_FORMAT, self.log_info_format)

        # set log_warning_format
        self.log_warning_format = parameters.get(
            Configurable.LOG_WARNING_FORMAT, self.log_warning_format)

        # set log_error_format
        self.log_error_format = parameters.get(
            Configurable.LOG_ERROR_FORMAT, self.log_error_format)

        # set log_critical_format
        self.log_critical_format = parameters.get(
            Configurable.LOG_CRITICAL_FORMAT, self.log_critical_format)

        from cconfiguration.manager import ConfigurationManager

        # set managers
        managers = parameters.get(Configurable.MANAGERS)
        if managers is not None:
            self.managers = list()
            managers = managers.split(',')
            for manager in managers:
                manager = ConfigurationManager.add_manager(manager)
                self.managers.add(manager)

    @staticmethod
    def _get_manager(
        conf_file, logger, managers
    ):
        """
        Get the first manager able to handle input conf_file.
        None if no manager is able to handle input conf_file.

        :return: first ConfigurationManager able to handle conf_file.
        :rtype: ConfigurationManager
        """

        result = None, None

        for manager in managers:
            manager = manager()

            handle = manager.handle(conf_file=conf_file, logger=logger)

            if handle:
                result = manager
                break

        return result
