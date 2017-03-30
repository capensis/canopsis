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

    def tearDown(self):
        self.entities_storage.remove_elements()
        self.organisations_storage.remove_elements()
        self.organisations_storage.remove_elements()


class TestManager(TestCase):

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

        self.entities_storage.put_element(
            element={
                '_id': 'r1/c1',
                'type': 'resource',
                'name': 'r1',
                'depends': [],
                'impact': [],
                'measurements': [],
                'infos': {}
            }
        )

        self.entities_storage.put_element(
            element={
                '_id': 'c1',
                'type': 'component',
                'name': 'c1',
                'depends': [],
                'impact': [],
                'measurements': [],
                'infos': {}
            }
        )

        self.entities_storage.put_element(
            element={
                '_id': 'conn1/conn1',
                'type': 'connector',
                'name': 'conn1',
                'depends': [],
                'impact': [],
                'measurements': [],
                'infos': {}
            }
        )

        self.entities_storage.put_element(
            element={
                '_id': 'r1/c1',
                'type': 'resource',
                'name': 'r1',
                'depends': [],
                'impact': [],
                'measurements': [],
                'infos': {}
            }
        )

        self.organisations_storage.put_element(
            element={
                '_id': 'org-1',
                'name': 'org-1',
                'parents': '',
                'children': [],
                'views': [],
                'users': []
            }
        )

        self.users_storage.put_element(
            element={
                '_id': 'user_1',
                'name': 'jean',
                'org': 'org-1',
                'access_to':  ['org-1']
            }
        )

    def tearDown(self):
        self.entities_storage.remove_elements()
        self.organisations_storage.remove_elements()
        self.organisations_storage.remove_elements()

    def test_create_entity(self):
        entity = {
            '_id': 'test_create',
            'name': 'test_create',
            'depends': [],
            'impact': [],
            'infos': {}
        }
        entity2 = {
            '_id': 'test_create2',
            'name': 'test_create2',
            'depends': [],
            'impact': [],
            'infos': {}
        }
        self.manager.create_entity(entity)
        self.manager.create_entity(entity2)
        self.assertDictEqual(
            entity,
            self.manager.get_entities(query={'_id': 'test_create'}[0])
        )

    def test_delete_entity(self):
        entity = {
            '_id': 'test_entity',
            'name': 'test_entity',
            'depends': [],
            'impact': ['test_entity2'],
            'infos': {}
        }
        entity2 = {
            '_id': 'test_entity2',
            'name': 'test_entity2',
            'depends': ['test_entity'],
            'impact': [],
            'infos': {}
        }
        self.manager.create_entity(entity)
        self.manager.create_entity(entity2)
        self.manager.delete_entity('test_entity')
        self.assertEqual(
            [],
            self.manager.get_entities(query={'_id': 'test_entity2'})[
                0]['depends']
        )

    def test_manage_comp_to_re_link(self):
        conn1_id = "mcr_conn1"
        re1_id = "mcr_re1"
        re2_id = "mcr_re2"
        re3_id = "mcr_re3"

        conn1 = create_conn(conn1_id, conn1_id, depends=[re1_id, re2_id])
        re1 = create_re(re1_id, re1_id, impact=[conn1_id])
        re2 = create_re(re2_id, re2_id, impact=[conn1_id])
        re3 = create_re(re3_id, re3_id)

        self.manager.put_entities(conn1)
        self.manager.put_entities(re1)
        self.manager.put_entities(re2)
        self.manager.put_entities(re3)

        self.manager.manage_comp_to_re_link(re3_id, conn1_id)

        doc = self.entities_storage.get_elements(ids=conn1_id)
        # FIXME : did we receive more than one element ?
        self.assertIn(re1_id, doc["depends"])
        self.assertIn(re2_id, doc["depends"])
        self.assertIn(re3_id, doc["depends"])

    def test_manage_re_to_conn_link(self):
        re1_id = "mrc_re1"
        conn1_id = "mrc_conn1"
        conn2_id = "mrc_conn2"
        conn3_id = "mrc_conn3"

        re1 = create_re(re1_id, re1_id, depends=[conn1_id, conn2_id])
        conn1 = create_conn(conn1_id, conn1_id, impact=[re1_id])
        conn2 = create_conn(conn2_id, conn2_id, impact=[re1_id])
        conn3 = create_conn(conn3_id, conn3_id)

        self.manager.put_entities(re1)
        self.manager.put_entities(conn1)
        self.manager.put_entities(conn2)
        self.manager.put_entities(conn3)

        self.manager.manage_comp_to_re_link(conn3_id, re1_id)

        doc = self.entities_storage.get_elements(ids=re1_id)
        # FIXME : did we receive more than one element ?
        self.assertIn(conn1_id, doc["depends"])
        self.assertIn(conn2_id, doc["depends"])
        self.assertIn(conn3_id, doc["depends"])

    def test_manage_comp_to_conn_link(self):
        comp1_id = "mrc_comp1"
        conn1_id = "mrc_conn1"
        conn2_id = "mrc_conn2"
        conn3_id = "mrc_conn3"

        comp1 = create_comp(comp1_id, comp1_id, depends=[conn1_id, conn2_id])
        conn1 = create_conn(conn1_id, conn1_id, impact=[comp1_id])
        conn2 = create_conn(conn2_id, conn2_id, impact=[comp1_id])
        conn3 = create_conn(conn3_id, conn3_id)

        self.manager.put_entities(comp1)
        self.manager.put_entities(conn1)
        self.manager.put_entities(conn2)
        self.manager.put_entities(conn3)

        self.manager.manage_comp_to_re_link(conn3_id, comp1_id)

        doc = self.entities_storage.get_elements(ids=comp1_id)
        # FIXME : did we receive more than one element ?
        self.assertIn(conn1_id, doc["depends"])
        self.assertIn(conn2_id, doc["depends"])
        self.assertIn(conn3_id, doc["depends"])

    def test_check_links(self):
        self.assertRaises(

            NotImplementedError,
            self.manager.check_links,
            None,
            None,
            None)


