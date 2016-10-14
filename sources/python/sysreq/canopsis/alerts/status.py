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

from canopsis.common.utils import ensure_iterable
from canopsis.check import Check

OFF = 0
ONGOING = 1
STEALTHY = 2
FLAPPING = 3
CANCELED = 4


def get_previous_step(alarm, steptypes, ts=None):
    """
    Get last step in alarm history.

    :param alarm: Alarm history
    :type alarm: dict

    :param steptypes: Step types wanted
    :type steptypes: str or list

    :param ts: Timestamp to look from (optional)
    :type ts: int

    :returns: Most recent step
    """

    if len(alarm['steps']) > 0:
        if ts is None:
            ts = alarm['steps'][-1]['t'] + 1

        steptypes = ensure_iterable(steptypes)

        for step in reversed(alarm['steps']):
            if step['t'] < ts and step['_t'] in steptypes:
                return step

    return None


def get_last_state(alarm, ts=None):
    """
    Get last alarm state.

    :param alarm: Alarm history
    :type alarm: dict

    :param ts: Timestamp to look from (optional)
    :type ts: int

    :returns: Most recent state
    """

    if alarm['state'] is not None:
        return alarm['state']['val']

    return Check.OK


def get_last_status(alarm, ts=None):
    """
    Get last alarm status.

    :param alarm: Alarm history
    :type alarm: dict

    :param ts: Timestamp to look from (optional)
    :type ts: int

    :returns: Most recent status
    """

    if alarm['status'] is not None:
        return alarm['status']['val']

    return OFF


def is_flapping(manager, alarm):
    """
    Check if alarm is flapping.

    :param manager: Alerts manager
    :type manager: canopsis.alerts.manager.Alerts

    :param alarm: Alarm history
    :type alarm: dict

    :returns: ``True`` if alarm is flapping, ``False`` otherwise
    """

    statestep = None
    freq = 0
    ts = alarm['state']['t']

    for step in reversed(alarm['steps']):
        if (ts - step['t']) > manager.flapping_interval:
            break

        if statestep is None and step['_t'] in ['stateinc', 'statedec']:
            statestep = step

        elif step['_t'] == 'stateinc' and statestep['_t'] == 'statedec':
            freq += 1
            statestep = step

        elif step['_t'] == 'statedec' and statestep['_t'] == 'stateinc':
            freq += 1
            statestep = step

        if freq >= manager.flapping_freq:
            return True

    return False


def is_stealthy(manager, alarm):
    """
    Check if alarm is stealthy.

    :param manager: Alerts manager
    :type manager: canopsis.alerts.manager.Alerts

    :param alarm: Alarm history
    :type alarm: dict

    :returns: ``True`` if alarm is stealthy, ``False`` otherwise
    """

    ts = alarm['state']['t']

    for step in reversed(alarm['steps']):
        if (ts - step['t']) > manager.stealthy_show_duration:
            break

        elif (ts - step['t']) > manager.stealthy_interval:
            break

        if step['_t'] in ['stateinc', 'statedec']:
            if step['val'] != Check.OK and alarm['state']['val'] == Check.OK:
                return True

    return False


def compute_status(manager, alarm):
    """
    Compute alarm status from its history.

    :param manager: Alerts manager
    :type manager: canopsis.alerts.manager.Alerts

    :param alarm: Alarm history
    :type alarm: dict

    :returns: Alarm status as int
    """

    if alarm['canceled'] is not None:
        return CANCELED

    if is_flapping(manager, alarm):
        return FLAPPING

    elif is_stealthy(manager, alarm):
        return STEALTHY

    elif alarm['state']['val'] != Check.OK:
        return ONGOING

    else:
        return OFF
