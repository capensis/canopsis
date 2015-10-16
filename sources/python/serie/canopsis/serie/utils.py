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


def build_filter_from_regex(self, regex):
    regex_parts = regex.split(' ')
    regex = {
        'component': [],
        'resource': [],
        'name': []
    }

    for part in regex_parts:
        if part.startswith('co:'):
            regex['component'].append({'$regex': part[3:]})

        elif part.startswith('re:'):
            regex['resource'].append({'$regex': part[3:]})

        elif part.startswith('me:'):
            regex['name'].append({'$regex': part[3:]})

        else:
            for key in regex.keys():
                regex[key].append({'$regex': part})

    mfilter = {'$and': []}

    for key in regex:
        if len(regex[key]) > 0:
            local_mfilter = {'$or': [
                subfilter for subfilter in regex[key]
            ]}

            if len(local_mfilter['$or']) == 1:
                local_mfilter = local_mfilter['$or'][0]

            mfilter['$and'].append(local_mfilter)

    return mfilter
