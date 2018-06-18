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
Adapter for view object.
"""

from __future__ import unicode_literals

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.views.enums import ViewField


class ViewAdapter(object):
    """
    Adapter for the view collection.
    """

    VIEWS_COLLECTION = 'views'

    def __init__(self):
        self.views_collection = MongoCollection(
            MongoStore.get_default().get_collection(self.VIEWS_COLLECTION))

    def get_by_id(self, view_id):
        """
        Get a view given its id.

        :param str view_id: the id of the view.
        """
        return self.views_collection.find_one({
            ViewField.id: view_id
        })

    def create(self, view):
        """
        Create a new view and return its id.

        :param Dict view:
        :rtype: str
        """
        return self.views_collection.insert(view)

    def update(self, view_id, view):
        """
        Update a view given its id
        """
        self.views_collection.update({
            ViewField.id: view_id
        }, view, upsert=False)

    def remove_with_id(self, view_id):
        """
        Remove a view given its id.
        """
        self.views_collection.remove({
            ViewField.id: view_id
        })
