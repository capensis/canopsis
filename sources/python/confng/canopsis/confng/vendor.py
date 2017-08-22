#!/usr/bin/env python
# -*- coding: utf-8 -*-

# https://github.com/Leryan/leryan.types/tree/v0.0.17

from __future__ import unicode_literals

from configparser import ConfigParser, ExtendedInterpolation
import json

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO


class Driver(object):
    """
    A generic driver class.
    """

    def __init__(self, fh=None, sconf=None, *args, **kwargs):
        super(Driver, self).__init__(*args, **kwargs)

        if fh is None and sconf is None:
            raise ValueError('pass either a file handler or a string')

        if sconf is not None and fh is None:
            fh = StringIO(sconf)

        self._fh = fh

    def export(self):
        """
        Must always return a dict() object.
        """
        raise NotImplementedError()


class Ini(Driver):
    """
    Reads ini file and returns configuration in a dict.

    Supports ExtendedInterpolation.
    """

    def __init__(self, fh=None, sconf=None, with_interpolation=False, *args, **kwargs):
        """
        :param fh: file-like object.
        :param sconf: string containing INI-formatted configuration.
        :param with_interpolation: enable ExtendedInterpolation. Default to False.
        """

        super(Ini, self).__init__(fh=fh, sconf=sconf, *args, **kwargs)

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


class Json(Driver):
    """
    Read a json paylod and returns configuration as a dict.
    """

    def export(self):
        return json.load(self._fh)


class SimpleConf(object):

    @staticmethod
    def export(driver, output_class=dict):
        return output_class(driver.export())
