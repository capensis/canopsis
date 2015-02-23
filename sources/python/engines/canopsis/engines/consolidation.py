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

from canopsis.engines import Engine
from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from canopsis.old.event import forger, get_routingkey
from canopsis.old.tools import roundSignifiantDigit
from canopsis.perfdata.manager import PerfData
from canopsis.timeserie import TimeSerie
from canopsis.timeserie.timewindow import TimeWindow, Period

from json import loads

from time import time
from datetime import datetime


class engine(Engine):
    etype = 'consolidation'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.storage = get_storage(
            namespace='events',
            account=Account(
                user="root",
                group="root"
            )
        )


    def pre_run(self):
        self.storage = get_storage(namespace='object',
            account=Account(user="root", group="root"))
        self.manager = PerfData()

    def publish_aggre_stats(self):

        series_event = forger(
            connector='engine',
            connector_name='engine',
            event_type='perf',
            source_type='resource',
            resource='series_events',
            state=0,
            perf_data_array=self.perf_data_array
        )

        rk = get_routingkey(series_event)

        self.logger.debug('Publishing {} : {}'.format(rk, series_event))

        self.amqp.publish(
            series_event,
            rk,
            self.amqp.exchange_name_events
        )

    def fetch(self, series, _from, _to):
        self.logger.debug("Je passe dans fetch \n\n\n\n\n")

    def consume_dispatcher(self, event, *args, **kargs):
        self.logger.debug("Start metrics consolidation")
        serie = self.get_ready_record(event)

        if not serie:
            # Show error message
            self.logger.error('No record found')

        self.logger.debug(serie)
        #self.fetch(serie)
        event_id = event['_id']
        # Update crecords informations
        self.crecord_task_complete(event_id)
