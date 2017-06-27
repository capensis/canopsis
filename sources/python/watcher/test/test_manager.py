#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware
from canopsis.watcher.manager import Watcher
from canopsis.context_graph.process import create_entity


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
        self.context_graph_manager[ContextGraph.ENTITIES_STORAGE] = (
            self.entities_storage
        )
        self.manager.context_graph = self.context_graph_manager
        self.manager[Watcher.WATCHER_STORAGE] = self.watcher_storage
        self.manager[Watcher.ALERTS_STORAGE] = self.alerts_storage

    def tearDown(self):
        self.watcher_storage.remove_elements()
        self.alerts_storage.remove_elements()
        self.entities_storage.remove_elements()

class GetWatcher(BaseTest):

    def test_get_watcher(self):
        self.assertIsNone(self.manager.get_watcher('watcher-one'))
        watcher_entity = create_entity(
            'watcher-one',
            'one',
            'watcher'
        )
        self.context_graph_manager.create_entity(watcher_entity)
        print(self.manager.get_watcher('watcher-one'))

    def create_watcher(self):
        body = {
            '_id': 'an_id',
            'mfilter': '{}',
            'display_name': 'a_name'
        }
        self.assertTrue(
            sortedDict(body) == (
                sortedDict(
                    self.watcher_storage.find_elements(query={'_id': 'an_id'})
                )
            )
        )
        entity = self.context_graph_manager.get_entities_by_id('an_id')
        entity['infos'].pop('enable_history', None)
        """
        self.assertTrue(
            sortedDict(
                {
                    u'impact': [],
                    u'name': u'one',
                    u'measurements': [],
                    u'depends': [],
                    u'infos': {
                        u'enabled': True
                    }
                    u'_id': u'watcher-one',
                    u'type': u'watcher'
                }
            )
        )
        """
if __name__ == "__main__":
    main()
