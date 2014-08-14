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

from canopsis.middleware import \
    Middleware, parse_scheme, SCHEME_SEPARATOR


class TestUnregisteredMiddleware(Middleware):

    __protocol__ = 'notprotocoltest'


class TestRegisteredMiddleware(Middleware):

    __register__ = True
    __protocol__ = 'protocoltest'


class TestRegisteredWithDataTypeMiddleware(TestRegisteredMiddleware):

    __datatype__ = 'datatypetest'


class MiddlewareTest(TestCase):

    def test_parse_scheme(self):

        uri = 'http://plop'

        protocol, data_type = parse_scheme(uri)

        self.assertEqual(protocol, 'http')
        self.assertEqual(data_type, None)

        uri = '%s%s%s://' % ('http', SCHEME_SEPARATOR, '')

        protocol, data_type = parse_scheme(uri)

        self.assertEqual(protocol, 'http')
        self.assertEqual(data_type, '')

        uri = '%s%s%s://' % ('http', SCHEME_SEPARATOR, 'ae')

        protocol, data_type = parse_scheme(uri)

        self.assertEqual(protocol, 'http')
        self.assertEqual(data_type, 'ae')

    def test_resolve_middleware(self):

        uri = '%s://' % TestUnregisteredMiddleware.__protocol__
        self.assertRaises(Middleware.Error, Middleware.resolve_middleware, uri)

        uri = '%s://' % TestRegisteredMiddleware.__protocol__
        middleware_class = Middleware.resolve_middleware(uri)
        self.assertEqual(middleware_class, TestRegisteredMiddleware)

        uri = '%s%s%s://' % (
            TestRegisteredWithDataTypeMiddleware.__protocol__,
            SCHEME_SEPARATOR,
            TestRegisteredWithDataTypeMiddleware.__datatype__)

        middleware_class = Middleware.resolve_middleware(uri)
        self.assertEqual(
            middleware_class, TestRegisteredWithDataTypeMiddleware)

if __name__ == '__main__':
    main()
