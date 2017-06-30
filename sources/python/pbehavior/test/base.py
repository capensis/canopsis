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

from unittest import TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware
from canopsis.pbehavior.manager import PBehaviorManager


class BaseTest(TestCase):
    def setUp(self):
        pbehavior_storage = Middleware.get_middleware_by_uri(
            'storage-default-testpbehavior://'
        )
        entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )

        self.pbm = PBehaviorManager()
        self.context = ContextGraph()
        self.context[ContextGraph.ENTITIES_STORAGE] = entities_storage
        self.pbm.context = self.context

        self.pbm[PBehaviorManager.PBEHAVIOR_STORAGE] = pbehavior_storage

    def tearDown(self):
        self.pbm.pbehavior_storage.remove_elements()
        self.context[ContextGraph.ENTITIES_STORAGE].remove_elements()
