from unittest import TestCase

from canopsis.webcore.services.pbehavior import check_values

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