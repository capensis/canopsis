# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2019 "Capensis" [http://www.capensis.fr]
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

from canopsis.common.collection import MongoCollection


class CanopsisWebhookManager(object):

    COLLECTION = "webhooks"

    def __init__(self, mongo_store):
        """
        :param collection: `pymongo.collection.Collection` object.
        """
        super(CanopsisWebhookManager, self).__init__()
        self.__mongo_store = mongo_store
        collection = self.__mongo_store.get_collection(self.COLLECTION)
        self.__collection = MongoCollection(collection)

    def get_webhook_from_id(self, wid):
        return self.__collection.find_one({'_id': wid})

    def create(self, webhook):
        return self.__collection.insert(webhook)

    def delete_webhook_from_id(self, wid):
        resp = self.__collection.remove({'_id': wid})

        return self.__collection.is_successfull(resp)
