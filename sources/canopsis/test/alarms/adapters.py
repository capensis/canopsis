#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from __future__ import unicode_literals

import logging
from unittest import TestCase, main

from canopsis.alarms.adapters import (
    Adapter, make_alarm_from_mongo, make_alarm_step_from_mongo
)
from canopsis.common.mongo_store import MongoStore


class AlarmsAdaptersTest(TestCase):

    @classmethod
    def setUp(self):
        self.logger = logging.getLogger('alarms')

        self.conf = {
            MongoStore.CONF_CAT: {
                'db': self.db_name
            }
        }
        self.cred_conf = {
            MongoStore.CRED_CAT: {
                'user': '',
                'pwd': ''
            }
        }
        self.ms = MongoStore(config=self.conf,
                             cred_config=self.cred_conf)
        self.collection_name = 'test_{}'.format(Adapter.COLLECTION)
        client = self.ms.get_collection(self.collection_name)
        self.adapter = Adapter(mongo_client=client)

    def test_adapter(self):
        pass
        #res = self.adapter.find_unresolved_snoozed_alarms()

    def test_make_alarm_from_mongo(self):
        pass

    def test_make_alarm_step_from_mongo(self):
        pass

if __name__ == '__main__':
    main()
