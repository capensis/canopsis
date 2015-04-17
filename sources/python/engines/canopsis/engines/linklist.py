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

from canopsis.engines.core import Engine
from datetime import datetime
from time import time
from json import loads
from canopsis.linklist.manager import Linklist
from canopsis.context.manager import Context
import pprint
pp = pprint.PrettyPrinter(indent=4)
# TODO remove
from canopsis.old.storage import get_storage
s = get_storage().get_backend('events')


class engine(Engine):
    etype = 'linklist'
    FILTER = 'mfilter'
    LINKS = 'filterlink'
    event_projection = {
        'resource': 1,
        'source_type': 1,
        'component': 1,
        'connector_name': 1,
        'connector': 1,
        'event_type': 1,
    }

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)
        self.context = Context()
        self.manager = Linklist()
        self.logger.debug = self.logger.info

    def consume_dispatcher(self, event, *args, **kargs):
        pass

    def pre_run(self):
        from time import sleep
        sleep(1)
        self.beat()

    def beat(self):

        links = {}

        # Computes links for all context elements
        # may cost some memory depending on filters and context size
        for linklist in self.manager.find():

            # condition to proceed a list link is they must be set
            name = linklist['name']
            l_filter = linklist.get(self.FILTER)
            l_list = linklist.get(self.LINKS)

            self.logger.debug('proceed linklist {}'.format(name))

            if not l_list or not l_filter:
                self.logger.info('Cannot proceed linklist for {}'.format(name))
            else:
                # Find context element ids matched by filter
                context_ids = self.get_ids_for_filter(l_filter)

                # Add all linklist to matched context element
                for context_id in context_ids:
                    if context_id not in links:
                        links[context_id] = []

                    # Append all links/labels to the context element
                    links[context_id] += l_list

        self.logger.debug('links')
        self.logger.debug(pp.pformat(links))

        entities = self.context.get_entities(links.keys())

        for entity in entities:
            self.update_context_with_links(
                entity,
                links[entity['_id']]
            )

    def update_context_with_links(self, entity, links):

        self.logger.debug(' + entity')
        self.logger.debug(entity)
        self.logger.debug(' + links')
        self.logger.debug(links)

        # element initialization
        if 'links' not in entity:
            entity['links'] = {}

        entity['links']['computed_links'] = links

        self.context.put(entity['type'], entity)

    def get_ids_for_filter(self, l_filter):

        context_ids = []
        a = l_filter
        try:
            l_filter = loads(l_filter)
        except Exception as e:
            self.logger.error(
                'Unable to parse mfilter, query aborted {}'.format(e)
            )
            return context_ids

        events = s.find(l_filter, self.event_projection)
        self.logger.debug('{} elements matches filter {}'.format(
            events.count(),
            a
        ))

        for event in events:
            entity = self.context.get_entity(event)
            entity_id = self.context.get_entity_id(entity)
            context_ids.append(entity_id)

        return context_ids
