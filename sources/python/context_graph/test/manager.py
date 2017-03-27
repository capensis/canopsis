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
                '_id': 'conn1',
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

    def trearDown(self):
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
            self.manager.get_entities(query={'_id': 'test_create'})[0]
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
        print(self.manager.get_entities(query={'_id': 'test_entity2'}))
        self.manager.delete_entity('test_entity')
        print(self.manager.get_entities(query={'_id': 'test_entity2'}))
        self.assertEqual(
            [],
            self.manager.get_entities(query={'_id': 'test_entity2'})[0]['depends']
        )

    def test_check_comp(self):
        self.assertEqual(self.manager.check_comp('c1'), True)
        self.assertEqual(self.manager.check_comp('c2'), False)

    def test_check_re(self):
        self.assertEqual(self.manager.check_re('r1/c1'), True)
        self.assertEqual(self.manager.check_re('r2/c1'), False)

    def test_check_conn(self):
        self.assertEqual(self.manager.check_conn('conn1'), True)
        self.assertEqual(self.manager.check_conn('conn2'), False)

    def test_add_comp(self):
        id = "comp1"
        comp = create_comp(id,
                           "comp1_name",
                           depends=["conn1", "conn2", "conn3", "conn4"],
                           impact=["res1", "res2", "res3", "res4"],
                           measurements=["m1", "m2", "m3"],
                           infos={"title": "Foo", "content": "bar"})

        self.manager.add_comp(comp)
        tmp_comp = self.entities_storage.get_elements(ids=id)
        self.assertEqual(comp, tmp_comp)

    def test_add_con(self):
        id = "conn1"
        conn = create_conn(id,
                           "conn1_name",
                           depends=["foo0", "bar0", "foo1", "bar1"],
                           impact=["comp1", "comp2", "comp3", "comp4"],
                           measurements=["m1", "m2", "m3"],
                           infos={"title": "Foo", "content": "bar"})
        self.manager.add_conn(conn)
        tmp_conn = self.entities_storage.get_elements(ids=id)
        self.assertEqual(conn, tmp_conn)

    def test_add_re(self):
        id = "re1"
        re = create_re(id,
                       "re1_name",
                       depends=["foo", "bar", "foo", "bar"],
                       impact=["comp1", "comp2", "comp3", "comp4"],
                       measurements=["m1", "m2", "m3"],
                       infos={"title": "Foo", "content": "bar"})
        self.manager.add_re(re)
        tmp_re = self.entities_storage.get_elements(ids=id)
        self.assertEqual(re, tmp_re)

    def test_manage_comp_to_re_link(self):
        conn1_id = "mcr_conn1"
        re1_id = "mcr_re1"
        re2_id = "mcr_re2"
        re3_id = "mcr_re3"

        conn1 = create_conn(conn1_id, conn1_id, depends=[re1_id, re2_id])
        re1 = create_re(re1_id, re1_id, impact=[conn1_id])
        re2 = create_re(re2_id, re2_id, impact=[conn1_id])
        re3 = create_re(re3_id, re3_id)

        self.manager.add_conn(conn1)
        self.manager.add_re(re1)
        self.manager.add_re(re2)
        self.manager.add_re(re3)

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

        self.manager.add_re(re1)
        self.manager.add_conn(conn1)
        self.manager.add_conn(conn2)
        self.manager.add_conn(conn3)

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

        self.manager.add_comp(comp1)
        self.manager.add_conn(conn1)
        self.manager.add_conn(conn2)
        self.manager.add_conn(conn3)

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

    def test_get_entity(self):
        pass

# Need to be adapted
class GetEvent(TestCase):
    """Test get_event method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context[ContextGraph.ENTITIES_STORAGE].remove_elements()
        self.entity_id = '/a/b/c/d/e/f'

    def tearDown(self):
        self.context[ContextGraph.ENTITIES_STORAGE].remove_elements()

    def test_get_check_event(self):

        entity_id = '/other/a/b/c'

        entity = self.context.get_entity_by_id(entity_id)

        event = self.context.get_event(
            entity, event_type='check', output='test'
        )

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
        error_desc = "Event type should be 'connector', 'resource' or\
            'component' not {0}.".format(self.event["source_type"])
        with self.assertRaisesRegexp(ValueError, error_desc):
            ContextGraph.get_id(self.event)


if __name__ == '__main__':
    main()
