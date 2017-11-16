# -*- coding: utf-8 -*-

import unittest

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.confng.simpleconf import Configuration
from canopsis.confng.vendor import Ini

from canopsis.activity.activity import (
    Activity, TimeUnits, DaysOfWeek, ActivityAggregate
)
from canopsis.activity.manager import ActivityManager, ActivityAggregateManager
from canopsis.activity.pbehavior import PBehaviorGenerator


class TestActivity(unittest.TestCase):

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

        ag = ActivityAggregate('agname')
        self.assertEqual(ag.name, 'agname')

        ac = Activity({}, DaysOfWeek.Monday, 0, 1)

        self.assertIsNone(ac.aggregate_name)

        ag.add(ac)

        self.assertEqual(ac.aggregate_name, ag.name)

    def test_aggregate_ac(self):
        ag = ActivityAggregate('agname')

        dow = DaysOfWeek.Monday
        ac1 = Activity({}, dow, 8 * TimeUnits.Hour, 17 * TimeUnits.Hour)

        ag.add(ac1)

        self.assertEqual(len(ag.activities), 1)
        self.assertEqual(ag.activities[0], ac1)

    def test_overlap(self):
        ag = ActivityAggregate('agname')

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

        acag = ActivityAggregate('testac')

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
            cls.store.get_collection(ActivityManager.ACTIVITY_COLLECTION))

    @classmethod
    def tearDownClass(self):
        self.col_activity.collection.drop()

    def setUp(self):
        self.col_activity.remove({})

    def test_store_get_then_equals(self):
        acm = ActivityManager(self.col_activity)
        agm = ActivityAggregateManager(acm)
        ag = ActivityAggregate('agtest')

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


if __name__ == '__main__':
    unittest.main()
