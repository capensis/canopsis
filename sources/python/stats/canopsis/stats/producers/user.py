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

from canopsis.stats.producers.metric import MetricProducer
from canopsis.configuration.configurable.decorator import add_category
from canopsis.configuration.configurable.decorator import conf_paths


CONF_PATH = 'stats/producers/user.conf'
CATEGORY = 'USER_METRIC_PRODUCER'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class UserMetricProducer(MetricProducer):
    def __init__(self, *args, **kwargs):
        super(UserMetricProducer, self).__init__(*args, **kwargs)

    def alarm_ack(self, event, user):
        event = {
            'connector': 'canopsis',
            'connector_name': 'stats',
            'event_type': 'perf',
            'source_type': 'resource',
            'component': user,
            'resource': 'alarm_ack',
            'perf_data_array': [
                {
                    'metric': filtername,
                    'value': 1,
                    'type': 'COUNTER'
                }
                for filtername in self.match(event)
            ]
        }

        return event

    def _delay(self, user, name, value):
        event = {
            'connector': 'canopsis',
            'connector_name': 'stats',
            'event_type': 'perf',
            'source_type': 'resource',
            'component': user,
            'resource': name,
            'perf_data_array': [
                {
                    'metric': 'sum',
                    'value': value,
                    'type': 'COUNTER'
                },
                {
                    'metric': 'last',
                    'value': value,
                    'type': 'GAUGE'
                }
            ]
        }

        for operator in ['min', 'max', 'average']:
            entity = self[MetricProducer.PERFDATA_MANAGER].get_metric_entity(
                'last', event
            )

            self.may_create_stats_serie(entity, operator)

        return event

    def alarm_ack_delay(self, user, delay):
        return self._delay(user, 'alarm_ack_delay', delay)

    def alarm_ack_solved(self, user, delay):
        return self._delay(user, 'alarm_ack_solved', delay)

    def alarm_solved(self, user, delay):
        return self._delay(user, 'alarm_solved', delay)
