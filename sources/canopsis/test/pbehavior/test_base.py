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
from canopsis.logger.logger import Logger, OutputNull
from canopsis.middleware.core import Middleware
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.watcher.manager import Watcher


class MockEngine():

    def __init__(self):
        self.logger = Logger.get('', None, output_cls=OutputNull)


class BaseTest(unittest.TestCase):
    def setUp(self):
        pbehavior_storage = Middleware.get_middleware_by_uri(
            'storage-default-testpbehavior://'
        )
        entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        alarm_storage = Middleware.get_middleware_by_uri(
            'storage-periodical-testalarm://'
        )
        watcher_manager = Watcher()

        logger = Logger.get('test_pb', None, output_cls=OutputNull)

        self.pbm = PBehaviorManager(
            logger=logger,
            pb_storage=pbehavior_storage,
            alarm_storage=alarm_storage,
            watcher_manager=watcher_manager,
        )
        self.context = ContextGraph(logger)
        self.context.ent_storage = entities_storage
        self.pbm.context = self.context

    def tearDown(self):
        self.pbm.pb_storage.remove_elements()
        self.context.ent_storage.remove_elements()
