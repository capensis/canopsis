# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from logging import Formatter, getLogger, FileHandler, Filter

from os.path import join, sep
from sys import prefix as sys_prefix

from inspect import isclass

from canopsis.common.init import basestring
from canopsis.configuration.parameters \
    import Configuration, Category, Parameter

from canopsis.configuration.driver import ConfigurationDriver


class MetaConfigurable(type):
    """
    Meta class for Configurable.
    """

    def __call__(cls, *args, **kwargs):
        """
        Get a new instance of input cls class, and if instance.auto_conf or
        instance.reconf_once, then call instance.apply_configuration().
        """

        result = type.__call__(cls, *args, **kwargs)

        if result.auto_conf or result.reconf_once:
            # get configuration
            conf = result.conf
            # add a last category which contains args and kwargs as parameters
            if kwargs:
                init_category = Category(Configurable.INIT_CAT)
                for name in kwargs:
                    param = Parameter(name=name, value=kwargs[name])
                    init_category += param
                conf += init_category
            # apply configuration
            result.apply_configuration(conf=conf)

        return result


class ConfigurableError(Exception):
    """
    Handle Configurable errors
    """


class Configurable(object):
    """
    Manages class conf synchronisation with conf resources.
    """

    class Error(Exception):
        """Handle Configurable errors.
        """

    __metaclass__ = MetaConfigurable

    DEFAULT_DRIVERS = '{0},{1}'.format(
        'canopsis.configuration.driver.file.json.JSONConfigurationDriver',
        'canopsis.configuration.driver.file.ini.INIConfigurationDriver')

    INIT_CAT = 'init_cat'  #: initialization category

    CONF_PATH = 'configuration/configurable.conf'

    CONF = 'CONFIGURATION'
    LOG = 'LOG'

    AUTO_CONF = 'auto_conf'
    RECONF_ONCE = 'reconf_once'
    CONF_PATHS = 'conf_paths'
    DRIVERS = 'drivers'

    LOG_NAME = 'log_name'
    LOG_LVL = 'log_lvl'
    LOG_DEBUG_FORMAT = 'log_debug_format'
    LOG_INFO_FORMAT = 'log_info_format'
    LOG_WARNING_FORMAT = 'log_warning_format'
    LOG_ERROR_FORMAT = 'log_error_format'
    LOG_CRITICAL_FORMAT = 'log_critical_format'

    DEBUG_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] \
