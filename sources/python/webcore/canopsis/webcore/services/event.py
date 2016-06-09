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

from canopsis.common.utils import ensure_iterable
from canopsis.common.ws import route
from canopsis import schema
from canopsis.event.eventslogmanager import EventsLog
from canopsis.common.utils import singleton_per_scope

from bottle import HTTPError
import requests
import json


def exports(ws):
    manager = singleton_per_scope(EventsLog)

    @route(ws.application.post, name='event', payload=['event', 'url'])
    @route(ws.application.put, name='event', payload=['event', 'url'])
    def send_event(event, url=None):
        if ws.enable_crossdomain_send_events and url is not None:
            payload = {
                'event': json.dumps(event)
            }

            response = requests.post(url, data=payload)

            if response.status_code != 200:
                api_response = json.loads(response.text)

                return (api_response['data'], api_response['total'])

            else:
                return HTTPError(response.status_code, response.text)

        else:
            events = ensure_iterable(event)
            exchange = ws.amqp.exchange_name_events

            for event in events:
                if schema.validate(event, 'cevent'):
                    sname = 'cevent.{0}'.format(event['event_type'])

                    if schema.validate(event, sname):
                        if event['event_type'] == 'eue':
                            sname = 'cevent.eue.{0}'.format(
                                event['type_message']
                            )

                            if not schema.validate(event, sname):
                                continue

                        rk = '{0}.{1}.{2}.{3}.{4}'.format(
                            u'{0}'.format(event['connector']),
                            u'{0}'.format(event['connector_name']),
                            u'{0}'.format(event['event_type']),
                            u'{0}'.format(event['source_type']),
                            u'{0}'.format(event['component'])
                        )

                        if event['source_type'] == 'resource':
                            rk = '{0}.{1}'.format(rk, event['resource'])

                        ws.amqp.publish(event, rk, exchange)

            return events

    @route(ws.application.get,
           name='eventslog/count',
           payload=['tstart', 'tstop', 'limit', 'select']
           )
    def get_event_count_per_day(tstart, tstop, limit=100, select={}):
        """ get eventslog log count for each days in a given period
            :param tstart: timestamp of the begin period
            :param tstop: timestamp of the end period
            :param limit: limit the count number per day
            :param select: filter for eventslog collection
            :return: list in which each item contains an interval and the
            related count
            :rtype: list
        """

        results = manager.get_eventlog_count_by_period(
            tstart, tstop, limit=limit, query=select
        )

        return results
