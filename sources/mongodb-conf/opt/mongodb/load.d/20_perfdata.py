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

from canopsis.context.manager import Context
from canopsis.perfdata.manager import PerfData

logger = None


def init():
    pass


def update():

    nonetonan()

def nonetonan():
    """change none values to nan"""

    context = Context()
    perfdata = PerfData()
    metrics = context.find(_type='metric')

    for metric in metrics:
        metric_id = context.get_entity_id(metric)
        points = perfdata.get(metric_id=metric_id, with_meta=False)

        nan = float('nan')

        nonepoints = list(
            (point[0], nan) for point in points if point[1] is None
        )

        perfdata.put(metric_id=metric_id, points=nonepoints)
