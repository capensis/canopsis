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

from canopsis.snmp.manager import SnmpManager
from canopsis.ctxinfo.funder import CTXInfoFunder


class SNMPFunder(CTXInfoFunder):
    """In charge of binding a snmp information to context entities.
    """

    __datatype__ = 'snmp'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(SNMPFunder, self).__init__(*args, **kwargs)

        self.manager = SnmpManager()

    def _do(self, cmd, entity_ids, *args, **kwargs):

        result = []

        if entity_ids is None:
            result = list(cmd(*args, **kwargs))
        else:
            for entity_id in entity_ids:
                cmdresult = cmd(oids=entity_id, *args, **kwargs)
                result.append(cmdresult)

        return result

    def _get(self, entity_ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.get, entity_ids=entity_ids)

    def _delete(self, entity_ids, query, *args, **kwargs):

        return self._do(
            cmd=self.manager.remove, entity_ids=entity_ids
        )
