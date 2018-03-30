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


class Filter(object):
    """
    Allow canopsis to generate custom mongo filters.
    """

    def make_filter(self, mfilter={}, includes=[], excludes=[]):

        """
        Generate an appropriate suitable filter for canopsis depending on
        inclusion elements, exclusion elements and an existing filter
        """

        final_filter = {}
        subfilter = {}

        minclude = self.add_filter_by_id(includes, True)
        mexclude = self.add_filter_by_id(excludes, False)

        if mfilter:
            if minclude:
                subfilter = {'$or': [mfilter, minclude]}
            else:
                subfilter = mfilter
        elif minclude:
            subfilter = minclude

        # When at least one exclusion exists
        if mexclude:
            if not mfilter and not minclude:
                final_filter = mexclude
            else:
                final_filter['$and'] = [subfilter, mexclude]
        elif subfilter:
            final_filter = subfilter

        return final_filter

    def add_filter_by_id(self, id_list, inclusion):
        """
        generate a filter depending on the given id list
        """

        multiop = '$in' if inclusion else '$nin'
        singleop = '$eq' if inclusion else '$ne'

        mfilter = {}
        if len(id_list) == 1:
            mfilter['_id'] = {singleop: id_list[0]}
        elif len(id_list) > 1:
            mfilter['_id'] = {multiop: id_list}

        return mfilter
