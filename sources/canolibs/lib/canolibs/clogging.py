#!/usr/bin/env python
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

# file for canopsis logging configuration file
LOGGING_CONFIGURATION_FILENAME = 'logging.conf'
# file for python logging configuration file
PYTHON_LOGGING_CONFIGURATION_FILENAME = 'python-logging.conf'


import os.path

LOG_DIRECTORY = os.path.expanduser('~/var/log/')

INFO_FORMAT = "%(asctime)s [%(name)s] [%(levelname)s] %(message)s"

DEBUG_FORMAT = "%(asctime)s [%(name)s] [%(levelname)s] [path: %(pathname)s]\
 [p: %(process)d] [t: %(thread)d] [f: %(funcName)s] [l: %(lineno)d]\
  %(message)s"

DATE_FORMAT = "%Y-%m-%d %H:%M:%S"

import logging

INFO_FORMATTER = logging.Formatter(fmt=INFO_FORMAT, datefmt=DATE_FORMAT)
DEBUG_FORMATTER = logging.Formatter(fmt=DEBUG_FORMAT, datefmt=DATE_FORMAT)

import sys
if hasattr(sys, 'frozen'):  # support for py2exe
    _srcfile = "logging%s__init__%s" % (os.sep, __file__[-4:])
elif str.lower(__file__[-4:]) in ['.pyc', '.pyo']:
    _srcfile = __file__[:-4] + '.py'
else:
    _srcfile = __file__
_srcfile = os.path.normcase(_srcfile)


class CanopsisLogger(logging.Logger):
    """
    Logger dedicated to Canopsis files.
    """

    # static namespace for global scope information
    __SCOPE__ = '__SCOPE__'

    def __init__(self, name, level=logging.INFO):
        """
        Create a default file handler where filename corresponds to input name.
        Name tree is preserved in log file tree.
        """

        super(CanopsisLogger, self).__init__(name, level)

        self.setScope()

    def setScope(self, scope=None):
        """
        Change scope, i.e. file handler target.
        """

        path = self.getLogPath(scope)

        # create log file if not exists
        if not os.path.exists(path):
            directory = os.path.dirname(path)
            if not os.path.exists(directory):
                os.makedirs(directory)

        self.handler = logging.FileHandler(path, 'a')
        self.addHandler()

    def getLogPath(self, scope=None):
        """
        Get log file path corresponding to its name.
        """

        if scope is None:
            scope = globals().get(CanopsisLogger.__SCOPE__, self.name)

        filename = scope.replace('.', os.path.sep) + '.log'
        result = os.path.join(LOG_DIRECTORY, filename)

        return result

    def debug(self, msg, *args, **kwargs):
        self.log(logging.DEBUG, msg, *args, **kwargs)

    def info(self, msg, *args, **kwargs):
        self.log(logging.INFO, msg, *args, **kwargs)

    def warning(self, msg, *args, **kwargs):
        self.log(logging.WARNING, msg, *args, **kwargs)

    def critical(self, msg, *args, **kwargs):
        self.log(logging.CRITICAL, msg, *args, **kwargs)

    def error(self, msg, *args, **kwargs):
        self.log(logging.ERROR, msg, *args, **kwargs)

    def log(self, level, msg, *args, **kwargs):
        """
        Change dynamically of formatter if no new handler has been requested.
        """

        if self.handler is not None:
            if self.isEnabledFor(level):
                if level <= logging.DEBUG:
                    self.handler.setFormatter(DEBUG_FORMATTER)
                else:
                    self.handler.setFormatter(INFO_FORMATTER)

        super(CanopsisLogger, self).log(level, msg, *args, **kwargs)

        # log debug message for this
        if self is not _logger:
            _logger.debug('log:%s, level: %s, msg: %s', self.name, level, msg)

    def findCaller(self):
        """
        Find the stack frame of the caller so that we can note the source
        file name, line number and function name.
        """

        f = logging.currentframe().f_back
        rv = "(unknown file)", 0, "(unknown function)"

        while hasattr(f, "f_code"):
            co = f.f_code
            filename = os.path.normcase(co.co_filename)
            # This line is modified.
            if filename in (_srcfile, logging._srcfile):
                f = f.f_back
                continue
            rv = (filename, f.f_lineno, co.co_name)
            break

        return rv

    def addHandler(self, handler=None):
        """
        Check if the call has been done during self initialization.
        """

        if handler is None:
            if self.handler is not None:
                self.removeHandler(self.handler)
            handler = self.handler
        else:
            self.handler = None

        super(CanopsisLogger, self).addHandler(handler)

logging.setLoggerClass(CanopsisLogger)

import inspect


def getLogger(name=None, scope=None):
    """
    Get a logger in a Canopsis environment.
    - name: name of new logger. If None, iname is callee module.
    - scope: output scope identity. If None, use the last defined.
    """

    if name is None:
        f_back = inspect.currentframe().f_back
        # get previous frame module name
        name = f_back.f_globals['__name__']
        if name == '__main__':
            # get filename in case of main process
            name = f_back.f_code.co_filename
            if name.endswith('.py'):
                name = name[:-3]
            elif name.endswith('.pyc'):
                name = name[:-4]

    result = logging.getLogger(name)

    if scope is not None:
        result.setScope(scope)

    return result

# instantiate a logger
_logger = getLogger()


def getChildLogger(name=None):
    """
    Get a child logger related to previous frame.
    """

    f_back = inspect.currentframe().f_back
    # get previous frame module name
    parent = f_back.f_back.f_globals['__name__']
    if parent == '__main__':
        # get filename in case of main process
        parent = f_back.f_back.f_code.co_filename
        if parent.endswith('.py'):
            parent = parent[:-3]
        elif parent.endswith('.pyc'):
            parent = parent[:-4]

    if name is not None:
        name = "{0}.{1}".format(parent, name)
    else:
        name = parent

    result = logging.getLogger(name)

    return result


def getRootLogger():
    """
    Get Root logger.
    """

    result = logging.getLogger()
    return result

# bind observers to both configuration files

# register file configuration changes into the global configuration file
import ConfigParser
LEVEL = 'level'


def loadConfigurationFile(src_path):
    """
    Reuse simple configuration file in order to parameterize loggers.

    Sections are logger names, and options are:
    - level: level_value or level_name.
    """

    _logger.debug('src_path: %s', src_path)

    config_parser = ConfigParser.RawConfigParser()
    config_parser.read(src_path)

    for section in config_parser.sections():
        logger = logging.getLogger(section)

        if config_parser.has_option(section, LEVEL):
            level = config_parser.get(section, LEVEL)

            if str.isdigit(level):
                level = int(level)

            logger.setLevel(level)

import logging.config


def loadPythonConfigurationFile(src_path):
    """
    Reuse python logging configuration file in order to parameterize loggers.
    """

    _logger.debug('src_path: %s', src_path)

    logging.config.fileConfig(src_path)
