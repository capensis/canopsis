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

from canopsis.engines.core import Engine, publish
from canopsis.old.storage import get_storage
from canopsis.old.account import Account
from canopsis.old.record import Record
from time import time

# Delay since the lock document is released in any cases
UNLOCK_DELAY = 60


class engine(Engine):
    etype = "crecord_dispatcher"

    def __init__(self, *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        self.crecords = []
        self.delays = {}
        self.beat_interval = 15
        self.nb_beat = 0
        self.crecords_types = [
            'selector',
            'topology',
            'derogation',
            'serie',
        ]

        self.beat_interval_trigger = {
            'eventstore': {
                'delay': 60,
                'elapsed_since_last_beat': 0
            },
            'datacleaner': {
                'delay': 3600,
                'elapsed_since_last_beat': 0
            },
            'stats': {
                'delay': 60,
                'elapsed_since_last_beat': 0
            },
        }

    def pre_run(self):
        # Load crecords from database
        self.storage = get_storage(
            namespace='object',
            account=Account(
                user="root",
                group="root"
            ))

        self.backend = self.storage.get_backend('object')

        self.logger.info('Release crecrord dispatcher lock')
        self.Lock.release('load_crecords', self.backend)

        self.ha_engine_triggers = {}

        self.beat()

    def load_crecords(self):
        crecords = []
        now = int(time())

        with self.Lock(self, 'load_crecords') as l:
            if l.own():
                # Crecord selection can be performed now as lock is ready.
                # Crecord are loaded again if last dispatch update
                # is not set or < now - 60 seconds
                crecords_json = self.storage.find({
                    'crecord_type': {'$in': self.crecords_types},
                    'enable': True,
                    '$or': [
                        {'next_ready_time':
                            {'$exists': False}},  # crecord is new
                        {'last_dispatch_update':
                            {'$exists': False}},  # crecord is new
                        {'last_dispatch_update':
                            {'$lte': now - 60}},  # unlock case
                        {'$and': [
                            {'next_ready_time':
                                {'$lte': now}},  # record is ready
                            {'loaded': False}
                        ]}
                    ]
                }, namespace="object")

                for crecord_json in crecords_json:
                    try:
                        self.storage.update(crecord_json._id, {
                            'loaded': True,
                            'last_dispatch_update': now}
                        )

                        crecord = Record(
                            storage=self.storage,
                            record=crecord_json
                        )

                        crecords.append(crecord)
                    except Exception as e:
                        self.logger.error(
                            'can\'t process crecord {}: {}'.format(
                                crecord_json._id,
                                e
                            ))
        return crecords

    def publish_record(self, event, crecord_type):

        try:
            if crecord_type == 'serie':
                rk = 'dispatcher.{0}'.format('consolidation')
            else:
                rk = 'dispatcher.{0}'.format(crecord_type)

            #rk = 'dispatcher.{0}'.format(crecord_type)

            self.amqp.get_exchange('media')
            publish(publisher=self.amqp, event=event, rk=rk, exchange='media')

            return True

        except Exception as e:
            # Will be reloaded on next beat
            self.logger.error('Unable to send crecord {} error : {}'.format(
                crecord_type,
                e
            ))
            return False

    def beat(self):

        # These events will run only once the list below consumer's
        # dispatch method.
        # This ensure those engine methods are run once in ha mode
        # Event are triggered only at engine's delay duration
        for trigger_engine in self.beat_interval_trigger:

            self.logger.debug('Processing engine trigger for {}'.format(
                trigger_engine
            ))

            tengine = self.beat_interval_trigger[trigger_engine]

            if tengine['elapsed_since_last_beat'] > tengine['delay']:

                self.logger.debug(u'triggering dispatch for {}'.format(
                    trigger_engine
                ))

                self.publish_record({
                    'event': 'engine process trigger'
                }, trigger_engine)

                # Update deplay
                tengine['elapsed_since_last_beat'] = 0

            # Reset delay
            tengine['elapsed_since_last_beat'] += self.beat_interval

        """
            Reinitialize crecords and may publish event related
            crecord targeted to other engines crecord queues
        """
        # Triggers consume dispatch method within engines.
        # This ensure those engine methods are run once in ha mode
        for ha_engine_trigger in self.ha_engine_triggers:

            tengine = self.ha_engine_triggers[ha_engine_trigger]

            if time() - tengine['last_update'] > tengine['delay']:

                self.publish_record(
                    {'crecord_type': ha_engine_trigger},
                    ha_engine_trigger
                )

                tengine['last_update'] = time()

        crecords = self.load_crecords()

        # Loop until list is empty
        self.logger.debug(u' + {} beat, {} crecords queued'.format(
            self.name,
            len(crecords)
        ))

        for crecord in crecords:

            # Every crecord is sent to rabbit mq queues
            # for each listening engines
            dump = crecord.dump()
            record_id = dump['_id']

            # Crecord is sent to other engines and is not kept anymore
            if '_id' in dump:
                dump['_id'] = str(dump['_id'])
                # Just sending key and type to build back object
                # From dedicated engines
                is_published = self.publish_record(
                    dump,
                    dump['crecord_type']
                )

                if not is_published:
                    self.storage.update(record_id, {'loaded': False})

        self.nb_beat += 1
