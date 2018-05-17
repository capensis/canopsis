# -*- coding: utf-8 -*-

from __future__ import unicode_literals

def compat_go_crop_states(alarm):
    """
    If a statecounter field exists in an alarm step,
    val key is replaced with the value of the statecounter field.

    :param dict alarm: dict with alarm data
    :returns: the alarm with modified steps
    :rtype dict:
    """
    steps = alarm.get('v', {}).get('steps', [])
    for i, step in enumerate(steps):
        if 'statecounter' in step:
            step['val'] = step['statecounter']
            alarm['v']['steps'][i] = step

    return alarm
