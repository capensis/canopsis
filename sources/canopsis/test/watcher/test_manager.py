#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest
from calendar import timegm
from datetime import datetime, timedelta
from json import loads
from time import sleep
from unittest import main, TestCase
from mock import Mock

from canopsis.alerts.manager import Alerts
from canopsis.common import root_path
from canopsis.common.ethereal_data import EtherealData
from canopsis.confng import Configuration, Ini
from canopsis.context_graph.manager import ContextGraph
# from canopsis.context_graph.process import create_entity
from canopsis.logger.logger import Logger, OutputNull
from canopsis.common.mongo_store import MongoStore
from canopsis.common.collection import MongoCollection
from canopsis.common.middleware import Middleware
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.statsng.event_publisher import StatEventPublisher
from canopsis.watcher.manager import Watcher
import xmlrunner

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
        logger = Logger.get('', None, output_cls=OutputNull)
        self.manager = Watcher()
        self.context_graph_manager = ContextGraph(logger)
        self.alerts_storage = Middleware.get_middleware_by_uri(
            'mongodb-periodical-testalarm://'
        )
        self.watcher_storage = Middleware.get_middleware_by_uri(
            'storage-default-testwatcher://'
        )
        self.entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )

        self.context_graph_manager.ent_storage = self.entities_storage
        self.manager.alert_storage = self.alerts_storage
        self.manager.context_graph = self.context_graph_manager
        self.manager.watcher_storage = self.watcher_storage

    def tearDown(self):
        self.watcher_storage.remove_elements()
        self.alerts_storage.remove_elements()
        self.entities_storage.remove_elements()


class GetWatcher(BaseTest):

    def test_get_watcher(self):
        self.assertIsNone(self.manager.get_watcher('watcher-one'))
        watcher_entity = ContextGraph.create_entity_dict(
            'watcher-one',
            'one',
            'watcher'
        )
        self.context_graph_manager.create_entity(watcher_entity)
        #print(self.manager.get_watcher('watcher-one'))


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
        entity['infos'].pop('enable_history', None)

        expected = {
            '_id': 'an_id',
            'impact': [],
            'name': 'a_name',
            'measurements': {},
            'depends': [],
            'mfilter': '{}',
            'state': 0,
            'enabled': True,
            'infos': {},
            'type': 'watcher'
        }

        self.assertEquals(expected['_id'], entity['_id'])
        self.assertEquals(expected['impact'], entity['impact'])
        self.assertEquals(expected['name'], entity['name'])
        self.assertEquals(expected['measurements'], entity['measurements'])
        self.assertEquals(expected['depends'], entity['depends'])
        self.assertEquals(expected['mfilter'], entity['mfilter'])
        self.assertEquals(expected['state'], entity['state'])
        self.assertEquals(expected['type'], entity['type'])
        self.assertEquals(expected['enabled'], entity['enabled'])


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
        self.assertEqual(get_w, None)


class GetWatcher2(BaseTest):

    def test_get_watcher(self):
        self.manager.create_watcher(watcher_example)
        get_w = self.manager.get_watcher('watcher-one')
        self.assertTrue(isinstance(get_w, dict))
        self.assertEqual(get_w['_id'], watcher_one)
        self.assertEqual(None, self.manager.get_watcher('watcher_two'))


class WorstState(BaseTest):

    def test_worst_state(self):
        self.assertEqual(self.manager.worst_state(0, 0, 0), 0)
        self.assertEqual(self.manager.worst_state(0, 0, 1), 1)
        self.assertEqual(self.manager.worst_state(0, 1, 0), 2)
        self.assertEqual(self.manager.worst_state(1, 0, 0), 3)


class ComputeState(BaseTest):

    def setUp(self):
        super(ComputeState, self).setUp()
        mongo = MongoStore.get_default()
        collection = mongo.get_collection("default_testpbehavior")
        pb_collection = MongoCollection(collection)

        filter_storage = Middleware.get_middleware_by_uri(
            'storage-default-testalarmfilter://'
        )
        config_storage = Middleware.get_middleware_by_uri(
            'storage-default-testconfig://'
        )
        config_storage.put_element(
            element={
                '_id': 'test_config',
                'crecord_type': 'statusmanagement',
                'bagot_time': 3600,
                'bagot_freq': 10,
                'stealthy_time': 300,
                'restore_event': True,
                'auto_snooze': False,
                'snooze_default_time': 300,
            },
            _id='test_config'
        )
        logger = Logger.get('test_pb', None, output_cls=OutputNull)

        config = Configuration.load(PBehaviorManager.CONF_PATH, Ini)

        self.pbm = PBehaviorManager(config=config,
                                    logger=logger,
                                    pb_collection=pb_collection)
        self.pbm.context = self.context_graph_manager
        self.manager.pbehavior_manager = self.pbm

        conf = Configuration.load(Alerts.CONF_PATH, Ini)
        filter_ = {'crecord_type': 'statusmanagement'}
        config_data = EtherealData(collection=config_storage._backend,
                                   filter_=filter_)

        event_publisher = Mock(spec=StatEventPublisher)

        self.alert_manager = Alerts(config=conf,
                                    logger=logger,
                                    alerts_storage=self.alerts_storage,
                                    config_data=config_data,
                                    filter_storage=filter_storage,
                                    context_graph=self.context_graph_manager,
                                    watcher=self.manager,
                                    event_publisher=event_publisher,
                                    pbehavior=self.pbm)

        # Creating entity
        self.type_ = 'resource'
        self.name = 'morticia'
        entity = ContextGraph.create_entity_dict(
            id=self.name,
            etype=self.type_,
            name=self.name
        )
        self.context_graph_manager.create_entity(entity)

        # Creating coresponding alarm
        event = {
            'connector': self.type_,
            'connector_name': 'connector_name',
            'component': self.name,
            'output': 'tadaTaDA tic tic',
            'timestamp': 0
        }
        alarm = self.alert_manager.make_alarm(self.name, event)
        self.state = 2
        alarm = self.alert_manager.update_state(alarm, self.state, event)
        new_value = alarm[self.alert_manager.alerts_storage.VALUE]
        self.alert_manager.update_current_alarm(alarm, new_value)

    def tearDown(self):
        super(ComputeState, self).tearDown()
        self.pbm.collection.remove({})


    def test_compute_state_issue427(self):
        # Aka: state desyncro
        watcher_id = 'addams'
        watcher = {
            '_id': watcher_id,
            'mfilter': '{"name": {"$in": ["morticia"]}}',
            'display_name': 'family'
        }
        self.assertTrue(self.manager.create_watcher(watcher))

        res = self.manager.get_watcher(watcher_id)
        self.assertEqual(res['state'], self.state)

        # Creating pbehavior on it
        now = datetime.utcnow()
        self.pbm.create(
            name='addam',
            filter=loads('{"name": "morticia"}'),
            author='addams',
            tstart=timegm(now.timetuple()),
            tstop=timegm((now + timedelta(seconds=2)).timetuple()),
            rrule=None,
            enabled=True
        )

        self.pbm.compute_pbehaviors_filters()

        res = self.manager.get_watcher(watcher_id)
        self.assertEqual(res['state'], self.state)

        self.manager.compute_watchers()

        res = self.manager.get_watcher(watcher_id)
        self.assertEqual(res['state'], 0)

        sleep(3)
        self.pbm.compute_pbehaviors_filters()
        self.manager.compute_watchers()

        res = self.manager.get_watcher(watcher_id)
        self.assertEqual(res['state'], self.state)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
