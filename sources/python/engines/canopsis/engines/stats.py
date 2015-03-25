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

from canopsis.engines.core import Engine, publish
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.event import forger


class engine(Engine):

    etype = 'stats'

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
        self.storage = get_storage(
            namespace='events',
            account=Account(
                user="root",
                group="root"
            )
        )

    def consume_dispatcher(self, event, *args, **kargs):
        self.logger.debug('Entered in stats consume dispatcher')

        self.perf_data_array = []

        self.compute_states()

        self.publish_states()

    def compute_states(self):

        # Event count computation by state
        for state in [0, 1, 2, 3]:
            # There is an index on state field in events collection,
            # that would keep the query efficient
            count = self.storage.get_backend().find(
                {'state': state}
            ).count()

            state_str = self.states_str[state]

            self.perf_data_array.append({
                'metric': 'cps_states_{}'.format(state_str),
                'value': count
            })

        for source_type in ['resource', 'component']:
            # Event count source type
            count = self.storage.get_backend().find(
                {'source_type': source_type}
            ).count()

            self.perf_data_array.append({
                'metric': 'cps_count_{}'.format(
                    source_type
                ),
                'value': count
            })

            # Event count by source type and state
            for state in [0, 1, 2, 3]:

                # There is an index on state and source_type field in
                # events collection, that would keep the query efficient
                count = self.storage.get_backend().find(
                    {'source_type': source_type, 'state': state}
                ).count()

                state_str = self.states_str[state]

                self.perf_data_array.append({
                    'metric': 'cps_states_{}_{}'.format(
                        source_type,
                        state_str
                    ),
                    'value': count
                })

        self.logger.debug('perf_data_array {}'.format(self.perf_data_array))

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
