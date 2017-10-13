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
OPERATORS = {
    '$lt': lt,
    '$lte': le,
    '$gt': gt,
    '$gte': ge,
    '$ne': ne,
    '$eq': eq}


def field_check(condition, event, key):

    for op in condition[key]:
        if op == '$exists':
            #check if key is in event
            if condition[key][op]:
                if key not in event:
                    return False
            #check if key is not in event
            else:
                if key in event:
                    return False

        elif op in {'$eq', '$ne', '$gt', '$gte', '$lt', '$lte'}:
            if not OPERATORS[op](event[key], condition[key][op]):
                return False

        elif op == '$regex' or (
                op == '$options' and "$regex" in condition[key]):
            if not regex_match(
                event[key], condition[key]["$regex"], condition[key].get(
                    "$options", None)):
                return False

        elif op == '$in':
            if event[key] not in condition[key][op]:
                return False

        elif op == '$nin':
            if event[key] in condition[key][op]:
                return False

        elif op == '$not':
            reverse_mfilter = {}
            reverse_mfilter[key] = condition[key][op]

            if field_check(reverse_mfilter, event, key):
                return False

        elif op == '$all':
            items = event[key]

            # If event[key] isn't a list, treat it as if it was
            if not isinstance(items, list):
                items = [items]

            # Check if all items from condition[key]['$all'] are in event[key]
            for item in condition[key][op]:
                if item not in items:
                    return False

        else:
            if event[key] != condition[key]:
                return False

    return True


def check(event, ctx, condition, **params):
    """
    Check if input event matches input condition.
    """

    # Check connector_name
    if 'connector_name' in condition and 'connector_name' in event:
        if condition['connector_name'] != event['connector_name']:
            return False

    # For each key of filter
    for key in condition:
        if key == '$and':
            # Check match for each elements in the list
            if isinstance(condition[key], list):
                result = True
                for sub_filter in condition[key]:
                    result = result and check(sub_filter, event)
                return result
            else:
                for element in condition[key]:
                    # If one does not match, then return False
                    if not check(element, event):
                        return False

        elif key == '$or':
            # Check match for each elements in the list
            if isinstance(condition[key], list):
                result = True
                #testing len of filter
                if len(condition[key]):
                    result = check(condition[key][0], event)
                    for sub_filter in condition[key][1:]:
                        result = result or check(sub_filter, event)
                return result
            else:
                for element in condition[key]:
                    # If one match, then return True
                    if check(element, event):
                        return True

            # Here nothing matched, then return False
            return False

        elif key == '$nor':
            # Check match for each elements in the list

            for element in condition[key]:
                # If one match, then return False
                if check(element, event):
                    return False

        # For each other case, just test the equality
        elif key in event and event[key]:
            if isinstance(condition[key], dict):
                if isinstance(event[key], dict) \
                    or isinstance(event[key], list) \
                        and '$in' in condition[key]:
                    if isinstance(event[key], list) \
                            and isinstance(event[key][0], dict):
                        #For each elem of event[key], check if it's in condition[key]['$in']
                        if not [x for x in event[key] if any(
                                y in x['name'] for y in condition[key]['$in'])
                                ]:
                            return False
                    else:
                        if not [x for x in event[key] if any(
                                y in x for y in condition[key]['$in'])]:
                            return False

                elif isinstance(event[key], dict) \
                    or isinstance(event[key], list) \
                        and '$nin' in condition[key]:
                    if isinstance(event[key], list) \
                            and isinstance(event[key][0], dict):
                        #For each elem of event[key], check if it's in condition[key]['$nin']
                        if [x for x in event[key] if any(
                                y in x['name'] for y in
                                condition[key]['$in'])]:
                            return False
                        else:
                            if [x for x in event[key] if any(
                                    y in x for y in condition[key]['$in'])]:
                                return False

                elif '$in' in condition[key]:
                    if event[key] not in condition[key]['$in']:
                        return False

                elif '$nin' in condition[key]:
                    if event[key] in condition[key]['$nin']:
                        return False

                else:
                    if field_check(condition, event, key):
                        continue

                    elif event[key] != condition[key]:
                        return False

            else:
                if event[key] != condition[key]:
                    return False

        else:
            return False

    # If we arrive here, everything matched
    return True


def regex_computeoptions(options):
    if isinstance(options, basestring):
        if "i" in options:
            return I

    return 0


def regex_match(phrase, pattern, options=None):
    options = regex_computeoptions(options)

    if not phrase or not pattern or not options:
        return False

    return bool(search(str(pattern), str(phrase), options))
