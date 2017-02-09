#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware


def create_entity(id, name, etype, depends=[], impact=[], measurements=[], infos={}):
    return {'_id': id,
            'type': etype,
            'name': name,
            'depends': depends,
            'impact': impact,
            'measurements': measurements,
            'infos': infos
            }


def create_comp(id, name, depends=[], impact=[], measurements=[], infos={}):
    return create_entity(id,
                         name,
                         "component",
                         depends,
                         impact,
                         measurements,
                         infos)


def create_conn(id, name, depends=[], impact=[], measurements=[], infos={}):
    return create_entity(id,
                         name,
                         "connector",
                         depends,
                         impact,
                         measurements,
                         infos)


def create_re(id, name, depends=[], impact=[], measurements=[], infos={}):
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


if __name__ == '__main__':
    main()
