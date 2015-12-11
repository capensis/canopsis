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

logger = None


def init():
    pass


def update():

    fixtimestampandnone()

def fixtimestampandnone():
    """Fix wrong timestamp when period uses a week."""

    perfdatatofix = PerfData()

    perfdata_migrating = 'perfdata_migrating'

    # start to rename perfdata collection
    perfdatatofix[PerfData.PERFDATA_STORAGE]._backend.rename(perfdata_migrating)
    perfdatatofix[PerfData.PERFDATA_STORAGE].table = perfdata_migrating

    fixedperfdata = PerfData()

    nan = float('nan')

    for document in perfdatatofix[PerfData.PERFDATA_STORAGE].find_elements():

        metric_id = document['i']

        values = document['v']
        t = document['t']

        points = list(
            (t + int(ts), nan if values[ts] is None else values[ts])
            for ts in values
        )

        fixedperfdata.put(metric_id=metric_id, points=points, cache=False)
