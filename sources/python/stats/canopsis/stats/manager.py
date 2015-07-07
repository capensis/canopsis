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
from canopsis.event.manager import Event
from canopsis.session.manager import Session
from canopsis.event import forger


class Stats(MiddlewareRegistry):

    event_manager = Event()
    session_manager = Session()

    """
    Manage stats in Canopsis
    """

    def __init__(self, *args, **kwargs):

        super(Stats, self).__init__(self, *args, **kwargs)
        self.perf_data_array = []

    def new_alert_event_count(self, event, devent):

        """
        Produce metric count for alert count.
        cps_new_alert +1 or -1 depends on previous event state

        :param: event is the current event passing through canopsis input
        :param: devent is the event from database matching event RK.
        it is the last state the event were
        """

        is_alert = self.event_manager.is_alert(event['state'])
        was_alert = self.event_manager.is_alert(devent['state'])

        metric = None
        # When alert
        if is_alert:
            # and event was not in alert
            if not was_alert:
                # Publish increment new_alarm count
                metric = {
                    'metric': 'cps_new_alert',
                    'value': 1,
                    'type': 'COUNTER'
                }
        elif was_alert:

            if not is_alert:

                # Publish decrement new_alarm count
                metric = {
                    'metric': 'cps_new_alert',
                    'value': -1,
                    'type': 'COUNTER'
                }

        self.logger.debug('new alert \n{}'.format(metric))

        return metric

    def solved_alarm_ack(self, devent):

        # Then produce metric
        was_ack = self.event_manager.is_ack(devent)
        if was_ack:
            metric_name = 'cps_solved_ack_alarms'
        else:
            metric_name = 'cps_solved_not_ack_alarms'

        metric = {
            'metric': metric_name,
            'value': 1,
            'type': 'COUNTER'
        }

        self.logger.debug('solved alarm ack \n{}'.format(metric))

        return metric

    def compute_ack_alerts(self, event, devent):

        """
        Compute ack [solved] alerts metrics
        :param: event is amqp message
        :param: devent is previous state
        """
        perf_data_array = []

        # Compute alert stats and publish metrics if any
        metric = self.new_alert_event_count(
            event,
            devent
        )
        if metric:

            perf_data_array.append(metric)

            solved = metric['value'] == -1
            if solved:
                # Compute alert ack depending on is ack
                metric = self.solved_alarm_ack(
                    devent
                )
                perf_data_array.append(metric)

        if len(perf_data_array):
            stats_event = forger(
                connector="Engine",
                connector_name='stats',
                event_type="perf",
                source_type="component",
                component="__canopsis__",
                perf_data_array=perf_data_array
            )
            return stats_event

    def add_metric(self, mname, mvalue, mtype='COUNTER'):

        """
        Add metric to the manager perf_data array property
        :param: mname is the metric name
        :param: mvalue is the metric value
        :param: mtype is the metric type
        """

        self.perf_data_array.append({
            'metric': mname,
            'value': mvalue,
            'type': mtype
        })

    def users_session_duration(self):

        """
        Produce user session duration statistics from the session manager
        """

        sessions = self.session_manager.get_new_inactive_sessions()
        metrics = self.session_manager.get_delta_session_time_metrics(sessions)
        self.perf_data_array += metrics

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

                state_str = self.event_manager.states_str[state]

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

            state_str = self.event_manager.states_str[state]

            self.add_metric(
                'cps_states_{}'.format(state_str),
                count
            )
