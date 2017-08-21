#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os
import sys

from canopsis.common import root_path
from canopsis.confng.vendor import SimpleConf

class ConfigurationFileNotFound(Exception):
    pass

class ConfigurationUnreachable(Exception):
    pass

class Configuration(SimpleConf):

    @staticmethod
    def load(conf_path, driver_cls, *args, **kwargs):
        """
        Load configuration file regarding available paths from sys.path
        when conf_path isn't an absolute path.

        :param str conf_path: the file location
        :param object driver_cls: the class to use to read the file
        :raises ConfigurationFileNotFound: configuration file is not present.
        :raises ConfigurationUnreachable: cannot open configuration file.
        :rtype: dict
        """
        conf_file = None

        if os.path.isabs(conf_path):
            conf_file = conf_path

        else:
            fpath = os.path.join(root_path, conf_path)

            if os.path.isfile(fpath):
                conf_file = fpath
            else:
                raise ConfigurationFileNotFound(fpath)

        conf = {}

        try:
            with open(conf_file, 'r') as fh:
                driver = driver_cls(fh=fh, *args, **kwargs)
                conf = SimpleConf.export(driver)
        except IOError, ex:
            raise ConfigurationUnreachable(str(ex))

        return conf
