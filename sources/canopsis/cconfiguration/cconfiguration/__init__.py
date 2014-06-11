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

from os.path import expanduser

import logging

from cconfiguration.manager import ConfigurationManager
from cconfiguration.watcher import add_configurable, remove_configurable


class Configurable(object):
    """
    Manages class configuration synchronisation with configuration files.
    """

    DEFAULT_CONFIGURATION_FILE = '~/etc/conf.conf'

    CONF = 'CONF'

    AUTO_CONF = 'auto_conf'
    MANAGERS = 'conf_managers'
    LOGGING_LEVEL = 'logging_level'

    PARSING_RULES = {
        CONF: {
            AUTO_CONF: bool,
            LOGGING_LEVEL: str,
            MANAGERS: str
        }
    }

    def __init__(
        self,
        configuration_files=None, auto_conf=True, logging_level=logging.INFO,
        managers=None, parsing_rules=None, _ready_to_conf=True, *args, **kwargs
    ):
        """
        :param configuration_files: configuration_files to parse.
        :type configuration_files: Iterable or str

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

        if configuration_files is None:
            configuration_files = [Configurable.DEFAULT_CONFIGURATION_FILE]

        self.configuration_files = configuration_files

        self.auto_conf = auto_conf

        self._logger = logging.getLogger(type(self).__name__)
        self._logger.setLevel(logging_level)

        self.parsing_rules = (Configurable.PARSING_RULES,) \
            if parsing_rules is None else parsing_rules

        # set managers
        self.managers = ConfigurationManager.get_managers() \
            if managers is None else managers

        if _ready_to_conf and self.auto_conf:
            self.apply_configuration()

    @property
    def configuration_files(self):
        """
        :return: self configuration files
        :rtype: tuple
        """

        result = tuple(getattr(self, '_configuration_files', list()))

        return result

    @configuration_files.setter
    def configuration_files(self, value):
        """
        Change of configuration_files in adding it in watching list.
        """

        # remove previous watching
        remove_configurable(self)
        self._configuration_files = tuple(value)
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
        self, parsing_rules=None, configuration_files=None,
        managers=None, *args, **kwargs
    ):
        """
        Apply configuration on a destination in 5 phases:

        1. identify the right manager to use with configuration_files to parse.
        2. for all configuration_files, get parameters which match
            with input parsing_rules.
        3. apply parsing rules on configuration_file parameters.
        4. put values and parsing errors in two different dictionaries.
        5. returns both dictionaries of parameter values and errors.

        :param parsing_rules: Iterable of parsing_rule.
        :type parsing_rules: Iterable

        :param configuration_files: configuration files to parse. If
            configuration_files is a str, it is automatically putted into an
            Iterable.
        :type configuration_files: Iterable(str) or str.
        """

        parameters, error_parameters = self.get_parameters(
            parsing_rules=parsing_rules,
            configuration_files=configuration_files,
            managers=managers,
            *args, **kwargs)

        self.configure(
            parameters=parameters, error_parameters=error_parameters,
            *args, **kwargs)

    def get_parameters(
        self,
        parsing_rules=None, configuration_files=None, logger=None,
        managers=None, *args, **kwargs
    ):
        """
        Get a dictionary of parameters by name from parsing_rules,
        configuration_files and conf_managers.

        :param parsing_rules: Iterable of parsing_rule. If None, use
            self.parsing_rules.
        :type parsing_rules:
            Iterable(dict(category, dict(parameter name, parser)))

        :param configuration_files: Iterable of configuration file path.
        :type configuration_files: Iterable(str)
        """

        # start to initialize input parameters
        if logger is None:
            logger = self._logger

        if parsing_rules is None:
            parsing_rules = self.parsing_rules

        if configuration_files is None:
            configuration_files = self._configuration_files

        if isinstance(configuration_files, str):
            configuration_files = [configuration_files]

        # clean configuration file list
        configuration_files = [
            expanduser(configuration_file) for configuration_file
            in configuration_files]

        parameters, error_parameters = dict(), dict()

        if managers is None:
            managers = self.managers

        # iterate on all configuration_files
        for configuration_file in configuration_files:

            config_manager = self._get_manager(
                configuration_file=configuration_file,
                logger=logger, managers=managers)

            # if a config_resource is not None
            if config_manager is not None:

                # get parameters from a good config_manager
                _parameters, _error_parameters = config_manager.get_parameters(
                    parsing_rules=parsing_rules,
                    configuration_file=configuration_file,
                    logger=logger)

                # update parameters and error_parameters
                parameters.update(_parameters)
                error_parameters.update(_error_parameters)

                # clean parameters and error_parameters in order to be
                # consistent with _parameters and _error_parameters
                for name in _error_parameters:
                    if name not in parameters:
                        del parameters[name]

                for name in _parameters:
                    if name not in error_parameters:
                        del error_parameters[name]

            else:
                # if no config_manager, display a warning log message
                logger.warning('No manager found among {0} for {1}'.format(
                    configuration_file))

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
        self, configuration_file, parameter_by_categories, config_manager=None,
        logger=None, *args, **kwargs
    ):
        """
        Set parameters on input configuration_file.

        Args:
            - configuration_files (str): configuration_file to udate with
                parameters
            - parameter_by_categories (dict(str: dict(str: object)):
            - logger (logging.Logger): logger to use to set parameters.
        """

        result = None

        if logger is None:
            logger = self._logger

        # try to find a good conf_manager if conf_manager is None
        if config_manager is None:
            config_manager = self._get_manager(
                configuration_file=configuration_file,
                logger=logger,
                managers=self.managers)

        if config_manager is not None:
            config_manager.set_parameters(
                configuration_file=configuration_file,
                parameter_by_categories=parameter_by_categories,
                logger=logger)

        else:
            logger.error(
                'No ConfigurationManager found for \
                configuration file {0}'.format(
                    configuration_file))

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
                ConfigurationManager.add_manager(manager)

    @staticmethod
    def _get_manager(
        configuration_file, logger, managers
    ):
        """
        Get the first manager able to handle input configuration_file.
        None if no manager is able to handle input configuration_file.

        :return: first ConfigurationManager able to handle configuration_file.
        :rtype: ConfigurationManager
        """

        result = None, None

        for manager in managers:
            manager = manager()

            handle = manager.handle(
                configuration_file=configuration_file, logger=logger)

            if handle:
                result = manager
                break

        return result
