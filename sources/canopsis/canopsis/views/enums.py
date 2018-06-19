#!/usr/bin/env python
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

"""
Enumerations for views.
"""

from __future__ import unicode_literals

from canopsis.common.enumerations import FastEnum


class ViewField(FastEnum):
    """
    The ViewField enumeration defines the names of the fields of a view.
    """
    id = "_id"
    group_id = "group_id"


class GroupField(FastEnum):
    """
    The GroupField enumeration defines the names of the fields of a group.
    """
    id = "_id"
