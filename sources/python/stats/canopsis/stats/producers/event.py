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


CONF_PATH = 'stats/producers/event.conf'
CATEGORY = 'USER_METRIC_PRODUCER'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class EventMetricProducer(MetricProducer):
    """
    Metric producer for event statistics.
    """

    def alarm(self, event):
        return self._counter('alarm', event)

    def alarm_ack(self, event):
        return self._counter('alarm_ack', event)

    def alarm_ack_solved(self, event):
        return self._counter('alarm_ack_solved', event)

    def alarm_solved(self, event):
        return self._counter('alarm_solved', event)

    def alarm_ack_delay(self, delay):
        return self._delay('alarm_ack_delay', delay)

    def alarm_ack_solved_delay(self, delay):
        return self._delay('alarm_ack_solved_delay', delay)

    def alarm_solved_delay(self, delay):
        return self._delay('alarm_solved_delay', delay)
