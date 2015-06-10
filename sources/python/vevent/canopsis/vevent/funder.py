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

    def _do(self, cmd, entity_ids):

        result = []

        for entity_id in entity_ids:
            cmdresult = cmd(sources=entity_id)
            result.append(cmdresult)

        return result

    def _get(self, entity_ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.values, entity_ids=entity_ids)

    def _delete(self, entity_ids, query, *args, **kwargs):

        return self._do(
            cmd=self.manager.remove_by_source, entity_ids=entity_ids
        )
