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

# MongoDB Operators:
# http://docs.mongodb.org/manual/reference/operator/

from re import I, search
from operator import lt, le, gt, ge, ne, eq

from canopsis.common.utils import get_sub_key_raise as get_sub_key

def field_check(mfilter, event, key):
    cond = {'$lt': lt,
            '$lte': le,
            '$gt': gt,
            '$gte': ge,
            '$ne': ne,
            '$eq': eq}

    for op in mfilter[key]:
        if op == '$exists':
            # check if key is in event
            if mfilter[key][op]:
                if key not in event:
                    return False
            # check if key is not in event
            else:
                if key in event:
                    return False

        elif op in ['$eq', '$ne', '$gt', '$gte', '$lt', '$lte']:
            if not cond[op](get_sub_key(event, key), mfilter[key][op]):
                return False

        elif op == '$regex':
            if not regex_match(
                    get_sub_key(event, key), mfilter[key]["$regex"],
                    mfilter[key].get("$options", None)
                ):
                return False

        elif op == '$notregex':
            if regex_match(
                get_sub_key(event, key), mfilter[key]["$notregex"],
                mfilter[key].get("$options", None)
            ):
                return False

        elif op == '$options' and (
            '$regex' in mfilter[key]
            or '$notregex' in mfilter[key]
        ):
            pass

        elif op == '$in':
            if get_sub_key(event, key) not in mfilter[key][op]:
                return False

        elif op == '$nin':
            if get_sub_key(event, key) in mfilter[key][op]:
                return False

        elif op == '$not':
            if isinstance(mfilter[key]['$not'], dict):
                reverse_mfilter = {}
                reverse_mfilter[key] = mfilter[key][op]

                if field_check(reverse_mfilter, event, key):
                    return False

            elif regex_match(
                get_sub_key(event, key), mfilter[key]["$not"],
                mfilter[key].get("$options", None)
            ):
                return False

        elif op == '$all':
            items = get_sub_key(event, key)

            # If get_sub_key(event, key) isn't a list, treat it as if it was
            if not isinstance(items, list):
                items = [items]

            # Check if all items from mfilter[key]['$all'] are in get_sub_key(event, key)
            for item in mfilter[key][op]:
                if item not in items:
                    return False

        else:
            if get_sub_key(event, key) != mfilter[key]:
                return False

    return True


def _check(mfilter, event):
    # For each key of filter
    for key in mfilter:
        if key == '$and':
            # Check match for each elements in the list
            if isinstance(mfilter[key], list):
                result = True

                for sub_filter in mfilter[key]:
                    result = result and check(sub_filter, event)

                return result

            else:
                for element in mfilter[key]:
                    # If one does not match, then return False
                    if not check(element, event):
                        return False

        elif key == '$or':
            # Check match for each elements in the list
            if isinstance(mfilter[key], list):
                result = True

                # testing len of filter
                if len(mfilter[key]):
                    result = check(mfilter[key][0], event)

                    for sub_filter in mfilter[key][1:]:
                        result = result or check(sub_filter, event)

                return result

            else:
                for element in mfilter[key]:
                    # If one match, then return True
                    if check(element, event):
                        return True

            # Here nothing matched, then return False
            return False

        elif key == '$nor':
            # Check match for each elements in the list
            for element in mfilter[key]:
                # If one match, then return False
                if check(element, event):
                    return False

        # For each other case, just test the equality
        else:
            event_value = None
            try:
                event_value = get_sub_key(event, key)
            except KeyError:
                pass

            if isinstance(mfilter[key], dict):
                if (
                    (
                        isinstance(event_value, dict)
                        or isinstance(event_value, list)
                    )
                    and '$in' in mfilter[key]
                ):
                    if (
                        isinstance(event_value, list)
                        and event_value
                        and isinstance(event_value[0], dict)
                    ):
                        l = len([
                            x
                            for x in event_value
                            if any(
                                y in x['name']
                                for y in mfilter[key]['$in']
                            )
                        ])

                        # For each elem of get_sub_key(event, key),
                        # check if it's in mfilter[key]['$in']
                        if not l:
                            return False

                    else:
                        l = len([
                            x
                            for x in event_value
                            if any(
                                y in x
                                for y in mfilter[key]['$in']
                            )
                        ])

                        if not l:
                            return False

                elif (
                    (
                        isinstance(event_value, dict)
                        or isinstance(event_value, list)
                    )
                    and '$nin' in mfilter[key]
                ):
                    if (
                        isinstance(event_value, list)
                        and isinstance(event_value[0], dict)
                    ):
                        l = len([
                            x
                            for x in event_value
                            if any(
                                y in x['name']
                                for y in mfilter[key]['$in']
                            )
                        ])

                        #For each elem of get_sub_key(event, key),
                        # check if it's in mfilter[key]['$nin']
                        if l:
                            return False

                        else:
                            l = len([
                                x
                                for x in event_value
                                if any(
                                    y in x
                                    for y in mfilter[key]['$in']
                                )
                            ])

                            if l:
                                return False

                elif '$in' in mfilter[key]:
                    ev = event_value
                    if ev is None:
                        ev = ''

                    if ev not in mfilter[key]['$in']:
                        return False

                elif '$nin' in mfilter[key]:
                    ev = event_value
                    if ev is None:
                        ev = ''

                    if ev in mfilter[key]['$nin']:
                        return False

                else:
                    if field_check(mfilter, event, key):
                        continue
                    elif event_value != mfilter[key]:
                        return False

            else:
                if event_value != mfilter[key]:
                    return False

    # If we arrive here, everything matched
    return True


def check(mfilter, event):
    try:
        return _check(mfilter, event)
    except KeyError:
        return False


def regex_computeoptions(options):
    if isinstance(options, basestring):
        if "i" in options:
            return I

    return 0


def regex_match(phrase, pattern, options=None):
    options = regex_computeoptions(options)
    if phrase is None or pattern is None or options is None:
        return False
    return bool(search(
        pattern.encode('utf-8'),
        phrase.encode('utf-8'),
        options
    ))
