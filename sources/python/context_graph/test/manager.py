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


# Test of deprecated functions
class ContextStorageTest(TestCase):
    """Test access to the context storage.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()

    def tearDown(self):
        pass
        #self.context.remove()

    def test_ctx_storage(self):
        """Test to access to the ctx storage.
        """
        context = self.context.context

        count_per_entity_type = 2

        # let's iterate on context items in order to create entities
        for n in range(1, len(context)):
            sub_context = context[:n]
            entity_context = {c: c for c in sub_context[1:]}
            entity = {}
            for i in range(count_per_entity_type):
                entity[ContextGraph.NAME] = str(i)
                self.context.put(
                    _type=context[n], context=entity_context, entity=entity
                )

        entities = self.context.find()

        self.assertEqual(
            len(entities),
            count_per_entity_type * (len(context) - 1) + len(context) - 2
        )

        for n in range(1, len(context)):
            sub_context = context[:n]
            entity_context = {c: c for c in sub_context[1:]}
            entities = self.context.find(
                _type=context[n], context=entity_context
            )

            self.assertEqual(
                len(entities),
                count_per_entity_type + (1 if n < (len(context) - 1) else 0)
            )

            _id = self.context.get_entity_id(entities[0])
            self.context.remove(ids=_id)
            entities = self.context.find(
                _type=context[n], context=entity_context
            )
            self.assertEqual(
                len(entities),
                count_per_entity_type - (0 if n < (len(context) - 1) else 1)
            )

            self.context.remove(_type=context[n], context=entity_context)
            entities = self.context.find(
                _type=context[n], context=entity_context
            )
            self.assertEqual(len(entities), 0)

    def test_incomplete_hierarchy(self):
        """Test to add elements where parents do not exists.
        """

        # first, ensure no entity exists
        # in constructing a context
        context = {}

        _context_keys = self.context.context[1:]
        # for all key in context.context keys
        for key in _context_keys:
            # check if entity does not exist
            entity = self.context.get(_type=key, names=key, context=context)
            context[key] = key  # update context with key
            self.assertIsNone(entity)
        # ensure entity does not exist
        entity = self.context.get(_type='test', names='test', context=context)
        self.assertIsNone(entity)

        # put entity in DB
        property_key = 'test'
        entity = {ContextGraph.NAME: 'test', property_key: 'test'}
        self.context.put(
            _type='test', entity=entity, context=context
        )

        # this time, check if parent have been putted
        context = {}
        # for all key in context.context keys
        for key in _context_keys:
            # check if entity does not exist
            entity = self.context.get(_type=key, names=key, context=context)
            context[key] = key  # update context with key
            self.assertIsNotNone(entity)
            self.assertNotIn(property_key, entity)

        # check if entity exists and if property key is in entity
        entity = self.context.get(_type='test', names='test', context=context)
        self.assertIsNotNone(entity)
        self.assertIn(property_key, entity)

        self.context.remove()

        del context['resource']
        # put entity in DB
        property_key = 'test'
        entity = {ContextGraph.NAME: 'test', property_key: 'test'}
        self.context.put(
            _type='test', entity=entity, context=context
        )
        # do the same with contex without resource

        # this time, check if parent have been putted
        context = {}
        # for all key in context.context keys
        for key in _context_keys[:-1]:
            # check if entity does not exist
            entity = self.context.get(_type=key, names=key, context=context)
            context[key] = key  # update context with key
            self.assertIsNotNone(entity)
            self.assertNotIn(property_key, entity)

        # check if entity exists and if property key is in entity
        entity = self.context.get(_type='test', names='test', context=context)
        self.assertIsNotNone(entity)
        self.assertIn(property_key, entity)

        # do the same in trying to put a resource
        self.context.remove()

        # put entity in DB
        property_key = 'test'
        entity = {ContextGraph.NAME: 'resource', property_key: 'test'}
        self.context.put(
            _type='resource', entity=entity, context=context
        )

        # this time, check if parent have been putted
        context = {}
        # for all key in context.context keys
        for key in _context_keys[:-1]:
            # check if entity does not exist
            entity = self.context.get(_type=key, names=key, context=context)
            self.assertIsNotNone(entity)
            self.assertNotIn(property_key, entity)
            context[key] = key  # update context with key

        # check if entity exists and if property key is in entity
        entity = self.context.get(
            _type='resource', names='resource', context=context
        )
        self.assertIsNotNone(entity)
        self.assertIn(property_key, entity)

        # do the same in trying to put a component
        self.context.remove()

        del context['component']
        # put entity in DB
        property_key = 'test'
        entity = {ContextGraph.NAME: 'component', property_key: 'test'}
        self.context.put(
            _type='component', entity=entity, context=context
        )

        # this time, check if parent have been putted
        context = {}
        # for all key in context.context keys
        for key in _context_keys[:-2]:
            # check if entity does not exist
            entity = self.context.get(_type=key, names=key, context=context)
            context[key] = key  # update context with key
            self.assertIsNotNone(entity)
            self.assertNotIn(property_key, entity)

        # check if entity exists and if property key is in entity
        entity = self.context.get(
            _type='component', names='component', context=context
        )
        self.assertIsNotNone(entity)
        self.assertIn(property_key, entity)


class GetEntityTest(TestCase):
    """Test get_entity method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()

    def tearDown(self):
        self.context.remove()

    def _assert_entity_id(self, event, entity):

        _entity = self.context.get_entity(event)
        self.assertEqual(entity, _entity)

    def test_noname_notype(self):
        """Test to get an entity from an event without both name and type.
        """

        event = {
            'event_type': 'test',
            'source_type': 'component',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'example'
        }
        entity = {
            'type': 'test',
            'connector': 'c',
            'connector_name': 'cn',
            ContextGraph.NAME: 'example'
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_connector(self):
        """Test get_entity_id from a connector.
        """

        event = {
            'event_type': 'test',
            'source_type': 'connector',
            'connector': 'c',
        }
        entity = {
            'type': 'test',
            ContextGraph.NAME: 'c',
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_connector_name(self):
        """Test get_entity_id from a connector_name.
        """

        # assert to get a connector name id
        event = {
            'event_type': 'test',
            'source_type': 'connector_name',
            'connector': 'c',
            'connector_name': 'cn'
        }
        entity = {
            'type': 'test',
            'connector': 'c',
            ContextGraph.NAME: 'cn',
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_component(self):
        """Test get_entity_id from a component.
        """

        event = {
            'event_type': 'test',
            'source_type': 'component',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k'
        }
        entity = {
            'type': 'test',
            'connector': 'c',
            'connector_name': 'cn',
            ContextGraph.NAME: 'k',
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_resource(self):
        """Test get_entity_id from a resource.
        """

        event = {
            'event_type': 'test',
            'source_type': 'resource',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            'resource': 'r'
        }
        entity = {
            'type': 'test',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'r',
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_other_component(self):
        """Test get_entity_id from a component other data.
        """

        event = {
            'event_type': 'test',
            'source_type': 'component',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k'
        }
        entity = {
            'type': 'test',
            'connector': 'c',
            'connector_name': 'cn',
            ContextGraph.NAME: 'k'
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_other_resource(self):
        """Test get_entity_id from a resource other data.
        """

        event = {
            'event_type': 'test',
            'source_type': 'resource',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            'resource': 'r'
        }
        entity = {
            'type': 'test',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'r',
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_other_check_type(self):
        """Test get_entity_id from a resource other data.
        """

        event = {
            'event_type': 'check',
            'source_type': 'resource',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            'resource': 'r'
        }
        entity = {
            'type': 'resource',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'r',
        }
        self._assert_entity_id(event, entity)

    def test_get_entity_id_other_other_check_type(self):
        """Test get_entity_id from a resource other data.
        """

        event = {
            'event_type': 'check',
            'source_type': 'other',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            'resource': 'r'
        }
        entity = {
            'type': 'other',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'r'
        }
        self._assert_entity_id(event, entity)


class ContextEntityTest(TestCase):
    """Test for playing with entities.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()

    def tearDown(self):
        self.context.remove()

    def _assert_entity_id(self, entity, entity_id):

        _entity_id = self.context.get_entity_id(entity)
        self.assertEqual(_entity_id, entity_id)

    def test_get_entity_id_connector(self):
        """Test get_entity_id from a connector.
        """

        # assert to get a connector id
        entity = {
            'type': 'connector',
            ContextGraph.NAME: 'c'
        }
        self._assert_entity_id(entity, '/connector/c')

    def test_get_entity_id_connector_name(self):
        """Test get_entity_id from a connector_name.
        """

        # assert to get a connector name id
        entity = {
            'type': 'connector_name',
            'connector': 'c',
            ContextGraph.NAME: 'cn'
        }
        self._assert_entity_id(entity, '/connector_name/c/cn')

    def test_get_entity_id_component(self):
        """Test get_entity_id from a component.
        """

        entity = {
            'type': 'component',
            'connector': 'c',
            'connector_name': 'cn',
            ContextGraph.NAME: 'k'
        }
        self._assert_entity_id(entity, '/component/c/cn/k')

    def test_get_entity_id_resource(self):
        """Test get_entity_id from a resource.
        """

        entity = {
            'type': 'resource',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'r'
        }
        self._assert_entity_id(entity, '/resource/c/cn/k/r')

    def test_get_entity_id_other_component(self):
        """Test get_entity_id from a component other data.
        """

        entity = {
            'type': 'other',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'o'
        }
        self._assert_entity_id(entity, '/other/c/cn/k/o')

    def test_get_entity_id_other_resource(self):
        """Test get_entity_id from a resource other data.
        """

        entity = {
            'type': 'other',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            'resource': 'r',
            'other': 'o',
            ContextGraph.NAME: 'o'
        }
        self._assert_entity_id(entity, '/other/c/cn/k/r/o')


class EntityIdContextTest(TestCase):
    """Test get_entity_id_context_name method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()

    def tearDown(self):
        self.context.remove()

    def _assert_entity_id(self, entity, entity_id, path, data_id):

        _entity_id, _path, _data_id = self.context.get_entity_id_context_name(
            entity
        )
        self.assertEqual(entity_id, _entity_id)
        self.assertEqual(path, _path)
        self.assertEqual(data_id, _data_id)

    def _get_path(self, entity):
        """Get entity path.
        """

        result = {}

        for ctx in self.context.context:
            if ctx in entity and entity[ContextGraph.NAME] != entity[ctx]:
                result[ctx] = entity[ctx]

        return result

    def test_get_entity_id_connector(self):
        """Test get_entity_id from a connector.
        """

        # assert to get a connector id
        entity = {
            'type': 'connector',
            ContextGraph.NAME: 'c'
        }
        path = self._get_path(entity)
        self._assert_entity_id(entity, '/connector/c', path, 'c')

    def test_get_entity_id_connector_name(self):
        """Test get_entity_id from a connector_name.
        """

        # assert to get a connector name id
        entity = {
            'type': 'connector_name',
            'connector': 'c',
            ContextGraph.NAME: 'cn'
        }
        path = self._get_path(entity)
        self._assert_entity_id(entity, '/connector_name/c/cn', path, 'cn')

    def test_get_entity_id_component(self):
        """Test get_entity_id from a component.
        """

        entity = {
            'type': 'component',
            'connector': 'c',
            'connector_name': 'cn',
            ContextGraph.NAME: 'k'
        }
        path = self._get_path(entity)
        self._assert_entity_id(entity, '/component/c/cn/k', path, 'k')

    def test_get_entity_id_resource(self):
        """Test get_entity_id from a resource.
        """

        entity = {
            'type': 'resource',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'r'
        }
        path = self._get_path(entity)
        self._assert_entity_id(entity, '/resource/c/cn/k/r', path, 'r')

    def test_get_entity_id_other_component(self):
        """Test get_entity_id from a component other data.
        """

        entity = {
            'type': 'other',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            ContextGraph.NAME: 'o'
        }
        path = self._get_path(entity)
        self._assert_entity_id(entity, '/other/c/cn/k/o', path, 'o')

    def test_get_entity_id_other_resource(self):
        """Test get_entity_id from a resource other data.
        """

        entity = {
            'type': 'other',
            'connector': 'c',
            'connector_name': 'cn',
            'component': 'k',
            'resource': 'r',
            ContextGraph.NAME: 'o'
        }
        path = self._get_path(entity)
        self._assert_entity_id(entity, '/other/c/cn/k/r/o', path, 'o')


class GetNameTest(TestCase):
    """Test get_name method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()
        self.entity_id = '/a/b/c/d/e/f'

    def tearDown(self):
        self.context.remove()

    def _assert_name(self, _type, result, entity_id=None):
        """Assert get_name(self.entity_id) result with input result.

        :param str _type: get_name parameter _type.
        :param str result: value to compare with get_name result.
        :param str entity_id: first get_name parameter. If None, use
            self.entity_id.
        """

        if entity_id is None:
            entity_id = self.entity_id

        name = self.context.get_name(entity_id, _type=_type)

        self.assertEqual(name, result)

    def test_type_none(self):
        """Test with _type is None.
        """

        self._assert_name(_type=None, result='f')

    def test_type(self):
        """Test with _type is type.
        """

        self._assert_name(_type='type', result='a')

    def test_connector(self):
        """Test with _type is connector.
        """

        self._assert_name(_type='connector', result='b')

    def test_connector_name(self):
        """Test with _type is connector_name.
        """

        self._assert_name(_type='connector_name', result='c')

    def test_component(self):
        """Test with _type is component.
        """

        self._assert_name(_type='component', result='d')

    def test_resource(self):
        """Test with _type is resource.
        """

        self._assert_name(_type='resource', result='e')

    def test_other(self):
        """Test with _type is other.
        """

        self._assert_name(_type='other', result='f')

    # TODO 4-01-2017
    # def test_error(self):
    #    """Test with _type is not in entity_id.
    #    """

    #    self._assert_name(_type='resource', result=None, entity_id='/a/b')


class GetEntityByIdTest(TestCase):
    """Test get_entity_by_id method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()
        self.entity_id = '/a/b/c/d/e/f'

    def tearDown(self):
        self.context.remove()

    def _assert_entity(self, _id, entity):
        """Assert get_entity_by_id() result with input result.

        :param str _id: value to use such as parameter of get_entity_by_id.
        :param str entity: value to compare with get_entity_by_id result.
        """

        entity_to_compare = self.context.get_entity_by_id(_id)

        self.assertEqual(entity_to_compare, entity)

    # TODO 4-01-2017
    # def test_entity_empty(self):
    #    """Test with empty name entity.
    #    """

    #    self._assert_entity(_id='', entity={})

    def test_connector(self):
        """Test with connector.
        """

        self._assert_entity(
            _id='/connector/a', entity={'name': 'a', 'type': 'connector'}
        )

    def test_connector_name(self):
        """Test with connector_name.
        """

        self._assert_entity(
            _id='/connector_name/a/b',
            entity={'name': 'b', 'connector': 'a', 'type': 'connector_name'}
        )

    def test_component(self):
        """Test with component.
        """

        self._assert_entity(
            _id='/component/a/b/c',
            entity={
                'name': 'c',
                'connector': 'a',
                'type': 'component',
                'connector_name': 'b'
            }
        )

    def test_resource(self):
        """Test with resource.
        """

        self._assert_entity(
            _id='/resource/a/b/c/d',
            entity={
                'name': 'd',
                'connector': 'a',
                'connector_name': 'b',
                'component': 'c',
                'type': 'resource'
            }
        )

    def test_component_other(self):
        """Test with other.
        """

        self._assert_entity(
            _id='/other/a/b/c/d',
            entity={
                'name': 'd',
                'connector': 'a',
                'connector_name': 'b',
                'component': 'c',
                'type': 'other'
            }
        )

    def test_resource_error(self):
        """Test with _type is not in entity_id.
        """

        self._assert_entity(
            _id='/other/a/b/c/d/e',
            entity={
                'name': 'e',
                'connector': 'a',
                'connector_name': 'b',
                'component': 'c',
                'resource': 'd',
                'type': 'other'
            }
        )


class GetEvent(TestCase):
    """Test get_event method.
    """

    def setUp(self):
        self.context = ContextGraph(data_scope='test_context')
        self.context.remove()
        self.entity_id = '/a/b/c/d/e/f'

    def tearDown(self):
        self.context.remove()

    def test_get_check_event(self):

        entity_id = '/other/a/b/c'

        entity = self.context.get_entity_by_id(entity_id)

        event = self.context.get_event(
            entity, event_type='check', output='test'
        )

        self.assertEqual(event['event_type'], 'check')
        self.assertEqual(event['output'], 'test')


if __name__ == '__main__':
    main()
