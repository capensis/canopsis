#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals

import unittest
from unittest import main, TestCase
from canopsis.common import root_path
from canopsis.context_graph.manager import ContextGraph
from canopsis.common.middleware import Middleware
from canopsis.logger import Logger, OutputNull

from test_infos_filter import Keys, SCHEMA, TEMPLATE_INFOS

import xmlrunner

logger = Logger.get("", None, output_cls=OutputNull)


def create_entity(
        id,
        name,
        etype,
        depends=[],
        impact=[],
        measurements=[],
        links={},
        infos={}):
    return {'_id': id,
            'type': etype,
            'name': name,
            'depends': depends,
            'impact': impact,
            'measurements': measurements,
            'links': links,
            'infos': infos
            }


def create_comp(id, name, depends=[], impact=[], measurements=[], infos={}):
    """Create a component with the given parameters"""
    return create_entity(id,
                         name,
                         "component",
                         depends,
                         impact,
                         measurements,
                         infos)


def create_conn(id, name, depends=[], impact=[], measurements=[], infos={}):
    """Create a connector with the given parameters"""
    return create_entity(id,
                         name,
                         "connector",
                         depends,
                         impact,
                         measurements,
                         infos)


def create_re(id, name, depends=[], impact=[], measurements=[], infos={}):
    """Create a resource with the given parameters"""
    return create_entity(id,
                         name,
                         "resource",
                         depends,
                         impact,
                         measurements,
                         infos)


class BaseTest(TestCase):

    def setUp(self):
        self.manager = ContextGraph(logger)
        self.entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        self.organisations_storage = Middleware.get_middleware_by_uri(
            'storage-default-testorganisations://'
        )
        self.users_storage = Middleware.get_middleware_by_uri(
            'storage-default-testusers://'
        )

        self.manager.ent_storage = self.entities_storage

        self.template = {'_id': None,
                         'type': 'connector',
                         'name': 'conn-name1',
                         'depends': [],
                         'impact': [],
                         'measurements': [],
                         'links': {},
                         'infos': {}}

    def tearDown(self):
        self.entities_storage.remove_elements()
        self.organisations_storage.remove_elements()
        self.users_storage.remove_elements()

    def assertEqualEntities(self, entity1, entity2):
        sorted(entity1["depends"])
        sorted(entity1["impact"])
        sorted(entity2["depends"])
        sorted(entity2["impact"])
        self.assertDictEqual(entity1, entity2)

    def _insertion_filter_test(self, function, expected=None):
        infos = {}

        entity = create_entity("id", "a name", "resource", infos=infos)
        if expected is None:
            expected = entity.copy()

        function(entity)
        result = self.manager.get_entities_by_id("id")[0]
        # pop non predictable fields
        try:
            result.pop("enable_history")
        except KeyError:
            pass

        self.assertEqualEntities(result, expected)

    def _insertion_filter_test_not_allowed_field(self, function, expected=None):

        infos = {"I am not allowed to be here": [1],
                 "me too": [1]}
        entity = create_entity("id", "a name", "resource", infos=infos)
        if expected is None:
            expected = entity.copy()

        function(entity)
        result = self.manager.get_entities_by_id("id")[0]
        # pop non predictable fields
        try:
            result.pop("enable_history")
        except KeyError:
            pass

        self.assertEqualEntities(result, expected)

class CreateEntityDict(BaseTest):

    def test_create_entity(self):
        id_ = "id_1"
        name = "name_1"
        etype = "resource"
        depends = ["id_2", "id_3", "id_4", "id_5"]
        impacts = ["id_6", "id_7", "id_8", "id_9"]
        measurements = {"tag_1": "data_1", "tag_2": "data_2"}
        infos = {"info_1": "foo_1", "info_2": "bar_2"}

        ent = ContextGraph.create_entity_dict(id_, name, etype, depends,
                                    impacts, measurements, infos)

        self.assertEqual(id_, ent["_id"])
        self.assertEqual(name, ent["name"])
        self.assertEqual(etype, ent["type"])
        self.assertEqual(depends, ent["depends"])
        self.assertEqual(impacts, ent["impact"])
        self.assertEqual(measurements, ent["measurements"])
        self.assertEqual(infos, ent["infos"])

    def test_create_component(self):
        id_ = "id_1"
        name = "name_1"
        etype = "component"
        depends = ["id_2", "id_3", "id_4", "id_5"]
        impacts = ["id_6", "id_7", "id_8", "id_9"]
        measurements = {"tag_1": "data_1", "tag_2": "data_2"}
        infos = {"info_1": "foo_1", "info_2": "bar_2"}

        ent = ContextGraph.create_entity_dict(id_, name, etype, depends,
                                    impacts, measurements, infos)

        self.assertEqual(id_, ent["_id"])
        self.assertEqual(name, ent["name"])
        self.assertEqual(etype, ent["type"])
        self.assertEqual(depends, ent["depends"])
        self.assertEqual(impacts, ent["impact"])
        self.assertNotIn("measurements", ent.keys())
        self.assertEqual(infos, ent["infos"])


