# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

from canopsis.middleware.core import \
    Middleware, parse_scheme, SCHEME_SEPARATOR, DEFAULT_DATA_SCOPE


class TestUnregisteredMiddleware(Middleware):

    __protocol__ = 'test0'


class TestRegisteredMiddleware(Middleware):

    __register__ = True
    __protocol__ = 'test1'


class TestRegisteredWithDataTypeMiddleware(TestRegisteredMiddleware):

    __datatype__ = 'dttest1'


class MiddlewareTest(TestCase):

    def test_parse_scheme(self):

        uri = 'http://plop'

        protocol, data_type, data_scope = parse_scheme(uri)

        self.assertEqual(protocol, 'http')
        self.assertEqual(data_type, None)

        uri = '%s%s%s://' % ('http', SCHEME_SEPARATOR, '')

        protocol, data_type, data_scope = parse_scheme(uri)

        self.assertEqual(protocol, 'http')
        self.assertEqual(data_type, '')

        uri = '%s%s%s://' % ('http', SCHEME_SEPARATOR, 'ae')

        protocol, data_type, data_scope = parse_scheme(uri)

        self.assertEqual(protocol, 'http')
        self.assertEqual(data_type, 'ae')

    def test_resolve_middleware(self):

        protocol = TestUnregisteredMiddleware.__protocol__
        self.assertRaises(
            Middleware.Error, Middleware.resolve_middleware, protocol=protocol)

        protocol = TestRegisteredMiddleware.__protocol__
        middleware_class = Middleware.resolve_middleware(protocol=protocol)
        self.assertEqual(middleware_class, TestRegisteredMiddleware)

        data_type = None
        middleware_class = Middleware.resolve_middleware(
            protocol=protocol, data_type=data_type)
        self.assertEqual(middleware_class, TestRegisteredMiddleware)

        data_type = TestRegisteredWithDataTypeMiddleware.__datatype__
        middleware_class = Middleware.resolve_middleware(
            protocol=protocol, data_type=data_type)
        self.assertEqual(
            middleware_class, TestRegisteredWithDataTypeMiddleware)

    def test_resolve_middleware_by_uri(self):

        uri = '%s://' % TestUnregisteredMiddleware.__protocol__
        self.assertRaises(
            Middleware.Error, Middleware.resolve_middleware_by_uri, uri)

        uri = '%s://' % TestRegisteredMiddleware.__protocol__
        middleware_class = Middleware.resolve_middleware_by_uri(uri)
        self.assertEqual(middleware_class, TestRegisteredMiddleware)

        uri = '%s%s%s://' % (
            TestRegisteredWithDataTypeMiddleware.__protocol__,
            SCHEME_SEPARATOR,
            TestRegisteredWithDataTypeMiddleware.__datatype__)

        middleware_class = Middleware.resolve_middleware_by_uri(uri)
        self.assertEqual(
            middleware_class, TestRegisteredWithDataTypeMiddleware)

    def test_get_middleware(self):

        protocol = TestUnregisteredMiddleware.__protocol__
        self.assertRaises(
            Middleware.Error, Middleware.get_middleware, protocol=protocol)

        protocol = TestRegisteredMiddleware.__protocol__
        middleware = Middleware.get_middleware(protocol=protocol)
        self.assertEqual(type(middleware), TestRegisteredMiddleware)
        self.assertEqual(middleware.data_scope, DEFAULT_DATA_SCOPE)

        data_type = None
        middleware = Middleware.get_middleware(
            protocol=protocol, data_type=data_type)
        self.assertEqual(type(middleware), TestRegisteredMiddleware)
        self.assertEqual(middleware.data_scope, DEFAULT_DATA_SCOPE)

        data_scope = 'test'
        data_type = TestRegisteredWithDataTypeMiddleware.__datatype__
        middleware = Middleware.get_middleware(
            protocol=protocol, data_type=data_type, data_scope=data_scope)
        self.assertEqual(middleware.data_scope, data_scope)
        self.assertEqual(
            type(middleware), TestRegisteredWithDataTypeMiddleware)

    def test_get_middleware_by_uri(self):

        uri = '%s://' % TestUnregisteredMiddleware.__protocol__
        self.assertRaises(
            Middleware.Error, Middleware.get_middleware_by_uri, uri)

        uri = '%s://' % TestRegisteredMiddleware.__protocol__
        middleware = Middleware.get_middleware_by_uri(uri)
        self.assertEqual(type(middleware), TestRegisteredMiddleware)
        self.assertEqual(middleware.data_scope, DEFAULT_DATA_SCOPE)

        uri = '%s%s%s://' % (
            TestRegisteredWithDataTypeMiddleware.__protocol__,
            SCHEME_SEPARATOR,
            TestRegisteredWithDataTypeMiddleware.__datatype__)
        middleware = Middleware.get_middleware_by_uri(uri)
        self.assertEqual(
            type(middleware), TestRegisteredWithDataTypeMiddleware)
        self.assertEqual(middleware.data_scope, DEFAULT_DATA_SCOPE)
        self.assertEqual(middleware.data_type,
            TestRegisteredWithDataTypeMiddleware.__datatype__)

        data_scope = 'test'
        uri = '%s%s%s%s%s://' % (
            TestRegisteredWithDataTypeMiddleware.__protocol__,
            SCHEME_SEPARATOR,
            TestRegisteredWithDataTypeMiddleware.__datatype__,
            SCHEME_SEPARATOR,
            data_scope)
        middleware = Middleware.get_middleware_by_uri(uri)
        self.assertEqual(
            type(middleware), TestRegisteredWithDataTypeMiddleware)
        self.assertEqual(middleware.data_scope, data_scope)
        self.assertEqual(middleware.data_type,
            TestRegisteredWithDataTypeMiddleware.__datatype__)

if __name__ == '__main__':
    main()
