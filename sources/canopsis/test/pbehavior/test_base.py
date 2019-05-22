#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
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

import unittest
from canopsis.context_graph.manager import ContextGraph
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.common.middleware import Middleware
from canopsis.confng import Configuration, Ini
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.logger.logger import Logger, OutputNull


class MockEngine():

    def __init__(self):
        self.logger = Logger.get('', None, output_cls=OutputNull)


class BaseTest(unittest.TestCase):
    def setUp(self):
        mongo = MongoStore.get_default()
        collection = mongo.get_collection("default_testpbehavior")
        pb_coll = MongoCollection(collection)
        entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )

        logger = Logger.get('test_pb', None, output_cls=OutputNull)
        conf = Configuration.load(PBehaviorManager.CONF_PATH, Ini)
        self.pbm = PBehaviorManager(config=conf,
                                    logger=logger,
                                    pb_collection=pb_coll)
        self.context = ContextGraph(logger)
        self.context.ent_storage = entities_storage
        self.pbm.context = self.context

    def tearDown(self):
        self.pbm.delete(_filter={})
        self.context.ent_storage.remove_elements()
