from unittest import TestCase

from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.watcher.manager import Watcher as WatcherManager
from canopsis.webcore.services.pbehavior import check_values, RouteHandlerPBehavior

class TestPbehavior(TestCase):

    INVALID_PB_STRINGS = {
        'name': 1,
        'author': 2,
        'rrule': 3,
        'component': 4,
        'connector': 5,
        'connector_name': 6
    }

    INVALID_PB_INT = {
        'tstart': '7',
        'tstop': '8'
    }

    INVALID_PB_TSTART_TSTOP = {
        'tstart': 1000,
        'tstop': 100
    }

    INVALID_PB_COMMENTS = {
        'comments': '9'
    }

    INVALID_PB_FILTER = {
        'filter': 'invalid_{}filter[]'
    }

    INVALID_PB_ENABLED = {
        'enabled': '10'
    }

    INVALID_PB_RRULE = {
        'rrule': 'FREQ=DAILY;INVALID'
    }

    VALID_PB = {
        'name': 'test_pb',
        'author': 'test_case_pb',
        'rrule': 'FREQ=DAILY;BYDAY=MO',
        'component': 'test_case_pb',
        'connector': 'unittest',
        'connector_name': 'unittest',
        'tstart': 1000,
        'tstop': 10000,
    }

    def test_check_values_invalid_pb(self):
        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_STRINGS)

        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_INT)

        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_TSTART_TSTOP)

        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_COMMENTS)

        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_FILTER)

        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_ENABLED)

        with self.assertRaises(ValueError):
            check_values(self.INVALID_PB_RRULE)

    def test_check_values_valid_pb(self):
        check_values(self.VALID_PB)

class TestPbehaviorWebservice(TestCase):

    INVALID_PB = {
        'name': 'test_bad_pb',
        'author': 'pb_author',
        'filter_': 'bad filter',
        'tstart': 1000,
        'tstop': 100,
        'rrule': 'blurp',
        'enabled': True,
        'comments': [],
        'connector': 'test_pb',
        'connector_name': 'test_pb'
    }

    VALID_PB = {
        'name': 'test_pb',
        'author': 'pb_author',
        'filter_': '{"nokey": "novalue"}',
        'tstart': 100,
        'tstop': 1000,
        'rrule': 'FREQ=DAILY;BYDAY=MO',
        'enabled': True,
        'comments': [],
        'connector': 'test_pb',
        'connector_name': 'test_pb'
    }

    @classmethod
    def setUpClass(cls):
        pb_manager = PBehaviorManager()
        watcher_manager = WatcherManager()
        cls.rhpb = RouteHandlerPBehavior(pb_manager, watcher_manager)

    def test_create_bad_pb(self):
        with self.assertRaises(ValueError):
            self.rhpb.create(**self.INVALID_PB)

    def test_create_pb(self):
        self.rhpb.create(**self.VALID_PB)