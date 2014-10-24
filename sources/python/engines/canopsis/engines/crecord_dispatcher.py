# -*- coding: utf-8 -*-
#--------------------------------
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
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.old.selector import Selector

from time import time

# Delay since the lock document is released in any cases
UNLOCK_DELAY = 60


class engine(Engine):
    etype = "crecord_dispatcher"

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.crecords = []
        self.delays = {}
        self.beat_interval = 5
        self.nb_beat = 0
        self.crecords_types = ['selector', 'topology', 'derogation', 'consolidation']
        self.beat_interval_trigger = {
            'downtime': {'delay': 60, 'elapsed_since_last_beat': 0},
            'perfstore2_rotate': {'delay': 60, 'elapsed_since_last_beat': 0}
        }

    def pre_run(self):
        #load crecords from database
        self.storage = get_storage(
            namespace='object', account=Account(user="root", group="root"))
        self.backend = self.storage.get_backend('object')

        self.ha_engine_triggers = {'downtime': {
            'last_update': time(),
            'delay': 60
        }}

        self.beat()

    def load_crecords(self):
        crecords = []
        now = int(time())

        with self.Lock(self, 'load_crecords') as l:
            if l.own():
                # Crecord selection can be performed now as lock is ready. Crecord are loaded again if last dispatch update is not set or < now - 60 seconds
                crecords_json = self.storage.find({
                    'crecord_type': {'$in': self.crecords_types},
                    'enable': True,
                    '$or': [
                        {'next_ready_time': {'$exists': False}},  # crecord is new
                        {'last_dispatch_update': {'$exists': False}},  # crecord is new
                        {'last_dispatch_update': {'$lte': now - 60}},  # unlock case
                        {'$and': [
                            {'next_ready_time': {'$lte': now}},  # record is ready
                            {'loaded': False}
                        ]}
                    ]
                }, namespace="object")

                for crecord_json in crecords_json:
                    # let say selector is loaded
                    try:
                        self.storage.update(crecord_json._id, {'loaded': True, 'last_dispatch_update': now})
                        crecord = Selector(storage=self.storage, record=crecord_json, logging_level=self.logging_level)
                        crecords.append(crecord)
                    except Exception as e:
                        self.logger.error('Unable to manage crecord with id {}: {}'.format(crecord_json._id, e))
        return crecords

    def publish_record(self, event, crecord_type):
        """
        Factorised code method
        """
        rk = 'dispatcher.{0}'.format(crecord_type)

        self.amqp.get_exchange('media')
        self.amqp.publish(event, rk, exchange_name='media')

    def beat(self):

        #These events will run only once the list below consumer's dispatch method.
        #This ensure those engine methods are run once in ha mode
        #Event are triggered only at engine's delay duration
        for trigger_engine in self.beat_interval_trigger:
            if self.beat_interval_trigger[trigger_engine]['delay'] > self.beat_interval_trigger[trigger_engine]['elapsed_since_last_beat']:
                self.logger.debug('triggering dispatch for ' + trigger_engine)
                self.publish_record({'event': 'engine process trigger'}, trigger_engine)
                #update deplay
                self.beat_interval_trigger[trigger_engine]['elapsed_since_last_beat'] = 0
            #reset delay
            self.beat_interval_trigger[trigger_engine]['elapsed_since_last_beat'] += self.beat_interval

        """ Reinitialize crecords and may publish event related credort targeted to other engines crecord queues"""
        #Triggers consume dispatch method within engines.
        #This ensure those engine methods are run once in ha mode
        for ha_engine_trigger in self.ha_engine_triggers:
            if time() - self.ha_engine_triggers[ha_engine_trigger]['last_update'] > self.ha_engine_triggers[ha_engine_trigger]['delay']:
                self.publish_record({'crecord_type': ha_engine_trigger}, ha_engine_trigger)
                self.ha_engine_triggers[ha_engine_trigger]['last_update'] = time()

        crecords = self.load_crecords()

        # Loop until list is empty
        self.logger.debug(' + %s beat, %s crecords queued to publish @ %s' % (self.name, len(crecords), int(time())))

        for crecord in crecords:

            # Every crecord is sent to rabbit mq queues for each listening engines
            dump = crecord.dump()
            record_id = dump['_id']

            #crecord is sent to other engines and is not kept anymore
            if '_id' in dump:
                dump['_id'] = str(dump['_id'])
                try:
                    # just sending key and type to build back object from dedicated engines
                    self.publish_record(dump, dump['crecord_type'])

                    #Special case: selector crecords targeted to SLA
                    if dump['crecord_type'] == 'selector' and 'rk' in dump and dump['rk'] and 'dosla' in dump and dump['dosla'] in [ True, 'on'] and 'dostate' in dump and dump['dostate'] in [ True, 'on']:
                        self.publish_record(dump, 'sla')

                except Exception as e:
                    #Crecord gets out of queue and will be reloaded on next beat
                    self.logger.error('Dispatcher was unable to send crecord_type error : %s' % (e))
                    self.storage.update(record_id, {'loaded': False})

        self.nb_beat += 1
