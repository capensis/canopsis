#!/usr/bin/env python
# -*- coding: utf-8 -*-

# TODO: rewrite thoses tests with BaseApiTest

#import argparse
#import json
#import re
#import requests
#import time
#import unittest
#from pymongo import MongoClient
from canopsis.context_graph.import_ctx import ImportKey, ContextGraphImport

#from test_base import BaseApiTest, Method, HTTP

URL_BASE = "http://{0}:8082/"
URL_IMPORT = "{0}api/contextgraph/import"
URL_STATUS = "{0}api/contextgraph/import/status/{1}"
URL_AUTH = "{0}/?authkey={1}"
URL_MONGO = 'mongodb://cpsmongo:canopsis@{0}:27017/canopsis'

ENTITIES_COL = "default_entities"
IMPORT_COL = "default_importgraph"

JSON_EMPTY = {}
JSON_CIS_LIST_ONE = {ContextGraphImport.K_CIS: [
    {ContextGraphImport.K_ID: 'id_test',
     ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
     ContextGraphImport.K_PROPERTIES: {},
     ContextGraphImport.K_DEPENDS: [],
     ContextGraphImport.K_IMPACT: [],
     ContextGraphImport.K_INFOS: {},
     ContextGraphImport.K_MEASUREMENTS: [],
     ContextGraphImport.K_NAME: 'id_test',
     ContextGraphImport.K_TYPE: 'connector'}]}

JSON_CIS_LIST = {ContextGraphImport.K_CIS: [
    {ContextGraphImport.K_ID: 'id_test',
     ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
     ContextGraphImport.K_PROPERTIES: {},
     ContextGraphImport.K_DEPENDS: [],
     ContextGraphImport.K_IMPACT: [],
     ContextGraphImport.K_INFOS: {},
     ContextGraphImport.K_MEASUREMENTS: [],
     ContextGraphImport.K_NAME: 'id_test',
     ContextGraphImport.K_TYPE: 'connector'},
    {ContextGraphImport.K_ID: 'id_test1',
     ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
     ContextGraphImport.K_PROPERTIES: {},
     ContextGraphImport.K_DEPENDS: [],
     ContextGraphImport.K_IMPACT: [],
     ContextGraphImport.K_INFOS: {},
     ContextGraphImport.K_MEASUREMENTS: [],
     ContextGraphImport.K_NAME: 'id_test1',
     ContextGraphImport.K_TYPE: 'connector'}]}

JSON_LINKS_LIST_ONE = {
    ContextGraphImport.K_CIS: [
        {ContextGraphImport.K_ID: 'id_test',
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_DEPENDS: [],
         ContextGraphImport.K_IMPACT: [],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_MEASUREMENTS: [],
         ContextGraphImport.K_NAME: 'id_test',
         ContextGraphImport.K_TYPE: 'connector'},
        {ContextGraphImport.K_ID: 'id_test1',
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_DEPENDS: [],
         ContextGraphImport.K_IMPACT: [],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_MEASUREMENTS: [],
         ContextGraphImport.K_NAME: 'id_test1',
         ContextGraphImport.K_TYPE: 'connector'}],
    ContextGraphImport.K_LINKS: [
        {ContextGraphImport.K_ID: "id_0",
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_FROM: ["id_test"],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_TO: "id_test1"}]}

JSON_LINKS_LIST = {
    ContextGraphImport.K_CIS: [
        {ContextGraphImport.K_ID: 'id_test',
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_DEPENDS: [],
         ContextGraphImport.K_IMPACT: [],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_MEASUREMENTS: [],
         ContextGraphImport.K_NAME: 'id_test',
         ContextGraphImport.K_TYPE: 'connector'},
        {ContextGraphImport.K_ID: 'id_test1',
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_DEPENDS: [],
         ContextGraphImport.K_IMPACT: [],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_MEASUREMENTS: [],
         ContextGraphImport.K_NAME: 'id_test1',
         ContextGraphImport.K_TYPE: 'connector'}],
    ContextGraphImport.K_LINKS: [
        {ContextGraphImport.K_ID: "id_0",
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_FROM: ["id_test"],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_TO: "id_test1"},
        {ContextGraphImport.K_ID: "id_1",
         ContextGraphImport.K_ACTION: ContextGraphImport.A_CREATE,
         ContextGraphImport.K_PROPERTIES: {},
         ContextGraphImport.K_FROM: ["id_test1"],
         ContextGraphImport.K_INFOS: {},
         ContextGraphImport.K_TO: "id_test"}]}


