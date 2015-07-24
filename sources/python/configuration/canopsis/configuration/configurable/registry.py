# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from . import Configurable
from canopsis.common.init import basestring
from canopsis.configuration.parameters import Configuration, Category
from canopsis.common.utils import lookup

from inspect import isclass


class Configurables(dict):
    """
    With a ConfigurableTypes, it is in charge of a ConfigurableRegistry
    sub-configurables.
    When a configurable is trying to be setted, the type is checked related
    to its ConfigurableRegistry ConfigurableTypes.
    """

    def __init__(self, registry, values=None, *args, **kwargs):
        """
        :param registry: related registry
        :type registry: ConfigurableRegistry

        :param values: default values if not None (default)
        :type values: dict
        """
        super(Configurables, self).__init__(*args, **kwargs)

        self.registry = registry

        if values is not None:
            for name, value in values:
                self[name] = value

    def __setitem__(self, name, value):
        """
        Set a new configurable value

        :param name: new configurable name
        :type name: str

        :param value: new configurable value
        :type value: str (path) or class or instance

        :param args: args for configurable instanciation if value is a path or
            a class
        :param kwargs: kwargs for configurable instanciation if value is a path
            or a class
        """

        # get configurable type. Configurable by default
        configurable_type = self.registry._configurable_types.get(
            name, Configurable)

        configurable = value

        # if value is a path
        if isinstance(configurable, basestring):
            # get related python object
            configurable = lookup(configurable)

        # if configurable is a class
        if isclass(configurable) and issubclass(
                configurable, configurable_type):
                # instantiate a new configurable with input args and kwargs
                configurable = configurable()

        # do nothing if configurable is not an instance of configurable_type
        if not isinstance(configurable, configurable_type):
            self.registry.logger.error(
                "Impossible to set configurable {}:{}. Not an instance of {}"
                .format(name, configurable, configurable_type)
            )

        else:
            # update self.configurables
            super(Configurables, self).__setitem__(name, configurable)


class ConfigurableTypes(dict):
    """
    With a Configurables, it is in charge of a set of configurable type.
    When a new type is setted but the old configurable value does not inherits
    from it, then the old value is removed automatically.
    """

    def __init__(self, registry, values=None, *args, **kwargs):
        """
        :param registry: related registry
        :type registry: ConfigurableRegistry

        :param values: default values if not None (default)
        :type values: dict
        """

        super(ConfigurableTypes, self).__init__(*args, **kwargs)

        self.registry = registry

        if values is not None:

            for name in values:
                value = values[name]
                self[name] = value

    def __setitem__(self, name, value):
        """
        Set a new configurable type.

        :param name: new configurable name.
        :type name: str

        :param value: new type value.
        :type value: str (path) or class
        """

        configurable_type = value

        # if configurable_type is a path
        if isinstance(configurable_type, basestring):
            # get related python object
            configurable_type = lookup(configurable_type)

        # check if configurable_type is a subclass of Configurable
        if not issubclass(configurable_type, Configurable):
            self.registry.logger.error(
                "Impossible to set configurable type {}: {}. Wrong type"
                .format(name, configurable_type)
            )

        else:
            # check if an old value exiss
            if name in self.registry._configurables \
                    and not isinstance(
                        self.registry._configurables[name], configurable_type):
                # if the old value is not an instance of newly type
                self.registry.logger.warning(
                    "Old configurable {} removed. Not an instance of {}"
                    .format(name, configurable_type)
                )
                # delete if
                del self.registry._configurables[name]

            # set the new type
            super(ConfigurableTypes, self).__setitem__(name, configurable_type)


