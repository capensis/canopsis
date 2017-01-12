# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)
from canopsis.middleware.registry import MiddlewareRegistry

# Might be useful when dealing with rrules
# from dateutil.rrule import rrulestr

CONF_PATH = 'pbehavior/pbehavior.conf'
CATEGORY = 'PBEHAVIOR'


@conf_paths(CONF_PATH)
@add_category(CATEGORY)
class PBehaviorManager(MiddlewareRegistry):
    def create(
            self,
            name, filter_, author,
            tstart, tstop, rrule=None,
            enabled=True,
            connector='canopsis', connector_name='canopsis'
    ):
        self['pbehavior_storage']._backend.insert({'test': name})

        return 'created'

    def get_behaviors(self, entity_id):
        """
        Return all pbehaviors related to an entity_id, sorted by descending
        dtstart.

        :param str entity_id: Id for which behaviors have to be returned

        :return: List of pbehaviors as dict, with name, dtstart, dtend, rrule
          and enabled keys
        :rtype: list of dict
        """

        return []

    def compute_pbehaviors_filters(self):
        return 'computing...'
