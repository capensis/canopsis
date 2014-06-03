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
from utils import resolve_element


class MetaConfigurationManager(type):
    """
    ConfigurationManager meta class which register all manager in a global
    set of managers.
    """
    def __init__(self, name, bases, attrs):

        super(MetaConfigurationManager, self).__init__(name, bases, attrs)

        # if the class claims to be registered
        if self.__register__:
            # add it among managers
            ConfigurationManager._MANAGERS.add(self)

from stat import ST_SIZE
from os import stat, isfile


class ConfigurationManager(object):
    """
    Base class for managing configuration.
    """

    """
    Apply meta class for registering automatically it among global managers
    if __register__ is True
    """
    __metaclass__ = MetaConfigurationManager

    """
    Static parameter which allows this class to be automatically registered
    among managers.
    """
    __register__ = False

    CONFIGURATION_FILE = 'CONFIGURATION_FILE'

    _MANAGERS = set()

    @staticmethod
    def get_managers():
        """
        Get global defined managers.
        """

        return tuple(ConfigurationManager._MANAGERS)

    @staticmethod
    def add_manager(path):
        """
        Add a configuration manager by its path definition.

        :param path: manager path to add. Must be a full path from a known
            package/module
        :type path: str
        """

        result = resolve_element(path)

        return result

    def _has_category(
        self, config_resource, category, logger, *args, **kwargs
    ):
        raise NotImplementedError()

    def _has_parameter(
        self,
        config_resource, category, parameter_name, logger,
        *args, **kwargs
    ):
        raise NotImplementedError()

    def _get_config_resource(
        self, logger, configuration_file=None, *args, **kwargs
    ):
        """
        Get config resource.

        :param configuration_file: if not None, the config resource is \
            configuration_file content.
        :type configuration_file: str

        :param logger: logger used to log processing information
        :type logger: logging.Logger

        :return: empty config resource if configuration_file is None, else \
            configuration_file content.
        """

        raise NotImplementedError()

    def _get_parameter(
        self, config_resource, category, parameter_name, *args, **kwargs
    ):
        raise NotImplementedError()

    def handle(self, configuration_file, logger, *args, **kwargs):
        """
        True iif input configuration_file can be handled by self.

        :return: True iif input configuration_file can be handled by self.
        :rtype: bool
        """

        config_resource = self._get_config_resource(
            configuration_file=configuration_file, logger=logger)

        result = config_resource is not None

        return result

    def get_parameters(
        self, configuration_file, parsing_rules, logger, *args, **kwargs
    ):
        """
        Parse a configuration_files with input parsing_rules and returns
        parameters and errors by parameter name.

        Args:
            - configuration_files list(str) or str:
                path to configuration_files.
            - parsing_rules list(dict(str: dict(str: callable))):
                a list of dictionaries of parsers  by parameter name and
                category.
            - logger (logging.Logger): logger to use instead of raising
                exceptions.

        Returns:
            tuple of 2 dictionaries containing respectively parameter name with
            - value.
            - parsing error.
        """

        # initialize return values
        parameters = dict()
        error_parameters = dict()

        config_resource = None

        # ensure configuration_file exists and is not empty.
        if isfile(configuration_file) and stat(configuration_file)[ST_SIZE]:

            try:
                # first, read configuration file
                config_resource = self._get_config_resource(
                    configuration_file=configuration_file, logger=logger)

            except Exception as e:
                # if an error occured, add it in error_parameters at
                # ConfigurationLanguage
                config_resource_error = error_parameters.setdefault(
                    ConfigurationManager.CONFIGURATION_FILE, list())
                config_resource_error.append(e)
                logger.error(
                    'Impossible to parse configuration_file {0}: {1}'.format(
                        configuration_file, e))

            else:  # else process configuration file

                # get generic logging message
                log_message = '{0}/{1}'.format(configuration_file, '{0}/{1}')

                # for each parsing rule in the ascending order
                for parsing_rule in parsing_rules:

                    # iterate on all category
                    for category, parsers_by_parameter in \
                            parsing_rule.iteritems():

                        # if parsing_rule category exists in
                        # configuration_files
                        if self._has_category(
                            config_resource=config_resource, category=category,
                                logger=logger):

                            # iterate on all parameter_name
                            for name, parser in \
                                    parsers_by_parameter.iteritems():

                                # if parameter_name exists
                                if self._has_parameter(
                                    config_resource=config_resource,
                                    category=category,
                                    parameter_name=name,
                                        logger=logger):

                                    # construct generic log message for each
                                    # name
                                    option_log_message = '{0} = {1}'.format(
                                        log_message.format(
                                            category, name),
                                        '{0}')

                                    # get sub_category_value
                                    sub_category_value = self._get_parameter(
                                        config_resource=config_resource,
                                        category=category,
                                        parameter_name=name,
                                        logger=logger)

                                    try:  # parse parameter_name
                                        value = parser(sub_category_value)

                                    # if an exception occured
                                    except Exception as e:
                                        # set error among errors result
                                        error_parameters[name] = e
                                        # remove value from parameters result
                                        parameters.pop(name, None)
                                        error_message = \
                                            option_log_message.format(
                                                e)
                                        logger.error(error_message)

                                    else:  # if parsing is ok
                                        # set value for parameters result
                                        parameters[name] = value
                                        # remove exception from errors result
                                        error_parameters.pop(name, None)
                                        info_message = \
                                            option_log_message.format(
                                                value)
                                        logger.info(info_message)

        # set the result with a tuple of parameters and error_parameters
        result = parameters, error_parameters

        return result

    def _set_category(
        self, config_resource, category, logger, *args, **kwargs
    ):
        raise NotImplementedError()

    def _set_parameter(
        self, config_resource, category, parameter_name, parameter, logger,
        *args, **kwargs
    ):
        raise NotImplementedError()

    def _write_config_resource(
        self, config_resource, configuration_file, *args, **kwargs
    ):
        raise NotImplementedError()

    def set_parameters(
        self, configuration_file, parameter_by_categories, logger,
        *args, **kwargs
    ):
        """
        Args:
            - configuration_files (str):
            - parameter_by_categories (dict(str: dict(str: object)):
            - logger (logging.Logger):
        """

        result = None
        config_resource = None

        try:  # get config_resource
            config_resource = self._get_config_resource(
                configuration_file=configuration_file,
                logger=logger)

        except Exception as e:
            # if an error occured, stop processing
            logger.error(
                'Impossible to parse configuration_file {0}: {1}'.format(
                    configuration_file, e))
            result = e

        # if configuration_file can not be loaded, get default config resource
        if config_resource is None:

            config_resource = self._get_config_resource(logger=logger)

        # iterate on all parameter_by_categories items
        for category, parameters in parameter_by_categories.iteritems():

            # set category
            self._set_category(
                config_resource=config_resource, category=category,
                logger=logger)

            # iterate on parameters
            for parameter_name, parameter_value in parameters.iteritems():

                # set parameter
                self._set_parameter(
                    config_resource=config_resource,
                    category=category,
                    parameter_name=parameter_name,
                    parameter=parameter_value,
                    logger=logger)

        # write conf_resource in configuration file
        self._write_config_resource(
            config_resource=config_resource,
            configuration_file=configuration_file)

        return result

"""
Load automatically all library managers.
"""
from cconfiguration.manager import *