class ConfigurableRegistry(Configurable):
    """
    Manage a set of configurables which are accessibles from self.configurables

    Each configurable can be defined in conf parameters where names are like
    {name}_configurable={configurable_path, configurable_class, configurable}

    And a configurable configuration are in categories {NAME}_CONF.

    Then all sub-configurables are accessibles from item accessors with name as
    item key.
    """

    class Error(Exception):
        """handle ConfigurableRegistry errors"""
        pass

    CONF_PATH = 'configuration/registry.conf'  #: default conf path

    CATEGORY = 'MANAGER'  #: default ConfigurableRegistry category name

    CONFIGURABLE_SUFFIX = '_value'  #: configurable configuration suffix
    CONFIGURABLE_TYPE_SUFFIX = '_type'  #: type config suffix

    def __init__(
        self, configurables=None, configurable_types=None, *args, **kwargs
    ):
        """
        :param configurables: dictionary of configurables by name.
        :type configurables: dict

        :param configurable_types: dictionary of configurable types by name
        :type configurable_types: dict
        """

        super(ConfigurableRegistry, self).__init__(*args, **kwargs)

        self._configurables = Configurables(self, configurables)
        self._configurable_types = ConfigurableTypes(self, configurable_types)

    def _get_category(self):
        """Get category.

        :rtype: Category
        """

        result = Category(ConfigurableRegistry.CATEGORY)

        return result

    def _conf(self, *args, **kwargs):

        result = super(ConfigurableRegistry, self)._conf(*args, **kwargs)

        result.add_unified_category(name=ConfigurableRegistry.CATEGORY)

        return result

    def _get_conf_paths(self, *args, **kwargs):

        result = super(ConfigurableRegistry, self)._get_conf_paths(
            *args, **kwargs
        )

        result.append(ConfigurableRegistry.CONF_PATH)

        return result

    def apply_configuration(
        self,
        conf=None, conf_paths=None, drivers=None, logger=None, override=True,
        *args, **kwargs
    ):

        super(ConfigurableRegistry, self).apply_configuration(
            conf=conf, conf_paths=conf_paths, drivers=drivers, logger=logger,
            override=override,
            *args, **kwargs
        )

        if conf_paths is None:
            conf_paths = self.conf_paths

        if conf is None:
            conf = self.conf

        if drivers is None:
            drivers = self.drivers

        # get self conf path
        conf_path = conf_paths[-1]

        configurables = self._configurables

        # apply configuration to all self configurables
        for name in configurables:
            configurable = configurables[name]
            # add self last conf paths to configurable conf paths
            configurable_conf_paths = list(configurable.conf_paths)
            configurable_conf_paths.append(conf_path)
            # get a copy of configurable configuration
            configurable_configuration = configurable.conf.copy()
            # add a unified category where name is {NAME}_CONF
            category_name = ConfigurableRegistry.get_configurable_category(
                name
            )
            configurable_configuration.add_unified_category(
                name=category_name, copy=True
            )
            # apply configurable configuration
            configurable.apply_configuration(
                conf=configurable_configuration,
                conf_paths=configurable_conf_paths,
                drivers=drivers, logger=logger, override=override
            )

    def _configure(self, unified_conf, *args, **kwargs):

        super(ConfigurableRegistry, self)._configure(
            unified_conf=unified_conf, *args, **kwargs
        )

        foreigns = unified_conf[Configuration.FOREIGNS]

        if foreigns:
            # get len of suffixes in order to extract sub configurable names
            lenconfsuffix = len(ConfigurableRegistry.CONFIGURABLE_SUFFIX)
            lenconftsuffix = len(ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX)

            # for all parameters among foreign parameters
            for parameter in foreigns:
                # if name matches with a configurable name
                if parameter.name.endswith(
                        ConfigurableRegistry.CONFIGURABLE_SUFFIX
                ):
                    name = parameter.name[:-lenconfsuffix]
                    # try update it
                    self._configurables[name] = parameter.value

                # if name matches with a configurable type name
                elif parameter.name.endswith(
                        ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX
                ):
                    name = parameter.name[:-lenconftsuffix]
                    # try to update it
                    self._configurable_types[name] = parameter.value

    def _is_local(self, to_configure, name, *args, **kwargs):

        result = super(ConfigurableRegistry, self)._is_local(
            to_configure, name, *args, **kwargs
        )

        if not result:

            result = (
                name.endswith(ConfigurableRegistry.CONFIGURABLE_SUFFIX)
                or name.endswith(ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX)
            )

        return result

    @property
    def configurables(self):
        """Configurable which manages sub-configurables.
        """

        return self._configurables

    @property
    def configurable_types(self):
        """ConfigurableTypes which manages restriction of sub-configurable
        types.
        """

        return self._configurable_types

    def __contains__(self, name):
        """Redirection to self.configurables.__contains__.
        """

        if name.endswith(ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX):
            return name in self._configurable_types

        return name in self._configurables

    def __getitem__(self, name):
        """Redirection to self.configurables.__getitem__.
        """

        if name.endswith(ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX):
            return self._configurables_types[name]

        return self._configurables[name]

    def __setitem__(self, name, value):
        """Redirection to self.configurables.__setitem__.
        """

        if name.endswith(ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX):
            self._configurables_types[name] = value

        else:
            self._configurables[name] = value

    def __delitem__(self, name):
        """Redirection to self.configurables.__delitem__.
        """

        if name.endswith(ConfigurableRegistry.CONFIGURABLE_TYPE_SUFFIX):
            del self._configurables_types[name]

        else:
            del self._configurables[name]

    def __iter__(self):
        """Redirection to iter(self.configurables).
        """
        return iter(self._configurables)

    @staticmethod
    def get_configurable_category(name):
        """Get generated sub-configurable category name.
        """

        return "{0}_CONF".format(name.upper())

    @staticmethod
    def get_configurable(configurable, *args, **kwargs):
        """Get a configurable instance from a configurable class/path/instance
        and args, kwargs, None otherwise.

        :param configurable: configurable path, class or instance
        :type configurable: str, class or Configurable

        :return: configurable instance or None if input configurable can not be
        solved such as a configurable.
        """

        result = configurable

        if isinstance(configurable, basestring):
            result = lookup(configurable)

        if issubclass(result, Configurable):
            result = result(*args, **kwargs)

        if not isinstance(result, Configurable):
            result = None

        return result
