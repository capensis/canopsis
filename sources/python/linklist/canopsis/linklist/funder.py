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
from canopsis.ctxinfo.funder import CTXInfoFunder


class LinklistFunder(CTXInfoFunder):
    """In charge of binding a linklist information to context entities.
    """

    __datatype__ = 'linklist'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(LinklistFunder, self).__init__(*args, **kwargs)

        self.manager = Linklist()

    def _do(self, cmd, entity_ids):

        result = []

        for entity_id in entity_ids:
            cmdresult = cmd(ids=entity_id)
            result.append(cmdresult)

        return result

    def _get(self, entity_ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.find, entity_ids=entity_ids)

    def _delete(self, entity_ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.remove, entity_ids=entity_ids)
