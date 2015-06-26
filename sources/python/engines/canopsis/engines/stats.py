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

from canopsis.event.manager import Event
from canopsis.organisation.rights import Rights
from canopsis.engines.core import Engine, publish
from canopsis.event import forger


class engine(Engine):

    etype = 'stats'

    event_manager = Event()
    right_manager = Rights()

    """
    This engine's goal is to compute canopsis internal statistics.
    Statistics are computed on each passing event and are updated
    in async way in order to manage performances issues.
    """

    def __init__(self, *args, **kargs):

        super(engine, self).__init__(*args, **kargs)

        self.states_str = {
            0: 'info',
            1: 'minor',
            2: 'major',
            3: 'critical'
        }

    def beat(self):
        users = right_manager.get_users()
        self.userlist = [user['_id'] for user in list(users)]

    def consume_dispatcher(self, event, *args, **kargs):
        self.logger.debug('Entered in stats consume dispatcher')

        self.perf_data_array = []

        self.compute_states()

        self.publish_states()

    def compute_states(self):

        # Allow individual stat computation management from ui.
        stats_to_compute = [
            'event_count_by_source',
            'event_count_by_source_and_state',
            'event_count_by_state',
            'ack_alerts_by_user'
        ]

        for stat in stats_to_compute:
            method = getattr(self, 'stats_to_compute')
            method()

        self.logger.debug('perf_data_array {}'.format(self.perf_data_array))

    def event_count_by_source(self):

        for source_type in ('resource', 'component'):
            # Event count source type
            cursor, count = self.event_manager.find(
                query={'source_type': source_type},
                with_count=True
            )

            self.perf_data_array.append({
                'metric': 'cps_count_{}'.format(
                    source_type
                ),
                'value': count
            })

    def event_count_by_source_and_state(self):

        for source_type in ('resource', 'component'):

            # Event count by source type and state
            for state in (0, 1, 2, 3):

                # There is an index on state and source_type field in
                # events collection, that would keep the query efficient
                cursor, count = self.event_manager.find(
                    query={
                        'source_type': source_type,
                        'state': state
                    }
                    with_count=True
                )

                state_str = self.states_str[state]

                self.perf_data_array.append({
                    'metric': 'cps_states_{}_{}'.format(
                        source_type,
                        state_str
                    ),
                    'value': count
                })

    def event_count_by_state(self):

        # Event count computation by state
        for state in (0, 1, 2, 3):
            # There is an index on state field in events collection,
            # that would keep the query efficient
            cursor, count = self.event_manager.find(
                query={'state': state},
                with_count=True
            )

            state_str = self.states_str[state]

            self.perf_data_array.append({
                'metric': 'cps_states_{}'.format(state_str),
                'value': count
            })

    def ack_alerts_by_user(self):

        # may be improved to get this stats by domain/perimeter
        for user in self.userlist:
            cursor, count = self.event_manager.find(
                query={
                    'ack.author': user,
                    'ack.isAck': True
                },
                with_count=True
            )

            self.perf_data_array.append({
                'type': 'COUNTER',
                'metric': 'cps_states_ack_alerts_by_user_{}'.format(user)
                'value': count
            })

    def publish_states(self):

        stats_event = forger(
            connector='engine',
            connector_name='engine',
            event_type='perf',
            source_type='resource',
            resource='Engine_stats',
            state=0,
            perf_data_array=self.perf_data_array
        )

        self.logger.debug('Publishing {}'.format(stats_event))

        publish(publisher=self.amqp, event=stats_event)
