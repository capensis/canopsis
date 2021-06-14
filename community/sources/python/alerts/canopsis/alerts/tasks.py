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

from canopsis.alerts.status import (
    compute_status, OFF, CANCELED, get_previous_step)

from canopsis.task.core import register_task


@register_task('alerts.useraction.ack')
def acknowledge(manager, alarm, author, message, event):
    """
    Called when a user adds an acknowledgment on an alarm.
    """

    step = {
        '_t': 'ack',
        't': event['timestamp'],
        'a': author,
        'm': message
    }

    alarm['ack'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.useraction.ackremove')
def unacknowledge(manager, alarm, author, message, event):
    """
    Called when a user removes an acknowledgment from an alarm.
    """

    step = {
        '_t': 'ackremove',
        't': event['timestamp'],
        'a': author,
        'm': message
    }

    alarm['ack'] = None
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.useraction.cancel')
def cancel(manager, alarm, author, message, event):
    """
    Called when an alarm is canceled by a user.
    """

    step = {
        '_t': 'cancel',
        't': event['timestamp'],
        'a': author,
        'm': message
    }

    alarm['canceled'] = step
    alarm['steps'].append(step)

    return alarm, CANCELED


@register_task('alerts.useraction.uncancel')
def restore(manager, alarm, author, message, event):
    """
    Called when an alarm is restored by a user.
    """

    step = {
        '_t': 'uncancel',
        't': event['timestamp'],
        'a': author,
        'm': message
    }

    canceled = alarm['canceled']
    alarm['canceled'] = None
    alarm['steps'].append(step)

    status = None

    if manager.restore_event:
        status = get_previous_step(
            alarm,
            ['statusinc', 'statusdec'],
            ts=canceled['t']
        )

        if status is not None:
            status = status['val']
        else:
            # This is not supposed to happen since a restored alarm
            # should have a status before its cancelation
            status = OFF

    else:
        status = compute_status(manager, alarm)

    return alarm, status


@register_task('alerts.useraction.declareticket')
def declare_ticket(manager, alarm, author, message, event):
    """
    Called when a user declares a ticket for an alarm.
    """

    step = {
        '_t': 'declareticket',
        't': event['timestamp'],
        'a': author,
        'm': message,
        'val': None
    }

    alarm['ticket'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.useraction.assocticket')
def associate_ticket(manager, alarm, author, message, event):
    """
    Called when a user associates a ticket to an alarm.
    """

    step = {
        '_t': 'assocticket',
        't': event['timestamp'],
        'a': author,
        'm': message,
        'val': event['ticket']
    }

    alarm['ticket'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.useraction.changestate')
def change_state(manager, alarm, author, message, event):
    """
    Called when a user manually changes the state of an alarm.
    """

    step = {
        '_t': 'changestate',
        't': event['timestamp'],
        'a': author,
        'm': message,
        'val': event['state']
    }

    alarm['state'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.useraction.snooze')
def snooze(manager, alarm, author, message, event):
    """
    Called when a user snoozes an alarm.
    """

    until = event['timestamp'] + event.get('duration', 300)

    step = {
        '_t': 'snooze',
        't': event['timestamp'],
        'a': author,
        'm': message,
        'val': until
    }

    alarm['snooze'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.systemaction.state_increase')
def state_increase(manager, alarm, state, event):
    """
    Called when the system detects a state escalation on an alarm.
    """

    step = {
        '_t': 'stateinc',
        't': event['timestamp'],
        'a': '{0}.{1}'.format(event['connector'], event['connector_name']),
        'm': event['output'],
        'val': state
    }

    if alarm['state'] is None or alarm['state']['_t'] != 'changestate':
        alarm['state'] = step

    alarm['steps'].append(step)
    status = compute_status(manager, alarm)

    return alarm, status


@register_task('alerts.systemaction.state_decrease')
def state_decrease(manager, alarm, state, event):
    """
    Called when the system detects a state decrease on an alarm.
    """

    step = {
        '_t': 'statedec',
        't': event['timestamp'],
        'a': '{0}.{1}'.format(event['connector'], event['connector_name']),
        'm': event['output'],
        'val': state
    }

    if alarm['state'] is None or alarm['state']['_t'] != 'changestate':
        alarm['state'] = step

    alarm['steps'].append(step)
    status = compute_status(manager, alarm)

    return alarm, status


@register_task('alerts.systemaction.status_increase')
def status_increase(manager, alarm, status, event):
    """
    Called when the system detects a status escalation on an alarm.
    """

    step = {
        '_t': 'statusinc',
        't': event['timestamp'],
        'a': '{0}.{1}'.format(event['connector'], event['connector_name']),
        'm': event['output'],
        'val': status
    }

    alarm['status'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.systemaction.status_decrease')
def status_decrease(manager, alarm, status, event):
    """
    Called when the system detects a status decrease on an alarm.
    """

    step = {
        '_t': 'statusdec',
        't': event['timestamp'],
        'a': '{0}.{1}'.format(event['connector'], event['connector_name']),
        'm': event['output'],
        'val': status
    }

    alarm['status'] = step
    alarm['steps'].append(step)

    return alarm


@register_task('alerts.systemaction.update_state_counter')
def update_state_counter(alarm, diff_counter):
    """
    Create or update alarm state counter related to last status change.

    :param dict alarm: Alarm value
    :param dict diff_counter: Update existing counter with this one (or
      start with this one if no exist yet)

    :return: Alarm with updated or newly created counter
    :rtype: dict
    """

    counter_i = alarm['steps'].index(alarm['status']) + 1

    if len(alarm['steps']) == counter_i:
        # The last step is the last change of status
        counter_template = {
            '_t': 'statecounter',
            'a': alarm['status']['a'],
            't': alarm['status']['t'],
            'm': '',
            'val': {}
        }

        alarm['steps'].append(counter_template)

    elif alarm['steps'][counter_i]['_t'] != 'statecounter':
        counter_template = {
            '_t': 'statecounter',
            'a': alarm['status']['a'],
            't': alarm['status']['t'],
            'm': '',
            'val': {}
        }

        alarm['steps'].insert(counter_i, counter_template)

    counter = alarm['steps'][counter_i]['val']

    for category, diff_count in diff_counter.items():
        counter[category] = counter.get(category, 0) + diff_count

    return alarm


@register_task('alerts.systemaction.hard_limit')
def hard_limit(manager, alarm):
    """
    Called when the system detects an hard limit overtake.
    """

    step = {
        '_t': 'hardlimit',
        't': int(time()),
        'a': '__canopsis__',
        'm': (
            'This alarm has reached an hard limit ({} steps recorded) : no '
            'more steps will be appended. Please cancel this alarm or extend '
            'hard limit value.'.format(manager.hard_limit)
        ),
        'val': manager.hard_limit
    }

    alarm['hard_limit'] = step
    alarm['steps'].append(step)

    return alarm