class GetEvent(BaseTest):
    """Test get_event method.
    """

    def setUp(self):
        super(GetEvent, self).setUp()
        self.manager.ent_storage.remove_elements()

    def test_get_check_event(self):

        # set of required field
        fields = sorted(['component', 'connector', 'connector_name',
                         'event_type', 'long_output', 'output', 'source_type',
                         'state', 'state_type', 'timestamp', 'links'])

        conn = "conn"
        conn_name = "conn-name"
        entity_id = conn + "/" + conn_name

        entity = create_conn(entity_id, conn_name, depends=[],
                             impact=[])

        event = self.manager.get_event(
            entity, event_type='check', output='test'
        )

        self.assertListEqual(sorted(event.keys()), fields)
        self.assertEqual(event['event_type'], 'check')
        self.assertEqual(event['output'], 'test')


class GetID(TestCase):

    def setUp(self):
        self.event = {
            'connector': 'connector',
            'connector_name': 'connector-0',
            'component': 'c',
            'output': '...',
            'timestamp': 0,
            'source_type': None,
            'resource': 'resource-0'
        }

    def test_get_id_component(self):
        self.event["source_type"] = 'component'
        expected_id = self.event["component"]
        result = ContextGraph.get_id(self.event)

        self.assertEquals(result, expected_id)

    def test_get_id_resource(self):
        self.event["source_type"] = 'resource'
        expected_id = "{0}/{1}".format(self.event["resource"],
                                       self.event["component"])
        result = ContextGraph.get_id(self.event)

        self.assertEquals(result, expected_id)

    def test_get_id_connector(self):
        self.event["source_type"] = 'connector'
        expected_id = "{0}/{1}".format(self.event["connector"],
                                       self.event["connector_name"])
        result = ContextGraph.get_id(self.event)

        self.assertEquals(result, expected_id)

    def test_get_id_error(self):
        self.event["source_type"] = 'something_else'
        error_desc = ("source_type should be one of 'connector', 'resource' or "
                      "'component' not {}.".format(self.event["source_type"]))
        with self.assertRaisesRegexp(ValueError, error_desc):
            ContextGraph.get_id(self.event)


class GetEntitiesByID(BaseTest):

    def test_get_entity_by_id_id(self):
        entity = self.template.copy()
        entity['_id'] = 'conn1/conn-name'
        entity['type'] = 'connector'
        entity['name'] = 'conn-name'

        self.manager._put_entities(entity)

        result = self.manager.get_entities_by_id(entity["_id"])
        self.assertIsInstance(result, type([]))
        self.assertDictEqual(result[0], entity)

    def test_get_entity_by_id_ids(self):
        entities = [self.template.copy(),
                    self.template.copy(),
                    self.template.copy()]

        entities[0]['_id'] = 'conn1/conn-name1'
        entities[0]['type'] = 'connector'
        entities[0]['name'] = 'conn-name1'

        entities[1]['_id'] = 'conn2/conn-name2'
        entities[1]['type'] = 'connector'
        entities[1]['name'] = 'conn-name2'

        entities[2]['_id'] = 'conn3/conn-name3'
        entities[2]['type'] = 'connector'
        entities[2]['name'] = 'conn-name3'

        sorted(entities)
        self.manager._put_entities(entities)

        ids = [x["_id"] for x in entities]

        result = self.manager.get_entities_by_id(ids)

        sorted(result)
        self.assertIsInstance(result, type([]))
        self.assertListEqual(result, entities)

    def test_get_entity_by_id(self):
        """Test the behaviour of the get_entity_by_id function with the id
        of a nonexistant entity"""

        result = self.manager.get_entities_by_id("id")

        self.assertIsInstance(result, type([]))
        self.assertEqual(len(result), 0)


