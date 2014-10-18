#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from unittest import TestCase, main

from canopsis.common.ws import route_name, route, response


class RouteNameTest(TestCase):

    def test_no_param(self):

        result = route_name('d')

        self.assertEqual(result, '/d')

    def test_param(self):

        result = route_name('d', 'a', 'b')

        self.assertEqual(result, '/d/:a/:b')


class RouteTest(TestCase):

    def op(self, route):

        if not hasattr(self, 'routes'):
            self.routes = []

        self.routes.append(route)

        return lambda function: None

    def assertRoutes(self, routes):

        self.assertEqual(self.routes, routes)

    def test_empty_params(self):

        @route(op=self.op)
        def a():
            pass

        self.assertRoutes(['/a'])

    def test_name(self):

        @route(op=self.op, name='b')
        def a():
            pass

        self.assertRoutes(['/b'])

    def test_header_param(self):

        @route(op=self.op)
        def a(a):
            pass

        self.assertRoutes(['/a/:a'])

    def test_header_params(self):

        @route(op=self.op)
        def a(a, b):
            pass

        self.assertRoutes(['/a/:a/:b'])

    def test_optional_param(self):

        @route(op=self.op)
        def a(a=None):
            pass

        self.assertRoutes(['/a', '/a/:a'])

    def test_optional_params(self):

        @route(op=self.op)
        def a(a=None, b=None):
            pass

        self.assertRoutes(['/a', '/a/:a', '/a/:a/:b'])

    def test_required_optional_param(self):

        @route(op=self.op)
        def a(a, b=None):
            pass

        self.assertRoutes(['/a/:a', '/a/:a/:b'])

    def test_required_optional_params(self):

        @route(op=self.op)
        def a(a, b, c=None, d=None):
            pass

        self.assertRoutes(['/a/:a/:b', '/a/:a/:b/:c', '/a/:a/:b/:c/:d'])

    def test_body_param(self):

        @route(op=self.op, payload='a')
        def a(a):
            pass

        self.assertRoutes(['/a'])

    def test_body_default_param(self):

        @route(op=self.op, payload='a')
        def a(a=None):
            pass

        self.assertRoutes(['/a'])

    def test_payload(self):

        @route(op=self.op, payload=['a', 'b'])
        def a(a, b=None):
            pass

        self.assertRoutes(['/a'])

    def test_required_payload(self):

        @route(op=self.op, payload=['a', 'c'])
        def a(a, b, c=None):
            pass

        self.assertRoutes(['/a/:b'])

    def test_optional_payload(self):

        @route(op=self.op, payload=['a', 'c'])
        def a(a, b=None, c=None):
            pass

        self.assertRoutes(['/a', '/a/:b'])

    def test_required_optional_payload(self):

        @route(op=self.op, payload=['b', 'd'])
        def a(a, b=None, c=None, d=None):
            pass

        self.assertRoutes(['/a/:a', '/a/:a/:c'])

    def test_already_defined(self):

        @route(op=self.op, name='a/:c/:b/test', payload=['b', 'd'])
        def a(a, b=None, c=None, d=None):
            pass

        self.assertRoutes(['/a/:c/:b/test/:a'])


class ResponseTest(TestCase):

    def test_none(self):

        result = response(None)

        self.assertEqual(result['total'], 0)
        self.assertTrue(result['success'])

    def test_simple(self):

        result = response(1)

        self.assertEqual(result['total'], 1)
        self.assertTrue(result['success'])

    def test_iterable(self):

        result = response([1])

        self.assertEqual(result['total'], 1)
        self.assertTrue(result['success'])

    def test_empty_iterable(self):

        result = response([])

        self.assertEqual(result['total'], 0)
        self.assertTrue(result['success'])

    def test_empty_dict(self):

        result = response({})

        self.assertEqual(result['total'], 1)
        self.assertTrue(result['success'])

    def test_dict(self):

        result = response({1: 1})

        self.assertEqual(result['total'], 1)
        self.assertTrue(result['success'])

    def test_count(self):

        total = 10
        result = response((0, total))

        self.assertEqual(result['total'], total)
        self.assertTrue(result['success'])

if __name__ == '__main__':
    main()
