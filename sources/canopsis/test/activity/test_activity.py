# -*- coding: utf-8 -*-

import unittest
import copy

from canopsis.logger import Logger, OutputNull
from canopsis.common.collection import MongoCollection, CollectionError
from canopsis.common.mongo_store import MongoStore
from canopsis.confng.simpleconf import Configuration
from canopsis.confng.vendor import Ini

from canopsis.activity.activity import (
    Activity, TimeUnits, DaysOfWeek, ActivityAggregate
)
from canopsis.activity.manager import ActivityManager, ActivityAggregateManager
from canopsis.activity.pbehavior import PBehaviorGenerator

from canopsis.pbehavior.manager import PBehaviorManager


class TestActivity(unittest.TestCase):

    def test_copy(self):
        aco = Activity({}, DaysOfWeek.Monday, 0, 1)
        acc = copy.copy(aco)

        self.assertEqual(aco, acc)
        self.assertNotEqual(id(aco), id(acc))

    def test_invalid_day_of_week(self):

        with self.assertRaises(ValueError):
            Activity({}, "1", 0, 1)

    def test_invalid_start_stop(self):

        startt = 24 * TimeUnits.Hour + TimeUnits.Second
        stopt = 6 * TimeUnits.Day
        dow = DaysOfWeek.Monday

        with self.assertRaises(ValueError):
            Activity({}, dow, startt, stopt)

        with self.assertRaises(ValueError):
            Activity({}, dow, 7 * TimeUnits.Day, stopt)

    def test_valid_start_stop(self):
        startt = 7 * TimeUnits.Hour - TimeUnits.Second
        stopt = 7 * TimeUnits.Day
        ac = Activity({}, DaysOfWeek.Monday, startt, stopt)

        self.assertEqual(ac.start_time_of_day, startt)
        self.assertEqual(ac.stop_after_time, stopt)

        with self.assertRaises(ValueError):
            ac.stop_after_time = 1 * TimeUnits.Hour

        with self.assertRaises(ValueError):
            ac.stop_after_time = 80 * 24 * TimeUnits.Hour

    def test_eq(self):

        ac1 = Activity(
            {}, DaysOfWeek.Monday, TimeUnits.Day, 2 * TimeUnits.Day)

        ac2 = Activity(
            {}, DaysOfWeek.Monday, TimeUnits.Day, 2 * TimeUnits.Day)

        self.assertEqual(ac1, ac2)

    def test_overlap_simple(self):
        """
        Simple cases: no time overlaps
        """

        ac1 = Activity(
            {}, DaysOfWeek.Monday, TimeUnits.Hour, 2 * TimeUnits.Hour)

        ac2 = Activity(
            {}, DaysOfWeek.Monday, TimeUnits.Hour, 2 * TimeUnits.Hour)

        ac3 = Activity(
            {}, DaysOfWeek.Thursday, TimeUnits.Hour, 2 * TimeUnits.Hour)

        ac4 = Activity(
            {}, DaysOfWeek.Monday, 3 * TimeUnits.Hour, 4 * TimeUnits.Hour)

        self.assertTrue(ac1.overlap(ac2))
        self.assertFalse(ac1.overlap(ac3))
        self.assertFalse(ac1.overlap(ac4))

    def test_overlap_advanced(self):
        """
        More complex cases: time overlaps on different days, etc...
        """

        # monday/00:00 -> tuesday/01:00
        ac1 = Activity(
            {}, DaysOfWeek.Monday,
            1 * TimeUnits.Hour, TimeUnits.Day + TimeUnits.Hour)

        # tuesday/00:00 -> tuesday/02:00
        ac2 = Activity(
            {}, DaysOfWeek.Tuesday,
            0 * TimeUnits.Hour, 2 * TimeUnits.Hour)

        # monday/00:00 -> wednesday/08:00
        ac3 = Activity(
            {}, DaysOfWeek.Monday,
            3 * TimeUnits.Hour, 2 * TimeUnits.Day + 8 * TimeUnits.Hour)

        self.assertTrue(ac1.overlap(ac2))
        self.assertTrue(ac3.overlap(ac1))
        self.assertTrue(ac1.overlap(ac3))

    def test_overlap_next_week(self):
        ac1 = Activity(
            {}, DaysOfWeek.Monday,
            2 * TimeUnits.Hour, 4 * TimeUnits.Hour)

        ac2 = Activity(
            {}, DaysOfWeek.Monday,
            5 * TimeUnits.Hour, 7 * TimeUnits.Day + 3 * TimeUnits.Hour)

        ac3 = Activity(
            {}, DaysOfWeek.Monday,
            5 * TimeUnits.Hour, 7 * TimeUnits.Day + 1 * TimeUnits.Hour)

        self.assertTrue(ac1.overlap(ac2))
        self.assertFalse(ac1.overlap(ac3))


