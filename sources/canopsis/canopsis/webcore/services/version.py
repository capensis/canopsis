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

from pymongo.errors import PyMongoError

from canopsis.common.mongo_store import MongoStore
from canopsis.version import CanopsisVersionManager
from canopsis.webcore.utils import (gen_json, gen_json_error,
                                    HTTP_NOT_FOUND, HTTP_ERROR)


def exports(ws):

    @ws.application.get(
        '/api/v2/version'
    )
    def get_canopsis_version():
        """
        Get Canopsis version.

        :returns: ``200 OK`` if success or ``404 Not Found``
                  if Canopsis version info not found or
                  ``400 Bad Request`` if database error.
        """
        try:
            store = MongoStore.get_default()
            collection = \
                store.get_collection(name=CanopsisVersionManager.COLLECTION)
            document = CanopsisVersionManager(collection).\
                find_canopsis_document()
        except PyMongoError:
            return gen_json_error(
                {"description": "can not retrieve the canopsis version from "
                                "database, contact your administrator."},
                HTTP_ERROR)

        if not document:
            return gen_json_error(
                {"description": "canopsis version info not found."},
                HTTP_NOT_FOUND)

        return gen_json(
            document, allowed_keys=[CanopsisVersionManager.VERSION_FIELD])
