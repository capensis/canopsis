# -*- coding: utf-8 -*-
from __future__ import unicode_literals

import os
import logging

from six import string_types

from logging import FileHandler, StreamHandler, NullHandler
from logging.handlers import MemoryHandler

from canopsis.common import root_path as canopath


class Output(object):

    @staticmethod
    def make(_):
        """
        Override this method, then return a subclass of logging.Handler.

        Like FileHandler, StreamHandler...
        """
        raise NotImplementedError

    @staticmethod
    def make_memory(capacity, flushlevel, target):
        """
        Same as make(), but with a memory buffer, flushed when full or
        on critical events.

        :param capacity: memory handler capacity
        :param flushlevel: memory handler flush level.
        :return: memory handler
        :rtype: MemoryHandler
        """
        return MemoryHandler(capacity, flushLevel=flushlevel, target=target)


class OutputStream(Output):

    @staticmethod
    def make(stream):
        """
        :param stream: any file-like object.
        :return: stream handler
        :rtype: logging.StreamHandler
        """
        return StreamHandler(stream)


class OutputFile(Output):

    @staticmethod
    def make(path, default_root_path=canopath):
        """
        :param str path: path to file. If relative, will use canopsis.common.root_path as prefix.
        :return: file handler
        :rtype: logging.FileHandler
        """

        if os.path.isabs(path):
            fpath = path

        else:
            fpath = os.path.join(default_root_path, path)

        return FileHandler(fpath)


class OutputNull(Output):

    @staticmethod
    def make(_):
        """
        :param _: unused required argument
        :return: null handler
        :rtype: logging.NullHandler
        """

        return NullHandler()


class Logger(object):

    DEFAULT_FORMAT = "[%(asctime)s] [%(levelname)s] [%(name)s] %(message)s"

    FORMATS = {
        logging.DEBUG: "[%(asctime)s] [%(levelname)s] [%(name)s] [%(process)d] [%(thread)d] [%(pathname)s] [%(lineno)d] %(message)s",
        logging.CRITICAL: DEFAULT_FORMAT,
        logging.ERROR: DEFAULT_FORMAT,
        logging.FATAL: DEFAULT_FORMAT,
        logging.INFO: DEFAULT_FORMAT,
        logging.WARN: DEFAULT_FORMAT,
        logging.WARNING: DEFAULT_FORMAT,
    }

    @staticmethod
    def _init(logger, output, output_cls, level, fmt,
              memory, memory_capacity, memory_flushlevel,
              driver_make_args):
        if isinstance(level, string_types):
            level = getattr(logging, level.upper(), logging.INFO)

        logger.setLevel(level)

        if fmt is None:
            fmt = Logger.FORMATS.get(level, Logger.DEFAULT_FORMAT)

        formatter = logging.Formatter(fmt=fmt)

        handler = output_cls.make(output, **driver_make_args)
        handler.setFormatter(formatter)
        handler.setLevel(level)

        if memory:
            memory_handler = output_cls.make_memory(memory_capacity,
                                                    memory_flushlevel,
                                                    handler)
            logger.addHandler(memory_handler)

        else:
            logger.addHandler(handler)

        return logger

    @staticmethod
    def _logger_exists(name):
        return name in logging.Logger.manager.loggerDict

    @staticmethod
    def get(name, output, output_cls=OutputFile, level=logging.INFO, fmt=None,
            memory=False, memory_capacity=100,
            memory_flushlevel=logging.WARNING, driver_make_args={}):
        """
        Return the right logger.

        :param str name: logger name
        :param str output: output given to output_cls. Can be anything from a file path, a stringio or a full URI. It only has to be supported by output_cls.
        :param class output_cls: canopsis.logger.Output<output> class.
        :param level: logging.<LEVEL>
        :param fmt: format to apply. If None, defaults to Logger.DEFAULT_FORMAT.
        :param memory: wrap logging handler with logging.handlers.MemoryHandler.
        :param memory_capacity: MemoryHandler log capacity.
        :param memory_flushlevel: if a log event is equal or greater than this level, force flush.
        :param dict driver_make_args: dict of arguments to pass to the Output<Driver>.make() function. Acts as the usual **kwargs.
        :return: python logger.
        :rtype: logging.Logger
        """
        if Logger._logger_exists(name):
            return logging.getLogger(name)

        logger = logging.getLogger(name)
        return Logger._init(logger, output, output_cls, level, fmt,
                            memory, memory_capacity, memory_flushlevel,
                            driver_make_args)
