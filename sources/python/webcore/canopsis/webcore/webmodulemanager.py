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

from canopsis.middleware.registry import MiddlewareRegistry
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_config
from canopsis.configuration.model import ParamList, Parameter


CONF_PATH = 'webserver.conf'
CONFIG = {
    'WEBMODULEMANAGER': [],
    'webmodules': ParamList(parser=Parameter.bool)
}


@conf_paths(CONF_PATH)
@add_config(CONFIG)
class WebModuleManager(MiddlewareRegistry):
    """
    Manage list of web modules.
    """

    CONFIG_STORAGE = 'config_storage'

    @property
    def webmodules(self):
        if not hasattr(self, '_webmodules'):
            self._webmodules = {}

        return self._webmodules

    def __init__(self, config_storage=None, *args, **kwargs):
        super(WebModuleManager, self).__init__(*args, **kwargs)

        if config_storage is not None:
            self[WebModuleManager.CONFIG_STORAGE] = config_storage

    @property
    def config(self):
        storage = self[WebModuleManager.CONFIG_STORAGE]
        config = storage.get_elements(ids='enabledmodules')
        return config

    @property
    def modules(self):
        return self.config['enabled']

    @modules.setter
    def modules(self, value):
        storage = self[WebModuleManager.CONFIG_STORAGE]
        storage.put_element(
            element={'enabled': value},
            _id='enabledmodules'
        )

    def enable_module(self, name):
        """
        Enable a module if not already registered.

        :param name: module's name
        :type name: str

        :returns: True if module was enabled, False otherwise.
        """

        modules = self.modules

        if name not in modules:
            modules.append(name)
            self.modules = modules

            return True

        return False

    def disable_module(self, name):
        """
        Disable a module if not already unregistered.

        :param name: module's name
        :type name: str

        :returns: True if module was disabled, False otherwise.
        """

        modules = self.modules

        if name in modules:
            modules.remove(name)
            self.modules = modules

            return True

        return False

    def init_modules(self):
        self.modules = [
            module
            for module in self.webmodules
            if self.webmodules[module]
        ]

    def clear_modules(self):
        self.modules = []
