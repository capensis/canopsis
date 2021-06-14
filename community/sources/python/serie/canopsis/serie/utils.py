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

"""Module of Serie utilities."""

NAME_BY_KEYS = {
    'co:': 'component',
    're:': 'resource',
    'me:': 'name'
}  #: entity name by regex entity header.


def build_filter_from_regex(regex):
    """
    Transform a metric filter into a MongoDB-like filter.

    co:<regex> --> {'component': {'$regex': '<regex>'}}
    re:<regex> --> {'resource': {'$regex': '<regex>'}}
    me:<regex> --> {'metric': {'$regex': '<regex>'}}

    co:<regex1> co:<regex2> --> {'$or': [...]}
    co:<regex> re:<regex> --> {'$and': [...]}

    :param regex: Metric filter to transform
    :type regex: str

    :returns: MongoDB-like filter as dict
    """

    regex_parts = regex.split(' ')

    regex = {NAME_BY_KEYS[key]: [] for key in NAME_BY_KEYS}

    for part in regex_parts:
        for key in NAME_BY_KEYS:
            if part.startswith(key):
                name = NAME_BY_KEYS[key]
                regex[name].append({'$regex': part[len(key):]})
                break

        else:
            for name in regex.keys():
                regex[name].append({'$regex': part})

    mfilter = {'$and': []}

    for key in regex:
        items = regex[key]
        if items:
            local_mfilter = {'$or': [
                {key: subfilter} for subfilter in items
            ]}

            if len(local_mfilter['$or']) == 1:
                local_mfilter = local_mfilter['$or'][0]

            mfilter['$and'].append(local_mfilter)

    if len(mfilter['$and']) == 0:
        mfilter = {}

    elif len(mfilter['$and']) == 1:
        mfilter = mfilter['$and'][0]

    return mfilter
