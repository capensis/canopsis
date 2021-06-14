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
from canopsis.ctxprop.registry import CTXPropRegistry


class CTXVEventRegistry(CTXPropRegistry):
    """In charge of binding a vevent information to context entities.
    """

    __datatype__ = 'vevent'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(CTXVEventRegistry, self).__init__(*args, **kwargs)

        self.manager = VEventManager()

    def _get_docs(self, ids, query, *args, **kwargs):

        result = []

        docs = self.manager.values(sources=ids, query=query)

        ctx_id_field = self._ctx_id_field()

        if ids is not None:
            ids = set(ids)

        for doc in docs:
            entity_id = doc[VEventManager.SOURCE]
            if ids is None or entity_id in ids:
                doc[ctx_id_field] = entity_id
                result.append(doc)

        return result

    def _get(self, ids, query, *args, **kwargs):

        return self._get_docs(
            ids=ids, query=query, *args, **kwargs
        )

    def _delete(self, ids, query, *args, **kwargs):

        result = self._get_docs(ids=ids, query=query)

        self.manager.remove_by_source(sources=ids, query=query)

        return result

    def ids(self, query=None):

        return self.manager.whois(query=query)