class PutEntities(BaseTest):

    def test_put_entities_entity(self):
        entity = self.template.copy()
        entity['_id'] = 'conn1/conn-name'
        entity['type'] = 'connector'
        entity['name'] = 'conn-name'

        self.manager._put_entities(entity)

        result = self.manager.get_entities_by_id(entity["_id"])[0]

        self.assertDictEqual(result, entity)

    def test_put_entities_entities(self):
        entities = [self.template.copy(),
                    self.template.copy(),
                    self.template.copy()]

        entities[0]['_id'] = 'conn1/conn-name1'
        entities[0]['type'] = 'connector'
        entities[0]['name'] = 'conn-name1'

        entities[1]['_id'] = 'conn2/conn-name2'
        entities[1]['type'] = 'connector'
        entities[1]['name'] = 'conn-name2'

        entities[2]['_id'] = 'conn3/conn-name3'
        entities[2]['type'] = 'connector'
        entities[2]['name'] = 'conn-name3'

        sorted(entities)

        self.manager._put_entities(entities)

        ids = [x["_id"] for x in entities]

        result = self.manager.get_entities_by_id(ids)
        sorted(result)
        self.assertListEqual(result, entities)

    def test_create_infos_filter(self):
        self._insertion_filter_test(self.manager._put_entities)

    def test_create_infos_filter_not_allowed_field(self):
        self._insertion_filter_test(self.manager._put_entities)


class DeleteEntities(BaseTest):

    def test_no_entities(self):
        id1 = "ent1"
        id2 = "not_an_entity_id"
        entity = self.template.copy()
        entity["_id"] = id1
        entity['type'] = 'connector'
        entity['name'] = 'conn-name'

        self.manager._put_entities(entity)

        self.manager._delete_entities(id2)

        result = self.manager.get_entities_by_id([id1, id2])

        self.assertEqual(len(result), 1)
        self.assertEqualEntities(entity, result[0])

    def test_entity(self):
        id_ = "ent1"
        entity = self.template.copy()
        entity["_id"] = id_
        entity['type'] = 'connector'
        entity['name'] = 'conn-name'

        self.manager._put_entities(entity)

        self.manager._delete_entities(id_)

        result = self.manager.get_entities_by_id(id_)

        self.assertEqual(len(result), 0)

    def test_entities(self):
        id1 = "ent1"
        entity1 = self.template.copy()
        entity1["_id"] = id1
        entity1['type'] = 'connector'
        entity1['name'] = 'conn-name'

        id2 = "ent2"
        entity2 = self.template.copy()
        entity2["_id"] = id2
        entity2['type'] = 'connector'
        entity2['name'] = 'conn-name'

        self.manager._put_entities([entity1, entity2])

        self.manager._delete_entities([id1, id2])

        result = self.manager.get_entities_by_id([id1, id2])

        self.assertEqual(len(result), 0)


class GetAllEntitiesId(BaseTest):

    def test_get_all_entities_id(self):
        entities = [self.template.copy(),
                    self.template.copy(),
                    self.template.copy()]

        entities[0]['_id'] = 'conn1/conn-name1'
        entities[0]['type'] = 'connector'
        entities[0]['name'] = 'conn-name1'

        entities[1]['_id'] = 'conn2/conn-name2'
        entities[1]['type'] ='connector'
        entities[1]['name'] = 'conn-name2'

        entities[2]['_id'] = 'conn3/conn-name3'
        entities[2]['type'] = 'connector'
        entities[2]['name'] = 'conn-name3'

        self.manager._put_entities(entities)
        expected = set()
        for entity in entities:
            expected.add(entity["_id"])

        result = self.manager.get_all_entities_id()
        self.assertSetEqual(result, expected)


