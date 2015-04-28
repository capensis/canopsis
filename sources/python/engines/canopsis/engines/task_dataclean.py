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

from canopsis.engines.core import TaskHandler
from canopsis.old.account import Account
from canopsis.old.storage import get_storage

from datetime import datetime
from time import time


class engine(TaskHandler):
    etype = "taskdataclean"

    def __init__(self, *args, **kwargs):
        super(engine, self).__init__(*args, **kwargs)

        self.clean_collection = {}

    def get_collection(self, collection):
        if collection not in self.clean_collection:
            self.clean_collection[collection] = get_storage(
                collection,
                account=Account(user='root')
            )

        return self.clean_collection[collection].get_backend()

    def get_retention_date(self, configuration):
        """
        Computes date beyond where date must be cleaned.
        Gets the retention duration from user GUI input.
        """

        # Alias
        datestr = datetime.utcfromtimestamp
        retention_duration = configuration['retention_duration']

        # Compute retention duration compare date (compare with security date)
        compare_duration = int(time() - retention_duration)

        # Just set security retention duration depending
        # on user sercure information
        if configuration['use_secure_delay']:
            security_duration_limit = int(time() - 3600 * 24 * 365)  # one year

            self.logger.info(
                'Secure delay is set to true, minimum delay set to 1 year'
            )

        else:
            security_duration_limit = int(time())  # now

            self.logger.info('Secure delay is set to false')

        self.logger.debug(
            'retention ts: {}, security ts: {}'.format(
                datestr(compare_duration),
                datestr(security_duration_limit)
            )
        )

        # When not secure, compare duration is at least gte security delay
        if compare_duration > security_duration_limit:
            self.logger.info(
                'Retention date too short, will use minimum delay'
            )

            compare_duration = security_duration_limit

        self.logger.debug('selected retention ts: {}'.format(
            datestr(compare_duration)
        ))

        return compare_duration

    def handle_task(self, job):
        self.logger.debug('taskdataclean.handle_task()')
        self.logger.debug('job: {0}'.format(job))

        # getting retention date limit
        retention_date_limit = self.get_retention_date(job)

        # formating query for deletion
        query = {
            'timestamp': {
                '$lte': retention_date_limit
            }
        }

        # iteration over collections to clean
        for collection in job['storages']:
            clean_collection = self.get_collection(collection)

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

        return (0, 'Collections were cleaned')
