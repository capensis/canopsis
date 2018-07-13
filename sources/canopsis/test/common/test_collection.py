from __future__ import unicode_literals

from pymongo.cursor import Cursor
from six import string_types
import unittest

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.common import root_path
from canopsis.common.collection import MongoCollection
from canopsis.logger import Logger, OutputStream
from canopsis.common.middleware import Middleware
import xmlrunner


class TestMongoCollection(unittest.TestCase):

    def setUp(self):
        output = StringIO()
        self.logger = Logger.get('test', output, OutputStream)

        self.storage = Middleware.get_middleware_by_uri(
            'storage-default-testmongocollection://'
        )

        self.collection = MongoCollection(collection=self.storage._backend,
                                          logger=self.logger)

        self.id_ = 'testid'

    def tearDown(self):
        """Teardown"""
        self.collection.remove()

    def test_insert(self):
        res = self.collection.insert(document={'_id': self.id_})
        self.assertEqual(res, self.id_)

        res2 = self.collection.insert(document={'up': 'down'})
        self.assertTrue(isinstance(res, string_types))
        self.assertNotEqual(res, res2)

    def test_update(self):
        res = self.collection.update(query={'_id': self.id_},
                                     document={'strange': 'charm'})
        self.assertTrue(MongoCollection.is_successfull(res))
        self.assertEqual(res['n'], 0)

        res = self.collection.update(query={'_id': self.id_},
                                     document={'yin': 'yang'},
                                     upsert=True)
        self.assertTrue(MongoCollection.is_successfull(res))
        self.assertEqual(res['n'], 1)

        res = self.collection.find_one(self.id_)
        self.assertEqual(res['yin'], 'yang')
        self.assertTrue('strange' not in res)

    def test_remove(self):
        res = self.collection.insert(document={'_id': self.id_, 'top': 'bottom'})
        self.assertIsNotNone(res)

        res = self.collection.remove(query={'_id': self.id_})
        self.assertTrue(MongoCollection.is_successfull(res))
        self.assertEqual(res['n'], 1)

        # Deleting non-existing object doesn't throw error
        res = self.collection.remove(query={})
        self.assertTrue(MongoCollection.is_successfull(res))
        self.assertEqual(res['n'], 0)

    def test_find(self):
        res = self.collection.insert(document={'_id': self.id_, 'yin': 'yang'})
        self.assertIsNotNone(res)

        res = self.collection.find(query={'_id': self.id_})
        self.assertTrue(isinstance(res, Cursor))
        res = list(res)
        self.assertTrue(isinstance(res, list))
        self.assertEqual(res[0]['yin'], 'yang')

    def test_find_one(self):
        res = self.collection.insert(document={'_id': self.id_, 'up': 'down'})
        self.assertIsNotNone(res)
        res = self.collection.insert(document={'strange': 'charm'})
        self.assertIsNotNone(res)

        res = self.collection.find_one(query={'_id': self.id_})
        self.assertTrue(isinstance(res, dict))
        self.assertEqual(res['up'], 'down')

    def test_is_successfull(self):
        dico = {'ok': 1.0, 'n': 2}
        self.assertTrue(MongoCollection.is_successfull(dico))

        dico = {'ok': 666.667, 'n': 1}
        self.assertFalse(MongoCollection.is_successfull(dico))

        dico = {'n': 2}
        self.assertFalse(MongoCollection.is_successfull(dico))

if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
