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

from os.path import expanduser, exists

import logging

from cconfiguration.manager import ConfigurationManager
from cconfiguration.watcher import add_configurable, remove_configurable


class Configuration(object):
    """
    Manage configuration thanks to a list of categories which contain
    properties.

    The order of categories permit to ensure parameter overriding.
    """

    def __init__(self, *categories, **kwargs):
        """
        :param categories: categories to configure.
        :type categories: list of Category.
        """

        super(Configuration, self).__init__(**kwargs)
        self.categories = categories

    def get_parameters(self, *args, **kwargs):
        """
        Get values and errors of parameters in respecting parameter overriding.

        :return: two dictionaries respectively with values and errors by name.
        :rtype: dict, dict
        """

        values = dict()
        errors = dict()

        result = values, errors

        for category in self.categories:

            for name, parameter in category.iteritems():

                if parameter.value is not None:

                    to_update, to_delete = errors, values if \
                        isinstance(parameter.value, Exception) \
                        else values, errors

                    to_update[name] = parameter.value

                    if name in to_delete:
                        del to_delete[name]

        return result


class Category(dict):
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
        super(Category, self).__init__(name=name, **kwargs)
        # set parameter by names.
        self.update({parameter.name: parameter for parameter in parameters})


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


class Configurable(object):
    """
    Manages class configuration synchronisation with configuration files.
    """

    DEFAULT_CONF_FILE = '~/etc/conf.conf'

    CONF = 'CONF'

    AUTO_CONF = 'auto_conf'
    MANAGERS = 'conf_managers'
    LOGGING_LEVEL = 'logging_level'

    CONFIGURATION = Configuration(
        Category(
            name=CONF,
            Parameter(name=AUTO_CONF, parser=bool),
            Parameter(name=MANAGERS),
            Parameter(name=LOGGING_LEVEL)))

    PARSING_RULES = {
        CONF: {
            AUTO_CONF: bool,
            LOGGING_LEVEL: str,
            MANAGERS: str
        }
    }

    def __init__(
        self,
        conf_files=None, auto_conf=True, logging_level=logging.INFO,
        managers=None, parsing_rules=None, _ready_to_conf=True, *args, **kwargs
    ):
        """
        :param conf_files: conf_files to parse.
        :type conf_files: Iterable or str

        :param auto_conf: true force auto conf as soon as parameter change
        :type auto_conf: bool

        :param parsing_rules: list of callable parser by parameter name
            and category name
        :type parsing_rules: list of dict of dict of callable

        :param logging_level: logging level
        :type logging_level: str

        :param _ready_to_conf: protected parameter permetting to deactivate
            auto_conf processing in this call.
        :type _ready_to_conf: bool
        """

        super(Configurable, self).__init__(*args, **kwargs)

        if conf_files is None:
            conf_files = [Configurable.DEFAULT_CONF_FILE]

        self.conf_files = conf_files

        self.auto_conf = auto_conf

        self._logger = logging.getLogger(type(self).__name__)
        self._logger.setLevel(logging_level)

        self.parsing_rules = (Configurable.PARSING_RULES,) \
            if parsing_rules is None else parsing_rules

        # set managers
        self.managers = set(ConfigurationManager.get_managers()) \
            if managers is None else set(managers)

        if _ready_to_conf and self.auto_conf:
            self.apply_configuration()

    @property
    def conf_files(self):
        """
        :return: self configuration files
        :rtype: tuple
        """

        result = tuple(getattr(self, '_conf_files', list()))

        return result

    @conf_files.setter
    def conf_files(self, value):
        """
        Change of conf_files in adding it in watching list.
        """

        # remove previous watching
        remove_configurable(self)
        self._conf_files = tuple(value)
        # add new watching
        add_configurable(self)

    @property
    def logging_level(self):
        """
        Get this logger.

        :return: self logger
        :rtype: logging.Logger
        """
        return self._logger.level

    @logging_level.setter
    def logging_level(self, value):
        """
        Change of logging level.

        :param value: new logging_level to set up.
        :type value: str
        """

        self._logger.setLevel(value)

    def apply_configuration(
        self, parsing_rules=None, conf_files=None,
        managers=None, *args, **kwargs
    ):
        """
        Apply configuration on a destination in 5 phases:

        1. identify the right manager to use with conf_files to parse.
        2. for all conf_files, get parameters which match
            with input parsing_rules.
        3. apply parsing rules on conf_file parameters.
        4. put values and parsing errors in two different dictionaries.
        5. returns both dictionaries of parameter values and errors.

        :param parsing_rules: Iterable of parsing_rule.
        :type parsing_rules: Iterable

        :param conf_files: configuration files to parse. If
            conf_files is a str, it is automatically putted into an
            Iterable.
        :type conf_files: Iterable(str) or str.
        """

        parameters, error_parameters = self.get_parameters(
            parsing_rules=parsing_rules,
            conf_files=conf_files,
            managers=managers,
            *args, **kwargs)

        self.configure(
            parameters=parameters, error_parameters=error_parameters,
            *args, **kwargs)

    def get_parameters(
        self,
        parsing_rules=None, conf_files=None, logger=None,
        managers=None, *args, **kwargs
    ):
        """
        Get a dictionary of parameters by name from parsing_rules,
        conf_files and conf_managers.

        :param parsing_rules: Iterable of parsing_rule. If None, use
            self.parsing_rules.
        :type parsing_rules:
            Iterable(dict(category, dict(parameter name, parser)))

        :param conf_files: Iterable of configuration file path.
        :type conf_files: Iterable(str)
        """

        # start to initialize input parameters
        if logger is None:
            logger = self._logger

        if parsing_rules is None:
            parsing_rules = self.parsing_rules

        if conf_files is None:
            conf_files = self._conf_files

        if isinstance(conf_files, str):
            conf_files = [conf_files]

        # clean configuration file list
        conf_files = [
            expanduser(conf_file) for conf_file
            in conf_files]

        parameters, error_parameters = dict(), dict()

        if managers is None:
            managers = self.managers

        # iterate on all conf_files
        for conf_file in conf_files:

            if not exists(conf_file):
                continue

            conf_manager = self._get_manager(
                conf_file=conf_file,
                logger=logger, managers=managers)

            # if a config_resource is not None
            if conf_manager is not None:

                # get parameters from a good conf_manager
                _parameters, _error_parameters = conf_manager.get_parameters(
                    parsing_rules=parsing_rules,
                    conf_file=conf_file,
                    logger=logger)

                # update parameters and error_parameters
                parameters.update(_parameters)
                error_parameters.update(_error_parameters)

                # clean parameters and error_parameters in order to be
                # consistent with _parameters and _error_parameters
                for name in _error_parameters:
                    if name in parameters:
                        del parameters[name]

                for name in _parameters:
                    if name in error_parameters:
                        del error_parameters[name]

            else:
                # if no conf_manager, display a warning log message
                logger.warning('No manager found among {0} for {1}'.format(
                    conf_file))

        result = parameters, error_parameters

        return result

    def get_parameters_by_categories(self):

        parameters = {
            ConfigurationManager.AUTO_CONF: self.auto_conf,
            Configurable.LOGGING_LEVEL: self.logger.level
        }

        result = {ConfigurationManager.CONF: parameters}

        return result

    def set_parameters(
        self, conf_file, parameter_by_categories, conf_manager=None,
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

        # try to find a good conf_manager if conf_manager is None
        if conf_manager is None:
            conf_manager = self._get_manager(
                conf_file=conf_file,
                logger=logger,
                managers=self.managers)

        elif issubclass(conf_manager, ConfigurationManager):
            conf_manager = conf_manager()

        if conf_manager is not None:
            conf_manager.set_parameters(
                conf_file=conf_file,
                parameter_by_categories=parameter_by_categories,
                logger=logger)

        else:
            logger.error(
                'No ConfigurationManager found for \
                configuration file {0}'.format(
                    conf_file))

        return result

    def configure(self, parameters, error_parameters, *args, **kwargs):
        """
        Configure this class with input parameters.
        """

        # set logging_level
        logging_level = parameters.get(Configurable.LOGGING_LEVEL)
        if logging_level is not None:
            self.logging_level = logging_level

        # set auto_conf
        auto_conf = parameters.get(Configurable.AUTO_CONF)
        if auto_conf is not None:
            self.auto_conf = auto_conf

        # set conf_managers
        conf_managers = parameters.get(Configurable.MANAGERS)
        if conf_managers is not None:
            self.conf_managers = list()
            conf_managers = conf_managers.split(',')
            for manager in conf_managers:
                manager = ConfigurationManager.add_manager(manager)
                self.managers.append(manager)

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