class TestActivityAggregate(unittest.TestCase):

    def test_aggregate(self):

        ag = ActivityAggregate('agname', {'filter': 'from_aggregate'})
        self.assertEqual(ag.name, 'agname')

        ac = Activity({}, DaysOfWeek.Monday, 0, 1)

        self.assertIsNone(ac.aggregate_name)

        ag.add(ac)

        self.assertEqual(
            ag.activities[0].aggregate_name,
            ag.name
        )
        self.assertEqual(
            ag.activities[0].entity_filter,
            {'filter': 'from_aggregate'}
        )

    def test_aggregate_ac(self):
        ag = ActivityAggregate('agname', {})

        dow = DaysOfWeek.Monday
        ac1 = Activity({}, dow, 8 * TimeUnits.Hour, 17 * TimeUnits.Hour)

        ag.add(ac1)

        self.assertEqual(len(ag.activities), 1)
        self.assertEqual(ag.activities[0], ac1)

    def test_overlap(self):
        ag = ActivityAggregate('agname', {})

        dow = DaysOfWeek.Monday
        ac1 = Activity({}, dow, 0 * TimeUnits.Hour, 2 * TimeUnits.Hour)
        ac2 = Activity({}, dow, 0 * TimeUnits.Hour, TimeUnits.Hour)

        ag.add(ac1)

        with self.assertRaises(ValueError):
            ag.add(ac2)


class TestActivityAggregateToPBehavior(unittest.TestCase):

    def test_pb(self):
        ac_pb_gen = PBehaviorGenerator()
        filter_ = {'_id': 'schtroumpf'}

        acag = ActivityAggregate('testac', {})

        ts1 = 8 * TimeUnits.Hour
        tS1 = 17 * TimeUnits.Hour

        ts2 = 9 * TimeUnits.Hour
        tS2 = 16 * TimeUnits.Hour

        ac1 = Activity(filter_, DaysOfWeek.Monday, ts1, tS1)
        ac4 = Activity(filter_, DaysOfWeek.Tuesday, ts2, tS2)
        ac2 = Activity(filter_, DaysOfWeek.Wednesday, ts1, tS1)
        ac5 = Activity(filter_, DaysOfWeek.Thursday, ts2, tS2)
        ac3 = Activity(filter_, DaysOfWeek.Friday, ts1, tS1)

        acag.add(ac1)
        acag.add(ac2)
        acag.add(ac3)
        acag.add(ac4)
        acag.add(ac5)

        start_date = ac_pb_gen._get_monday()
        pbehaviors = ac_pb_gen.activities_to_pbehaviors(acag, start_date)

        pb1, pb2, pb3 = pbehaviors

        self.assertEqual(pb1.rrule, 'FREQ=DAILY;BYDAY=FR')
        self.assertEqual(pb2.rrule, 'FREQ=DAILY;BYDAY=TU,TH')
        self.assertEqual(pb3.rrule, 'FREQ=DAILY;BYDAY=MO,WE')


