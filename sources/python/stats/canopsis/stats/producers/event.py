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

    def alarm_opened(self, extra_fields={}):
        self._count(
            'alarm_opened_count',
            extra_fields=extra_fields
        )

    def alarm_ack_solved_delay(self, delay, extra_fields={}):
        self._delay(
            'alarm_ack_solved_delay',
            delay,
            extra_fields=extra_fields
        )

    def alarm_solved_delay(self, delay, extra_fields={}):
        self._delay(
            'alarm_solved_delay',
            delay,
            extra_fields=extra_fields
        )
