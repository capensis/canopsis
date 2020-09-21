from __future__ import unicode_literals

import unittest

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.action.manager import ActionManager
from canopsis.common import root_path
from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.logger import Logger, OutputStream
from canopsis.models.action import Action
import xmlrunner


class TestActionManager(unittest.TestCase):

    def setUp(self):
        output = StringIO()
        self.logger = Logger.get('test', output, OutputStream)

        store = MongoStore.get_default()
        self.collection = store.get_collection(name='default_test')
        self.mongo_collection = MongoCollection(
            collection=self.collection,
            logger=self.logger
        )
        # Cleanup
        self.tearDown()

        self.manager = ActionManager(
            logger=self.logger,
            mongo_collection=self.mongo_collection
        )

        self.id_ = 'testid'
        self.action = {
            "_id": self.id_,
            "hook": None,
            "type": "pbehavior",
            "fields": ["Resource"],
            "regex": ".*wine.*",
            "parameters": {
                "author": "Matho",
                "name": "Salammbo",
                "reason": "Madness",
                "type": "Mercenary War",
                "rrule": ""
            },
            "delay": ""
        }

    def tearDown(self):
        """Teardown"""
        self.mongo_collection.remove({})

    def test_crud(self):
        res = self.manager.create(action=self.action)
        self.assertTrue(res)

        res = self.manager.get_id(self.id_)
        self.assertIsNotNone(res)
        self.assertDictEqual(res.to_dict(), self.action)

        action2 = self.action.copy()
        action2[Action.FIELDS] = ['Component']
        res = self.manager.update_id(id_=self.id_, action=action2)
        self.assertTrue(res)

        res = self.manager.get_id(self.id_)
        self.assertIsNotNone(res)
        self.assertDictEqual(res.to_dict(), action2)

        res = self.manager.delete_id(id_=self.id_)
        self.assertTrue(res)

        res = self.manager.get_id(self.id_)
        self.assertIsNone(res)

    def test_is_delay_valid(self):
        self.assertTrue(self.manager.is_delay_valid(""))
        self.assertTrue(self.manager.is_delay_valid("30s"))
        self.assertTrue(self.manager.is_delay_valid("1.0s"))
        self.assertTrue(self.manager.is_delay_valid("1m"))
        self.assertTrue(self.manager.is_delay_valid("10h"))
        self.assertTrue(self.manager.is_delay_valid("1m3s"))
        self.assertFalse(self.manager.is_delay_valid("1m3"))
        self.assertFalse(self.manager.is_delay_valid("1k"))
        self.assertFalse(self.manager.is_delay_valid("1h3m5p"))


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
