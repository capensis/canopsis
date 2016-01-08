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

"""This module is dedicated to CTXPropRegistry interface.

It aims to execute ctxprop manager read and delete methods related to its own
information managers such as the perfdata manager for example.
"""

from copy import deepcopy

from canopsis.storage.core import Cursor
from canopsis.mongo.core import MongoStorage
from canopsis.ctxprop.registry import CTXPropRegistry


class CTXOldRegistry(CTXPropRegistry):
    """In charge of binding an old collection (not generated from the Storage
    project) to context entities.

    Such old collections contain all entity fields.
    """

    def __init__(self, table=None, *args, **kwargs):
        """
        :param Storage storage: event storage to use.
        """

        super(CTXOldRegistry, self).__init__(*args, **kwargs)

        self.storage = MongoStorage(table=table)

    def _do(self, command, ids, query, queryname, *args, **kwargs):
        """Execute storage command related to input ids and query.

        :param command: storage command.
        :param list ids: entity id(s).
        :param dict query: storage command query.
        :param str queryname: storage command query parameter.
        :return: storage command execution documents.
        :rtype: list
        """

        result = []
        # initialize query
        query = deepcopy(query) if query else {}
        # get entity id field name
        ctx_id_field = self._ctx_id_field()

        if ids is None:
            # get all existing entity ids
            ids = []
            events = self.storage.find_elements()
            for event in events:
                entity = self.context.get_entity(event)
                entity_id = self.context.get_entity_id(entity)
                ids.append(entity_id)

        for entity_id in ids:
            # get entity
            entity = self.context.get_entity_by_id(entity_id)
            cleaned_entity = self.context.clean(entity)
            cleaned_entity['source_type'] = cleaned_entity.pop('type')
            for ctx in self.context.context[1:]:
                if ctx in cleaned_entity:
                    continue
                else:
                    cleaned_entity[ctx] = cleaned_entity.pop(Context.NAME)
                    break
            # update query with entity information
            _query = deepcopy(query)
            _query.setdefault('$and', []).append(cleaned_entity)
            # update kwargs with query and queryname
            kwargs[queryname] = _query
            # execute the storage command
            documents = command(*args, **kwargs)
            if isinstance(documents, Cursor):
                documents = list(documents)
                # update entity_id in documents
                for document in documents:
                    document[ctx_id_field] = entity_id
            else:
                documents = [
                    {ctx_id_field: entity_id, 'result': documents}
                ]
            # add all documents into the result
            result += documents

        return result

    def _get(self, ids, query, *args, **kwargs):

        return self._do(
            command=self.storage.find_elements,
            ids=ids, query=query, queryname='query',
            *args, **kwargs
        )

    def _count(self, ids, query, *args, **kwargs):

        return self._do(
            command=self.storage.count_elements,
            ids=ids, query=query, queryname='query',
            *args, **kwargs
        )

    def _delete(self, ids, query, *args, **kwargs):

        return self._do(
            command=self.storage.remove_elements,
            ids=ids, query=query, queryname='_filter',
            *args, **kwargs
        )

    def ids(self, query=None, *args, **kwargs):

        result = []

        events = self.storage.find_elements(query=query)

        for event in events:
            entity = self.context.get_entity(event)
            entity_id = self.context.get_entity_id(entity)
            result.append(entity_id)

        return result


class CTXEventRegistry(CTXOldRegistry):
    """Provider bound to the old events collection.
    """

    __datatype__ = 'events'  #: default datatype name

    def __init__(self, table=__datatype__, *args, **kwargs):

        super(CTXEventRegistry, self).__init__(table=table, *args, **kwargs)


class CTXEventLogRegistry(CTXOldRegistry):
    """Provider bound to the old events_log collection.
    """

    __datatype__ = 'eventslog'  #: default datatype name

    def __init__(self, table='events_log', *args, **kwargs):

        super(CTXEventLogRegistry, self).__init__(table=table, *args, **kwargs)
