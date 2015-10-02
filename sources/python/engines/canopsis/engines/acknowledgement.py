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
from canopsis.event import get_routingkey, forger, is_host_acknowledged
from canopsis.old.account import Account
from canopsis.old.storage import get_storage
from copy import deepcopy
from time import time
from json import dumps


class engine(Engine):
    etype = "acknowledgement"

    def __init__(self, acknowledge_on='canopsis.events', *args, **kargs):
        super(engine, self).__init__(*args, **kargs)

        account = Account(user="root", group="root")

        self.storage = get_storage(namespace='ack', account=account)
        self.events_collection = self.storage.get_backend('events')
        self.stbackend = self.storage.get_backend('ack')
        self.objects_backend = self.storage.get_backend('object')
        self.acknowledge_on = acknowledge_on

    def pre_run(self):

        self.ack_event = forger(
            connector="Engine",
            connector_name=self.etype,
            event_type="perf",
            source_type='component',
            resource='ack',
            state=0,
            state_type=1,
        )

        self.beat()

    def beat(self):
        self.reload_ack_cache()
        self.reload_ack_comments()

    def reload_ack_comments(self):

        # reload comment for ack comparison
        self.comments = []
        query = self.objects_backend.find(
            {'crecord_type': 'comment'}, {'comment': 1, '_id': 1})
        for comment in query:
            self.comments.append(comment)

    def reload_ack_cache(self):
        query = self.stbackend.find({
            'solved': False,
            'ackts': {'$gt': -1}
        }, {'rk': 1})

        # dictionary is faster than list for key test existance
        self.cache_acks = {}
        for ack in query:
            self.cache_acks[ack['rk']] = 1
            self.logger.debug(' + ack cache key > ' + ack['rk'])

    def get_metric_name_adp(self, event):

        domain = event.get('domain', '')
        perimeter = event.get('perimeter', '')

        domain_perimeter = u''

        # If no domain, information domain perimeter is useless
        if domain:
            domain_perimeter = u'_d-{}p-{}'.format(domain, perimeter)

        metric_name_adp = u'{}'.format(domain_perimeter)

        return metric_name_adp

    def work(self, event, *args, **kargs):
        logevent = None

        ackremove = False
        state = event.get('state', 0)
        state_type = event.get('state_type', 1)

        if event['event_type'] == 'ackremove':
            # remove ack from event
            # Ack remove information exists when ack is just removed
            # And deleted if event is ack again
            rk = event['ref_rk']
            self.events_collection.update(
                {'_id': rk},
                {
                    '$set': {
                        'ack_remove': {
                            'author': event['author'],
                            'comment': event['output'],
                            'timestamp': time()
                        }
                    },
                    '$unset': {
                        'ack': '',
                        'ticket_declared_author': '',
                        'ticket_declared_date': '',
                        'ticket': '',
                        'ticket_date': ''
                    }
                }
            )
            ackremove = True

        # If event is of type ack, then ack reference event
        if event['event_type'] == 'ack':
            self.logger.debug('Ack event found, will proceed ack.')

            rk = event.get('referer', event.get('ref_rk', None))

            author = event['author']

            self.logger.debug(dumps(event, indent=2))

            if not rk:
                self.logger.error(
                    'Cannot get acknowledged event, missing referer or ref_rk'
                )
                return event

            for comment in self.comments:
                if comment['comment'] in event['output']:
                    # An ack comment is contained into a defined comment
                    # Then let save referer key to the comment
                    # Set referer rk to last update date
                    self.objects_backend.update(
                        {'_id': comment['_id']},
                        {"$addToSet": {'referer_event_rks': {'rk': rk}}},
                        upsert=True)
                    self.logger.info(
                        'Added a referer rk to the comment {}'.format(
                            comment['comment']
                        )
                    )

            ackts = int(time())
            ack_info = {
                'timestamp': event['timestamp'],
                'ackts': ackts,
                'rk': rk,
                'author': author,
                'comment': event['output']
            }

            # add rk to acknowledged rks
            response = self.stbackend.find_and_modify(
                query={'rk': rk, 'solved': False},
                update={'$set': ack_info},
                upsert=True,
                full_response=True,
                new=True
            )

            self.logger.debug(
                u'Updating event {} with author {} and comment {}'.format(
                    rk,
                    author,
                    ack_info['comment']
                )
            )

            ack_info['isAck'] = True
            # Useless information for event ack data
            del ack_info['ackts']

            # clean eventual previous ack remove information
            self.events_collection.update(
                {
                    '_id': rk
                }, {
                    '$set': {
                        'ack': ack_info,
                    },
                    '$unset': {
                        'ack_remove': '',
                    }
                }
            )

            # When an ack status is changed

            # Emit an event log
            referer_event = self.storage.find_one(
                mfilter={'_id': rk},
                namespace='events'
            )

            if referer_event:

                referer_event = referer_event.dump()

                # Duration between event last state and acknolegement date
                duration = ackts - referer_event.get(
                    'last_state_change', event['timestamp']
                )
                logevent = forger(
                    connector="Engine",
                    connector_name=self.etype,
                    event_type="log",
                    source_type=referer_event['source_type'],
                    component=referer_event['component'],
                    resource=referer_event.get('resource', None),
                    state=0,
                    state_type=1,
                    ref_rk=event['rk'],
                    output=u'Event {0} acknowledged by {1}'.format(
                        rk, author),
                    long_output=event['output'],
                )

                # Now update counters
                ackhost = is_host_acknowledged(event)
                # Cast response to ! 0|1
                cvalues = int(not ackhost)

                ack_event = deepcopy(self.ack_event)
                ack_event['component'] = author
                ack_event['perf_data_array'] = [
                    {
                        'metric': 'alerts_by_host',
                        'value': cvalues,
                        'type': 'COUNTER'
                    },
                    {
                        'metric': 'alerts_count{}'.format(
                            self.get_metric_name_adp(event)
                        ),
                        'value': 1,
                        'type': 'COUNTER'
                    },
                    {
                        'metric': 'delay',
                        'value': duration,
                        'type': 'COUNTER'
                    }
                ]

                publish(
                    publisher=self.amqp, event=ack_event,
                    exchange=self.acknowledge_on
                )

                self.logger.debug('Ack internal metric sent. {}'.format(
                    dumps(ack_event['perf_data_array'], indent=2)
                ))

            for hostgroup in event.get('hostgroups', []):
                ack_event = deepcopy(self.ack_event)
                ack_event['perf_data_array'] = [
                    {
                        'metric': 'alerts',
                        'value': cvalues,
                        'type': 'COUNTER'
                    }
                ]

                publish(
                    publisher=self.amqp, event=ack_event,
                    exchange=self.acknowledge_on
                )

            self.logger.debug('Reloading ack cache')
            self.reload_ack_cache()

        # If event is acknowledged, and went back to normal, remove the ack
        # This test concerns most of case
        # And could not perform query for each event
        elif state == 0 and state_type == 1:
            solvedts = int(time())

            if event['rk'] in self.cache_acks:
                self.logger.debug(
                    'Ack exists for this event, and has to be recovered.'
                )

                # We have an ack to process for this event
                query = {
                    'rk': event['rk'],
                    'solved': False,
                    'ackts': {'$gt': -1}
                }

                ack = self.stbackend.find_one(query)

                if ack:

                    ackts = ack['ackts']

                    self.stbackend.update(
                        query,
                        {
                            '$set': {
                                'solved': True,
                                'solvedts': solvedts
                            }
                        }
                    )

                    logevent = forger(
                        connector="Engine",
                        connector_name=self.etype,
                        event_type="log",
                        source_type=event['source_type'],
                        component=event['component'],
                        resource=event.get('resource', None),

                        state=0,
                        state_type=1,

                        ref_rk=event['rk'],
                        output=u'Acknowledgement removed for event {0}'.format(
                            event['rk']),
                        long_output=u'Everything went back to normal'
                    )

                    logevent['acknowledged_connector'] = event['connector']
                    logevent['acknowledged_source'] = event['connector_name']
                    logevent['acknowledged_at'] = ackts
                    logevent['solved_at'] = solvedts

                    # Metric for solved alarms
                    ack_event = deepcopy(self.ack_event)
                    ack_event['component'] = 'solved_alert'
                    ack_event['perf_data_array'] = [
                        {
                            'metric': 'delay',
                            'value': solvedts - ackts,
                            'unit': 's'
                        },
                        {
                            'metric': 'count',
                            'value': 1,
                            'type': 'COUNTER'
                        }
                    ]

                    publish(
                        publisher=self.amqp, event=ack_event,
                        exchange=self.acknowledge_on
                    )


        # If the event is in problem state,
        # update the solved state of acknowledgement
        elif ackremove or (state != 0 and state_type == 1):
            self.logger.debug('Alert on event, preparing ACK statement.')

            self.stbackend.find_and_modify(
                query={'rk': event['rk'], 'solved': True},
                update={'$set': {
                    'solved': False,
                    'solvedts': -1,
                    'ackts': -1,
                    'timestamp': -1,
                    'author': '',
                    'comment': ''
                }}
            )

        if logevent:
            self.logger.debug('publishing log event {}'.format(
                dumps(logevent, indent=2)
            ))
            publish(
                publisher=self.amqp, event=logevent,
                exchange=self.acknowledge_on
            )

        return event
