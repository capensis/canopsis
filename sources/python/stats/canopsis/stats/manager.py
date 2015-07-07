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
from canopsis.event import forger


class Stats(MiddlewareRegistry):

    event_manager = Event()
    # alias for easier testing purposes
    """
    Manage stats in Canopsis
    """

    def new_alert_event_count(self, event, devent):

        """
        Produce metric count for alert count.
        cps_new_alert +1 or -1 depends on previous event state

        :param: event is the current event passing through canopsis input
        :param: devent is the event from database matching event RK.
        it is the last state the event were
        """

        perf_data_array = []
        is_alert = self.event_manager.is_alert(event['state'])
        was_alert = self.event_manager.is_alert(devent['state'])

        # When alert
        if is_alert:
            # and event was not in alert
            if not was_alert:
                # Publish increment new_alarm count
                perf_data_array.append({
                    'metric': 'cps_new_alert',
                    'value': 1,
                    'type': 'COUNTER'
                })
        elif was_alert:

            if not is_alert:

                # Publish decrement new_alarm count
                perf_data_array.append({
                    'metric': 'cps_new_alert',
                    'value': -1,
                    'type': 'COUNTER'
                })

        # Do not generate event if there is no metric
        if len(perf_data_array) == 0:
            return None

        stats_event = forger(
            connector="Engine",
            connector_name='stats',
            event_type="perf",
            source_type="component",
            component="__canopsis__",
        )

        stats_event['perf_data_array'] = perf_data_array

        return stats_event
