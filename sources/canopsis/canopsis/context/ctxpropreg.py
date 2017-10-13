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

from canopsis.context.manager import Context
from canopsis.ctxprop.registry import CTXPropRegistry


class CTXContextRegistry(CTXPropRegistry):
    """In charge of contextual context properties.
    """

    __datatype__ = 'context'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(CTXContextRegistry, self).__init__(*args, **kwargs)

        self.manager = Context()

    def _get_documents(self, ids, query):

        query[Context.DATA_ID] = ids

        return self.manager.find(_filter=query)

    def _get(self, ids, query, *args, **kwargs):

        return self._get_documents(ids=ids, query=query)

    def _delete(self, ids, query, *args, **kwargs):

        docs = self._get_documents(ids=ids, query=query)

        ids = [doc[Context.DATA_ID] for doc in docs]

        self.manager.remove(ids=ids)

        return docs

    def ids(self, query=None):

        result = set()

        elts = self.manager.find(_filter=query)

        for elt in elts:
            result.add(elt[Context.DATA_ID])

        return list(result)
