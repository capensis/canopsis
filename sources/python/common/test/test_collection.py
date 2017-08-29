from __future__ import unicode_literals

from pymongo.cursor import Cursor
import unittest

try:
    from io import StringIO
except ImportError:
    from StringIO import StringIO

from canopsis.common.collection import MongoCollection
from canopsis.logger import Logger, OutputStream
from canopsis.middleware.core import Middleware


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
        self.storage.remove_elements()

    def test_update(self):
        res = self.collection.update(query={'_id': self.id_},
                                     document={'yin': 'yang'})
        self.assertTrue(MongoCollection.is_mongo_successfull(res))
        self.assertEqual(res['n'], 0)

        res = self.collection.update(query={'_id': self.id_},
                                     document={'yin': 'yang'},
                                     upsert=True)
        self.assertTrue(MongoCollection.is_mongo_successfull(res))
        self.assertEqual(res['n'], 1)

    def test_find(self):
        res = self.collection.update(query={'_id': self.id_},
                                     document={'yin': 'yang'},
                                     upsert=True)
        self.assertTrue(MongoCollection.is_mongo_successfull(res))
        self.assertEqual(res['n'], 1)

        res = self.collection.find(query={'_id': self.id_})
        self.assertTrue(isinstance(res, Cursor))
        res = list(res)
        self.assertTrue(isinstance(res, list))
        self.assertEqual(res[0]['yin'], 'yang')

    def test_find_one(self):
        res = self.collection.update(query={'_id': self.id_},
                                     document={'yin': 'yang'},
                                     upsert=True)
        res = self.collection.update(query={'_id': 'anotherid'},
                                     document={'bambou': 'arbre'},
                                     upsert=True)
        self.assertTrue(MongoCollection.is_mongo_successfull(res))
        self.assertEqual(res['n'], 1)

        res = self.collection.find_one(query={'_id': self.id_})
        self.assertTrue(isinstance(res, dict))
        self.assertEqual(res['yin'], 'yang')

    def test_is_mongo_successfull(self):
        dico = {'ok': 1.0, 'n': 2}
        self.assertTrue(MongoCollection.is_mongo_successfull(dico))

        dico = {'ok': 666.66, 'n': 1}
        self.assertFalse(MongoCollection.is_mongo_successfull(dico))

        dico = {'n': 2}
        self.assertFalse(MongoCollection.is_mongo_successfull(dico))

if __name__ == '__main__':
    unittest.main()
