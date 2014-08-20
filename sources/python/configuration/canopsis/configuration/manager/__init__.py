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

from canopsis.common.utils import resolve_element, path

from canopsis.configuration.parameters import \
    Configuration, Parameter, Category


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
            ConfigurationManager._MANAGERS[path(self)] = self


class ConfigurationManager(object):
    """
    Base class for managing conf.
    """

    """
    Apply meta class for registering automatically it among global managers
    if __register__ is True
    """
    __metaclass__ = MetaConfigurationManager

    """
    Static param which allows this class to be automatically registered
    among managers.
    """
    __register__ = False

    CONF_FILE = 'CONF_FILE'

    """
    Private set of shared manager types.
    """
    _MANAGERS = {}

    def handle(self, conf_path, logger):
        """
        True iif input conf_path can be handled by self.

        :return: True iif input conf_path can be handled by self.
        :rtype: bool
        """

        conf_resource = self._get_conf_resource(
            conf_path=conf_path, logger=logger)

        result = conf_resource is not None

        return result

    def exists(self, conf_path):
        """
        True if conf_path exist related to this manager behaviour.
        """

        raise NotImplementedError()

    def get_configuration(
        self, conf_path, logger, conf=None, fill=False, override=True
    ):
        """
        Parse a configuration_files with input conf and returns
        parameters and errors by param name.

        :param conf_path: conf file to parse and from get parameters
        :type conf_path: str

        :param conf: conf to fill with conf_path values and
            conf param names.
        :type conf: canopsis.configuration.Configuration

        :param logger: logger to use in order to trace information/error
        :type logger: logging.Logger

        :param fill: result is all conf_path content
        :type fill: bool

        :param override: if True (by default), override self configuration
        :type override: bool
        """

        conf_resource = None

        result = None

        # ensure conf_path exists and is not empty.
        if self.exists(conf_path):

            try:
                # first, read conf file
                conf_resource = self._get_conf_resource(
                    conf_path=conf_path, logger=logger)

            except Exception as e:
                # if an error occured, log it
                logger.error(
                    'Impossible to parse conf_path %s with %s: %s' % (
                        conf_path, type(self), e))

            else:  # else process conf file

                if conf_resource is None:
                    return result

                result = Configuration() if conf is None \
                    else conf

                # get generic logging message
                log_message = '{0}/{{0}}/{{1}}'.format(conf_path)

                if fill:

                    for category in self._get_categories(
                            conf_resource=conf_resource, logger=logger):

                        category = result.setdefault(
                            category, Category(category))

                        for name in self._get_parameters(
                                conf_resource=conf_resource,
                                category=category,
                                logger=logger):

                            param = category.setdefault(name, Parameter(name))

                            value = self._get_value(
                                conf_resource=conf_resource,
                                category=category,
                                param=param,
                                logger=logger)

                            if value not in (None, ''):
                                if override or param.value in (None, ''):
                                    param.value = value

                else:

                    # for each parsing rule in the ascending order
                    for category in result:

                        if self._has_category(
                            conf_resource=conf_resource,
                            category=category,
                                logger=logger):

                            for param in category:

                                name = param.name

                                # if parameter_name exists
                                if self._has_parameter(
                                    conf_resource=conf_resource,
                                    category=category,
                                    param=param,
                                    logger=logger
                                ):

                                    # construct generic log message for each
                                    # name
                                    option_log_message = '{0} = {{0}}'.format(
                                        log_message.format(
                                            category.name, name))

                                    # get sub_category_value
                                    value = self._get_value(
                                        conf_resource=conf_resource,
                                        category=category,
                                        param=param,
                                        logger=logger)

                                    if override \
                                            or param.value not in (None, ''):
                                        param.value = value

                                    # if an exception occured
                                    if isinstance(param.value, Exception):
                                        # set error among errors result
                                        error_message = \
                                            option_log_message.format(
                                                param.value)
                                        logger.error(error_message)

                                    info_message = option_log_message.format(
                                        param.value)
                                    logger.info(info_message)

        return result

    def set_configuration(self, conf_path, conf, logger):
        """
        Set input conf in input conf_path.

        :param conf_path:
        :type conf_path: str

        :param conf: conf to write to conf_path.
        :type conf: canopsis.configuration.Configuration

        :param logger: used to log info/errors
        :type logger: logging.Logger
        """

        result = None
        conf_resource = None

        try:  # get conf_resource
            conf_resource = self._get_conf_resource(
                conf_path=conf_path,
                logger=logger)

        except Exception as e:
            # if an error occured, stop processing
            logger.error(
                'Impossible to parse conf_path {0}: {1}'.format(
                    conf_path, e))
            result = e

        # if conf_path can not be loaded, get default config conf_resource
        if conf_resource is None:

            conf_resource = self._get_conf_resource(logger=logger)

        # iterate on all conf items
        for category in conf:

            # set category
            self._set_category(
                conf_resource=conf_resource, category=category,
                logger=logger)

            # iterate on parameters
            for param in category:

                # set param
                self._set_parameter(
                    conf_resource=conf_resource,
                    category=category,
                    param=param,
                    logger=logger)

        # write conf_resource in conf file
        self._update_conf_resource(
            conf_resource=conf_resource, conf_path=conf_path)

        return result

    @staticmethod
    def get_managers():
        """
        Get global defined managers.
        """

        return set(ConfigurationManager._MANAGERS.values())

    @staticmethod
    def get_manager(path):
        """
        Add a conf manager by its path definition.

        :param path: manager path to add. Must be a full path from a known
            package/module
        :type path: str
        """

        # try to get if from global definition
        result = ConfigurationManager._MANAGERS.get(path)

        # if not already added
        if result is None:
            # resolve it and add it in global definition
            result = resolve_element(path)
            ConfigurationManager._MANAGERS[path] = result

        return result

    def _get_categories(self, conf_resource, logger):
        """
        Get a list of category names in conf_resource
        """

        raise NotImplementedError()

    def _get_parameters(self, conf_resource, category, logger):
        """
        Get a list of param names in conf_resource related to category
        """

        raise NotImplementedError()

    def _has_category(self, conf_resource, category, logger):
        """
        True iif input conf_resource contains input category.
        """

        raise NotImplementedError()

    def _has_parameter(self, conf_resource, category, param, logger):
        """
        True iif input conf_resource has input parameter_name in input category
        """

        raise NotImplementedError()

    def _get_conf_resource(self, logger, conf_path=None):
        """
        Get config conf_resource.

        :param conf_path: if not None, the config conf_resource is \
            conf_path content.
        :type conf_path: str

        :param logger: logger used to log processing information
        :type logger: logging.Logger

        :return: empty config conf_resource if conf_path is None, else \
            conf_path content.
        """

        raise NotImplementedError()

    def _get_value(self, conf_resource, category, param):
        """
        Get a param related to input conf_resource, category and param
        """

        raise NotImplementedError()

    def _set_category(self, conf_resource, category, logger):
        """
        Set category on conf_resource.
        """

        raise NotImplementedError()

    def _set_parameter(self, conf_resource, category, param, logger):
        """
        Set param on conf_resource.
        """

        raise NotImplementedError()

    def _update_conf_resource(self, conf_resource, conf_path):
        """
        Write conf_resource into conf_path.
        """

        raise NotImplementedError()