class TestActivityAggregateManager(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        conf = Configuration.load(MongoStore.CONF_PATH, Ini)
        cred_conf = Configuration.load(MongoStore.CRED_CONF_PATH, Ini)

        cls.store = MongoStore(config=conf, cred_config=cred_conf)

        cls.col_activity = MongoCollection(
            cls.store.get_collection('test_activities'))
        cls.col_ag = MongoCollection(
            cls.store.get_collection('test_activity_aggregates'))
        cls.col_pb = MongoCollection(
            cls.store.get_collection('default_testactivitiespb'))

    @classmethod
    def tearDownClass(self):
        self.col_ag.collection.drop()
        self.col_activity.collection.drop()
        self.col_pb.collection.drop()

    def setUp(self):
        self.col_ag.remove({})
        self.col_activity.remove({})
        self.col_pb.remove({})

    def test_pb_attach(self):
        acm = ActivityManager(self.col_activity)
        agm = ActivityAggregateManager(self.col_ag, self.col_pb, acm)
        ag = ActivityAggregate('agtest', {})

        ac1 = Activity({}, DaysOfWeek.Monday, 0, 1)
        pb_ids = ['je', 'ne', 'suis', 'pas', 'un', 'robot', 'robot']

        ag.add(ac1)

        with self.assertRaises(CollectionError):
            agm.attach_pbehaviors(ag, pb_ids)

        agm.store(ag)
        agm.attach_pbehaviors(ag, pb_ids)

        agg = agm.get(ag.name)

        self.assertSetEqual(
            set(agg.pb_ids),
            set(pb_ids)
        )

    def test_store_get_then_equals(self):
        acm = ActivityManager(self.col_activity)
        agm = ActivityAggregateManager(self.col_ag, self.col_pb, acm)
        ag = ActivityAggregate('agtest', {})

        ac1 = Activity({}, DaysOfWeek.Monday, 0, 1)
        ac2 = Activity({}, DaysOfWeek.Saturday, 0, 1)

        self.assertTrue(ag.add(ac1))
        self.assertTrue(ag.add(ac2))
        self.assertFalse(ag.add(ac1))

        agm.store(ag)

        res = acm.get_by_aggregate_name(ag.name)
        self.assertEqual(len(res), 2)
        self.assertIsInstance(res[0], Activity)
        self.assertIsInstance(res[1], Activity)

        self.assertEqual(ac1, res[0])
        self.assertEqual(ac2, res[1])

        # Testing ActivityManager
        res = acm.get_all()
        self.assertEqual(len(res), 2)
        res = acm.del_by_aggregate_name('agtest')
        self.assertEqual(res['n'], 2)

    def test_store_get_delete_with_pb(self):
        from canopsis.middleware.core import Middleware
        ac_pb_gen = PBehaviorGenerator()
        acm = ActivityManager(self.col_activity)
        agm = ActivityAggregateManager(self.col_ag, self.col_pb, acm)
        ag = ActivityAggregate('agtest', {})
        pb_man = PBehaviorManager(
            Logger.get('test', None, output_cls=OutputNull),
            Middleware.get_middleware_by_uri(
                'mongodb-default-testactivitiespb://'))

        ac1 = Activity({}, DaysOfWeek.Monday, 0, 1)
        ac2 = Activity({}, DaysOfWeek.Thursday, 0, 1)

        self.assertTrue(ag.add(ac1))
        self.assertTrue(ag.add(ac2))

        agm.store(ag)

        agc = agm.get(ag.name)

        start_date = ac_pb_gen._get_monday()
        pbehaviors = ac_pb_gen.activities_to_pbehaviors(agc, start_date)

        self.assertTrue(agc.name, ag.name)
        self.assertEqual(len(pbehaviors), 2)

        pb_ids = []
        for pb in pbehaviors:
            pb_id = pb_man.create(
                name=pb.name,
                filter=pb.filter_,
                author=pb.author,
                tstart=pb.tstart,
                tstop=pb.tstop,
                rrule=pb.rrule,
                enabled=pb.enabled,
                comments=pb.comments,
            )
            pb_ids.append(pb_id)

        agm.attach_pbehaviors(agc, pb_ids)

        agc2 = agm.get(ag.name)
        self.assertEqual(agc2.pb_ids, agc.pb_ids)

        res = pb_man.pb_storage._backend.find({'_id': {'$in': pb_ids}})
        self.assertEqual(len(list(res)), 2)

        agm.delete(agc2)
        res = pb_man.pb_storage._backend.find({'_id': {'$in': pb_ids}})
        self.assertEqual(len(list(res)), 0)

    def test_store_get_dbid(self):
        acm = ActivityManager(self.col_activity)
        agm = ActivityAggregateManager(self.col_ag, self.col_pb, acm)

        ag = ActivityAggregate('agtest', {})

        ac1 = Activity({}, DaysOfWeek.Monday, 0, 1)
        ac2 = Activity({}, DaysOfWeek.Saturday, 0, 1)

        ag.add(ac1)
        ag.add(ac2)

        store_ids = agm.store(ag)
        activities = acm.get_by_aggregate_name(ag.name)

        self.assertEqual(len(store_ids), len(activities))
        self.assertEqual(str(store_ids[0]), str(activities[0].dbid))
        self.assertEqual(str(store_ids[1]), str(activities[1].dbid))


if __name__ == '__main__':
    unittest.main()
