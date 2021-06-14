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

from canopsis.linklist.manager import Linklist
from canopsis.ctxprop.registry import CTXPropRegistry
from canopsis.mongo.core import MongoStorage
from canopsis.context.manager import Context

from json import loads


class CTXLinklistRegistry(CTXPropRegistry):
    """In charge of ctx linklist information.
    """

    __datatype__ = 'linklist'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(CTXLinklistRegistry, self).__init__(*args, **kwargs)

        self.manager = Linklist()
        self.events = MongoStorage(table='events')
        self.context = Context()

    def _get_documents(self, ids, query):
        """Get documents related to input ids and query.

        :param list ids: entity ids. If None, get all documents.
        :param dict query: additional selection query.
        :return: list of documents.
        :rtype: list
        """
        result = []
        # get entity id field name
        ctx_id_field = self._ctx_id_field()
        # get a set of entity ids for execution speed reasons
        if ids is not None:
            ids = set(ids)
        # get documents
        docs = self.manager.find(_filter=query)
        for doc in docs:
            try:
                mfilter = loads(doc['mfilter'])
            except Exception:
                pass
            else:  # get entities from events
                events = self.events.find_elements(query=mfilter)
                for event in events:
                    entity = self.context.get_entity(event)
                    entity_id = self.context.get_entity_id(entity)
                    if ids is None or entity_id in ids:
                        doc[ctx_id_field] = entity_id  # add eid to the doc
                        result.append(doc)

        return result

    def _get(self, ids, query, *args, **kwargs):

        return self._get_documents(ids=ids, query=query)

    def _delete(self, ids, query, *args, **kwargs):

        result = self._get_documents(ids=ids, query=query)

        ids = [doc['_id'] for doc in result]

        self.manager.remove(ids=ids)

        return result

    def ids(self, query=None):

        result = []

        documents = self.manager.find(_filter=query)

        for document in documents:
            try:
                mfilter = loads(document['mfilter'])
            except Exception:
                pass
            else:
                # get entities from events
                events = self.events.find_elements(query=mfilter)
                for event in events:
                    entity = self.context.get_entity(event)
                    entity_id = self.context.get_entity_id(entity)
                    result.append(entity_id)

        return result
