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
from canopsis.common.utils import lookup

from j1bz.expression import parse_dsl, DSLException
from b3j0f.requester.request.crud import Read
from link.utils.filter import Filter
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
    ROUTING_STORAGE = 'routing_storage'
    SYSREQ = 'sysreq'

    DEFAULT_EXCHANGE = 'canopsis.queues'
    DEFAULT_CONFIG = {
        'default_action': 'pass',
        'at_least': 0
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
                query={'crecord_type': 'defaultrule'}
            )

        if value:
            value = value[0]

        else:
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
            value = self[RouterManager.CONFIG_STORAGE].get_elements(
                query={'crecord_type': 'filter', 'enabled': True},
                sort='priority'
            )

        self._filters = value

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

    def match_filters(self, event, publisher=None):
        does_match = True

        for doc in self.filters:
            event_filter = Filter(json.loads(doc['filter']))

            if event_filter.match(event):
                self.logger.debug(u'Event {0} matched filter {1}'.format(
                    event['rk'],
                    doc['crecord_name']
                ))

                event = self.apply_actions(
                    event,
                    doc['actions'],
                    publisher=publisher
                )

                if event is None:
                    # Event dropped
                    break

                if doc.get('break', False):
                    self.logger.debug(u'Filter {0} breaking chain'.format(
                        doc['crecord_name']
                    ))

                    break

            requests = doc.get('sysreq', [])

            if requests:
                nmatch = 0

                for request in requests:
                    self.logger.debug(u'Request: {0}'.format(request))

                    try:
                        cruds = parse_dsl(request)

                    except DSLException as err:
                        self.logger.error(
                            u'An error occured while parsing request: {0}'
                            .format(err)
                        )
                        cruds = None

                    if cruds is not None:
                        sysreq = self[RouterManager.SYSREQ]
                        reads = [
                            crud
                            for crud in cruds
                            if isinstance(crud, Read)
                        ]

                        ctx = sysreq(cruds)

                        for read in reads:
                            if ctx[read]:
                                nmatch += 1
                                break

                if nmatch <= self.config['at_least']:
                    does_match = False
                    break

        if event is not None and not does_match:
            self.logger.debug(u'Event {0} did not match any filter'.format(
                event['rk']
            ))

            event = self.apply_actions(event, {
                'type': self.config['default_action']
            })

        return event

    def apply_actions(self, event, actions, publisher=None):
        for action in actions:
            atype = action.pop('type')

            self.logger.debug(u'Apply action {0} on event {1}'.format(
                atype,
                event['rk']
            ))

            try:
                handler = lookup('router.actions.{0}'.format(atype))

            except ImportError:
                handler = None

            if not callable(handler):
                self.logger.error(u'Unknown action {0}, ignoring'.format(
                    atype
                ))

            else:
                event = handler(event, publisher=publisher, **action)

            if event is None:
                # Event dropped
                break

        return event

    def get_routing_key(self, event):
        for doc in self.routes:
            event_filter = Filter(json.loads(doc['filter']))

            if event_filter.match(event):
                self.logger.debug(u'Event {0} matched route {1}'.format(
                    event['rk'],
                    doc['rk']
                ))

                return doc['rk']

        return None

    def reload(self):
        self.logger.debug('Reload filters and routes')

        self.filters = None
        self.routes = None
