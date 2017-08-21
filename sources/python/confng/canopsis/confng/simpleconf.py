#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import os
import sys

from canopsis.confng.vendor import SimpleConf


class Configuration(SimpleConf):

    @staticmethod
    def load(conf_path, driver_cls, *args, **kwargs):
        """
        Load configuration file regarding available paths from sys.path
        when conf_path isn't an absolute path.

        :param str conf_path: the file location
        :param object driver_cls: the class to use to read the file
        :rtype: dict
        """
        conf_file = None

        if os.path.isabs(conf_path):
            conf_file = conf_path

        else:
            for path in sys.path:

                fpath = os.path.join(path, conf_path)

                if os.path.isfile(fpath):
                    conf_file = fpath
                    break

        conf = {}
        with open(conf_file, 'r') as fh:
            driver = driver_cls(fh=fh, *args, **kwargs)
            conf = SimpleConf.export(driver)

        return conf
