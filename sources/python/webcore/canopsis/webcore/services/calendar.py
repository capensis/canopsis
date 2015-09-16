# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 'Capensis' [http://www.capensis.com]
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

from canopsis.common.ws import route
from canopsis.ccalendar.manager import Calendar

from dateutil.rrule import rrulestr
from datetime import datetime
from time import mktime


calendar_manager = Calendar()


def exports(ws):
    rest = ws.require('rest')

    @route(ws.application.get)
    def cal(source, interval_start, interval_end):
        params = {
            'filter': {
                '$and': [
                    {'event_type': 'calendar'},
                    {'component': source},
                    {'rrule': {'$exists': False}},
                    {'$or': [
                        {'$and': [
                            {'start': {'$gt': int(interval_start)}},
                            {'start': {'$lt': int(interval_end)}}
                        ]},
                        {'$and': [
                            {'end': {'$gt': int(interval_start)}},
                            {'end': {'$lt': int(interval_end)}}
                        ]}
                    ]}
                ]
            }
        }

        events = rest.get_records(ws, 'events', **params)

        params['filter'] = {
            '$and': [
                {'event_type': 'calendar'},
                {'component': source},
                {'rrule': {'$exists': True}}
            ]
        }

        recurrent_events = rest.get_records(ws, 'events', **params)

        for event in recurrent_events['data']:
            try:
                dtstart = datetime.fromtimestamp(float(interval_start))
                dtend = datetime.fromtimestamp(float(interval_end))

                evstart = datetime.fromtimestamp(float(event['start']))
                evend = datetime.fromtimestamp(float(event['end']))

                occurences = list(
                    rrulestr(
                        event['rrule'],
                        dtstart=evstart
                    ).between(dtstart, dtend)
                )

                event_duration = evend - evstart

                n_occur = 0

                # instantiate an event occurence for each found date
                for occurence in occurences:
                    n_occur += 1

                    new_event = event.copy()
                    occur_start = mktime(occurence.timetuple())
                    occur_end = occurence + event_duration
                    occur_end = mktime(occur_end.timetuple())

                    new_event['start'] = int(occur_start)
                    new_event['end'] = int(occur_end)
                    events['data'].append(new_event)

            except Exception as e:
                ws.logger.error('Error parsing rrule for event: {0}'.format(e))

        return events

    @route(ws.application.delete, payload=['ids'])
    def calendar(ids):
        calendar_manager.remove(ids)
        ws.logger.info('Delete : {}'.format(ids))
        return True

    @route(
        ws.application.post,
        payload=['document'],
        name='calendar/put'
    )
    def calendar(document):
        ws.logger.debug({
            'document': document,
            'type': type(document)
        })

        calendar_manager.put(document)

        return True

    @route(ws.application.post, payload=['limit', 'start', 'sort', 'filter'])
    def calendar(limit=0, start=0, sort=None, filter={}):
        result = calendar_manager.find(
            limit=limit,
            skip=start,
            query=filter,
            sort=sort,
            with_count=True
        )
        return result