[%(process)d] [%(thread)d] [%(pathname)s] [%(lineno)d] %(message)s"
    INFO_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s"
    WARNING_FORMAT = INFO_FORMAT
    ERROR_FORMAT = WARNING_FORMAT
    CRITICAL_FORMAT = ERROR_FORMAT

    def __init__(
        self,
        unified_category=None,
        to_configure=None,
        conf_paths=None, drivers=DEFAULT_DRIVERS,
        auto_conf=True, reconf_once=False,
        log_lvl='INFO', log_name=None, log_info_format=INFO_FORMAT,
        log_debug_format=DEBUG_FORMAT, log_warning_format=WARNING_FORMAT,
        log_error_format=ERROR_FORMAT, log_critical_format=CRITICAL_FORMAT
    ):
        """
        :param str unified_category: if not None, used such as a unified
            category

        :param to_configure: object to reconfigure. Such object may implements
            the methods configure apply_configuration and configure
        :type to_configure: class instance

        :param conf_paths: conf_paths to parse
        :type conf_paths: Iterable or str

        :param bool auto_conf: true force auto conf as soon as param change

        :param bool reconf_once: true force auto conf reconf_once as soon as
            param change

        :param str log_lvl: logging level
        """

        super(Configurable, self).__init__()

        self.unified_category = unified_category

        self._to_configure = self if to_configure is None else to_configure

        self.auto_conf = auto_conf
        self.reconf_once = reconf_once

        # set conf files
        self._init_conf_paths(conf_paths)

        # set drivers
        self.drivers = drivers

        # set logging properties
        self._log_lvl = log_lvl
        self._log_name = log_name if log_name is not None else \
            type(self).__name__.lower()
        self._log_debug_format = log_debug_format
        self._log_info_format = log_info_format
        self._log_warning_format = log_warning_format
        self._log_error_format = log_error_format
        self._log_critical_format = log_critical_format

        self._logger = self.newLogger()

    @property
    def drivers(self):
        return self._drivers

    @drivers.setter
    def drivers(self, value):
        self._drivers = value

    def newLogger(self):
        """
        Get a new logger related to self properties.
        """

        result = getLogger(self.log_name)
        result.setLevel(self.log_lvl)

        def setHandler(logger, lvl, path, _format):
            """
            Set right handler related to input lvl, path and format.

            :param Logger logger: logger on which add an handler.
            :param str lvl: logging level.
            :param str path: file path.
            :param str _format: logging message format.
            """

            class _Filter(Filter):
                """
                Ensure message will be given for specific lvl
                """
                def filter(self, record):
                    return record.levelname == lvl

            # get the rights formatter and filter to set on a file handler
            handler = FileHandler(path)
            handler.addFilter(_Filter())
            handler.setLevel(lvl)
            formatter = Formatter(_format)
            handler.setFormatter(formatter)

            # if an old handler exist, remove it from logger
            if hasattr(logger, lvl):
                old_handler = getattr(logger, lvl)
                logger.removeHandler(old_handler)

            logger.addHandler(handler)
            setattr(logger, lvl, handler)

        filename = self.log_name.replace('.', sep)
        path = join(sys_prefix, 'var', 'log', '{0}.log'.format(filename))

        setHandler(result, 'DEBUG', path, self.log_debug_format)
        setHandler(result, 'INFO', path, self.log_info_format)
        setHandler(result, 'WARNING', path, self.log_warning_format)
        setHandler(result, 'ERROR', path, self.log_error_format)
        setHandler(result, 'CRITICAL', path, self.log_critical_format)

        return result

    @property
    def conf(self):
        """
        Get conf with parsers and self property values
        """
        result = self._conf()

        # add a last unified category if asked by self
        if self.unified_category is not None:
            result.add_unified_category(self.unified_category)

        return result

    def _conf(self):
        """
        Protected method to override in order to specify which conf
        return with parsers and default values
        """

        result = Configuration(
            Category(
                Configurable.CONF,
                Parameter(Configurable.AUTO_CONF, Parameter.bool),
                Parameter(Configurable.DRIVERS),
                Parameter(Configurable.RECONF_ONCE, Parameter.bool),
                Parameter(Configurable.CONF_PATHS, Parameter.array)
            ),
            Category(
                Configurable.LOG,
                Parameter(Configurable.LOG_NAME, critical=True),
                Parameter(Configurable.LOG_LVL, critical=True),
                Parameter(Configurable.LOG_DEBUG_FORMAT, critical=True),
                Parameter(Configurable.LOG_INFO_FORMAT, critical=True),
                Parameter(Configurable.LOG_WARNING_FORMAT, critical=True),
                Parameter(Configurable.LOG_ERROR_FORMAT, critical=True),
                Parameter(Configurable.LOG_CRITICAL_FORMAT, critical=True)
            )
        )

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

        :param str value: new log_lvl to set up.
        """

        self._log_lvl = value
        self._logger.setLevel(self._log_lvl)

    @property
    def logger(self):

        return self._logger

    @logger.setter
    def logger(self, value):

        self._logger = value

    @property
    def conf_paths(self):
        """
        Get all type conf files and user files.

        :return: self conf files
        :rtype: tuple
        """

        if not hasattr(self, '_conf_paths'):
            self._conf_paths = []

        result = self._conf_paths

        return result

    @conf_paths.setter
    def conf_paths(self, value):
        """
        Change of conf_paths in adding it in watching list.

        .. TODO:: add watchers here
        """

        #from canopsis.configuration.watcher import add_configurable,\
        #    remove_configurable

        # remove previous watching
        #remove_configurable(self)
        self._conf_paths = tuple(value)
        # add new watching
        #add_configurable(self)

    @property
    def auto_conf(self):
        return self._auto_conf

    @auto_conf.setter
    def auto_conf(self, value):
        self._auto_conf = value

    @property
    def reconf_once(self):
        return self._reconf_once

    @reconf_once.setter
    def reconf_once(self, value):
        self._reconf_once = value

    def apply_configuration(
        self, conf=None, conf_paths=None, drivers=None, logger=None,
        override=True, to_configure=None
    ):
        """
        Apply conf on a destination in 5 phases:

        1. identify the right driver to use with conf_paths to parse.
        2. for all conf_paths, get conf which match
            with input conf.
        3. apply parsing rules on conf_path params.
        4. put values and parsing errors in two different dictionaries.
        5. returns both dictionaries of param values and errors.

        :param Configuration conf: conf from where get conf

        :param conf_paths: conf files to parse. If
            conf_paths is a str, it is automatically putted into a list
        :type conf_paths: list of str

        :param bool override: if True (by default), override self configuration

        :param to_configure: object to configure. self by default.
        """

        conf = self.get_configuration(
            conf=conf, conf_paths=conf_paths, logger=logger,
            drivers=drivers, override=override
        )

        self.configure(conf=conf, to_configure=to_configure)

    def get_configuration(
        self,
        conf=None, conf_paths=None, drivers=None, logger=None,
        override=True
    ):
        """
        Get a dictionary of params by name from conf,
        conf_paths and drivers

        :param Configuration conf: conf to update. If None, use self.conf

        :param conf_paths: list of conf files. If None, use self.conf_paths
        :type conf_paths: list of str

        :param Logger logger: logger to use for logging info/error messages.
            If None, use self.logger

        :param drivers: conf drivers to use. If None, use self.drivers
        :type drivers: list of ConfigurationDriver

        :param bool override: if True (by default), override self configuration

        :param to_configure: object to configure. self by default.
        """

        # start to initialize input params
        if logger is None:
            logger = self.logger

        if conf is None:
            conf = self.conf

        if conf_paths is None:
            conf_paths = self._conf_paths

        if isinstance(conf_paths, basestring):
            conf_paths = [conf_paths]

        if drivers is None:
            drivers = self.drivers

        # iterate on all conf_paths
        for conf_path in conf_paths:

            conf_driver = self._get_driver(
                conf_path=conf_path, logger=logger, drivers=drivers
            )

            # if a config_resource is not None
            if conf_driver is not None:

                conf = conf_driver.get_configuration(
                    conf=conf, logger=logger,
                    conf_path=conf_path, override=override
                )

            else:
                # if no conf_driver, display a warning log message
                self.logger.warning(
                    'No driver found among {0} for processing {1}'.format(
                        drivers, conf_path
                    )
                )

        return conf

    def set_configuration(self, conf_path, conf, driver=None, logger=None):
        """
        Set params on input conf_path.

        :param str conf_paths: conf_path to udate with params
        :param dict parameter_by_categories: (dict(str: dict(str: object))
        :param Logger logger: logger to use to set params.
        """

        result = None

        if logger is None:
            logger = self.logger

        # first get content of input conf_path
        prev_driver = self._get_driver(
            conf_path=conf_path,
            logger=logger,
            drivers=self.drivers
        )

        if prev_driver is not None:
            prev_conf = prev_driver.get_configuration(
                conf_path=conf_path, logger=logger
            )

        # try to find a good driver if driver is None
        if driver is None:
            driver = self._get_driver(
                conf_path=conf_path,
                logger=logger,
                drivers=self.drivers
            )

        elif isclass(driver):
            driver = driver()

        else:
            driver = self._get_driver(
                conf_path=None,
                logger=logger,
                drivers=driver)

        # if prev driver is not the new driver
        if prev_conf is not None \
                and type(driver) is not type(prev_driver):
            # update prev_conf with input conf
            prev_conf.update(conf)
            conf = prev_conf

        if driver is not None:
            driver.set_configuration(
                conf_path=conf_path,
                conf=conf,
                logger=logger
            )

        else:
            self.logger.error(
                'No ConfigurationDriver found for conf resource {0}'.format(
                    conf_path
                )
            )

        return result

    def configure(self, conf, logger=None, to_configure=None):
        """
        Update self properties with input params only if:
        - self.configure is True
        - self.auto_conf is True
        - param conf 'configure' is True
        - param conf 'auto_conf' is True

        This method may not be overriden. see _configure instead

        :param Configuration conf: object from where get paramters
        :param to_configure: object to configure. self if equals None.
        """

        if logger is None:
            logger = self.logger

        unified_conf = conf.unify()

        values = unified_conf[Configuration.VALUES]

        # set configure
        reconf_once = values.get(Configurable.RECONF_ONCE)

        if reconf_once is not None:
            self.reconf_once = reconf_once.value

        # set auto_conf
        auto_conf_parameter = values.get(Configurable.AUTO_CONF)

        if auto_conf_parameter is not None:
            self.auto_conf = auto_conf_parameter.value

        if self.reconf_once or self.auto_conf:
            self._configure(unified_conf=unified_conf, logger=logger)
            # when conf succeed, deactive reconf_once
            self.reconf_once = False

    def _is_local(self, to_configure, name):
        """
        True iif input name parameter can be handled by to_configure.
        """

        return hasattr(to_configure, name)

    def _configure(self, unified_conf, logger=None, to_configure=None):
        """
        Configure this class with input conf only if auto_conf or
        configure is true.

        This method should be overriden for specific conf

        :param Configuration unified_conf: Configuration with two categories
            VALUES and ERRORS

        :param bool configure: if True, force full self conf
        :param to_configure: object to configure. self if equals None.
        """

        if to_configure is None:
            to_configure = self

        values = [p for p in unified_conf[Configuration.VALUES]]
        foreigns = [p for p in unified_conf[Configuration.FOREIGNS]]

        criticals = []  # list of critical parameters

        for parameter in values + foreigns:
            name = parameter.name
            if not parameter.asitem:
                # if parameter is local, to_configure must a related name
                if self._is_local(to_configure, name):
                    if hasattr(to_configure, name):
                        param_value = parameter.value

                        # in case of a critical parameter
                        if parameter.critical:
                            # check if current value is not the same as new val
                            value = getattr(to_configure, name)

                            if value != parameter.value:
                                # add it to list of criticals
                                criticals.append(parameter)
                                # set private name
                                privname = '_%s' % name

                                # if private name exists
                                if hasattr(to_configure, privname):
                                    # change of value of private name
                                    setattr(
                                        to_configure, privname, param_value)

                                else:  # update public value
                                    setattr(to_configure, name, param_value)

                        else:  # change public value
                            setattr(to_configure, name, param_value)

                else:  # else log the warning
                    message = 'Parameter {} is not bound to an attribute in {}'
                    self.logger.warning(message.format(name, to_configure))

            else:
                value = getattr(to_configure, parameter.asitem.name)
                param_value = parameter.value
                param_name = parameter.name

                if param_name not in value or param_value != value[param_name]:
                    criticals.append(parameter)

                value[parameter.name] = parameter.value

        # if criticals
        if criticals:
            self.restart(to_configure=to_configure, criticals=criticals)

    def restart(self, criticals, to_configure=None):
        """
        Restart a configurable object with critical parameters.

        :param to_configure: object to configure with critical parameters
        :param tuple criticals: list of critical parameters
        """

        if to_configure is None:
            to_configure = self

        if self._is_critical_category(Configurable.LOG, criticals):
            self._logger = self.newLogger()
            to_configure.logger = self.logger

    def _is_critical_category(self, category, criticals):
        """
        Check if input category parameters are among criticals
        """

        result = False

        properties = (param.name for param in self.conf[category])

        # if property is among criticals
        for prop in properties:
            for critical in criticals:
                if critical.name == prop:
                    result = True
                    break

        return result

    def _init_conf_paths(self, conf_paths):

        self.conf_paths = self._get_conf_paths() \
            if conf_paths is None else conf_paths

    def _get_conf_paths(self):

        result = [Configurable.CONF_PATH]

        return result

    @staticmethod
    def _get_driver(conf_path, drivers, logger):
        """
        Get the first driver able to handle input conf_path.
        None if no driver is able to handle input conf_path.

        :return: first ConfigurationDriver able to handle conf_path.
        :rtype: ConfigurationDriver
        """

        result = None

        for driver in drivers.split(','):
            driver = ConfigurationDriver.get_driver(driver)
            driver = driver()

            handle = conf_path is None \
                or driver.handle(conf_path=conf_path, logger=logger)

            if handle:
                result = driver
                break

        return result
