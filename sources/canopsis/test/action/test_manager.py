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
            "field": "resource",
            "regex": ".*wine.*",
            "parameters": {
                "author": "Matho",
                "name": "Salammbo",
                "reason": "Madness",
                "type": "Mercenary War",
                "rrule": ""
            }
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
        action2[Action.FIELD] = 'component'
        res = self.manager.update_id(id_=self.id_, action=action2)
        self.assertTrue(res)

        res = self.manager.get_id(self.id_)
        self.assertIsNotNone(res)
        self.assertDictEqual(res.to_dict(), action2)

        res = self.manager.delete_id(id_=self.id_)
        self.assertTrue(res)

        res = self.manager.get_id(self.id_)
        self.assertIsNone(res)

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
