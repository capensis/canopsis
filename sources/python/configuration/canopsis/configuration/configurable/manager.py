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

from canopsis.configuration import Configurable
from canopsis.configuration.parameters import Configuration, Category
from canopsis.common.utils import resolve_element

from inspect import isclass


class Manager(Configurable):
    """
    Manage a set of configurables which are accessibles
    from self.configurables.

    Each configurable can be defined in conf parameters where names are like
    {name}_configurable={configurable_path, configurable_class, configurable}

    And a configurable configuration are in categories {NAME}_CONF.
    """

    class Error(Exception):
        """handle manager errors"""
        pass

    CONF_PATH = 'configuration/manager.conf'

    CATEGORY = 'MANAGER'

    CONFIGURABLE_SUFFIX = '_configurable'
    CONFIGURABLE_TYPE_SUFFIX = '_configurable_type'

    def __init__(
        self, configurables=None, configurable_types=None, *args, **kwargs
    ):
        """
        :param configurables: dictionary of configurables by name.
        :type configurables: dict

        :param configurable_types: dictionary of configurable types by name
        :type configurable_types: dict
        """

        super(Manager, self).__init__(*args, **kwargs)

        self.configurables = {} if configurables is None else configurables
        self.configurable_types = {} if configurable_types is None \
            else configurable_types

    def _get_category(self):
        """
        Get category.
        """

        result = Category(Manager.CATEGORY)

        return result

    def _conf(self, *args, **kwargs):

        result = super(Manager, self)._conf(*args, **kwargs)

        result.add_unified_category(name=Manager.CATEGORY)

        return result

    def _get_conf_paths(self, *args, **kwargs):

        result = super(Manager, self)._get_conf_paths(*args, **kwargs)

        result.append(Manager.CONF_PATH)

        return result

    def apply_configuration(
        self, conf=None, conf_paths=None, managers=None, logger=None
    ):

        super(Manager, self).apply_configuration(
            conf=conf, conf_paths=conf_paths, managers=managers, logger=logger)

        if conf_paths is None:
            conf_paths = self.conf_paths

        if conf is None:
            conf = self.conf

        if managers is None:
            managers = self.managers

        # get self conf path
        conf_path = self.conf_paths[-1]

        # apply configuration to all self configurables
        for name, configurable in self.configurables.iteritems():
            # add self last conf paths to configurable conf paths
            configurable_conf_paths = list(configurable.conf_paths)
            configurable_conf_paths.append(conf_path)
            # get a copy of configurable configuration
            configurable_configuration = configurable.conf.copy()
            # add a unified category where name is {NAME}_CONF
            category_name = Manager.get_configurable_category(name)
            configurable_configuration.add_unified_category(
                name=category_name, copy=True)
            # apply configurable configuration
            configurable.apply_configuration(
                conf=configurable_configuration,
                conf_paths=configurable_conf_paths,
                managers=managers, logger=logger)

    def _configure(self, unified_conf, *args, **kwargs):

        super(Manager, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        values = unified_conf[Configuration.VALUES]

        for parameter in values:
            if parameter.name.endswith(Manager.CONFIGURABLE_SUFFIX):
                self[parameter.name] = parameter.value

            elif parameter.name.endswith(Manager.CONFIGURABLE_TYPE_SUFFIX):
                self[parameter.name] = parameter.value

    @property
    def configurables(self):
        """
        Dictionary of configurable by attribute name
        """
        return self._configurables

    @configurables.setter
    def configurables(self, value):
        self._configurables = value

    @property
    def configurable_types(self):
        """
        Dictionary of configurable types by attribute name
        """
        return self._configurable_types

    @configurable_types.setter
    def configurable_types(self, value):
        self._configurable_types = value

    def __contains__(self, name):
        """
        Check whatever or not if a configurable name is in self configurables
        """

        if name.endswith(Manager.CONFIGURABLE_SUFFIX):
            name = name[:-len(Manager.CONFIGURABLE_SUFFIX)]

        return name in self.configurables

    def __getitem__(self, name):
        """
        Redirect to configurables if name in configurables.

        Example: let a configurable named foo.
        configurableManager = Manager()
        configurableManager.foo_configurable = Configurable()
        assert self.foo == self.foo_configurable
        """

        result = None

        # if name matches with a configurable name
        if name in self.configurables:
            # allow to matches self.foo_configurable with self.foo
            result = self.configurables[name]

        # if name matches a configurable, returns it
        elif name.endswith(Manager.CONFIGURABLE_SUFFIX):
            configurable_name = name[:-len(
                Manager.CONFIGURABLE_SUFFIX)]

            if configurable_name in self._configurables:
                result = self.configurables[configurable_name]

            else:
                raise Manager.Error(
                    "Configurable %s does not exist" % name)

        # if name matches a configurable type, returns it
        elif name.endswith(Manager.CONFIGURABLE_TYPE_SUFFIX):
            configurable_type_name = name[:-len(
                Manager.CONFIGURABLE_TYPE_SUFFIX)]

            if configurable_type_name in self._configurable_types:
                result = self._configurable_types[configurable_type_name]

            else:
                raise Manager.Error(
                    "Configurable class %s does not exist" % name)

        if result is None:
            raise AttributeError(
                "'%s' manager has no configurable '%s'" % (self, name))

        return result

    def __setitem__(self, name, value):
        """
        Check if value isinstance of registered configurable type if name in
        configurables.

        Example: let a configurable named foo.
        cm = Manager()
        cm.foo_configurable = Configurable()
        raisedError = False
        try:
            cm.foo = 1
        except Manager:
            raisedError = True
        assert raisedError
        """

        # if name matches to a Configurable name, modify it
        if name in self.configurables:
            name = '%s%s' % (name, Manager.CONFIGURABLE_SUFFIX)

        # if name corresponds to a Configurable
        if name.endswith(Manager.CONFIGURABLE_SUFFIX):

            name = name[:-len(Manager.CONFIGURABLE_SUFFIX)]

            configurable_type = self.configurable_types.get(
                name, Configurable)
            configurable = value

            # in case of value is str
            if isinstance(configurable, str):
                configurable = resolve_element(value)

            # in case of value is configurable type
            if isclass(configurable) \
                    and issubclass(configurable, configurable_type):
                configurable = configurable()

            # in case of value is an instance of configurable_type
            if isinstance(configurable, configurable_type):
                self._configurables[name] = configurable
            # else, raise an Error
            else:
                raise Manager.Error(
                        "Configurable %s:%s must be of type %s" % (
                            name, configurable, configurable_type))

        # if name corresponds to a configurable type
        elif name.endswith(Manager.CONFIGURABLE_TYPE_SUFFIX):

            name = name[:-len(Manager.CONFIGURABLE_TYPE_SUFFIX)]

            configurable_type = value

            if isinstance(value, str):
                configurable_type = resolve_element(value)

            if not issubclass(configurable_type, Configurable):
                raise Manager.Error(
                    "Configurable type %s %s must inherits from Configurable" %
                    (name, value))

            self._configurable_types[name] = configurable_type

            # check if old configurable inherits from new type
            if name in self._configurables:
                # if it does not, delete it from self.configurables
                configurable = self._configurables[name]
                if not isinstance(configurable, configurable_type):
                    self.logger.info(
                        "remove old incompatible configurable %s %s with %s" %
                            (name, configurable, configurable_type))
                    del self.configurables[name]

    def __delitem__(self, name):
        """
        Redirect to configurables if name in configurables.

        Example:
        cm = Manager()
        cm.foo_configurable = Configurable()
        del cm.foo
        raisedError = False
        try:
            del cm.foo_configurable
        except Manager.Error:
            raisedError = True
        assert raisedError
        """

        if name in self._configurables:
            name = '%s%s' % (name, Manager.CONFIGURABLE_SUFFIX)

        # if name matches a configurable name, modify it
        if name.endswith(Manager.CONFIGURABLE_SUFFIX):
            try:
                del self._configurables[name]
            except KeyError:
                raise Manager.Error(
                    "Configurable %s does not exist" % name)

        # if name matches a configurable type name, delete related type
        elif name.endswith(Manager.CONFIGURABLE_TYPE_SUFFIX):
            name = name[:-len(Manager.CONFIGURABLE_TYPE_SUFFIX)]
            if name in self._configurable_types:
                try:
                    del self._configurable_types[name]
                except KeyError:
                    raise Manager.Error(
                        "Configurable type %s does not exist" % name)

    @staticmethod
    def get_configurable_category(name):

        return "%s_CONF" % name.upper()

    @staticmethod
    def get_configurable(self, configurable, *args, **kwargs):
        """
        Get a configurable instance from a configurable class/path and
        args, kwargs

        :param configurable: configurable path or class
        :type configurable: str or Configurable
        """
        result = configurable

        if isinstance(configurable, str):
            result = resolve_element(configurable)

        if issubclass(result, Configurable):
            result = result(*args, **kwargs)

        return result
