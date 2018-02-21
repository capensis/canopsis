#!/usr/bin/env python
# -*- coding: utf-8 -*-

# https://github.com/Leryan/leryan.types/tree/v0.0.17

from __future__ import unicode_literals

from configparser import ConfigParser, ExtendedInterpolation
import json
import os

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.common import root_path


class Driver(object):
    """
    A generic driver class.
    """

    def __init__(self, url, *args, **kwargs):
        super(Driver, self).__init__(*args, **kwargs)
        self._url = url

    def export(self):
        """
        Must always return a dict() object.
        """
        raise NotImplementedError()


class FileDriver(Driver):

    def __init__(self, path, sconf=None, fh=None, default_root=root_path):
        """
        Used conf follow this preference order:

            path
            sconf
            fh

        :param str path: absolute or relative path. if relative, will be concatenated to default_root: we do not support relative files from current workdir.
        :param str sconf: configuration as string
        :param file fh: file handler
        :param str default_root: default root for relative paths
        :raises ValueError: path, sconf and fh are None
        """
        super(FileDriver, self).__init__(path)

        if path is not None:
            conf_file = path

            if not os.path.isabs(path):
                conf_file = os.path.join(default_root, path)

            self._fh = open(conf_file, 'r')

        elif sconf is not None:
            self._fh = StringIO(sconf)

        elif fh is not None:
            self._fh = fh

        else:
            raise ValueError('pass either a file path, a file handler or a string')


class Ini(FileDriver):
    """
    Reads ini file and returns configuration in a dict.

    Supports ExtendedInterpolation.
    """

    def __init__(self, path=None, fh=None, sconf=None, with_interpolation=False, *args, **kwargs):
        """
        See FileDriver.__init__ doc for parameters.
        :param with_interpolation: enable ExtendedInterpolation. Default to False.
        """

        super(Ini, self).__init__(path=path, fh=fh, sconf=sconf, *args, **kwargs)

        self._with_interpolation = with_interpolation

    def export(self):
        if self._with_interpolation:
            config = ConfigParser(interpolation=ExtendedInterpolation())
        else:
            config = ConfigParser()

        config.read_file(self._fh)

        conf = {}

        for section in config.sections():
            conf[section] = {}

            for k, v in config.items(section=section):
                conf[section][k] = v

        return conf


class Json(FileDriver):
    """
    Read a json paylod and returns configuration as a dict.
    """

    def export(self):
        return json.load(self._fh)


class SimpleConf(object):

    @staticmethod
    def export(driver, output_class=dict):
        return output_class(driver.export())