# class ImportContextTest(BaseApiTest):

#     # def __init__(self, ent_col, imp_col):
#     #     super(ImportContextTest, self).__init__("all_tests")
#     #     self.session = None
#     #     self.cookies = None
#     #     self.ent_col = ent_col
#     #     self.imp_col = imp_col

#     def setUp(self):
#         self.ent_col.remove({})

#         self.session = requests.Session()
#         response = self.session.get(URL_AUTH)

#         if re.search("<title>Canopsis | Login</title>", response.text)\
#            is not None:
#             self.fail("Authentication error.")

#         self.cookies = response.cookies

#     def _launch_import(self, import_):
#         data = json.dumps({"json": import_})
#         headers = {"Content-type": "application/json", "Accept": "text/plain"}

#         response = self.session.put(URL_IMPORT,
#                                     data=data,
#                                     headers=headers,
#                                     cookies=self.cookies)

#         if re.search('"success": false', response.text) is not None:
#             self.fail("Error during import : {0}".format(response.text))

#         else:
#             response = json.loads(response.text)
#             id_ = response["data"][0]["import_id"]
#             status = ImportKey.ST_ONGOING

#             while status not in [ImportKey.ST_FAILED, ImportKey.ST_DONE]:
#                 time.sleep(1)
#                 response = self.session.get(URL_STATUS.format(URL_BASE, id_))
#                 response = json.loads(response.text)
#                 status = response[ImportKey.F_STATUS]

#             if status == ImportKey.ST_FAILED:
#                 self.fail("The import failed : {}".format(response[ImportKey.F_INFO]))

#         return id_

#     def _check_import_report(self, report, status, deleted, updated):
#         self.assertEqual(report[ImportKey.F_STATUS], status)
#         self.assertEqual(report[ImportKey.F_STATS][ImportKey.F_DELETED],
#                          deleted)
#         self.assertEqual(report[ImportKey.F_STATS][ImportKey.F_UPDATED],
#                          updated)

#     def all_tests(self):
#         self.test_empty_json()
#         self.ent_col.remove({"_id": {"$in": ["id_test", "id_test1"]}})
#         self.test_json_cis_list_item()
#         self.ent_col.remove({"_id": {"$in": ["id_test", "id_test1"]}})
#         self.test_json_cis_list_items()
#         self.ent_col.remove({"_id": {"$in": ["id_test", "id_test1"]}})
#         self.test_json_links_list_item()
#         self.ent_col.remove({"_id": {"$in": ["id_test", "id_test1"]}})
#         self.test_json_links_list_items()
#         self.ent_col.remove({"_id": {"$in": ["id_test", "id_test1"]}})

#     def test_empty_json(self):
#         uuid = self._launch_import(JSON_EMPTY)
#         report = list(self.imp_col.find({ImportKey.F_ID: uuid}))[0]
#         self._check_import_report(report, ImportKey.ST_DONE, 0, 0)

#         entities = list(self.ent_col.find({"_id":
#                                            {"$in": ["id_test", "id_test1"]}}))

#         self.assertListEqual(entities, [])

#     def assertListEntitiesEquals(self, l1, l2):
#         self.assertEqual(len(l1), len(l2))
#         l1 = sorted(l1)
#         l2 = sorted(l2)
#         for i in range(len(l1)):
#             self.assertDictEqual(l1[i], l2[i])

#     def test_json_cis_list_item(self):
#         uuid = self._launch_import(JSON_CIS_LIST_ONE)
#         report = list(self.imp_col.find({ImportKey.F_ID: uuid}))[0]

#         self._check_import_report(report, ImportKey.ST_DONE, 0, 1)

#         expected = JSON_CIS_LIST_ONE[ContextGraphImport.K_CIS]
#         expected[0].pop(ContextGraphImport.K_ACTION)
#         expected[0].pop(ContextGraphImport.K_PROPERTIES)

#         entities = list(self.ent_col.find({"_id":
#                                            {"$in": ["id_test"]}}))

#         self.assertListEntitiesEquals(entities, expected)

