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

import pprint
pp = pprint.PrettyPrinter(indent=4)


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

    def pre_run(self):

        self.beat()
        #TODO remove consume dispatcher call
        self.consume_dispatcher({})

    def consume_dispatcher(self, event, *args, **kargs):

        """
        Entry point for stats computation. Triggered by the dispatcher
        engine for distributed processing puroses.
        Following methods will generate metrics that are finally embeded
        in a metric event.
        """

        self.logger.debug('Entered in stats consume dispatcher')

        self.perf_data_array = []

        # Some queries may be used twice, so cache them for performance
        self.prepare_queries()

        self.compute_states()

        self.publish_states()

    def prepare_queries(self):

        """
        Stats are computed from database information. This methods
        let perform cached results queries that are available in
        all stats methods. Cached result should not be altered while 
        processing stats computation
        """

        users = self.right_manager.get_users()
        self.userlist = [user['_id'] for user in list(users)]

        # Query for current ack events
        cursor = self.event_manager.find(
            query={
                'ack.isAck': True
            },
            projection={
                'ack.author': 1,
                'ack.timestamp': 1,
                'last_state_change': 1,
                'domain': 1,
                'perimeter': 1
            }
        )
        self.ack_events = list(cursor)

    def compute_states(self):

        """
        Entry point for dynamic stats method triggering
        Dynamic triggering allow greated control on which
        stats are computed and can be activated/deactivated
        from frontend.
        """

        # Allow individual stat computation management from ui.
        stats_to_compute = [
            'event_count_by_source',
            'event_count_by_source_and_state',
            'event_count_by_state',
            'ack_alerts_by_user',
            'delta_alert_ack_by_user'
        ]

        for stat in stats_to_compute:
            method = getattr(self, stat)
            method()

    def add_metric(self, mname, mvalue, mtype='COUNTER'):
        self.perf_data_array.append({
            'metric': mname,
            'value': mvalue,
            'type': mtype
        })

    def delta_alert_ack_by_user(self):

        """
        Computes time sum between an alert and a user ack.
        metric is generated for each user.
        """

        metrics = {}
        for event in self.ack_events:

            ack = event.get('ack', {})
            last_state_change = event['last_state_change']
            ack_timestamp = ack.get('timestamp')
            user = ack.get('author')

            if last_state_change and ack_timestamp and user:
                # Start delta time aggregation by user
                delta = ack_timestamp - last_state_change
                if user not in metrics:
                    metrics[user] = 0
                metrics[user] += delta

            else:
                self.logger.warning(
                    'Stat delta_alert_ack_by_user error {} {} {}'.format(
                        last_state_change,
                        ack_timestamp,
                        user
                    )
                )

        for user in metrics:
            self.add_metric(
                'cps_delta_alert_ack_by_user_{}'.format(user),
                metrics[user]
            )
        self.add_metric(
            'cps_delta_alert_ack_all',
            sum(metrics.values())
        )

    def ack_alerts_by_user(self):

        """
        Counts how many alerts are ack by each user. It also
        depends on event domain and perimeter
        """

        metrics = {}
        for event in self.ack_events:
            ack = event.get('ack', {})
            user = ack.get('author')
            domain = event['domain']
            perimeter = event['perimeter']

            metric_name = 'cps_ack_alerts_by_user_{}_on_{}_{}'.format(
                user,
                domain,
                perimeter
            )

            if metric_name not in metrics:
                metrics[metric_name] = 0

            metrics[metric_name] += 1

        for metric_name in metrics:
            self.add_metric(
                metric_name,
                metrics[metric_name]
            )
        self.add_metric(
            'cps_ack_alerts_all',
            sum(metrics.values())
        )

    def event_count_by_source(self):

        """
        Counts and produces metrics for events depending on source type
        """

        for source_type in ('resource', 'component'):
            # Event count source type
            cursor, count = self.event_manager.find(
                query={'source_type': source_type},
                with_count=True
            )

            self.add_metric(
                'cps_count_{}'.format(
                    source_type
                ),
                count
            )

    def event_count_by_source_and_state(self):

        """
        Counts and produces metrics for events depending on source type,
        by state
        """

        for source_type in ('resource', 'component'):

            # Event count by source type and state
            for state in (0, 1, 2, 3):

                # There is an index on state and source_type field in
                # events collection, that would keep the query efficient
                cursor, count = self.event_manager.find(
                    query={
                        'source_type': source_type,
                        'state': state
                    },
                    with_count=True
                )

                state_str = self.states_str[state]

                self.add_metric(
                    'cps_states_{}_{}'.format(
                        source_type,
                        state_str
                    ),
                    count
                )

    def event_count_by_state(self):

        """
        Counts and produces metrics for events depending on state
        """

        # Event count computation by state
        for state in (0, 1, 2, 3):
            # There is an index on state field in events collection,
            # that would keep the query efficient
            cursor, count = self.event_manager.find(
                query={'state': state},
                with_count=True
            )

            state_str = self.states_str[state]

            self.add_metric(
                'cps_states_{}'.format(state_str),
                count
            )

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

        metrics = []
        for m in self.perf_data_array:
            metrics.append('{}\t{}\t{}'.format(
                m['value'],
                m['type'],
                m['metric']
            ))

        self.logger.debug('Generated perfdata\n{}'.format(
            '\n'.join(metrics)
        ))

        publish(publisher=self.amqp, event=stats_event)
