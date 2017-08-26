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
from __future__ import unicode_literals

from canopsis.engines.core import Engine
from canopsis.entitylink.manager import Entitylink
from canopsis.linklist.dbconfigurationmanager import DBConfiguration


class engine(Engine):

    etype = 'linklist'

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)
        self.entity_link_manager = Entitylink()
        self.db_configuration_manager = DBConfiguration()
        self.link_field = []
        self.reload_configuration()

    def beat(self):
        self.reload_configuration()

    def reload_configuration(self):
        configuration = self.db_configuration_manager.find_one(
            query={'crecord_type': 'linklistfieldsurl'}
        )
        if configuration is not None:
            self.link_field = configuration['field_list']
            self.logger.debug(
                'searching for urls in events fields : {}'.format(
                    self.link_field
                )
            )

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
