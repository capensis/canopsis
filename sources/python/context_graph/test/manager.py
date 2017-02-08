#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware

class TestManager(TestCase):

    def setUp(self):
        self.manager = ContextGraph()
        self.entities_storage = Middleware.get_middleware_by_uri(
            'storage-default-testentities'
        )
        self.organisations_storage = Middleware.get_middleware_by_uri(
            'storage-default-testorganisations'
        )
        self.users_storage = Middleware.get_middleware_by_uri(
            'storage-default-testusers'
        )

        self.manager[ContextGraph.ENTITIES_STORAGE] = self.entities_storage
        self.manager[ContextGraph.ORGANISATIONS_STORAGE] = self.organisations_storage
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
            self.assertEqual(self.manager.chek_comp('c1'), True)
            self.assertEqual(self.manager.chek_comp('c2'), False)

        def test_check_re(self):
            self.assertEqual(self.manager.check_re('r1/c1'), True)
            self.assertEqual(self.manager.check_re('r2/c1'), False)

        def test_check_conn(self):
            self.assertEqual(self.manager.check_conn('conn1'), True)
            self.assertEqual(self.manager.check_conn('conn2'), False)

        def test_check_links(self):
            pass


if __name__ == '__main__':
    main()