class UpdateEntity(BaseTest):

    def setUp(self):
        super(UpdateEntity, self).setUp()

        self.ent1 = self.template.copy()
        self.ent2 = self.template.copy()
        self.ent3 = self.template.copy()
        self.ent4 = self.template.copy()

        self.ent1["_id"] = "ent1"
        self.ent2["_id"] = "ent2"
        self.ent3["_id"] = "ent3"
        self.ent4["_id"] = "ent4"

        self.manager._put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

    def test_update_entity_wrong_id(self):
        fake_entity = {"_id": "wrong id"}
        error_desc = "The _id {0} does not match any entity in database."\
                     .format(fake_entity["_id"])
        with self.assertRaisesRegexp(ValueError, error_desc):
            self.manager.update_entity(fake_entity)

    def __test_multiple(self, from_, to):
        self.ent1[from_] = ["ent2", "ent3"]
        self.manager.update_entity(self.ent1)

        entity = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(self.ent1, entity)

        entity = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.ent2[to] = ["ent1"]
        self.assertEqualEntities(self.ent2, entity)

        entity = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.ent3[to] = ["ent1"]
        self.assertEqualEntities(self.ent3, entity)

        entity = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(self.ent4, entity)

    def __test_single(self, from_, to):
        self.ent1[from_] = ["ent2"]
        self.manager.update_entity(self.ent1)

        entity = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(self.ent1, entity)

        entity = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.ent2[to] = ["ent1"]
        self.assertEqualEntities(self.ent2, entity)

        entity = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.assertEqualEntities(self.ent3, entity)

        entity = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(self.ent4, entity)

    def test_update_entity_update_depends_insert_multiple(self):
        self.__test_multiple("depends", "impact")

    def test_update_entity_update_depends_insert_single(self):
        self.__test_single("depends", "impact")

    def test_update_entity_update_impact_insert_multiple(self):
        self.__test_multiple("impact", "depends")

    def test_update_entity_update_impact_insert_single(self):
        self.__test_single("impact", "depends")

    def test_update_other_fields(self):
        new_entity = self.ent1.copy()
        new_entity['infos'] = {"Change": True}

        self.manager.update_entity(new_entity)
        result = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(result, new_entity)

    def test_update_no_update(self):
        new_entity = self.ent1.copy()

        self.manager.update_entity(new_entity)
        result = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(result, new_entity)


class UpdateDependancies(BaseTest):

    def test_update_dependancies_wrong_type(self):

        desc = "Dependancy_type should be depends or impact not {0}."

        with self.assertRaisesRegexp(ValueError, desc.format(None)):
            self.manager._ContextGraph__update_dependancies(None, None, None)

    def __test(self, from_, to, delete):

        self.ent1 = self.template.copy()
        self.ent2 = self.template.copy()
        self.ent3 = self.template.copy()
        self.ent4 = self.template.copy()

        self.ent1["_id"] = "ent1"
        self.ent1[from_] = ["ent2", "ent3"]
        self.ent2["_id"] = "ent2"
        self.ent3["_id"] = "ent3"
        self.ent4["_id"] = "ent4"
        self.ent4[to] = ["dummy"]

        if delete:
            self.ent2[to] = ["ent1"]
            self.ent3[to] = ["ent1", "dummy"]
        else:
            self.ent3[to] = ["dummy"]

        self.manager._put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

        if delete:
            status = {"deletions": ["ent2", "ent3"],
                      "insertions": []}
        else:
            status = {"deletions": [],
                      "insertions": ["ent2", "ent3"]}

        result = self.manager._ContextGraph__update_dependancies(
            self.ent1["_id"], status, from_)

        for entity in result:
            if entity["_id"] == "ent2":
                expected = self.template.copy()
                expected["_id"] = "ent2"

                if delete is False:
                    expected[to] = ["ent1"]

                self.assertEqualEntities(entity, expected)

            elif entity["_id"] == "ent3":
                expected = self.template.copy()
                expected["_id"] = "ent3"

                if delete is False:
                    expected[to] = ["dummy", "ent1"]
                else:
                    expected[to] = ["dummy"]

                self.assertEqualEntities(entity, expected)

            else:
                self.fail(
                    "The entity of id {0} should not be here".format(
                        entity["_id"]))

    def test_update_dependancies_entity_depends_delete(self):
        self.__test("depends", "impact", True)

    def test_update_dependancies_entity_impact_delete(self):
        self.__test("impact", "depends", True)

    def test_update_dependancies_entity_depends_insert(self):
        self.__test("depends", "impact", False)

    def test_update_dependancies_entity_impact_insert(self):
        self.__test("impact", "depends", False)


