#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware


def create_entity(
        id,
        name,
        etype,
        depends=[],
        impact=[],
        measurements=[],
        infos={}):
    return {'_id': id,
            'type': etype,
            'name': name,
            'depends': depends,
            'impact': impact,
            'measurements': measurements,
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
        self.manager = ContextGraph()
        self.entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities://'
        )
        self.organisations_storage = Middleware.get_middleware_by_uri(
            'storage-default-testorganisations://'
        )
        self.users_storage = Middleware.get_middleware_by_uri(
            'storage-default-testusers://'
        )

        self.manager[ContextGraph.ENTITIES_STORAGE] = self.entities_storage
        self.manager[
            ContextGraph.ORGANISATIONS_STORAGE] = self.organisations_storage
        self.manager[ContextGraph.USERS_STORAGE] = self.users_storage

        self.template = {'_id': None,
                         'type': 'connector',
                         'name': 'conn-name1',
                         'depends': [],
                         'impact': [],
                         'measurements': [],
                         'infos': {}}

    def tearDown(self):
        self.entities_storage.remove_elements()
        self.organisations_storage.remove_elements()
        self.organisations_storage.remove_elements()

    def assertEqualEntities(self, entity1, entity2):
        sorted(entity1["depends"])
        sorted(entity1["impact"])
        sorted(entity2["depends"])
        sorted(entity2["impact"])
        self.assertDictEqual(entity1, entity2)

class GetEvent(TestCase):
    """Test get_event method.
    """

    def setUp(self):
        super(GetEvent, self).setUp()
        self.context = ContextGraph(data_scope='test_context')
        self.context[ContextGraph.ENTITIES_STORAGE].remove_elements()

    def test_get_check_event(self):

        # set of required field
        fields = sorted(['component', 'connector', 'connector_name',
                         'event_type', 'infos', 'long_output', 'measurements',
                         'output', 'source_type', 'state', 'state_type',
                         'timestamp'])

        conn = "conn"
        conn_name = "conn-name"
        entity_id = conn + "/" + conn_name
        measurements = ["measurements1", "measurements2"]
        infos = {"info1": "data1",
                 "info2": "data2",
                 "info3": "data3"}

        entity = create_conn(entity_id, conn_name, depends=[],
                             impact=[], measurements=measurements, infos=infos)

        event = self.context.get_event(
            entity, event_type='check', output='test'
        )

        self.assertListEqual(sorted(event.keys()), fields)
        self.assertEqual(event['event_type'], 'check')
        self.assertEqual(event['output'], 'test')
        self.assertListEqual(event['measurements'], measurements)
        self.assertDictEqual(event['infos'], infos)


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
        error_desc = "Event type should be 'connector', 'resource' or\
            'component' not {0}.".format(self.event["source_type"])
        with self.assertRaisesRegexp(ValueError, error_desc):
            ContextGraph.get_id(self.event)


class GetEntitiesByID(BaseTest):

    def test_get_entity_by_id_id(self):
        entity = self.template.copy()
        entity['_id'] = 'conn1/conn-name'
        entity['type'] = 'connector'
        entity['name'] = 'conn-name'

        self.manager.put_entities(entity)

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
        entities[1]['type'] ='connector'
        entities[1]['name'] = 'conn-name2'

        entities[2]['_id'] = 'conn3/conn-name3'
        entities[2]['type'] = 'connector'
        entities[2]['name'] = 'conn-name3'

        sorted(entities)
        self.manager.put_entities(entities)

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

        self.manager.put_entities(entity)

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
        entities[1]['type'] ='connector'
        entities[1]['name'] = 'conn-name2'

        entities[2]['_id'] = 'conn3/conn-name3'
        entities[2]['type'] = 'connector'
        entities[2]['name'] = 'conn-name3'

        sorted(entities)

        self.manager.put_entities(entities)

        ids = [x["_id"] for x in entities]

        result = self.manager.get_entities_by_id(ids)
        sorted(result)
        self.assertListEqual(result, entities)


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

        self.manager.put_entities(entities)
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

        self.manager.put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

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

    def test_update_entity_update_impact_missing_ids(self):
        not_id = "no_an_entity_id"
        self.ent1["depends"] = [not_id]
        desc = "Could not find some entity in database."
        with self.assertRaisesRegexp(ValueError, desc):
            self.manager.update_entity(self.ent1)

        entity = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.assertEqualEntities(self.ent2, entity)

        entity = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.assertEqualEntities(self.ent3, entity)

        entity = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(self.ent4, entity)

    def test_update_entity_update_depends_missing_ids(self):
        not_id = "no_an_entity_id"
        self.ent1["depends"] = [not_id]

        desc = "Could not find some entity in database."
        with self.assertRaisesRegexp(ValueError, desc):
            self.manager.update_entity(self.ent1)

        entity = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.assertEqualEntities(self.ent2, entity)

        entity = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.assertEqualEntities(self.ent3, entity)

        entity = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(self.ent4, entity)


class UpdateDependancies(BaseTest):

    def test_update_dependancies_wrong_type(self):

        desc = "Dependancy_type should be depends or impact not {0}."

        with self.assertRaisesRegexp(ValueError, desc.format(None)):
            self.manager._ContextGraph__update_dependancies(None, None, None)

    def test_update_dependancies_entity_not_found_in_db(self):

        status = {"deletions": ["ent2", "ent3"]}

        desc = "Could not find some entity in database."
        with self.assertRaisesRegexp(ValueError, desc):
            self.manager._ContextGraph__update_dependancies(None, status,
                                                            "impact")

        status = {"insertions": ["ent2", "ent3"],
                  "deletions": []}
        with self.assertRaisesRegexp(ValueError, desc):
            self.manager._ContextGraph__update_dependancies(None, status,
                                                            "impact")
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

        self.manager.put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

        if delete:
            status = {"deletions": ["ent2", "ent3"],
                      "insertions":[]}
        else:
            status = {"deletions": [],
                      "insertions":["ent2", "ent3"]}

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

        self.manager.put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

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

    def test_create_entity_entity_exists(self):
        entity = self.template.copy()
        entity['_id'] = "I am here"
        entity['type'] = 'connector'
        entity['name'] = 'conn-name1'

        self.manager.create_entity(entity)

        desc = "An entity  id {0} already exist".format(entity["_id"])
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

        self.manager.put_entities(self.ent2)
        self.manager.put_entities(self.ent3)
        self.manager.put_entities(self.ent4)

        self.manager.create_entity(self.ent1)

        result = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(result, self.ent1)

        result = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.ent2["impact"].append(self.ent1["_id"])
        self.assertEqualEntities(result, self.ent2)

        result = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.ent3["impact"].append(self.ent1["_id"])
        self.assertEqualEntities(result, self.ent3)

        result = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(result, self.ent4)

    def test_create_entity_entity_depends(self):
        self.__test("depends", "impact")

    def test_create_entity_entity_impact(self):
        self.__test("depends", "impact")


if __name__ == '__main__':
    main()
