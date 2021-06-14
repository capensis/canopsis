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

from canopsis.perfdata.manager import PerfData
from canopsis.ctxprop.registry import CTXPropRegistry


class CTXPerfDataRegistry(CTXPropRegistry):
    """In charge of ctx perfdata properties.
    """

    __datatype__ = 'perfdata'  #: default datatype name

    def __init__(self, *args, **kwargs):

        super(CTXPerfDataRegistry, self).__init__(*args, **kwargs)

        self.manager = PerfData()

    def _do(self, cmd, ids, *args, **kwargs):

        result = []

        if ids is None:
            metrics = self.manager.context.find(_type='metric')
            ids = [metric['_id'] for metric in metrics]

        entity_id_field = self._entity_id_field()

        for entity_id in ids:
            cmdresult = cmd(metric_id=entity_id, **kwargs)
            if isinstance(cmdresult, list):
                result += [
                    {entity_id_field: entity_id, 'point': point}
                    for point in cmdresult
                ]
            else:
                item = {entity_id_field: entity_id, 'result': cmdresult}
                result.append(item)

        return result

    def _get(self, ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.get, ids=ids, with_tags=False)

    def _count(self, ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.count, ids=ids)

    def _delete(self, ids, query, *args, **kwargs):

        return self._do(cmd=self.manager.remove, ids=ids)

    def ids(self, query=None):

        result = self.manager.get_metrics(query=query)

        return result
