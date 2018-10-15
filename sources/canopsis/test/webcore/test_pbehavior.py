from unittest import TestCase

from canopsis.pbehavior.manager import PBehaviorManager
from canopsis.watcher.manager import Watcher as WatcherManager
from canopsis.webcore.services.pbehavior import check_values, RouteHandlerPBehavior
import unittest
from canopsis.common import root_path
import xmlrunner


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
        config, logger, storage = PBehaviorManager.provide_default_basics()
        pb_manager = PBehaviorManager(config, logger, storage)
        watcher_manager = WatcherManager()
        cls.rhpb = RouteHandlerPBehavior(pb_manager, watcher_manager)

    def test_create_bad_pb(self):
        with self.assertRaises(ValueError):
            self.rhpb.create(**self.INVALID_PB)

    def test_create_pb(self):
        self.rhpb.create(**self.VALID_PB)

    def test_read_pb(self):
        pb_id = self.rhpb.create(**self.VALID_PB)
        self.assertIsInstance(pb_id, str)

        pbehavior = self.rhpb.read(pb_id)

        self.assertIsInstance(pbehavior, dict)
        self.assertEquals(pbehavior.get('name'), self.VALID_PB.get('name'))

    def test_update_pb(self):
        pb_id = self.rhpb.create(**self.VALID_PB)

        pbehavior = self.rhpb.read(pb_id)

        self.assertEquals(pbehavior.get('name'), self.VALID_PB.get('name'))

        updated_pb = self.VALID_PB.copy()
        updated_pb['name'] = 'pb_new_name'
        updated_pb['author'] = 'pb_new_author'
        updated_pb['filter_'] = '{"new": "filter"}'
        updated_pb['tstart'] = pbehavior.get('tstart') * 10
        updated_pb['tstop'] = pbehavior.get('tstop') * 10
        updated_pb['rrule'] = 'FREQ=DAILY;BYDAY=SU'
        updated_pb['enabled'] = "true"
        updated_pb['comments'] = [{'author': 'test', 'message': 'test'}]
        updated_pb['connector'] = 'test_pb_new'
        updated_pb['connector_name'] = 'test_pb_new'

        res = self.rhpb.update(pb_id, **updated_pb)
        self.assertIsInstance(res, dict)


        updated_pb = self.rhpb.read(pb_id)
        self.assertEquals(updated_pb.get('name'), 'pb_new_name')
        self.assertEquals(updated_pb.get('author'), 'pb_new_author')
        self.assertEquals(updated_pb.get('filter_'), '{"new": "filter"}')
        self.assertEquals(updated_pb.get('tstart'), pbehavior.get('tstart') * 10)
        self.assertEquals(updated_pb.get('tstop'), pbehavior.get('tstop') * 10)
        self.assertEquals(updated_pb.get('rrule'), 'FREQ=DAILY;BYDAY=SU')
        self.assertEquals(updated_pb.get('enabled'), True)
        self.assertListEqual(updated_pb.get('comments'), [{'author': 'test', 'message': 'test'}])
        self.assertEquals(updated_pb.get('connector'), 'test_pb_new')
        self.assertEquals(updated_pb.get('connector_name'), 'test_pb_new')

    def test_delete_pb(self):
        pb_id = self.rhpb.create(**self.VALID_PB)

        pbehavior = self.rhpb.read(pb_id)

        self.assertEquals(pbehavior.get('name'), self.VALID_PB.get('name'))

        delres = self.rhpb.delete(pb_id)
        self.assertEquals(delres.get('deletedCount'), 1)

        pbehavior = self.rhpb.read(pb_id)
        self.assertIsNone(pbehavior)

    def test_create_comments_pb(self):
        pb_id = self.rhpb.create(**self.VALID_PB)

        c1 = self.rhpb.create_comment(pb_id, author='pb_test', message='pb_comment_msg')
        c2 = self.rhpb.create_comment(pb_id, author='pb_test', message='pb_comment_msg2')
        c3 = self.rhpb.create_comment(pb_id, author='pb_test2', message='pb_comment_msg')
        c4 = self.rhpb.create_comment(pb_id, author='pb_test2', message='pb_comment_msg2')

        self.assertIsInstance(c1, str)
        self.assertIsInstance(c2, str)
        self.assertIsInstance(c3, str)
        self.assertIsInstance(c4, str)

        pbehavior = self.rhpb.read(pb_id)

        comments = pbehavior.get('comments')

        self.assertEquals(comments[0].get('author'), 'pb_test')
        self.assertEquals(comments[0].get('message'), 'pb_comment_msg')

        self.assertEquals(comments[1].get('author'), 'pb_test')
        self.assertEquals(comments[1].get('message'), 'pb_comment_msg2')

        self.assertEquals(comments[2].get('author'), 'pb_test2')
        self.assertEquals(comments[2].get('message'), 'pb_comment_msg')

        self.assertEquals(comments[3].get('author'), 'pb_test2')
        self.assertEquals(comments[3].get('message'), 'pb_comment_msg2')

    def test_update_comment_pb(self):
        pb_id = self.rhpb.create(**self.VALID_PB)

        c1 = self.rhpb.create_comment(pb_id, author='pb_test', message='pb_comment_msg')

        self.assertIsInstance(c1, str)

        self.rhpb.update_comment(pb_id, c1, author='pb_test_new', message='pb_comment_new')

        pbehavior = self.rhpb.read(pb_id)
        comments = pbehavior.get('comments')

        self.assertEquals(len(comments), 1)
        self.assertEquals(comments[0].get('author'), 'pb_test_new')
        self.assertEquals(comments[0].get('message'), 'pb_comment_new')

    def test_delete_comment_pb(self):
        pb_id = self.rhpb.create(**self.VALID_PB)

        c1 = self.rhpb.create_comment(pb_id, author='pb_test', message='pb_comment_msg')

        self.assertIsInstance(c1, str)

        pbehavior = self.rhpb.read(pb_id)
        comments = pbehavior.get('comments')

        self.assertEquals(len(comments), 1)

        res = self.rhpb.delete_comment(pb_id, c1)

        self.assertEquals(res.get('deletedCount'), 1)

        pbehavior = self.rhpb.read(pb_id)
        comments = pbehavior.get('comments')

        self.assertEquals(len(comments), 0)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
