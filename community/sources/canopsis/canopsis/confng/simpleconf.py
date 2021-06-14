# -*- coding: utf-8 -*-

from __future__ import unicode_literals

from canopsis.confng.vendor import SimpleConf


class ConfigurationUnreachable(Exception):
    """
    Raise/catch this exception when configuration cannot be fetched
    for any reason.
    """
    pass


class Configuration(SimpleConf):
    """
    Lightweigh configuration loader.
    """

    @staticmethod
    def load(path, driver_cls, *args, **kwargs):
        """
        :param str path: configuration location. path format depends on the driver used
        :param object driver_cls: the class to use to read the file
        :param args: positional parameters for driver_cls
        :param kwargs: named parameters for driver_cls
        :raises ConfigurationUnreachable: cannot open configuration
        :rtype: dict
        """
        conf = {}

        try:
            driver = driver_cls(path, *args, **kwargs)
            conf = SimpleConf.export(driver)
        except IOError as ex:
            raise ConfigurationUnreachable(str(ex))

        return conf
