import unittest

from canopsis.common.collection import MongoCollection
from canopsis.common.mongo_store import MongoStore
from canopsis.dynamic_infos.manager import DynamicInfosManager, DynamicInfosRule
from canopsis.logger.logger import Logger, OutputNull


class DynamicInfosManagerTest(unittest.TestCase):
    def setUp(self):
        mongo = MongoStore.get_default()
        collection = mongo.get_collection("test_dynamic_infos")
        self.dynamic_coll = MongoCollection(collection)

        logger = Logger.get('test_dynamic_infos', None, output_cls=OutputNull)
        self.dynamic_infos = DynamicInfosManager(logger=logger,
                                                 mongo_collection=self.dynamic_coll)
        self.dynamic_coll.drop()
        self.dynamic_infos_doc = {
            "_id": "rule2",
            "name": "Test",
            "author": "billy",
            "creation_date": 1576260000,
            "last_modified_date": 1576260000,
            "description": "Freedom !",
            "infos": [
                {
                    "name": "info",
                    "value": "value"
                },
                {
                    "name": "info2",
                    "value": "value2"
                }
            ],
            "entity_patterns": [
                {
                    "_id": "cpu/billys-laptop"
                }
            ],
            "alarm_patterns": [
                {
                    "v": {
                        "state": {
                            "val": 3
                        }
                    }
                }
            ]
        }

    def tearDown(self):
        self.dynamic_coll.drop()

    def test_count(self):
        rule = DynamicInfosRule.new_from_dict(self.dynamic_infos_doc, "test_author", 1583301306)
        self.dynamic_infos.create(rule)
        count = self.dynamic_infos.count()
        self.assertEqual(count, 1)
        count = self.dynamic_infos.count(search="test_author", search_fields=["author"])
        self.assertEqual(count, 1)
        count = self.dynamic_infos.count(search="test_author", search_fields=["description"])
        self.assertEqual(count, 0)
        count = self.dynamic_infos.count(search="Test", search_fields=["name"])
        self.assertEqual(count, 1)
