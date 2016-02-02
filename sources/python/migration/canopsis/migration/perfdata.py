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

from canopsis.migration.manager import MigrationModule
from canopsis.configuration.configurable.decorator import conf_paths
from canopsis.configuration.configurable.decorator import add_category
from canopsis.perfdata.manager import PerfData


CONF_PATH = 'migration/perfdata.conf'
CATEGORY = 'PERFDATA'
CONTENT = []


@conf_paths(CONF_PATH)
@add_category(CATEGORY, content=CONTENT)
class PerfdataModule(MigrationModule):
    def __init__(self, *args, **kwargs):
        super(PerfdataModule, self).__init__(*args, **kwargs)

        self.manager = PerfData()

    def init(self):
        pass

    def update(self):
        if self.get_version('perfdata') < 1:
            self.logger.info('Migrating to version 1')

            self.update_to_version_1()
            self.set_version('perfdata', 1)

    def update_to_version_1(self):
        storage = self.manager[PerfData.PERFDATA_STORAGE]
        nan = float('nan')

        oneweek = 3600 * 24 * 7

        for document in storage.find_elements():

            metric_id = document['i']

            values = document['v']
            t = document['t']

            points = list(
                (t + int(ts), nan if values[ts] is None else values[ts])
                for ts in values
            )

            rightvalues = {
                key: values[key] for key in values if int(key) < oneweek
            }
            document['v'] = rightvalues

            storage.put_element(
                element=document, cache=False
            )

            self.manager.put(metric_id=metric_id, points=points, cache=False)
