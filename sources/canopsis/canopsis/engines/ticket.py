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

from __future__ import unicode_literals

from canopsis.engines.core import Engine
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from copy import deepcopy
from time import time
from time import sleep


class engine(Engine):
    etype = "ticket"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.store = get_storage('object', account=Account(user='root'))

    def pre_run(self):
        self.beat()

    def beat(self):
        try:
            self.config = self.store.get('cservice.ticket').dump()

        except KeyError:
            if not hasattr(self, 'config'):
                self.config = {}

    def work(self, event, *args, **kwargs):
        if 'job' in self.config:

            if 'rrule' in self.config['job']:
                del self.config['job']['rrule']

            if event['event_type'] == 'declareticket':
                self.logger.debug(u'Declare Ticket')

                try:
                    refevt = self.store.get(
                        event['ref_rk'], namespace='events'
                    )
                    refevt = refevt.dump()

                except KeyError:
                    refevt = {}

                if refevt.get('ack', {}) == {}:
                    sleep(2)
                    try:
                        refevt = self.store.get(
                            event['ref_rk'], namespace='events'
                        )
                        refevt = refevt.dump()

                    except KeyError:
                        refevt = {}

                job = deepcopy(self.config['job'])
                job['_id'] = self.config['_id']
                job['context'] = refevt

                try:
                    self.work_amqp_publisher.direct_event(
                        job, 'Engine_scheduler')
                except Exception as e:
                    self.logger.exception("Unable to send event to next queue")

                self.logger.info(
                    'Setting ticked received for {}'
                    .format(event['ref_rk'])
                )

                self.store.get_backend('events').update({
                    'rk': event['ref_rk']
                }, {
                    '$set': {
                        'ticket_declared_author': event['author'],
                        'ticket_declared_date': int(time()),
                    }
                })

            elif (event['event_type'] in ['ack', 'assocticket']
                    and 'ticket' in event):

                self.logger.info(
                    'Associate ticket for event type {}'
                    .format(event['event_type'])
                )

                events = self.store.get_backend('events')

                self.logger.info(
                    'Update events with rk {0}'
                    .format(event['ref_rk'])
                )
                events.update({
                    'rk': event['ref_rk']
                }, {
                    '$set': {
                        'ticket': event['ticket'],
                        'ticket_date': int(time())
                    }
                })

        return event
