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

from canopsis.mongo.core import MongoStorage

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

    def update_mongo2influxdb(self):
        """Convert mongo data to influxdb data."""
        mongostorage = MongoStorage(table='periodic_perfdata')

        count = mongostorage._backend.count()

        if count:  # if data have to be deleted
            for document in mongostorage._backend.find():

                data_id = document['i']
                timestamp = int(document['t'])
                values = document['v']

                points = [(timestamp + int(ts), values[ts]) for ts in values]

                perfdata.put(metric_id=data_id, points=points)

                mongostorage.remove_elements(ids=document['_id'])

    def update(self):
        # FIXME
        if self.get_version('perfdata') < 1 and False:
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
