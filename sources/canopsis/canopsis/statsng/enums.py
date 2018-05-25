# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

"""
The canopsis.statsng.enums module contains enumerations related to the
statsng engine.
"""

from __future__ import unicode_literals


# Inheriting from Enum would prevent the members from being JSON serializable
class StatEvents(object):
    """
    The StatEvents enumeration defines the types of events handled by the
    statsng engine.
    """
    statcounterinc = 'statcounterinc'
    statduration = 'statduration'


class StatEventFields(object):
    """
    The StatEventFields enumeration defines the names of the field of a stat*
    events.
    """
    alarm = 'alarm'
    entity = 'entity'
    counter_name = 'counter_name'
    duration_name = 'duration_name'
    duration = 'duration'


class StatMeasurements(object):
    """
    The StatMeasurements enumeration defines the names of the influxdb
    measurements.
    """
    counters = 'statcounters'
    durations = 'statdurations'


class StatCounters(object):
    """
    The StatCounters enumeration defines the names of the counters.
    """
    alarms_created = 'alarms_created'
    alarms_canceled = 'alarms_canceled'
    alarms_resolved = 'alarms_resolved'


class StatDurations(object):
    """
    The StatDurations enumeration defines the names of the durations.
    """
    ack_time = 'ack_time'
    resolve_time = 'resolve_time'
