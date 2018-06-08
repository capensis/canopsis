# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals


class StatsAPIError(Exception):
    """
    A StatsAPIError is an Exception that can be raised by the statistics API.

    It should be handled in `webcore/services/statsng.py`, and returned as a
    JSON object in the response.
    """
    def __init__(self, message):
        super(StatsAPIError, self).__init__(message)
        self.message = message


class UnknownStatNameError(StatsAPIError):
    """
    A UnknownStatNameError is an Exception that can be raised by a StatsAPI
    object when requesting an unknown statistic.
    """
    pass
