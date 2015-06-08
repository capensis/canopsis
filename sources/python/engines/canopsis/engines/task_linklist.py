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

from canopsis.engines.core import TaskHandler
from canopsis.linklist.manager import Linklist
from canopsis.context.manager import Context
from canopsis.event.manager import Event
from canopsis.entitylink.manager import Entitylink
from json import loads


class engine(TaskHandler):

    etype = 'tasklinklist'

    event_projection = {
        'resource': 1,
        'source_type': 1,
        'component': 1,
        'connector_name': 1,
        'connector': 1,
        'event_type': 1,
    }

    def handle_task(self, job):

        self.link_list_manager = Linklist()
        self.context = Context()
        self.event = Event()
        self.entity_link_manager = Entitylink()

        """
        This task computes all links associated to an entity.
        Link association are managed by entity link system.
        """

        for entity_id in self.context.iter_ids():
            self.entity_link_manager.put(entity_id, {
                'computed_links': []
            })

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
        self.logger.debug(links)

        entities = self.context.get_entities(links.keys())

        for entity in entities:
            self.update_context_with_links(
                entity,
                links[entity['_id']]
            )

        return (0, 'Link list computation complete')

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
