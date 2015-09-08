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

"""
Task condition functions such as duration/rrule condition, switch, all and any.
"""

from canopsis.common.init import basestring
from canopsis.task.core import register_task, run_task

from time import time
from datetime import datetime

from dateutil.relativedelta import relativedelta
from dateutil.rrule import rrule as rrule_class


@register_task
def during(rrule, duration=None, timestamp=None, **kwargs):
    """
    Check if input timestamp is in rrule+duration period

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
def _any(confs=None, **kwargs):
    """
    True iif at least one input condition is equivalent to True.

    :param confs: confs to check.
    :type confs: list or dict or str
    :param kwargs: additional task kwargs.

    :return: True if at least one condition is checked (compared to True, but
            not an strict equivalence to True). False otherwise.
    :rtype: bool
    """

    result = False

    if confs is not None:
        # ensure confs is a list
        if isinstance(confs, (basestring, dict)):
            confs = [confs]
        for conf in confs:
            result = run_task(conf, **kwargs)
            if result:  # leave function as soon as a result if True
                break

    return result


@register_task
def _all(confs=None, **kwargs):
    """
    True iif all input confs are True.

    :param confs: confs to check.
    :type confs: list or dict or str
    :param kwargs: additional task kwargs.

    :return: True if all conditions are checked. False otherwise.
    :rtype: bool
    """

    result = False

    if confs is not None:
        # ensure confs is a list
        if isinstance(confs, (basestring, dict)):
            confs = [confs]
        # if at least one conf exists, result is True by default
        result = True
        for conf in confs:
            result = run_task(conf, **kwargs)
            # stop when a result is False
            if not result:
                break

    return result


STATEMENT = 'statement'


@register_task
def _not(condition=None, **kwargs):
    """
    Return the opposite of input condition.

    :param condition: condition to process.

    :result: not condition.
    :rtype: bool
    """

    result = True

    if condition is not None:
        result = not run_task(condition, **kwargs)

    return result


@register_task
def condition(condition=None, statement=None, _else=None, **kwargs):
    """
    Run an statement if input condition is checked and return statement result.

    :param condition: condition to check.
    :type condition: str or dict
    :param statement: statement to process if condition is checked.
    :type statement: str or dict
    :param _else: else statement.
    :type _else: str or dict
    :param kwargs: condition and statement additional parameters.

    :return: statement result.
    """

    result = None

    checked = False

    if condition is not None:
        checked = run_task(condition, **kwargs)

    if checked:  # if condition is checked
        if statement is not None:  # process statement
            result = run_task(statement, **kwargs)
    elif _else is not None:  # else process _else statement
        result = run_task(_else, **kwargs)

    return result


@register_task
def switch(
        confs=None, remain=False, all_checked=False, _default=None, **kwargs
):
    """
    Execute first statement among conf where task result is True.
    If remain, process all statements conf starting from the first checked
    conf.

    :param confs: task confs to check. Each one may contain a task action at
        the key 'action' in conf.
    :type confs: str or dict or list
    :param bool remain: if True, execute all remaining actions after the
        first checked condition.
    :param bool all_checked: execute all statements where conditions are
        checked.
    :param _default: default task to process if others have not been checked.
    :type _default: str or dict

    :return: statement result or list of statement results if remain.
    :rtype: list or object
    """

    # init result
    result = [] if remain else None

        # check if remain and one task has already been checked.
    remaining = False

    if confs is not None:
        if isinstance(confs, (basestring, dict)):
            confs = [confs]
        for conf in confs:
            # check if task has to be checked or not
            check = remaining
            if not check:
                # try to check current conf
                check = run_task(conf=conf, **kwargs)
            # if task is checked or remaining
            if check:
                if STATEMENT in conf:  # if statements exist, run them
                    statement = conf[STATEMENT]
                    statement_result = run_task(statement, **kwargs)
                    # save result
                    if not remain:  # if not remain, result is statement_result
                        result = statement_result
                    else:  # else, add statement_result to result
                        result.append(statement_result)
                # if remain
                if remain:
                    # change of remaining status
                    if not remaining:
                        remaining = True
                elif all_checked:
                    pass
                else:  # leave execution if one statement has been executed
                    break

    # process _default statement if necessary
    if _default is not None and (remaining or (not result) or all_checked):
        last_result = run_task(_default, **kwargs)
        if not remain:
            result = last_result
        else:
            result.append(last_result)

    return result