class DeleteEntity(BaseTest):

    def test_delete_entity_wrong_id(self):
        id_ = "this_is_not_an_id"
        desc = "No entity found for the following id : {0}".format(id_)

        with self.assertRaisesRegexp(ValueError, desc):
            self.manager.delete_entity(id_)

    def __test(self, from_, to):
        self.ent1 = self.template.copy()
        self.ent2 = self.template.copy()
        self.ent3 = self.template.copy()
        self.ent4 = self.template.copy()

        self.ent1["_id"] = "ent1"
        self.ent1[from_] = ["ent2", "ent3"]
        self.ent2["_id"] = "ent2"
        self.ent2[to] = ["ent1"]
        self.ent3["_id"] = "ent3"
        self.ent3[to] = ["ent1", "dummy"]
        self.ent4["_id"] = "ent4"
        self.ent4[from_] = ["dummy"]

        self.manager._put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

        self.manager.delete_entity(self.ent1["_id"])

        result = self.manager.get_entities_by_id([self.ent1,
                                                  self.ent2,
                                                  self.ent3,
                                                  self.ent4])

        for entity in result:
            if entity["_id"] == "ent1":
                self.fail("The entity ent1 should be deleted")

            if entity["_id"] == "ent2":
                expected = self.template.copy()
                expected["_id"] = "ent2"

                self.assertEqualEntities(entity, expected)

            elif entity["_id"] == "ent3":
                expected = self.template.copy()
                expected["_id"] = "ent3"
                expected[to] = ["dummy"]

                self.assertEqualEntities(entity, expected)

            elif entity["_id"] == "ent4":
                expected = self.template.copy()
                expected["_id"] = "ent4"
                expected[to] = ["dummy"]

                self.assertEqualEntities(entity, expected)

    def test_delete_entity_depends(self):
        self.__test("depends", "impact")

    def test_delete_entity_impact(self):
        self.__test("impact", "depends")


class CreateEntity(BaseTest):

    SCHEMA = {"schema": {
        "properties": {
            "disable_history": "string",
            "enable_history": "string",
            "enabled": "boolean"
        }
    }}

    def test_create_entity_entity_exists(self):
        entity = self.template.copy()
        entity['_id'] = "I am here"
        entity['type'] = 'connector'
        entity['name'] = 'conn-name1'

        self.manager.create_entity(entity)

        desc = "An entity id {0} already exist".format(entity["_id"])
        with self.assertRaisesRegexp(ValueError, desc):
            self.manager.create_entity(entity)

    def __test(self, from_, to):
        self.ent1 = self.template.copy()
        self.ent2 = self.template.copy()
        self.ent3 = self.template.copy()
        self.ent4 = self.template.copy()

        self.ent1["_id"] = "ent1"
        self.ent1[from_] = ["ent2", "ent3"]
        self.ent2["_id"] = "ent2"
        self.ent2[to] = []
        self.ent3["_id"] = "ent3"
        self.ent3[to] = ["dummy"]
        self.ent4["_id"] = "ent4"
        self.ent4[from_] = ["dummy"]
        self.ent4[to] = ["dummy"]

        self.manager._put_entities(self.ent2)
        self.manager._put_entities(self.ent3)
        self.manager._put_entities(self.ent4)

        self.manager.create_entity(self.ent1)

        result = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(result, self.ent1)

        result = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.ent2[to].append(self.ent1["_id"])
        self.assertEqualEntities(result, self.ent2)

        result = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.ent3[to].append(self.ent1["_id"])
        self.assertEqualEntities(result, self.ent3)

        result = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(result, self.ent4)

    def test_create_entity_entity_depends(self):
        self.__test("depends", "impact")

    def test_create_entity_entity_impact(self):
        self.__test("impact", "depends")

    def test_create_infos_filter(self):
        self.manager.filter_._schema = self.SCHEMA
        infos = {}
        expected = create_entity("id", "a name", "resource", infos=infos)
        expected["enabled"] = True
        self._insertion_filter_test(
            self.manager.create_entity, expected
        )

    def test_create_infos_filter_not_allowed_field(self):
        self.manager.filter_._schema = self.SCHEMA
        infos = {}
        expected = create_entity("id", "a name", "resource", infos=infos)
        expected["enabled"] = True
        self._insertion_filter_test_not_allowed_field(
            self.manager.create_entity, expected
        )


