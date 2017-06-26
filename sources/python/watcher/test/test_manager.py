#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware
from canopsis.watcher.manager import Watcher


class BaseTest(TestCase):

    def setUp(self):
        self.manager = Watcher()
        self.context_graph_manager = ContextGraph()
        self.watcher_storage = Middleware.get_middleware_by_uri(
            'storage-default-testwatcher://'
        )
        self.alerts_storage = Middleware.get_middleware_by_uri(
            'storage-default-testalerts://'
        )
        self.entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        self.context_graph_manager[ContextGraph.ENTITIES_STORAGE] = self.entities_storage

        self.manager[Watcher.WATCHER_STORAGE] = self.watcher_storage
        self.manager[Watcher.ALERTS_STORAGE] = self.alerts_storage

    def tearDown(self):
        self.watcher_storage.remove_elements()
        self.alerts_storage.remove_elements()

class GetWatcher(BaseTest):

    def test_get_watcher():
        self.assertIsNone(self.manager.get_watcher('watcher-one'))