#     def test_json_cis_list_items(self):
#         uuid = self._launch_import(JSON_CIS_LIST)
#         report = list(self.imp_col.find({ImportKey.F_ID: uuid}))[0]

#         self._check_import_report(report, ImportKey.ST_DONE, 0, 2)

#         expected = JSON_CIS_LIST[ContextGraphImport.K_CIS]
#         for i in range(len(expected)):
#             expected[i].pop(ContextGraphImport.K_ACTION)
#             expected[i].pop(ContextGraphImport.K_PROPERTIES)
#         entities = list(self.ent_col.find({"_id":
#                                            {"$in": ["id_test", "id_test1"]}}))

#         self.assertListEntitiesEquals(entities, expected)

#     def test_json_links_list_item(self):
#         uuid = self._launch_import(JSON_LINKS_LIST_ONE)
#         report = list(self.imp_col.find({ImportKey.F_ID: uuid}))[0]

#         self._check_import_report(report, ImportKey.ST_DONE, 0, 2)

#         expected = JSON_LINKS_LIST_ONE[ContextGraphImport.K_CIS]
#         for i in range(len(expected)):
#             expected[i].pop(ContextGraphImport.K_ACTION)
#             expected[i].pop(ContextGraphImport.K_PROPERTIES)
#             if expected[i][ContextGraphImport.K_ID] == "id_test":
#                 expected[i][ContextGraphImport.K_IMPACT] = ["id_test1"]
#             if expected[i][ContextGraphImport.K_ID] == "id_test1":
#                 expected[i][ContextGraphImport.K_DEPENDS] = ["id_test"]

#         entities = list(self.ent_col.find({"_id":
#                                            {"$in": ["id_test", "id_test1"]}}))

#         self.assertListEntitiesEquals(entities, expected)

#     def test_json_links_list_items(self):
#         uuid = self._launch_import(JSON_LINKS_LIST)
#         report = list(self.imp_col.find({ImportKey.F_ID: uuid}))[0]

#         self._check_import_report(report, ImportKey.ST_DONE, 0, 2)

#         expected = JSON_LINKS_LIST[ContextGraphImport.K_CIS]
#         for i in range(len(expected)):
#             expected[i].pop(ContextGraphImport.K_ACTION)
#             expected[i].pop(ContextGraphImport.K_PROPERTIES)
#             if expected[i][ContextGraphImport.K_ID] == "id_test":
#                 expected[i][ContextGraphImport.K_IMPACT] = ["id_test1"]
#                 expected[i][ContextGraphImport.K_DEPENDS] = ["id_test1"]
#             if expected[i][ContextGraphImport.K_ID] == "id_test1":
#                 expected[i][ContextGraphImport.K_DEPENDS] = ["id_test"]
#                 expected[i][ContextGraphImport.K_IMPACT] = ["id_test"]
#         entities = list(self.ent_col.find({"_id":
#                                            {"$in": ["id_test", "id_test1"]}}))

#         self.assertListEntitiesEquals(entities, expected)



# ====================================== Not needed with BaseApiTest

# def parse_args():
#     parser = argparse.ArgumentParser()
#     parser.add_argument('-a', type=str, dest="authkey", help='authkey')
#     parser.add_argument('-w', type=str, dest="web_host", default="locahost",
#                         help='The webserver address.')
#     parser.add_argument('-m', type=str, dest="mongo_host", default="locahost",
#                         help='The mongodb address.')
#     return parser.parse_args()


# def setup(args):
#     global URL_BASE, URL_IMPORT, URL_AUTH

#     client = MongoClient(URL_MONGO.format(args.mongo_host))
#     db = client.canopsis
#     ent_col = db[ENTITIES_COL]
#     imp_col = db[IMPORT_COL]

#     URL_BASE = URL_BASE.format(args.web_host)
#     URL_AUTH = URL_AUTH.format(URL_BASE, args.authkey)
#     URL_IMPORT = URL_IMPORT.format(URL_BASE)

#     return ent_col, imp_col


# def main():
#     args = parse_args()
#     ent_col, imp_col = setup(args)
#     suite = unittest.TestSuite()
#     suite.addTest(ImportContextTest(ent_col, imp_col))
#     t = unittest.TextTestRunner()
#     t.run(suite)

# if __name__ == '__main__':
#     main()
