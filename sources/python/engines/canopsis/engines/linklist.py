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
from canopsis.entitylink.manager import Entitylink
from canopsis.context.manager import Context
from canopsis.event.manager import Event

import pprint
pp = pprint.PrettyPrinter(indent=4)


class engine(Engine):

    etype = 'linklist'

    event_projection = {
        'resource': 1,
        'source_type': 1,
        'component': 1,
        'connector_name': 1,
        'connector': 1,
        'event_type': 1,
    }

    # Search fields in the event for possible url information
    link_field = [
        'action_url'
    ]

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)
        #TODO REMOVE BEAT INERVAL AND SWITCH TO SCHEDULED JOB THE BEAT PROCESS
        self.beat_interval = 3600

        self.context = Context()
        self.event = Event()
        self.link_list_manager = Linklist()
        self.entity_link_manager = Entitylink()

    def work(self, event, *args, **kwargs):
        """
        Check if event contains any url interesting keys.
        When one of self.link_field found, it is added to the event entity link
        """

        link_key = 'event_links'

        event_links = {}
        for key in self.link_field:
            if key in event and event[key]:
                event_links[key] = event[key]

        if event_links:

            self.logger.debug('found links into the event {}: {}'.format(
                event['rk'],
                event_links
            ))

            # Proceed link update in case some link found in the event
            entity = self.entity_link_manager.get_or_create_from_event(event)
            _id = entity['_id']

            links_to_integrate = {}

            # integrate previous links if any
            if link_key in entity:
                for link in entity[link_key]:
                    links_to_integrate[link['label']] = link['url']

            # New event link definition upsert (override)
            links_to_integrate.update(event_links)

            # Define new data context as in upsert mode
            context = {link_key: []}
            for key in links_to_integrate:
                context[link_key].append({
                    'label': key,
                    'url': links_to_integrate[key]
                })

            # Push changes to db
            self.entity_link_manager.put(_id, context)

    def beat(self):

        """
        This task computes all links associated to an entity.
        Link association are managed by entity link system.
        """

        links = {}

        # Computes links for all context elements
        # may cost some memory depending on filters and context size
        for linklist in self.link_list_manager.find():

            # condition to proceed a list link is they must be set
            name = linklist['name']
            l_filter = linklist.get('mfilter')
            l_list = linklist.get('filterlink')

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

        """
        Upsert computed links to the entity link storage
        """

        self.logger.debug(' + entity')
        self.logger.debug(entity)
        self.logger.debug(' + links')
        self.logger.debug(links)

        context = {
            'computed_links': links
        }

        _id = self.context.get_entity_id(entity)

        self.entity_link_manager.put(_id, context)

    def get_ids_for_filter(self, l_filter):

        """
        Retrieve a list of id from event collection.
        Can be performance killer as matching mfilter
        is only available on the event collection at the moment
        """

        context_ids = []

        try:
            l_filter = loads(l_filter)
        except Exception as e:
            self.logger.error(
                'Unable to parse mfilter, query aborted {}'.format(e)
            )
            return context_ids

        events = self.event.find(
            query=l_filter,
            projection=self.event_projection
        )

        for event in events:
            self.logger.debug('rk : {}'.format(event['_id']))
            entity = self.context.get_entity(event)
            entity_id = self.context.get_entity_id(entity)
            context_ids.append(entity_id)

        return context_ids
