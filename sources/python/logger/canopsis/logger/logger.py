import os
import logging

from canopsis.common import root_path as canopath


class Output(object):

    @staticmethod
    def make(arg):
        raise NotImplementedError

class OutputStream(object):

    @staticmethod
    def make(stream):
        return logging.StreamHandler(stream)

class OutputFile(object):

    @staticmethod
    def make(path):

        if os.path.isabs(path):
            fpath = path

        else:
            fpath = os.path.join(canopath, path)

        return logging.FileHandler(fpath)


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
    def _init(logger, output, output_cls, level, fmt):
        logger.setLevel(level)

        if fmt is None:
            fmt = Logger.FORMATS.get(level, Logger.DEFAULT_FORMAT)

        formatter = logging.Formatter(fmt)

        handler = output_cls.make(output)
        handler.setLevel(level)
        handler.setFormatter(formatter)

        logger.addHandler(handler)

        return logger

    @staticmethod
    def _logger_exists(name):
        return name in logging.Logger.manager.loggerDict

    @staticmethod
    def get(name, output, output_cls=OutputFile, level=logging.INFO, fmt=None):
        """
        :param name: logger name
        :param output: output given to output_cls. Can be anything from a file path, a stringio or a full URI. It only has to be supported by output_cls.
        :param output_cls: canopsis.logger.Output<output> class.
        :param level: logging.<LEVEL>
        :param fmt: format to apply. If None, defaults to Logger.DEFAULT_FORMAT.
        """
        if Logger._logger_exists(name):
            return logging.getLogger(name)

        logger = logging.getLogger(name)
        return Logger._init(logger, output, output_cls, level, fmt)