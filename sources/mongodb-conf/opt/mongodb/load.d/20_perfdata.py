#!/usr/bin/env python
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

from canopsis.perfdata.manager import PerfData
from time import time

logger = None


def init():
    pass


def update():

    fixtimestampandnone()

def fixtimestampandnone():
    """Fix wrong timestamp when period uses a week."""

    perfdata = PerfData()

    nan = float('nan')

    oneweek = 3600 * 24 * 7

    cursor = perfdata[PerfData.PERFDATA_STORAGE].find_elements(
        {'p': {'$exists': True}}
    )

    count = cursor.count()

    i = 0
    prct = 0

    start = time()

    for document in cursor:

        _prct = round(float(i) / count, 2) * 100
        if prct != _prct:
            prct = _prct
            print('{0}% . Elapsed time: {1}s'.format(prct, time() - start))

        i += 1

        del document['p']

        metric_id = document['i']

        values = document['v']
        t = document['t']

        points = list(
            (t + int(ts), nan if values[ts] is None else values[ts])
            for ts in values
        )

        rightvalues = {
            key: values[key] for key in values if int(key) < oneweek
        }
        document['v'] = rightvalues

        perfdata[PerfData.PERFDATA_STORAGE]._update(
            spec={'_id': document['_id']}, document={'$unset': {'p': ''}},
            multi=False, cache=False
        )

        perfdata[PerfData.PERFDATA_STORAGE].put_element(
            element=document, cache=False
        )

        perfdata.put(metric_id=metric_id, points=points, cache=False)
