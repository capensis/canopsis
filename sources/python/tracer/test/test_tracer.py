import unittest
from unittest import TestCase

from canopsis.tracer.manager import Trace, TraceSetError, TraceNotFound
from canopsis.tracer.manager import TracerManager

class Unencodable(object):
    pass

class TestTracerManager(TestCase):

    def setUp(self):
        self.manager = TracerManager('mongodb-default-testtracer://')

    def test_set_trace(self):
        self.manager.set_trace('test_trace', 'unittest')

    def test_get_trace(self):
        self.manager.set_trace('test_trace_get1', 'unittest')

        trace = self.manager.get_by_id('test_trace_get1')

        self.assertEqual(trace[Trace.ID], 'test_trace_get1')
        self.assertEqual(trace[Trace.TRIGGERED_BY], 'unittest')
        self.assertEqual(trace[Trace.IMPACT_ENTITIES], [])
        self.assertEqual(trace[Trace.EXTRA], {})

    def test_trace_not_found(self):
        with self.assertRaises(TraceNotFound):
            self.manager.get_by_id('test_trace_doesnt_exists')

    def test_del_trace(self):
        self.manager.set_trace('test_trace_del', 'unittest')
        trace = self.manager.get_by_id('test_trace_del')
        self.manager.del_trace('test_trace_del')
        with self.assertRaises(TraceNotFound):
            self.manager.get_by_id('test_trace_del')

    def test_set_trace_error(self):
        with self.assertRaises(TraceSetError):
            self.manager.set_trace('test_trace2', 'unittest', [Unencodable])

    def tearDown(self):
        self.manager.storage._backend.drop()

if __name__ == '__main__':
    unittest.main()