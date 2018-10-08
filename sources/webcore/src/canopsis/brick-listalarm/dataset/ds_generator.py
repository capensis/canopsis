#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from uuid import uuid4
from random import choice, randint
from string import lowercase
from datetime import date
from json import dumps

N_ALARMS = 50


def randword(n_chars):
    return ''.join(choice(lowercase) for i in range(n_chars))


def randts(start, stop):
    return int(date.fromordinal(randint(start, stop)).strftime('%s'))


def randrrule():
    return 'FREQ={}'.format(choice(['DAILY', 'WEEKLY', 'YEARLY']))


def generate():
    offset = randint(1, 3) * N_ALARMS

    ds = {
        'alarms': [],
        'first': offset + 1,
        'last': offset + N_ALARMS,
        'total': offset + randint(1, 3) * N_ALARMS + randint(0, N_ALARMS - 1)
    }

    tstart = date.today().replace(day=1, month=1).toordinal()
    tstop = date.today().toordinal()

    connectors = {}
    for i in range(randint(3, 5)):
        connectors[randword(3)] = []

    for k, v in connectors.items():
        v.append(randword(3))

    domains = [randword(3) for i in range(4)]
    perimeters = [randword(3) for i in range(4)]

    for i in range(N_ALARMS):
        connector = choice(connectors.keys())
        connector_name = choice(connectors[connector])
        component = randword(3)
        resource = choice([None, randword(3)])

        if resource is None:
            entity_id = '/component/{}/{}/{}'.format(
                connector,
                connector_name,
                component
            )
        else:
            entity_id = '/resource/{}/{}/{}/{}'.format(
                connector,
                connector_name,
                component,
                resource
            )

        alarm = {
            '_id': str(uuid4()),
            'v': {
                'connector': connector,
                'connector_name': connector_name,
                'component': component,
                'resource': resource,
                'state': {
                    'a': '{}.{}'.format(connector, connector_name),
                    '_t': choice(['statedec', 'stateinc']),
                    't': randts(tstart, tstop),
                    'm': ' '.join(randword(3) for i in range(4)),
                    'val': randint(0, 3),
                },
                'status': {
                    'a': '{}.{}'.format(connector, connector_name),
                    '_t': choice(['statedec', 'stateinc']),
                    't': randts(tstart, tstop),
                    'm': ' '.join(randword(3) for i in range(4)),
                    'val': randint(0, 4),
                },
                'ack': choice([
                    None,
                    {
                        'a': randword(3),
                        '_t': 'ack',
                        't': randts(tstart, tstop),
                        'm': ' '.join(randword(3) for i in range(4)),
                    }
                ]),
                'cancel': choice([
                    None,
                    {
                        'a': randword(3),
                        '_t': 'cancel',
                        't': randts(tstart, tstop),
                        'm': ' '.join(randword(3) for i in range(4)),
                    }
                ]),
                'ticket': choice([
                    None,
                    {
                        'a': randword(3),
                        '_t': choice(['declareticket', 'assocticket']),
                        't': randts(tstart, tstop),
                        'm': ' '.join(randword(3) for i in range(4)),
                        'val': str(uuid4()),
                    }
                ]),
                'snooze': choice([
                    None,
                    {
                        'a': randword(3),
                        '_t': 'snooze',
                        't': randts(tstart, tstop),
                        'm': ' '.join(randword(3) for i in range(4)),
                        'val': randts(tstart, tstop)
                    }
                ]),
                'output': ' '.join(randword(3) for i in range(4)),
                'resolved': None,
                'extra': {
                    'domain': choice(domains),
                    'perimeter': choice(perimeters)
                },
                'pbehaviors': choice([None, []]),
                'hard_limit': 'to be ignored...',
                'steps': ['to be ignored...']
            },
            'd': entity_id,
            't': randts(tstart, tstop)
        }

        if choice([True, False, False]):
            alarm['v']['linklist']['event_links'] = []

            for i in range(randint(2, 4)):
                randurl = 'http{}://{}'.format(choice(['', 's']), uuid4().hex)
                alarm['v']['linklist']['event_links'].append(
                    {
                        'url': randurl,
                        'label': randword(3)
                    }
                )

        if alarm['v']['pbehaviors'] == []:
            if choice([True, False]):
                alarm['v']['pbehaviors'].append(
                    {
                        'name': 'downtime',
                        'enabled': choice([True, False]),
                        'tstart': randts(tstart, tstop),
                        'tstop': randts(tstart, tstop),
                        'rrule': randrrule(),
                    }
                )

            if choice([True, False]):
                alarm['v']['pbehaviors'].append(
                    {
                        'name': 'custom pbehavior',
                        'enabled': choice([True, False]),
                        'tstart': randts(tstart, tstop),
                        'tstop': randts(tstart, tstop),
                        'rrule': randrrule(),
                    }
                )

        ds['alarms'].append(alarm)

    return ds


if __name__ == '__main__':
    dataset = generate()

    print(dumps(dataset, separators=(',', ': '), indent=2))
