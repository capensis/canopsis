# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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

from copy import deepcopy


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
                self.logger.debug('Declare Ticket')

                try:
                    refevt = self.store.get(event['ref_rk'], namespace='events')
                    refevt = refevt.dump()

                except KeyError:
                    refevt = {}

                job = deepcopy(self.config['job'])
                job['_id'] = self.config['_id']
                job['context'] = refevt

                self.amqp.publish(job, 'Engine_scheduler', 'amq.direct')

            elif event['event_type'] == 'ack' and 'ticket' in event:
                self.logger.info('Associate ticket')

                events = self.store.get_backend('events')

                self.logger.info('Update events with rk {0}'.format(event['ref_rk']))
                events.update({
                    'rk': event['ref_rk']
                }, {
                    '$set': {
                        'ticket': event['ticket']
                    }
                })

        return event