class TestGraphRequests(BaseTest):

    def setUp(self):
        super(TestGraphRequests, self).setUp()
        self.manager.collection_name = 'default_testentities'

    def test(self):
        ent1 = create_entity(
            id='e1',
            name='e1',
            etype='component',
            depends=[],
            impact=['e2']
        )
        ent2 = create_entity(
            id='e2',
            name='e2',
            etype='component',
            depends=['e1'],
            impact=['e3']
        )
        ent3 = create_entity(
            id='e3',
            name='e3',
            etype='component',
            depends=['e2'],
            impact=['e4']
        )
        ent4 = create_entity(
            id='e4',
            name='e4',
            etype='component',
            depends=['e3'],
            impact=[]
        )

        self.manager._put_entities(ent1)
        self.manager._put_entities(ent2)
        self.manager._put_entities(ent3)
        self.manager._put_entities(ent4)

        res = self.manager.get_graph_depends('e4')['graph']
        res.sort()
        self.assertEqual(res, [
            {
                'impact': ['e2'],
                'depth': 3L,
                '_id': 'e1',
                'name': 'e1',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': []
            }, {
                'impact': ['e3'],
                'depth': 2L,
                '_id': 'e2',
                'name': 'e2',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e1']
            }, {
                'impact': ['e4'],
                'depth': 1L,
                '_id': 'e3',
                'name': 'e3',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e2']
            }, {
                'impact': [],
                'depth': 0L,
                '_id': 'e4',
                'name': 'e4',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e3']
            }
        ])

        res = self.manager.get_graph_depends('e4', deepness=0)['graph']
        res.sort()
        self.assertEqual(res, [{
            'impact': [],
            'depth': 0L,
            '_id': 'e4',
            'name': 'e4',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e3']
        }])

        res = self.manager.get_graph_depends('e4', deepness=1)['graph']
        res.sort()
        self.assertEqual(res, [
            {
                'impact': ['e4'],
                'depth': 1L,
                '_id': 'e3',
                'name': 'e3',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e2']
            }, {
                'impact': [],
                'depth': 0L,
                '_id': 'e4',
                'name': 'e4',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e3']
            }
        ])

        res = self.manager.get_graph_depends('e4', deepness=2)['graph']
        res.sort()
        self.assertEqual(res, [
            {
                'impact': ['e3'],
                'depth': 2L,
                '_id': 'e2',
                'name': 'e2',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e1']
            }, {
                'impact': ['e4'],
                'depth': 1L,
                '_id': 'e3',
                'name': 'e3',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e2']
            }, {
                'impact': [],
                'depth': 0L,
                '_id': 'e4',
                'name': 'e4',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e3']
            }])

        res = self.manager.get_graph_depends('e4', deepness=3)['graph']
        res.sort()
        self.assertEqual(res, [
            {
                'impact': ['e2'],
                'depth': 3L,
                '_id': 'e1',
                'name': 'e1',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': []
            }, {
                'impact': ['e3'],
                'depth': 2L,
                '_id': 'e2',
                'name': 'e2',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e1']
            }, {
                'impact': ['e4'],
                'depth': 1L,
                '_id': 'e3',
                'name': 'e3',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e2']
            }, {
                'impact': [],
                'depth': 0L,
                '_id': 'e4',
                'name': 'e4',
                'infos': {},
                'measurements': [],
                'type': 'component',
                'links': {},
                'depends': ['e3']
            }
        ])

        res = self.manager.get_graph_impact('e1')['graph']
        res.sort()
        self.assertEqual(res, [{
            'impact': ['e2'],
            'depth': 0L,
            '_id': 'e1',
            'name': 'e1',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': []
        }, {
            'impact': ['e3'],
            'depth': 1L, '_id':
            'e2', 'name': 'e2',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e1']
        }, {
            'impact': ['e4'],
            'depth': 2L,
            '_id': 'e3',
            'name': 'e3',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e2']
        }, {
            'impact': [],
            'depth': 3L, '_id':
            'e4', 'name':
            'e4', 'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e3']
        }])

        res = self.manager.get_graph_impact('e1', deepness=0)['graph']
        res.sort()
        self.assertEqual(res, [{
            'impact': ['e2'],
            'depth': 0L,
            '_id': 'e1',
            'name': 'e1',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': []
        }])

        res = self.manager.get_graph_impact('e1', deepness=1)['graph']
        res.sort()
        self.assertEqual(res, [{
            'impact': ['e2'],
            'depth': 0L,
            '_id': 'e1',
            'name': 'e1',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': []
        }, {
            'impact': ['e3'],
            'depth': 1L, '_id':
            'e2', 'name': 'e2',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e1']
        }])

        res = self.manager.get_graph_impact('e1', deepness=2)['graph']
        res.sort()
        self.assertEqual(res, [{
            'impact': ['e2'],
            'depth': 0L,
            '_id': 'e1',
            'name': 'e1',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': []
        }, {
            'impact': ['e3'],
            'depth': 1L, '_id':
            'e2', 'name': 'e2',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e1']
        }, {
            'impact': ['e4'],
            'depth': 2L,
            '_id': 'e3',
            'name': 'e3',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e2']
        }])

        res = self.manager.get_graph_impact('e1', deepness=3)['graph']
        res.sort()
        self.assertEqual(res, [{
            'impact': ['e2'],
            'depth': 0L,
            '_id': 'e1',
            'name': 'e1',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': []
        }, {
            'impact': ['e3'],
            'depth': 1L, '_id':
            'e2', 'name': 'e2',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e1']
        }, {
            'impact': ['e4'],
            'depth': 2L,
            '_id': 'e3',
            'name': 'e3',
            'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e2']
        }, {
            'impact': [],
            'depth': 3L, '_id':
            'e4', 'name':
            'e4', 'infos': {},
            'measurements': [],
            'type': 'component',
            'links': {},
            'depends': ['e3']
        }])


