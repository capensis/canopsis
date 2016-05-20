# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.model import Parameter

from link.utils.filter import Filter
from jsonpatch import JsonPatch
import json


CONF_PATH = 'router/manager.conf'
CATEGORY = 'ROUTER'
CONTENT = [
    Parameter('exchange')
]


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class RouterManager(MiddlewareRegistry):

    CONFIG_STORAGE = 'config_storage'
    FILTER_STORAGE = 'filter_storage'
    PATCH_STORAGE = 'patch_storage'
    ROUTING_STORAGE = 'routing_storage'

    DEFAULT_EXCHANGE = 'canopsis.queues'
    DEFAULT_CONFIG = {
        'is_blacklist': False
    }

    @property
    def exchange(self):
        if not hasattr(self, '_exchange'):
            self.exchange = None

        return self._exchange

    @exchange.setter
    def exchange(self, value):
        if value is None:
            value = RouterManager.DEFAULT_EXCHANGE

        self._exchange = value

    @property
    def config(self):
        if not hasattr(self, '_config'):
            self.config = None

        return self._config

    @config.setter
    def config(self, value):
        if value is None:
            value = self[RouterManager.CONFIG_STORAGE].get_elements(
                ids=self.__class__.__name__.lower()
            )

        if value is None:
            value = RouterManager.DEFAULT_CONFIG

        self._config = value

    @property
    def filters(self):
        if not hasattr(self, '_filters'):
            self.filters = None

        return self._filters

    @filters.setter
    def filters(self, value):
        if value is None:
            value = self[RouterManager.FILTER_STORAGE].get_elements()

        self._filters = value

    @property
    def patchs(self):
        if not hasattr(self, '_patchs'):
            self.patchs = None

        return self._patchs

    @patchs.setter
    def patchs(self, value):
        if value is None:
            value = self[RouterManager.PATCH_STORAGE].get_elements()

        self._patchs = value

    @property
    def routes(self):
        if not hasattr(self, '_routes'):
            self.routes = None

        return self._routes

    @routes.setter
    def routes(self, value):
        if value is None:
            value = self[RouterManager.ROUTING_STORAGE].get_elements()

        self._routes = value

    def match_filters(self, event):
        does_match = True
        blacklist = self.config['is_blacklist']

        for doc in self.filters:
            event_filter = Filter(json.loads(doc['filter']))

            conditions = [
                blacklist and event_filter.match(event),
                not blacklist and not event_filter.match(event)
            ]

            if any(conditions):
                does_match = False
                break

        return does_match

    def apply_patchs(self, event):
        for doc in self.patchs:
            event_filter = Filter(json.loads(doc['filter']))

            if event_filter.match(event):
                patch = JsonPatch(json.loads(doc['patch']))
                event = patch.apply(event)

        return event

    def get_routing_key(self, event):
        for doc in self.routes:
            event_filter = Filter(json.loads(doc['filter']))

            if event_filter.match(event):
                return doc['rk']

        return None

    def reload(self):
        self.filters = None
        self.patchs = None
        self.routes = None
