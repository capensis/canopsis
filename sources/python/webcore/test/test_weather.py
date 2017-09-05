import unittest
from unittest import TestCase

import time

from datetime import datetime
from canopsis.webcore.services.weather import get_pb_range, dt_to_ts
from canopsis.webcore.services.weather import RouteHandlerWeather

from canopsis.middleware.core import Middleware

from canopsis.alerts.reader import AlertsReader
from canopsis.context_graph.manager import ContextGraph
from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.tracer.manager import TracerManager
from canopsis.watcher.manager import Watcher as WatcherManager

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

class TestServiceWeather(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.rhw = RouteHandlerWeather()

        watcher_storage = Middleware.get_middleware_by_uri('mongodb-default-testwatcher://')

        cls.context_manager = ContextGraph()
        cls.alarmreader_manager = AlertsReader()
        cls.pbehavior_manager = PBehaviorManager()
        cls.tracer_manager = TracerManager()
        cls.watcher_manager = WatcherManager(
            watcher_storage=watcher_storage,
            context_graph_manager=cls.context_manager,
            pb_manager=cls.pbehavior_manager
        )

    def test_get_watchers(self):
        """
        FIXIT: code isn't testable for now, waiting for configurable deletion
         + dependency injection.
        """
        filter_ = {}
        start = 0
        limit = None
        sort = False

        watchers = self.rhw.get_weather_watchers(
            filter_, start, limit, sort, self.context_manager,
            self.tracer_manager, self.pbehavior_manager,
            self.watcher_manager, self.alarmreader_manager
        )

if __name__ == '__main__':
    unittest.main()