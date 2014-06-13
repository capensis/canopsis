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

from ccommon.utils import resolve_element

from stat import ST_SIZE

from os import stat
from os.path import exists

from cconfiguration import Configuration, Parameter, Category


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

    CONF_FILE = 'CONF_FILE'

    """
    Private set of shared manager types.
    """
    _MANAGERS = set()

    def handle(self, conf_file, logger, *args, **kwargs):
        """
        True iif input conf_file can be handled by self.

        :return: True iif input conf_file can be handled by self.
        :rtype: bool
        """

        conf_resource = self._get_conf_resource(
            conf_file=conf_file, logger=logger)

        result = conf_resource is not None

        return result

    def get_configuration(
        self, conf_file, logger, configuration=None, fill=False,
        *args, **kwargs
    ):
        """
        Parse a configuration_files with input configuration and returns
        parameters and errors by parameter name.

        :param conf_file: configuration file to parse and from get parameters
        :type conf_file: str

        :param configuration: configuration to fill with conf_file values and
            configuration parameter names.
        :type configuration: cconfiguration.Configuration

        :param logger: logger to use in order to trace information/error
        :type logger: logging.Logger

        :param fill: result is all conf_file content
        :type fill: bool
        """

        conf_resource = None

        result = None

        # ensure conf_file exists and is not empty.
        if exists(conf_file) and stat(conf_file)[ST_SIZE]:

            try:
                # first, read configuration file
                conf_resource = self._get_conf_resource(
                    conf_file=conf_file, logger=logger)

            except Exception as e:
                # if an error occured, log it
                logger.error(
                    'Impossible to parse conf_file {0}: {1}'.format(
                        conf_file, e))

            else:  # else process configuration file

                if conf_resource is None:
                    return result

                result = Configuration() if configuration is None \
                    else configuration

                # get generic logging message
                log_message = '{0}/{{0}}/{{1}}'.format(conf_file)

                if fill:

                    for category in self._get_categories(
                            conf_resource=conf_resource, logger=logger):

                        category = result.setdefault(
                            category, Category(category))

                        for parameter in self._get_parameters(
                                conf_resource=conf_resource,
                                category=category,
                                logger=logger):

                            parameter = category.setdefault(
                                parameter, Parameter(parameter))

                            value = self._get_value(
                                conf_resource=conf_resource,
                                category=category,
                                parameter=parameter,
                                logger=logger)

                            parsed_value = parameter.parse(
                                value, logger=logger)

                            parameter.value = parsed_value

                else:

                    # for each parsing rule in the ascending order
                    for category in result:

                        if self._has_category(
                            conf_resource=conf_resource,
                            category=category,
                                logger=logger):

                            for parameter in category:

                                name = parameter.name

                                # if parameter_name exists
                                if self._has_parameter(
                                    conf_resource=conf_resource,
                                    category=category,
                                    parameter=parameter,
                                        logger=logger):

                                    # construct generic log message for each
                                    # name
                                    option_log_message = '{0} = {{0}}'.format(
                                        log_message.format(
                                            category.name, name))

                                    # get sub_category_value
                                    value = self._get_value(
                                        conf_resource=conf_resource,
                                        category=category,
                                        parameter=parameter,
                                        logger=logger)

                                    parsed_value = parameter.parse(
                                        value, logger)

                                    # if an exception occured
                                    if isinstance(parsed_value, Exception):
                                        # set error among errors result
                                        error_message = \
                                            option_log_message.format(
                                                parsed_value)
                                        logger.error(error_message)

                                    # set value on parameter
                                    parameter.value = parsed_value
                                    info_message = option_log_message.format(
                                        parsed_value)
                                    logger.info(info_message)

        return result

    def set_configuration(
        self, conf_file, configuration, logger,
        *args, **kwargs
    ):
        """
        Set input configuration in input conf_file.

        :param conf_file:
        :type conf_file: str

        :param configuration: configuration to write in conf_file.
        :type configuration: cconfiguration.Configuration

        :param logger: used to log info/errors
        :type logger: logging.Logger
        """

        result = None
        conf_resource = None

        try:  # get conf_resource
            conf_resource = self._get_conf_resource(
                conf_file=conf_file,
                logger=logger)

        except Exception as e:
            # if an error occured, stop processing
            logger.error(
                'Impossible to parse conf_file {0}: {1}'.format(
                    conf_file, e))
            result = e

        # if conf_file can not be loaded, get default config resource
        if conf_resource is None:

            conf_resource = self._get_conf_resource(logger=logger)

        # iterate on all configuration items
        for category in configuration:

            # set category
            self._set_category(
                conf_resource=conf_resource, category=category,
                logger=logger)

            # iterate on parameters
            for parameter in category:

                # set parameter
                self._set_parameter(
                    conf_resource=conf_resource,
                    category=category,
                    parameter=parameter,
                    logger=logger)

        # write conf_resource in configuration file
        self._update_conf_file(
            conf_resource=conf_resource,
            conf_file=conf_file)

        return result

    @staticmethod
    def get_managers():
        """
        Get global defined managers.
        """

        return set(ConfigurationManager._MANAGERS)

    @staticmethod
    def add_manager(path):
        """
        Add a configuration manager by its path definition.

        :param path: manager path to add. Must be a full path from a known
            package/module
        :type path: str
        """

        result = resolve_element(path)

        # add it to _MANAGERS
        ConfigurationManager._MANAGERS.add(result)

        return result

    def _get_categories(self, conf_resource, logger, *args, **kwargs):
        """
        Get a list of category names in conf_resource
        """

        raise NotImplementedError()

    def _get_parameters(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        """
        Get a list of parameter names in conf_resource related to category
        """

        raise NotImplementedError()

    def _has_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        """
        True iif input conf_resource contains input category.
        """

        raise NotImplementedError()

    def _has_parameter(
        self,
        conf_resource, category, parameter, logger,
        *args, **kwargs
    ):
        """
        True iif input conf_resource has input parameter_name in input category
        """

        raise NotImplementedError()

    def _get_conf_resource(
        self, logger, conf_file=None, *args, **kwargs
    ):
        """
        Get config resource.

        :param conf_file: if not None, the config resource is \
            conf_file content.
        :type conf_file: str

        :param logger: logger used to log processing information
        :type logger: logging.Logger

        :return: empty config resource if conf_file is None, else \
            conf_file content.
        """

        raise NotImplementedError()

    def _get_value(
        self, conf_resource, category, parameter, *args, **kwargs
    ):
        """
        Get a parameter related to input conf_resource, category and parameter
        """

        raise NotImplementedError()

    def _set_category(
        self, conf_resource, category, logger, *args, **kwargs
    ):
        """
        Set category on conf_resource.
        """

        raise NotImplementedError()

    def _set_parameter(
        self, conf_resource, category, parameter, logger,
        *args, **kwargs
    ):
        """
        Set parameter on conf_resource.
        """

        raise NotImplementedError()

    def _update_conf_file(
        self, conf_resource, conf_file, *args, **kwargs
    ):
        """
        Write conf_resource into conf_file.
        """

        raise NotImplementedError()

"""
Load automatically all library managers.
"""
from cconfiguration.manager import *
