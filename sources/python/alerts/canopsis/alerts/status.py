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
from canopsis.check import archiver, Check


def get_previous_step(alarm, steptypes, ts=None):
    if len(alarm['steps']) > 0:
        if ts is None:
            ts = alarm['steps'][-1]['t'] + 1

        steptypes = ensure_iterable(steptypes)

        for step in reversed(alarm['steps']):
            if step['t'] < ts and step['_t'] in steptypes:
                return step

    return None


def get_last_state(alarm):
    step = get_previous_step(alarm, ['stateinc', 'statedec'])

    if step is not None:
        return step['val']

    return Check.OK


def is_flapping(manager, alarm):
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
    if alarm['canceled'] is not None:
        return archiver.CANCELED

    if is_flapping(manager, alarm):
        return archiver.BAGOT

    elif is_stealthy(manager, alarm):
        return archiver.STEALTHY

    elif alarm['state']['val'] != Check.OK:
        return archiver.ONGOING

    else:
        return archiver.OFF
