#!/usr/bin/env/python
# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from __future__ import unicode_literals
from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.context_graph.process import create_entity
from canopsis.middleware.core import Middleware
from canopsis.watcher.manager import Watcher
from collections import OrderedDict

watcher_one = "watcher-one"
watcher_example = {
    "_id": watcher_one,
    "alert_level": "minor",
    "crecord_creation_time": None,
    "crecord_name": None,
    "crecord_type": "selector",
    "crecord_write_time": None,
    "display_name": "a_displayed_name",
    "description": "a_description",
    "dosla": True,
    "dostate": True,
    "downtimes_as_ok": True,
    "enable": True,
    "exclude_ids": [],
    "include_ids": [],
    "last_dispatcher_update": None,
    "loaded": False,
    "mfilter": "{}",
    "output_tpl": "Off: [OFF], Minor: [MINOR], Major: [MAJOR], Critical: [CRITICAL], Ack count [ACK], Total: [TOTAL]",
    "sla_critical": 75,
    "sla_output_tpl": "Available: [P_AVAIL]%, Off: [OFF]%, Minor: [MINOR]%, Major: [MAJOR]%, Critical: [CRITICAL]%, Alerts [ALERTS]%, sla start [TSTART],  time available [T_AVAIL], time alert [T_ALERT]",
    "sla_timewindow": {
        "seconds": 12,
        "durationType": "second",
        "value": 12
    },
    "sla_warning": 90,
    "state": None,
    "state_algorithm": None,
    "state_when_all_ack": "worststate",
}


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


class CreateWatcher(BaseTest):

    def test_create_watcher(self):
        body = {
            '_id': 'an_id',
            'mfilter': '{}',
            'display_name': 'a_name'
        }
        self.manager.create_watcher(body)
        watcher = list(
            self.watcher_storage.find_elements(query={'_id': 'an_id'})
        )[0]
        self.assertTrue(sorted(list(body.values())) == sorted(list(watcher.values())))
        entity = self.context_graph_manager.get_entities_by_id('an_id')[0]
        print('------------------')
        print(entity)
        print('------------------')
        del entity['infos']['enable_history']
        self.assertTrue(
            sorted(list(
                {
                    '_id': 'watcher-one',
                    'impact': [],
                    'name': 'one',
                    'measurements': [],
                    'depends': [],
                    'infos': {
                        'enabled': True
                    },
                    'type': 'watcher'
                }
            )) == sorted(list(entity.values()))
        )
        entity['infos'].pop('enable_history', None)


class DeleteWatcher(BaseTest):

    def test_delete_watcher(self):
        self.manager.create_watcher(watcher_example)
        get_w = self.manager.get_watcher('watcher-one')
        self.assertTrue(isinstance(get_w, dict))
        self.assertEqual(get_w['_id'], watcher_one)

        del_w = self.manager.delete_watcher('watcher-one')
        self.assertTrue(isinstance(del_w, dict))
        self.assertEqual(del_w['n'], 1)

        get_w = self.manager.get_watcher('watcher-one')
        self.assertTrue(isinstance(get_w, dict))
        self.assertFalse(get_w['infos']['enabled'])

class GetWatcher(BaseTest):

    def test_get_watcher(self):
        self.manager.create_watcher(watcher_example)
        get_w = self.manager.get_watcher('watcher-one')
        self.assertTrue(isinstance(get_w, dict))
        self.assertEqual(get_w['_id'], watcher_one)
        self.assertEqual(None, self.manager.get_watcher('watcher_two'))

class WorstState(BaseTest):

    def test_worst_state(self):
        self.assertEqual(self.manager.worst_state(0,0,0), 0)
        self.assertEqual(self.manager.worst_state(0,0,1), 1)
        self.assertEqual(self.manager.worst_state(0,1,0), 2)
        self.assertEqual(self.manager.worst_state(1,0,0), 3)


if __name__ == "__main__":
    main()
