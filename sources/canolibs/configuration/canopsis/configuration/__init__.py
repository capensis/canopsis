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

from os.path import expanduser, exists, sep, abspath

from os import stat

from collections import OrderedDict, Iterable

from inspect import isclass


class Configuration(object):
    """
    Manage conf such as a list of Categories.

    The order of categories permit to ensure param overriding.
    """

    ERRORS = 'ERRORS'
    VALUES = 'VALUES'

    def __init__(self, *categories, **kwargs):
        """
        :param categories: categories to configure.
        :type categories: list of Category.
        """

        super(Configuration, self).__init__()

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

    def __iadd__(self, other):
        """
        Add categories or conf categories in self
        """

        # if other is a conf add a copy of all other categories
        if isinstance(other, Configuration):
            for category in other:
                self += category

        else:  # in case of category
            category = self.get(other.name)

            if category is None:
                self.put(other)

            else:
                for param in other:
                    category.put(param)

        return self

    def __repr__(self):

        return 'Configuration({0})'.format(self.categories)

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

    def unify(self, copy=False, *args, **kwargs):
        """
        Get a conf which contains only two categories:
        - VALUES where params are all self params where values are not
            exceptions.contains all values
        - ERRORS where params are all self params where values are
            exceptions

        :param copy: copy self params (default False)
        :type côpy: bool

        :return: two categories named respectivelly VALUES and ERRORS and
            contain respectivelly self param values and parsing errors
        :rtype: Configuration
        """

        result = Configuration()

        values = Category(Configuration.VALUES)
        errors = Category(Configuration.ERRORS)

        for category in self:

            for param in category:

                if param.value is not None:

                    to_update, to_delete = (errors, values) if \
                        isinstance(param.value, Exception) \
                        else (values, errors)

                    to_update.put(param.copy() if copy else param)

                    if param.name in to_delete:
                        del to_delete[param.name]

        result += values
        result += errors

        return result

    def get_unified_category(self, name, copy=False, *args, **kwargs):
        """
        Add a category with input name which takes all params provided
        by other categories

        :param name: new category name
        :type name: str

        :param copy: copy self params (default False)
        :type copy: bool
        """

        result = Category(name)

        for category in self:
            for param in category:
                result.put(param.copy() if copy else param)

        return result

    def add_unified_category(
        self, name, copy=False, new_content=None, *args, **kwargs
    ):
        """
        Add a unified category to self and add new_content if not None
        """
        category = self.get_unified_category(name=name, copy=copy)

        if new_content is not None:
            category += new_content

        self += category

    def clean(self, *args, **kwargs):
        """
        Clean this params in setting value to None.
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

    def update(self, conf, *args, **kwargs):
        """
        Update this content with input conf
        """

        for category in conf:
            category = self.setdefault(
                category.name, category.copy())

            for param in category:
                param = category.setdefault(
                    param.name, param.copy())


class Category(object):
    """
    Parameter category which contains a dictionary of params.
    """

    def __init__(self, name, *params, **kwargs):
        """
        :param name: unique in a conf.
        :type name: str

        :param params: Parameters
        :type params: list of Parameter
        """
        super(Category, self).__init__()

        self.name = name
        # set param by names.
        self.params = {
            param.name: param for param in params}

    def __iter__(self, *args, **kwargs):

        return iter(self.params.values())

    def __delitem__(self, param_name, *args, **kwargs):

        del self.params[param_name]

    def __getitem__(self, param_name, *args, **kwargs):

        return self.params[param_name]

    def __contains__(self, param_name, *args, **kwargs):

        return param_name in self.params

    def __len__(self):

        return len(self.params)

    def __eq__(self, other):

        return isinstance(other, Category) and other.name == self.name

    def __hash__(self):

        return hash(self.name)

    def __repr__(self):

        return 'Category({0}, {1})'.format(self.name, self.params)

    def __iadd__(self, value):

        if isinstance(value, Category):
            self += value.params.values()

        elif isinstance(value, Iterable):
            for content in value:
                self += content

        elif isinstance(value, Parameter):
            self.put(value)

        else:
            raise Exception('Wrong type to add {0} to {1}. \
Must be a Category, a Parameter or a list of {Parameter, Category}'.format(
                value, self))

        return self

    def setdefault(self, param_name, param, *args, **kwargs):

        return self.params.setdefault(param_name, param)

    def get(self, param_name, default=None, *args, **kwargs):

        return self.params.get(param_name, default)

    def put(self, param, *args, **kwargs):
        """
        Put a param and return the previous one if exist
        """

        result = self.get(param.name)
        self.params[param.name] = param
        return result

    def clean(self, *args, **kwargs):
        """
        Clean this params in setting value to None.
        """

        for param in self.params.values():

            param.clean()

    def copy(self, name=None, *args, **kwargs):

        if name is None:
            name = self.name

        result = Category(name)

        for param in self:
            result.put(param.copy())

        return result


class Parameter(object):
    """
    Parameter identified among a category by its name.
    Provide a value (None by default) and a parser (str by default).
    """

    @staticmethod
    def bool(value):
        return value == 'True' or value == 'true' or value == '1'

    def __init__(self, name, value=None, parser=str, *args, **kwargs):
        """
        :param name: unique by category
        :type name: str

        :param value: param value. None if not given.
        :type value: object

        :param parser: param test deserializer which takes in param
            a str.
        :type parser: callable
        """

        super(Parameter, self).__init__()

        self.name = name
        self._value = value
        self.parser = parser

    def __eq__(self, other):

        return isinstance(other, Parameter) and other.name == self.name

    def __hash__(self):

        return hash(self.name)

    def __repr__(self):

        return 'Parameter({0}, {1}, {2})'.format(
            self.name, self.value, self.parser)

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

    def clean(self, *args, **kwargs):
        """
        Clean this param in removing values
        """

        self._value = None


class MetaConfigurable(type):
    """
    Meta class for Configurable.
    """

    def __call__(cls, *args, **kwargs):
        """
        Get a new instance of input cls class, and if instance.auto_conf or
        instance.once, then call instance.apply_configuration().
        """

        result = type.__call__(cls, *args, **kwargs)

        if result.auto_conf or result.once:
            result.apply_configuration()

        return result


class Configurable(object):
    """
    Manages class conf synchronisation with conf files.
    """

    __metaclass__ = MetaConfigurable

    DEFAULT_MANAGERS = '%s,%s' % (
        'canopsis.configuration.manager.json.ConfigurationManager',
        'canopsis.configuration.manager.ini.ConfigurationManager')

    CONF_FILE = '~/etc/global.conf'

    CONF = 'CONFIGURATION'
    LOG = 'LOG'

    AUTO_CONF = 'auto_conf'
    ONCE = 'once'
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
        conf_files=None, managers=DEFAULT_MANAGERS,
        auto_conf=True, once=False,
        log_lvl='INFO', log_name=None, log_info_format=INFO_FORMAT,
        log_debug_format=DEBUG_FORMAT, log_warning_format=WARNING_FORMAT,
        log_error_format=ERROR_FORMAT, log_critical_format=CRITICAL_FORMAT,
        *args, **kwargs
    ):
        """
        :param conf_files: conf_files to parse
        :type conf_files: Iterable or str

        :param auto_conf: true force auto conf as soon as param change
        :type auto_conf: bool

        :param once: true force auto conf once as soon as param change
        :type once: bool

        :param log_lvl: logging level
        :type log_lvl: str
        """

        super(Configurable, self).__init__()

        self.auto_conf = auto_conf
        self.once = once

        # set conf files
        self._init_conf_files(conf_files)

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

    def _conf(self, *args, **kwargs):
        """
        Protected method to override in order to specify which conf
        return with parsers and default values
        """

        result = Configuration(
            Category(Configurable.CONF,
                Parameter(
                    Configurable.AUTO_CONF, self.auto_conf, Parameter.bool),
                Parameter(Configurable.MANAGERS, self.managers),
                Parameter(Configurable.ONCE, self.once, Parameter.bool)),
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
    def conf_files(self):
        """
        Get all type conf files and user files.

        :return: self conf files
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

        from canopsis.configuration.watcher import add_configurable,\
            remove_configurable

        # remove previous watching
        remove_configurable(self)
        self._conf_files = tuple(value)
        # add new watching
        add_configurable(self)

    def apply_configuration(
        self, conf=None, conf_files=None,
        managers=None, *args, **kwargs
    ):
        """
        Apply conf on a destination in 5 phases:

        1. identify the right manager to use with conf_files to parse.
        2. for all conf_files, get conf which match
            with input conf.
        3. apply parsing rules on conf_file params.
        4. put values and parsing errors in two different dictionaries.
        5. returns both dictionaries of param values and errors.

        :param conf: conf from where get conf
        :type conf: Configuration

        :param conf_files: conf files to parse. If
            conf_files is a str, it is automatically putted into a list
        :type conf_files: list of str
        """

        if conf is None:
            conf = self.conf

        conf = self.get_configuration(
            conf=conf, conf_files=conf_files,
            managers=managers, *args, **kwargs)

        self.configure(conf=conf, *args, **kwargs)

    def get_configuration(
        self,
        conf=None, conf_files=None, logger=None,
        managers=None, fill=False, *args, **kwargs
    ):
        """
        Get a dictionary of params by name from conf,
        conf_files and conf_managers

        :param conf: conf to update. If None, use \
            self.conf
        :type conf: Configuration

        :param conf_files: list of conf files. If None, use \
            self.conf_files
        :type conf_files: list of str

        :param logger: logger to use for logging info/error messages.
            If None, use self.logger
        :type logger: logging.Logger

        :param managers: conf managers to use. If None, use self.managers
        :type managers: list of ConfigurationManager

        :param fill: if True (False by default) load in conf all \
            conf_files content
        :type fill: bool
        """

        # start to initialize input params
        if logger is None:
            logger = self._logger

        if conf is None:
            conf = self.conf

        # remove values from conf
        conf.clean()

        if conf_files is None:
            conf_files = self._conf_files

        if isinstance(conf_files, str):
            conf_files = [conf_files]

        # clean conf file list
        conf_files = [
            abspath(expanduser(conf_file)) for conf_file
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

                conf = conf_manager.get_configuration(
                    conf=conf, fill=fill,
                    conf_file=conf_file, logger=logger)

            else:
                # if no conf_manager, display a warning log message
                logger.warning('No manager found among {0} for {1}'.format(
                    conf_file))

        return conf

    def set_configuration(
        self, conf_file, conf, manager=None,
        logger=None, *args, **kwargs
    ):
        """
        Set params on input conf_file.

        Args:
            - conf_files (str): conf_file to udate with
                params
            - parameter_by_categories (dict(str: dict(str: object)):
            - logger (logging.Logger): logger to use to set params.
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

        else:
            manager = self._get_manager(
                conf_file=None,
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
                conf_file=conf_file,
                conf=conf,
                logger=logger)

        else:
            logger.error(
                'No ConfigurationManager found for \
                conf file {0}'.format(
                    conf_file))

        return result

    def configure(self, conf, *args, **kwargs):
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

        unified_conf = conf.unify()

        values = unified_conf[Configuration.VALUES]

        # set configure
        once_parameter = values.get(Configurable.ONCE)
        if once_parameter is not None:
            self.once = once_parameter.value

        # set auto_conf
        auto_conf_parameter = values.get(Configurable.AUTO_CONF)
        if auto_conf_parameter is not None:
            self.auto_conf = auto_conf_parameter.value

        if self.once or self.auto_conf:
            self._configure(unified_conf=unified_conf, *args, **kwargs)
            # when conf succeed, deactive once
            self.once = False

    def _configure(self, unified_conf, *args, **kwargs):
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
            unified_conf, Configurable.MANAGERS, public_property=True)

    def _update_property(
        self, unified_conf, param_name, public_property=False,
        *args, **kwargs
    ):
        """
        True if a property update is required and do it.

        Check if a param exist in paramters where name is param_name.
        Then update self property depending on input public_property:
        - True => name is param_name
        - False => name is '_{param_name}'

        :param unified_conf: unified conf
        :type params: Configuration

        :param param_name: param name to find in params
        :type param_name: str

        :param public_property: If False (default), update directly private
            property, else update public property in using the property.setter
        :type property_name: bool
        """

        result = False

        param = unified_conf[Configuration.VALUES].get(
            param_name)
        if param is not None:
            property_name = '{0}{1}'.format(
                '' if public_property else '_', param_name)
            setattr(self, property_name, param.value)
            result = True

        return result

    def _init_conf_files(self, conf_files, *args, **kwargs):

        self.conf_files = self._get_conf_files(*args, **kwargs) \
            if conf_files is None else conf_files

    def _get_conf_files(self, *args, **kwargs):

        result = [Configurable.CONF_FILE]

        return result

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

        result = None

        from canopsis.configuration.manager import ConfigurationManager

        for manager in managers.split(','):
            manager = ConfigurationManager.get_manager(manager)
            manager = manager()

            handle = conf_file is None \
                or manager.handle(conf_file=conf_file, logger=logger)

            if handle:
                result = manager
                break

        logger.warning(
            'No manager found among {0} for processing file {1}'.format(
                managers, conf_file))

        return result