class GetEvent(TestCase):
    """Test get_event method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context[ContextGraph.ENTITIES_STORAGE].remove_elements()

    def tearDown(self):
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


class CheckComp(BaseTest):

    def setUp(self):
        super(CheckComp, self).setUp()
        entity = {'_id': 'c1',
                  'type': 'component',
                  'name': 'c1',
                  'depends': [],
                  'impact': [],
                  'measurements': [],
                  'infos': {}}
        self.manager.put_entities(entity)

    def test_check_comp(self):
        self.assertEqual(self.manager.check_comp('c1'), True)
        self.assertEqual(self.manager.check_comp('c2'), False)


class CheckRe(BaseTest):

    def setUp(self):
        super(CheckRe, self).setUp()
        entity = {'_id': 'r1/c1',
                  'type': 'resource',
                  'name': 'r1',
                  'depends': [],
                  'impact': [],
                  'measurements': [],
                  'infos': {}}
        self.manager.put_entities(entity)

    def test_check_re(self):
        self.assertEqual(self.manager.check_re('r1/c1'), True)
        self.assertEqual(self.manager.check_re('r2/c1'), False)


class CheckConn(BaseTest):

    def setUp(self):
        super(CheckConn, self).setUp()
        entity = {'_id': 'conn1/conn-name',
                  'type': 'connector',
                  'name': 'conn-name',
                  'depends': [],
                  'impact': [],
                  'measurements': [],
                  'infos': {}}
        self.manager.put_entities(entity)

    def test_check_conn(self):
        self.assertEqual(self.manager.check_conn('conn1/conn-name'), True)
        self.assertEqual(self.manager.check_conn('conn2'), False)


class GetEntitiesByID(BaseTest):

    def test_get_entity_by_id_id(self):
        entity = {'_id': 'conn1/conn-name',
                  'type': 'connector',
                  'name': 'conn-name',
                  'depends': [],
                  'impact': [],
                  'measurements': [],
                  'infos': {}}

        self.manager.put_entities(entity)

        result = self.manager.get_entities_by_id(entity["_id"])[0]
        self.assertIsInstance(result, type({}))
        self.assertDictEqual(result, entity)

    def test_get_entity_by_id_ids(self):
        entities = [{'_id': 'conn1/conn-name1',
                     'type': 'connector',
                     'name': 'conn-name1',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn2/conn-name2',
                     'type': 'connector',
                     'name': 'conn-name2',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn3/conn-name3',
                     'type': 'connector',
                     'name': 'conn-name3',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}}]

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

        result = self.manager.get_entities_by_id("id")[0]

        sorted(result)
        self.assertIsInstance(result, type({}))
        self.assertDictEqual(result, {})


class PutEntities(BaseTest):

    def test_put_entities_entity(self):
        entity = {'_id': 'conn1/conn-name',
                  'type': 'connector',
                  'name': 'conn-name',
                  'depends': [],
                  'impact': [],
                  'measurements': [],
                  'infos': {}}

        self.manager.put_entities(entity)

        result = self.manager.get_entities_by_id(entity["_id"])[0]

        self.assertDictEqual(result, entity)

    def test_put_entities_entities(self):
        entities = [{'_id': 'conn1/conn-name1',
                     'type': 'connector',
                     'name': 'conn-name1',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn2/conn-name2',
                     'type': 'connector',
                     'name': 'conn-name2',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn3/conn-name3',
                     'type': 'connector',
                     'name': 'conn-name3',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}}]
        sorted(entities)

        self.manager.put_entities(entities)

        ids = [x["_id"] for x in entities]

        result = self.manager.get_entities_by_id(ids)
        sorted(result)
        self.assertListEqual(result, entities)


class GetAllEntitiesId(BaseTest):

    def test_get_all_entities_id(self):
        entities = [{'_id': 'conn1/conn-name1',
                     'type': 'connector',
                     'name': 'conn-name1',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn2/conn-name2',
                     'type': 'connector',
                     'name': 'conn-name2',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn3/conn-name3',
                     'type': 'connector',
                     'name': 'conn-name3',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}},
                    {'_id': 'conn4/conn-name4',
                     'type': 'connector',
                     'name': 'conn-name3',
                     'depends': [],
                     'impact': [],
                     'measurements': [],
                     'infos': {}}]

        expected = set()
        for entity in entities:
            expected.add(entity["_id"])

        result = self.manager.get_all_entities()
        self.assertSetEqual(result, expected)


class CheckLinks(BaseTest):

    def test_check_links(self):
        with self.assertRaises(NotImplementedError):
            self.manager.check_links(None, None, None)


class UpdateEntity(BaseTest):

    def setUp(self):
        super(UpdateEntity, self).setUp()
        template = {'_id': None,
                    'type': 'connector',
                    'name': 'conn-name1',
                    'depends': [],
                    'impact': [],
                    'measurements': [],
                    'infos': {}}
        self.ent1 = template.copy()
        self.ent2 = template.copy()
        self.ent3 = template.copy()
        self.ent4 = template.copy()

        self.ent1["_id"] = "ent1"
        self.ent2["_id"] = "ent2"
        self.ent3["_id"] = "ent3"
        self.ent4["_id"] = "ent4"

        self.manager.put_entities([self.ent1, self.ent2, self.ent3, self.ent4])

    def assertEqualEntities(self, entity1, entity2):
        sorted(entity1["depends"])
        sorted(entity1["impact"])
        sorted(entity2["depends"])
        sorted(entity2["impact"])
        self.assertDictEqual(entity1, entity2)

    def test_update_entity_wrong_id(self):
        fake_entity = {"_id": "wrong id"}
        error_desc = "The _id {0} does not match any entity in database."\
                     .format(fake_entity["_id"])
        with self.assertRaisesRegexp(ValueError, error_desc):
            self.manager.update_entity(fake_entity)

    def test_update_entity_update_depends_insert_multiple(self):
        self.ent1["depends"] = ["ent2", "ent3"]
        self.manager.update_entity(self.ent1)

        entity = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(self.ent1, entity)

        entity = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.ent2["impact"] = ["ent1"]
        self.assertEqualEntities(self.ent2, entity)

        entity = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.ent3["impact"] = ["ent1"]
        self.assertEqualEntities(self.ent3, entity)

        entity = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(self.ent4, entity)

    def test_update_entity_update_depends_insert_single(self):
        self.ent1["depends"] = ["ent2"]
        self.manager.update_entity(self.ent1)

        entity = self.manager.get_entities_by_id(self.ent1["_id"])[0]
        self.assertEqualEntities(self.ent1, entity)

        entity = self.manager.get_entities_by_id(self.ent2["_id"])[0]
        self.ent2["impact"] = ["ent1"]
        self.assertEqualEntities(self.ent2, entity)

        entity = self.manager.get_entities_by_id(self.ent3["_id"])[0]
        self.assertEqualEntities(self.ent3, entity)

        entity = self.manager.get_entities_by_id(self.ent4["_id"])[0]
        self.assertEqualEntities(self.ent4, entity)

    def test_update_entity_update_impact_insert_multiple(self):
        pass

    def test_update_entity_update_impact_insert_single(self):
        pass

    def test_update_entity_update_impact_missing_ids(self):
        pass

    def test_update_entity_update_depends_missing_ids(self):
        pass


if __name__ == '__main__':
    main()
