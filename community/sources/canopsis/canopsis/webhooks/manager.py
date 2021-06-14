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

from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
import re
import time


class WebhookManager(object):
    """
    Manager for the webhooks
    """

    COLLECTION = "webhooks"

    def __init__(self, mongo_collection):
        """
        :param mongo_collection: `pymongo.collection.Collection` object.
        """
        super(WebhookManager, self).__init__()
        self.__collection = mongo_collection

    @classmethod
    def default_collection(cls):
        """
        Returns the default collection for the manager.

        ! Do not use in tests !

        :rtype: canopsis.common.collection.MongoCollection
        """
        store = MongoStore.get_default()
        collection = store.get_collection(name=cls.COLLECTION)
        return MongoCollection(collection)

    def get_webhook_list(self):
        """
        Return a list of all the webhooks.

        :rtype: List[Dict[str, Any]]
        """
        return list(self.__collection.find({}))

    def get_webhook_by_id(self, wid):
        """
        Get a webhook given its id.

        :param str rule_id: the id of the webhook.
        :rtype: Dict[str, Any]
        """
        return self.__collection.find_one({'_id': wid})

    def create_webhook(self, webhook):
        """
        Create a webhook and return its id.

        :param Dict[str, Any] webhook:
        :rtype: str
        :raises: CollectionError if the creation fails.
        """
        return self.__collection.insert(webhook)

    def update_webhook_by_id(self, webhook, wid):
        """
        Update a webhook given its id.
        Return a boolean if the operation is successful.

        :param str wid: the id of the webhook.
        :param Dict[str, Any] webhook:
        :rtype: bool
        :raises: CollectionError if the update fails.
        """
        current_webhook = self.__collection.find_one({'_id': wid})
        if 'creation_date' in current_webhook:
            webhook['creation_date'] = current_webhook['creation_date']
        now = int(time.time())
        webhook['last_update_date'] = now
        resp = self.__collection.update(query={'_id': wid}, document=webhook)
        return self.__collection.is_successfull(resp)

    def delete_webhook_by_id(self, wid):
        """
        Remove a webhook given its id.
        Return a boolean if the operation is successful.

        :param str wid: the id of the rule.
        :rtype: bool
        :raises: CollectionError if the deletion fails.
        """
        resp = self.__collection.remove({'_id': wid})
        return self.__collection.is_successfull(resp)

    def read(self, _id, search, limit, skip):
        """
        Retrieve webhook by search phrase
        :param _id:
        :param search:
        :param limit:
        :param skip:
        :return:
        """
        pipeline = []
        if _id is None:
            if search is not None:
                regex_search = re.compile(str(search), re.IGNORECASE)
                or_query = [
                    {"_id": regex_search},
                    {"request.url": regex_search},
                    {"request.payload": regex_search},
                    {"author": regex_search}
                ]
                pipeline.append({"$match": {"$or": or_query}})
            else:
                pipeline.append({"$match": {}})
        else:
            pipeline.append({"$match": {"_id": _id}})

        total_count_data = list(self.__collection.aggregate(
            pipeline + [{'$count': 'total_count'}]))

        if len(total_count_data) == 1:
            try:
                total_count = total_count_data[0]["total_count"]
            except (IndexError, KeyError):
                return {"total_count": 0, "count": 0, "data": []}
        else:
            return {"total_count": 0, "count": 0, "data": []}

        if _id is None:
            if skip is not None:
                pipeline.append({"$skip": skip})
            if limit is not None:
                pipeline.append({"$limit": limit})

        webhooks = list(self.__collection.aggregate(pipeline))

        return {"total_count": total_count,
                "count": len(webhooks),
                "data": webhooks}
