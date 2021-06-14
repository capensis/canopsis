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

from canopsis.check.manager import CheckManager
from canopsis.ctxprop.registry import CTXPropRegistry


class CTXCheckRegistry(CTXPropRegistry):
    """In charge of contextual state properties.
    """

    __datatype__ = 'check'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(CTXCheckRegistry, self).__init__(*args, **kwargs)

        self.manager = CheckManager()

    def _get_documents(self, ids, query):

        result = []

        ctx_id_field = self._ctx_id_field()

        docs = self.manager.state(ids=ids, query=query)
        for doc in docs:
            doc[ctx_id_field] = doc[CheckManager.ID]
            result.append(doc)

        return result

    def _get(self, ids, query, *args, **kwargs):

        return self._get_documents(ids=ids, query=query)

    def _delete(self, ids, query, *args, **kwargs):

        docs = self._get_documents(ids=ids, query=query)

        ids = [doc[CheckManager.ID] for doc in docs]

        self.manager.del_state(ids=ids, query=query)

        return docs

    def ids(self, query=None):

        result = set()

        elts = self.manager.state(query=query)

        for elt in elts:
            result.add(elt[CheckManager.ID])

        return list(result)
