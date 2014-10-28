# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.engines import Engine
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from datetime import datetime
from time import time

class engine(Engine):
    etype = "datacleaner"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.clean_collection = {}
        collections_to_clean = ['events', 'events_log']

        for collection in collections_to_clean:
            self.clean_collection[collection] = get_storage(
                collection,
                account=Account(user='root')
            ).get_backend()

        #task ran daily
        self.beat_interval = 3600 * 24


    def pre_run(self):

        self.beat()


    def get_retention_date(self):

        """
            Computes date beyond where date must be cleaned.
            Gets the retention duration from user GUI input.
        """

        retention_duration = 1 # TODO get it from record

        compare_duration = int(time() - retention_duration)

        security_duration_limit = int(time() - 3600 * 24 * 365) # one year

        self.logger.debug(
            'retention ts {}, security ts {}'.format(
                datetime.utcfromtimestamp(compare_duration),
                datetime.utcfromtimestamp(security_duration_limit)
            )
        )

        if compare_duration > security_duration_limit:
            self.logger.info(
                'Retention date too short, ' +
                'prevent data deletion by setting retention to one year'
            )
            compare_duration = security_duration_limit

        self.logger.debug('selected retention ts {}'.format(
            datetime.utcfromtimestamp(compare_duration)
        ))

        return compare_duration


    def beat(self):

        # getting retention date limit
        retention_date_limit = self.get_retention_date()

        # formating query for deletion
        query = {
            'timestamp': {
                '$lte': retention_date_limit
            }
        }

        # iteration over collections to clean
        for collection in self.clean_collection:

            clean_collection = self.clean_collection[collection]

            count = clean_collection.find(query).count()

            total = clean_collection.find().count()

            self.logger.info(
                'Clean {}/{} documents in collection {} starts.'.format(
                    count,
                    total,
                    collection
                )
            )

            # effective collection clean
            count = clean_collection.remove(query)


            self.logger.info('Clean complete for collection {}'.format(
                collection
            ))

