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
from json import dumps


class Stats(MiddlewareRegistry):

    """
    Manage stats in Canopsis
    """

    def __init__(self, *args, **kwargs):

        super(Stats, self).__init__(self, *args, **kwargs)
        self.set_perf_data_array([])

        self.event_manager = Event()
        self.session_manager = Session()

    def set_perf_data_array(self, perf_data_array):

        """
        Property perf_data_array setter

        :param: value perf_data_array set the perf_data_array value
        """

        self.perf_data_array = perf_data_array

    def get_perf_data_array(self):

        """
        Property perf_data_array getter
        """

        return self.perf_data_array

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

        perf_data_array = []

        value = None

        # When alert and event was not in alert
        if is_alert and not was_alert:
            # Publish increment new_alarm count
            value = 1

        if was_alert and not is_alert:
            # Publish decrement new_alarm count
            value = -1

        self.logger.debug('Alerts count {}'.format(value))

        if value is not None:
            perf_data_array.append({
                'metric': 'alert_count',
                'value': value,
                'type': 'COUNTER'
            })

        return perf_data_array

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

    def compute_stats(self, event, devent, new_event):

        """
        Computes general statistics that are computable
        from event and it's previous state

        :param event: the current event within canopsis engines
        :param devent: previous event state from event storage
        :param new_event: boolean information True if event is new

        :returns: a new event to be published if any metrics
        """
        perf_data_array = []
        perf_data_array += self.new_alert_event_count(
            event, devent
        )
        perf_data_array += self.compute_ack_alerts(event, devent)
        perf_data_array += self.compute_by_states_and_sources(
            event,
            devent,
            new_event
        )
        self.logger.info('compute stats generated perfdata')
        self.logger.info(dumps(perf_data_array, indent=2))

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

    def compute_by_states_and_sources(self, event, devent, new_event):

        metrics = []

        metrics += self.event_add_by_source(event, new_event)
        metrics += self.event_add_by_state(event, new_event)
        metrics += self.event_count_by_source_and_state(
            event,
            devent,
            new_event
        )

        return metrics

    def compute_ack_alerts(self, event, devent):

        """
        Compute ack [solved] alerts metrics
        :param: event is amqp message
        :param: devent is previous state
        """
        perf_data_array = []

        # Compute alert stats and publish metrics if any
        new_alert = self.new_alert_event_count(
            event,
            devent
        )
        self.logger.info('New alert metric {}'.format(new_alert))
        if new_alert:

            solved = new_alert[0]['value'] == -1
            self.logger.debug('solved alert {}'.format(solved))
            if solved:
                # Compute alert ack depending on is ack
                metric = self.solved_alarm_ack(devent)
                perf_data_array.append(metric)

        return perf_data_array

    def users_session_duration(self):

        """
        Produce user session duration statistics from the session manager
        """

        sessions = self.session_manager.get_new_inactive_sessions()
        metrics = self.session_manager.get_delta_session_time_metrics(sessions)
        self.perf_data_array += metrics

    def event_add_by_source(self, event, new_event):

        """
        Counts and produces metrics for events depending on source type
        """
        metrics = []

        if new_event:
            metrics.append({
                'metric': 'cps_count_{}'.format(event['source_type']),
                'value': 1,
                'type': 'COUNTER'
            })

        return metrics

    def event_count_by_source_and_state(self, event, devent, new_event):

        """
        Counts and produces metrics for events depending on source type,
        by state
        """

        metrics = []

        if not new_event and devent['state'] != event['state']:

            state_src = self.event_manager.states_str[devent['state']]
            state_dest = self.event_manager.states_str[event['state']]

            metrics.append({
                'metric': 'cps_states_{}_{}'.format(
                    event['source_type'],
                    state_src
                ),
                'value': -1,
                'type': 'COUNTER'
            })

            metrics.append({
                'metric': 'cps_states_{}_{}'.format(
                    event['source_type'],
                    state_dest
                ),
                'value': 1,
                'type': 'COUNTER'
            })

        return metrics

    def event_add_by_state(self, event, new_event):

        """
        Counts and produces metrics for events depending on state
        """

        metrics = []

        if new_event:
            state_str = self.event_manager.states_str[event['state']]
            metrics.append({
                'metric': 'cps_states_{}'.format(state_str),
                'value': 1,
                'type': 'COUNTER'
            })

        return metrics
