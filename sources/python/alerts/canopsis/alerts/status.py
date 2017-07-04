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

from time import time

from canopsis.alerts import AlarmField, States
from canopsis.common.utils import ensure_iterable
from canopsis.check import Check

OFF = 0
ONGOING = 1
STEALTHY = 2
FLAPPING = 3
CANCELED = 4


def get_last_state(alarm, ts=None):
    """
    Get last alarm state.

    :param alarm: Alarm history
    :type alarm: dict

    :param ts: Timestamp to look from (optional)
    :type ts: int

    :returns: Most recent state
    """

    if alarm[AlarmField.state.value] is not None:
        return alarm[AlarmField.state.value]['val']

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

    if alarm[AlarmField.status.value] is not None:
        return alarm[AlarmField.status.value]['val']

    return OFF


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

    if len(alarm[AlarmField.steps.value]) > 0:
        if ts is None:
            ts = alarm[AlarmField.steps.value][-1]['t'] + 1

        steptypes = ensure_iterable(steptypes)

        for step in reversed(alarm[AlarmField.steps.value]):
            if step['t'] < ts and step['_t'] in steptypes:
                return step

    return None


def is_flapping(manager, alarm):
    """
    Check if alarm is flapping.

    :param manager: Alerts manager
    :type manager: canopsis.alerts.manager.Alerts

    :param alarm: Alarm history
    :type alarm: dict

    :returns: ``True`` if alarm is flapping, ``False`` otherwise
    :rtype: bool
    """

    statestep = None
    freq = 0
    ts = alarm[AlarmField.state.value]['t']

    for step in reversed(alarm[AlarmField.steps.value]):
        if (ts - step['t']) > manager.flapping_interval:
            break

        if statestep is None and step['_t'] in ['stateinc', 'statedec']:
            freq += 1
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


def is_keeped_state(alarm):
    """
    Check if an alarm state must be keeped.

    :param manager: Alerts manager
    :type manager: canopsis.alerts.manager.Alerts

    :param alarm: Alarm history
    :type alarm: dict

    :returns: ``True`` if alarm state is forced, ``False`` otherwise
    :rtype: bool
    """
    state = alarm[AlarmField.state.value]

    return state is not None and '_t' in state and state['_t'] == States.changestate.value


def is_stealthy(manager, alarm):
    """
    Check if alarm is stealthy.

    :param manager: Alerts manager
    :type manager: canopsis.alerts.manager.Alerts

    :param alarm: Alarm history
    :type alarm: dict

    :returns: ``True`` if alarm is supposed to be stealthy, ``False`` otherwise
    :rtype: bool
    """

    ts = alarm[AlarmField.state.value]['t']

    for step in reversed(alarm[AlarmField.steps.value]):
        delta1 = ts - step['t']  # delta from last state change
        delta2 = int(time()) - step['t']  # delta from now
        if delta1 > manager.stealthy_show_duration or \
           delta1 > manager.stealthy_interval or \
           delta2 > manager.stealthy_show_duration or \
           delta2 > manager.stealthy_interval:
            break

        if step['_t'] in ['stateinc', 'statedec']:
            if step['val'] != Check.OK and alarm[AlarmField.state.value]['val'] == Check.OK:
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
    :rtype: int
    """

    if alarm[AlarmField.canceled.value] is not None:
        return CANCELED

    if is_flapping(manager, alarm):
        return FLAPPING

    elif is_stealthy(manager, alarm):
        return STEALTHY

    elif alarm[AlarmField.state.value]['val'] != Check.OK:
        return ONGOING

    else:
        return OFF
