import unittest
from unittest import TestCase

import time

from datetime import datetime
from canopsis.webcore.services.weather import get_pb_range

def dt_to_ts(dt):
    """
    datetime to timestamp

    :param datetime dt: datetime
    :rtype: int
    """
    return int(time.mktime(dt.timetuple()))

class TestWeather(TestCase):

    def test_pb_range(self):
        pb1 = {
            'rrule': 'FREQ=DAILY;BYDAY=MO,TU',
            'tstart': dt_to_ts(datetime(2017, 01, 01, 18, 15, 0)),
            'tstop': dt_to_ts(datetime(2017, 01, 02, 8, 0, 0))
        }

        pb2 = {
            'rrule': 'FREQ=DAILY;BYDAY=MO,TU,WE',
            'tstart': dt_to_ts(datetime(2017, 01, 01, 17, 15, 0)),
            'tstop': dt_to_ts(datetime(2017, 01, 02, 9, 0, 0))
        }

        rset, tod_start, tod_stop = get_pb_range([pb1, pb2])

        self.assertEquals(tod_start.hour, 17)
        self.assertEquals(tod_start.minute, 15)
        self.assertEquals(tod_stop.hour, 9)
        self.assertEquals(tod_stop.minute, 0)

if __name__ == '__main__':
    unittest.main()