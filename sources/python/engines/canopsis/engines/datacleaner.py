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

        self.object = get_storage(
            'object',
            account=Account(user='root')
        ).get_backend()

    def get_configuration(self):
        """ Retrieve engine configuration from database """
        self.logger.info('Reloading configuration from database')
        # Get engine configuration
        return self.object.find_one(
            {'crecord_type': 'datacleaner'}
        )

    def get_retention_date(self):

        """
            Computes date beyond where date must be cleaned.
            Gets the retention duration from user GUI input.
        """
        # Alias
        datestr = datetime.utcfromtimestamp

        # Easier to test this way
        configuration = self.get_configuration()

        if configuration is None:
            self.logger.warning(
                'No configuration found, cannot process data cleaning'
            )
            return None

        self.logger.info('configuration reloaded')
        self.logger.debug(configuration)

        retention_duration = configuration['retention_duration']

        # Compute retention duration compare date (compare with security date)
        compare_duration = int(time() - retention_duration)

        # Just set security retention duration depending
        # on user sercure information
        if configuration['use_secure_delay']:

            security_duration_limit = int(time() - 3600 * 24 * 365)  # one year

            self.logger.info(
                'Secure delay is set to true.' +
                ' datacleaner will use one year delay security'
            )

        else:

            security_duration_limit = int(time())  # now

            self.logger.info(
                'Secure delay is set to false.' +
                ' datacleaner will use user delay'
            )

        self.logger.debug(
            'retention ts {}, security ts {}'.format(
                datestr(compare_duration),
                datestr(security_duration_limit)
            )
        )

        # When not secure, compare duration is at least gte security delay
        if compare_duration > security_duration_limit:
            self.logger.info(
                'Retention date too short, ' +
                'prevent data deletion by setting retention to one year'
            )
            compare_duration = security_duration_limit

        self.logger.debug('selected retention ts {}'.format(
            datestr(compare_duration)
        ))

        return compare_duration

    def consume_dispatcher(self, event, *args, **kargs):
        self.logger.warning('enter datacleaner beat')
        # getting retention date limit
        retention_date_limit = self.get_retention_date()

        if retention_date_limit is None:
            return

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
