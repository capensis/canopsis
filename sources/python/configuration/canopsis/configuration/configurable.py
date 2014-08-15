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

from logging import Formatter, getLogger, FileHandler, Filter

from stat import ST_SIZE

from os import stat
from os.path import expanduser, exists, sep, abspath

from inspect import isclass

from canopsis.configuration.parameters \
    import Configuration, Category, Parameter


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
            result.apply_configuration()

        return result


class ConfigurableError(Exception):
    """
    Handle Configurable errors
    """
    pass


class Configurable(object):
    """
    Manages class conf synchronisation with conf resources.
    """

    __metaclass__ = MetaConfigurable

    DEFAULT_MANAGERS = '%s,%s' % (
        'canopsis.configuration.manager.file.json.JSONConfigurationManager',
        'canopsis.configuration.manager.file.ini.INIConfigurationManager')

    CONF_RESOURCE = 'configuration/configurable.conf'

    CONF = 'CONFIGURATION'
    LOG = 'LOG'

    AUTO_CONF = 'auto_conf'
    RECONF_ONCE = 'reconf_once'
    CONF_PATHS = 'conf_paths'
    MANAGERS = 'conf_managers'

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
        to_configure=None,
        conf_paths=None, managers=DEFAULT_MANAGERS,
        auto_conf=True, reconf_once=False,
        log_lvl='INFO', log_name=None, log_info_format=INFO_FORMAT,
        log_debug_format=DEBUG_FORMAT, log_warning_format=WARNING_FORMAT,
        log_error_format=ERROR_FORMAT, log_critical_format=CRITICAL_FORMAT
    ):
        """
        :param to_configure: object to reconfigure. Such object may implements
            the methods configure apply_configuration and configure
        :type to_configure: class instance

        :param conf_paths: conf_paths to parse
        :type conf_paths: Iterable or str

        :param auto_conf: true force auto conf as soon as param change
        :type auto_conf: bool

        :param reconf_once: true force auto conf reconf_once as soon as param
            change
        :type reconf_once: bool

        :param log_lvl: logging level
        :type log_lvl: str
        """

        super(Configurable, self).__init__()

        self._to_configure = self if to_configure is None else to_configure

        self.auto_conf = auto_conf
        self.reconf_once = reconf_once

        # set conf files
        self._init_conf_files(conf_paths)

        # set managers
        self.managers = managers

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

    @property
    def managers(self):
        return self._managers

    @managers.setter
    def managers(self, value):
        self._managers = value

    def newLogger(self):
        """
        Get a new logger related to self properties.
        """

        result = getLogger(self.log_name)
        result.setLevel(self.log_lvl)

        def setHandler(logger, lvl, path, format):
            """
            Set right handler related to input lvl, path and format
            """

            class filter(Filter):
                """
                Ensure message will be given for specific lvl
                """
                def filter(self, record):
                    return record.levelname == lvl

            # get the rights formatter and filter to set on a file handler
            handler = FileHandler(path)
            handler.addFilter(Filter())
            handler.setLevel(lvl)
            formatter = Formatter(format)
            handler.setFormatter(formatter)
            # if an old handler exist, remove it from logger
            if hasattr(logger, lvl):
                old_handler = getattr(logger, lvl)
                logger.removeHandler(old_handler)
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
    def conf(self):
        """
        Get conf with parsers and self property values
        """
        return self._conf()

    def _conf(self):
        """
        Protected method to override in order to specify which conf
        return with parsers and default values
        """

        result = Configuration(
            Category(Configurable.CONF,
                Parameter(
                    Configurable.AUTO_CONF, self.auto_conf, Parameter.bool),
                Parameter(Configurable.MANAGERS, self.managers),
                Parameter(
                    Configurable.RECONF_ONCE, self.reconf_once,
                    Parameter.bool),
                Parameter(
                    Configurable.CONF_PATHS, self.conf_paths,
                    Parameter.array)),
            Category(Configurable.LOG,
                Parameter(Configurable.LOG_NAME, self.log_name),
                Parameter(Configurable.LOG_LVL, self.log_lvl),
                Parameter(
                    Configurable.LOG_DEBUG_FORMAT, self.log_debug_format),
                Parameter(Configurable.LOG_INFO_FORMAT, self.log_info_format),
                Parameter(
                    Configurable.LOG_WARNING_FORMAT, self.log_warning_format),
                Parameter(
                    Configurable.LOG_ERROR_FORMAT, self.log_warning_format),
                Parameter(
                    Configurable.LOG_CRITICAL_FORMAT,
                    self.log_critical_format)
            ))

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

        self._log_lvl = value
        self._logger.setLevel(self._log_lvl)

    @property
    def logger(self):

        return self._logger

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
        self, conf=None, conf_paths=None, managers=None, logger=None
    ):
        """
        Apply conf on a destination in 5 phases:

        1. identify the right manager to use with conf_paths to parse.
        2. for all conf_paths, get conf which match
            with input conf.
        3. apply parsing rules on conf_path params.
        4. put values and parsing errors in two different dictionaries.
        5. returns both dictionaries of param values and errors.

        :param conf: conf from where get conf
        :type conf: Configuration

        :param conf_paths: conf files to parse. If
            conf_paths is a str, it is automatically putted into a list
        :type conf_paths: list of str
        """

        if logger is None:
            logger = self.logger

        if conf is None:
            conf = self.conf

        conf = self.get_configuration(
            conf=conf, conf_paths=conf_paths, logger=logger,
            managers=managers)

        self.configure(conf=conf)

    def get_configuration(
        self,
        conf=None, conf_paths=None, managers=None, fill=False, logger=None
    ):
        """
        Get a dictionary of params by name from conf,
        conf_paths and conf_managers

        :param conf: conf to update. If None, use \
            self.conf
        :type conf: Configuration

        :param conf_paths: list of conf files. If None, use \
            self.conf_paths
        :type conf_paths: list of str

        :param logger: logger to use for logging info/error messages.
            If None, use self.logger
        :type logger: logging.Logger

        :param managers: conf managers to use. If None, use self.managers
        :type managers: list of ConfigurationManager

        :param fill: if True (False by default) load in conf all \
            conf_paths content
        :type fill: bool
        """

        # start to initialize input params
        if logger is None:
            logger = self.logger

        if conf is None:
            conf = self.conf

        # remove values from conf
        conf.clean()

        if conf_paths is None:
            conf_paths = self._conf_paths

        if isinstance(conf_paths, str):
            conf_paths = [conf_paths]

        # clean conf file list
        conf_paths = [
            abspath(expanduser(conf_path)) for conf_path
            in conf_paths]

        if managers is None:
            managers = self.managers

        # iterate on all conf_paths
        for conf_path in conf_paths:

            if not exists(conf_path) or stat(conf_path)[ST_SIZE] == 0:
                continue

            conf_manager = self._get_manager(
                conf_path=conf_path, logger=logger, managers=managers)

            # if a config_resource is not None
            if conf_manager is not None:

                conf = conf_manager.get_configuration(
                    conf=conf, fill=fill, logger=logger,
                    conf_path=conf_path)

            else:
                # if no conf_manager, display a warning log message
                self.logger.warning(
                    'No manager found among {0} for {1}'.format(
                        conf_path))

        return conf

    def set_configuration(self, conf_path, conf, manager=None, logger=None):
        """
        Set params on input conf_path.

        Args:
            - conf_paths (str): conf_path to udate with
                params
            - parameter_by_categories (dict(str: dict(str: object)):
            - logger (logging.Logger): logger to use to set params.
        """

        result = None

        if logger is None:
            logger = self.logger

        # first get content of input conf_path
        prev_manager = self._get_manager(
            conf_path=conf_path,
            logger=logger,
            managers=self.managers)

        if prev_manager is not None:
            prev_conf = prev_manager.get_configuration(
                conf_path=conf_path, logger=logger)

        # try to find a good manager if manager is None
        if manager is None:
            manager = self._get_manager(
                conf_path=conf_path,
                logger=logger,
                managers=self.managers)

        elif isclass(manager):
            manager = manager()

        else:
            manager = self._get_manager(
                conf_path=None,
                logger=logger,
                managers=manager)

        # if prev manager is not the new manager
        if prev_conf is not None \
                and type(manager) is not type(prev_manager):
            # update prev_conf with input conf
            prev_conf.update(conf)
            conf = prev_conf

        if manager is not None:
            manager.set_configuration(
                conf_path=conf_path,
                conf=conf,
                logger=logger)

        else:
            self.logger.error(
                'No ConfigurationManager found for \
                conf file {0}'.format(
                    conf_path))

        return result

    def configure(self, conf, logger=None):
        """
        Update self properties with input params only if:
        - self.configure is True
        - self.auto_conf is True
        - param conf 'configure' is True
        - param conf 'auto_conf' is True

        This method may not be overriden. see _configure instead

        :param conf: object from where get paramters
        :type conf: Configuration
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

    def _configure(self, unified_conf, logger=None):
        """
        Configure this class with input conf only if auto_conf or
        configure is true.

        This method should be overriden for specific conf

        :param unified_conf: Configuration with two categories
            VALUES and ERRORS
        :type params: Configuration

        :param configure: if True, force full self conf
        :type configure: bool
        """

        new_logger = False

        log_properties = [param.name for param in
            self.conf[Configurable.LOG]]

        for log_property in log_properties:
            new_logger = self._update_property(
                unified_conf, log_property) or new_logger

        # if needed, renew the logger
        if new_logger:
            self._logger = self.newLogger()

        # set managers
        self._update_property(
            unified_conf, Configurable.MANAGERS, public=True)

    def _update_property(
        self, unified_conf, param_name, public=False
    ):
        """
        True if a property update is required and do it.

        Check if a param exist in paramters where name is param_name.
        Then update self property depending on input public:
        - True => name is param_name
        - False => name is '_{param_name}'

        The idea of the public argument permits to avoid to run an auto_conf in
        changing a private attribute in using its setter method.

        :param unified_conf: unified conf
        :type params: Configuration

        :param param_name: param name to find in params
        :type param_name: str

        :param public: If False (default), update directly private
            property, else update public property in using the property.setter
        :type property_name: bool
        """

        result = False

        param = unified_conf[Configuration.VALUES].get(
            param_name)
        if param is not None:
            property_name = '{0}{1}'.format(
                '' if public else '_', param_name)
            setattr(self, property_name, param.value)
            result = True

        return result

    def _init_conf_files(self, conf_paths):

        self.conf_paths = self._get_conf_files() \
            if conf_paths is None else conf_paths

    def _get_conf_files(self):

        result = [Configurable.CONF_RESOURCE]

        return result

    @staticmethod
    def _get_manager(conf_path, managers, logger):
        """
        Get the first manager able to handle input conf_path.
        None if no manager is able to handle input conf_path.

        :return: first ConfigurationManager able to handle conf_path.
        :rtype: ConfigurationManager
        """

        result = None

        from canopsis.configuration.manager import ConfigurationManager

        for manager in managers.split(','):
            manager = ConfigurationManager.get_manager(manager)
            manager = manager()

            handle = conf_path is None \
                or manager.handle(conf_path=conf_path, logger=logger)

            if handle:
                result = manager
                break

        logger.warning(
            'No manager found among {0} for processing file {1}'.format(
                managers, conf_path))

        return result


def conf_paths(*conf_paths):
    """
    Configurable decorator which adds conf_path paths to a Configurable.

    :param paths: conf resource pathes to add to a Configurable
    :type paths: list of str

    Example:
    >>>@conf_paths('myexample/example.conf')
    >>>class MyConfigurable(Configurable):
    >>>    pass
    >>>MyConfigurable().conf_paths == (Configurable().conf_paths + ['myexample/example.conf'])
    """

    def _get_conf_files(self):
        # get super result and append conf_paths
        result = super(type(self), self)._get_conf_files()
        result += conf_paths

        return result

    def add_conf_paths(cls):
        # add _get_conf_files method to configurable classes
        if issubclass(cls, Configurable):
            cls._get_conf_files = _get_conf_files

        else:
            raise Configurable.Error(
                "class %s is not a Configurable class" % cls)

        return cls

    return add_conf_paths
