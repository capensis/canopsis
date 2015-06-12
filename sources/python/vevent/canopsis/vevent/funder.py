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

from canopsis.vevent.manager import VEventManager
from canopsis.ctxinfo.funder import CTXInfoFunder


class VEventFunder(CTXInfoFunder):
    """In charge of binding a vevent information to context entities.
    """

    __datatype__ = 'vevent'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(VEventFunder, self).__init__(*args, **kwargs)

        self.manager = VEventManager()

    def _get_docs(self, entity_ids, query, *args, **kwargs):

        result = []

        docs = self.manager.values(sources=entity_ids, query=query)

        entity_id_field = self._entity_id_field()

        if entity_ids is not None:
            entity_ids = set(entity_ids)

        for doc in docs:
            entity_id = doc[VEventManager.SOURCE]
            if entity_ids is None or entity_id in entity_ids:
                doc[entity_id_field] = entity_id
                result.append(doc)

        return result

    def _get(self, entity_ids, query, *args, **kwargs):

        return self._get_docs(
            entity_ids=entity_ids, query=query, *args, **kwargs
        )

    def _delete(self, entity_ids, query, *args, **kwargs):

        result = self._get_docs(entity_ids=entity_ids, query=query)

        self.manager.remove_by_source(sources=entity_ids, query=query)

        return result

    def entity_ids(self, query=None):

        return self.manager.whois(query=query)
