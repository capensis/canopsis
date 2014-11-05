# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
Rule condition functions
"""

from canopsis.task import get_task_with_params, register_task

from time import time
from datetime import datetime

from dateutil.relativedelta import relativedelta
from dateutil.rrule import rrule as rrule_class


@register_task
def during(event, ctx, rrule, duration=None, timestamp=None, **kwargs):
    """
    Check if input timestamp is in rrule+duration period

    :param dict event: event to process
    :param dict ctx: rule context
    :param rrule: rrule to check
    :type rrule: str or dict
        (freq, dtstart, interval, count, wkst, until, bymonth, byminute, etc.)
    :param dict duration: time duration from rrule step. Ex:{'minutes': 60}
    :param float timestamp: timestamp to check between rrule+duration. If None,
        use now
    """

    result = False

    # if rrule is a string expression
    if isinstance(rrule, basestring):
        rrule_object = rrule_class.rrulestr(rrule)
    else:
        rrule_object = rrule_class(**rrule)

    # if timestamp is None, use now
    if timestamp is None:
        timestamp = time()

    # get now object
    now = datetime.fromtimestamp(timestamp)

    # get delta object
    duration_delta = now if duration is None else relativedelta(**duration)

    # get last date
    last_date = rrule_object.before(now, inc=True)

    # if a previous date exists
    if last_date is not None:
        next_date = last_date + duration_delta

        # check if now is between last_date and next_date
        result = last_date <= now <= next_date

    return result


@register_task
def any(event, ctx, conditions=None, **kwargs):
    """
    True if at least one input condition is True
    """

    result = False

    if conditions is not None:

        for condition in conditions:
            condition_task, params = get_task_with_params(condition)

            result = condition_task(event=event, ctx=ctx, **params)

            if result:
                break

    return result


@register_task
def all(event, ctx, conditions=None, **kwargs):
    """
    True iif all input conditions is True
    """

    result = False

    if conditions is not None:

        result = True

        for condition in conditions:
            condition_task, params = get_task_with_params(condition)

            result = condition_task(event=event, ctx=ctx, **params)

            if not result:
                break

    return result