class TestLeavesRequests(BaseTest):

    def setUp(self):
        super(TestLeavesRequests, self).setUp()
        self.manager.collection_name = 'default_testentities'

    def test(self):
        ent1 = create_entity(
            id='e1',
            name='e1',
            etype='component',
            depends=[],
            impact=['e2']
        )
        ent2 = create_entity(
            id='e2',
            name='e2',
            etype='component',
            depends=['e1'],
            impact=['e3']
        )
        ent3 = create_entity(
            id='e3',
            name='e3',
            etype='component',
            depends=['e2'],
            impact=['e4']
        )
        ent4 = create_entity(
            id='e4',
            name='e4',
            etype='component',
            depends=['e3'],
            impact=[]
        )

        self.manager._put_entities(ent1)
        self.manager._put_entities(ent2)
        self.manager._put_entities(ent3)
        self.manager._put_entities(ent4)

        res = self.manager.get_leaves_depends('e4')
        self.assertEqual(res,
                         [{
                             'impact': ['e2'],
                             'name': 'e1',
                             'measurements': [],
                             'depth': 3L,
                             'depends': [],
                             'infos': {},
                             'links': {},
                             '_id': 'e1',
                             'type': 'component'
                         }]
                         )

        res = self.manager.get_leaves_depends('e4', deepness=0)
        self.assertEqual(res,
                         [{
                             'impact': [],
                             'name': 'e4',
                             'measurements': [],
                             'depth': 0L,
                             'depends': ['e3'],
                             'infos': {},
                             'links': {},
                             '_id': 'e4',
                             'type': 'component'
                         }])

        res = self.manager.get_leaves_depends('e4', deepness=1)
        self.assertEqual(res,
                         [{
                             'impact': ['e4'],
                             'name': 'e3',
                             'measurements': [],
                             'depth': 1L,
                             'depends': ['e2'],
                             'infos': {},
                             'links': {},
                             '_id': 'e3',
                             'type': 'component'
                         }])

        res = self.manager.get_leaves_depends('e4', deepness=2)
        self.assertEqual(res,
                         [{
                             'impact': ['e3'],
                             'name': 'e2',
                             'measurements': [],
                             'depth': 2L,
                             'depends': ['e1'],
                             'infos': {},
                             'links': {},
                             '_id': 'e2',
                             'type': 'component'
                         }])

        res = self.manager.get_leaves_depends('e4', deepness=3)
        self.assertEqual(res,
                         [{
                             'impact': ['e2'],
                             'name': 'e1',
                             'measurements': [],
                             'depth': 3L,
                             'depends': [],
                             'infos': {},
                             'links': {},
                             '_id': 'e1',
                             'type': 'component'
                         }]
                         )

        res = self.manager.get_leaves_impact('e1')
        self.assertEqual(res,
                         [{
                             'impact': [],
                             'name': 'e4',
                             'measurements': [],
                             'depth': 3L,
                             'depends': ['e3'],
                             'infos': {},
                             'links': {},
                             '_id': 'e4',
                             'type': 'component'
                         }])


if __name__ == '__main__':
    output = root_path + "/tmp/tests_report"
    unittest.main(
        testRunner=xmlrunner.XMLTestRunner(output=output),
        verbosity=3)
