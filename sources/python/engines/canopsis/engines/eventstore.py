# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
from canopsis.old.archiver import Archiver
from canopsis.old.downtime import Downtime
from canopsis.old.storage import CONFIG

from csv import reader


class engine(Engine):
    etype = 'eventstore'

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.archiver = Archiver(
            namespace='events', confnamespace='object',
            autolog=False, logging_level=self.logging_level)

        self.event_types = reader([CONFIG.get('events', 'types')]).next()
        self.check_types = reader([CONFIG.get('events', 'checks')]).next()
        self.log_types = reader([CONFIG.get('events', 'logs')]).next()
        self.comment_types = reader([CONFIG.get('events', 'comments')]).next()

        self.cdowntime = Downtime()
        self.beat()

    def beat(self):
        self.archiver.beat()
        self.cdowntime.reload(self.beat_interval)

    def store_check(self, event):
        _id = self.archiver.check_event(event['rk'], event)

        if event.get('downtime', False):
            event['previous_state_change_ts'] = \
                self.cdowntime.get_downtime_end_date(
                    event['component'], event.get('resource', ''))

        if _id:
            event['_id'] = _id
            event['event_id'] = event['rk']
            ## Event to Alert
            self.amqp.publish(
                event, event['rk'], self.amqp.exchange_name_alerts)

    def store_log(self, event, store_new_event=True):
        ## passthrough

        if store_new_event:
            self.archiver.store_new_event(event['rk'], event)

        _id = self.archiver.log_event(event['rk'], event)
        event['_id'] = _id
        event['event_id'] = event['rk']

        ## Event to Alert
        self.amqp.publish(event, event['rk'], self.amqp.exchange_name_alerts)

    def work(self, event, *args, **kargs):

        if 'exchange' in event:
            del event['exchange']

        event_type = event['event_type']

        if event_type not in self.event_types:
            self.logger.warning("Unknown event type '%s', id: '%s', event:\n%s" % (event_type, event['rk'], event))
            return event

        elif event_type in self.check_types:
            self.store_check(event)

        elif event_type in self.log_types:
            self.store_log(event)

        elif event_type in self.comment_types:
            self.store_log(event, store_new_event=False)

        return event
