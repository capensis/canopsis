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

from canopsis.stats.manager import Stats
from canopsis.engines.core import Engine, publish
from canopsis.event import forger

import pprint
pp = pprint.PrettyPrinter(indent=4)


class engine(Engine):

    etype = 'stats'

    """
    This engine's goal is to compute canopsis internal statistics.
    Statistics are computed on each passing event and are updated
    in async way in order to manage performances issues.
    """

    def pre_run(self):

        self.stats_manager = Stats()

        self.beat()

    def consume_dispatcher(self, event, *args, **kargs):

        """
        Entry point for stats computation. Triggered by the dispatcher
        engine for distributed processing puroses.
        Following methods will generate metrics that are finally embeded
        in a metric event.
        """

        self.logger.debug('Entered in stats consume dispatcher')

        self.stats_manager.set_perf_data_array([])

        self.compute_states()

        self.publish_states()

    def compute_states(self):

        """
        Entry point for dynamic stats method triggering
        Dynamic triggering allow greated control on which
        stats are computed and can be activated/deactivated
        from frontend.
        """

        # Allow individual stat computation management from ui.
        stats_to_compute = [
            'users_session_duration',
        ]

        for stat in stats_to_compute:
            method = getattr(self.stats_manager, stat)
            method()

    def publish_states(self):

        perfdatas = self.stats_manager.get_perf_data_array()

        stats_event = forger(
            connector='engine',
            connector_name='engine',
            event_type='perf',
            source_type='resource',
            component='__canopsis__',
            resource='Engine_stats',
            state=0,
            perf_data_array=perfdatas
        )

        # Just log information
        metrics = []
        for m in perfdatas:
            metrics.append('{}\t{}\t{}'.format(
                m['value'],
                m['type'],
                m['metric']
            ))

        self.logger.debug('-- Generated perfdata --\n{}'.format(
            '\n'.join(metrics)
        ))

        publish(publisher=self.amqp, event=stats_event)
