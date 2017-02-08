#!/usr/bin/env/python
# -*- coding: utf-8 -*-

from unittest import main, TestCase

from canopsis.context_graph.manager import ContextGraph
from canopsis.middleware.core import Middleware

class TestManager(TestCase)

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

    def test_config(self):
        self.assertEqual(self.manager)

if __name__ == '__main__':
    main()
