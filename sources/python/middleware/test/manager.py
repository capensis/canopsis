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

from canopsis.middleware import Middleware, SCHEME_SEPARATOR
from canopsis.middleware.manager import Manager


class TestUnregisteredMiddleware(Middleware):

    __protocol__ = 'notprotocoltest'


class TestRegisteredMiddleware(Middleware):

    __register__ = True
    __protocol__ = 'protocoltest'


class TestRegisteredWithDataTypeMiddleware(TestRegisteredMiddleware):

    __datatype__ = 'datatypetest'


class TestManager(Manager):
    pass


class ManagerTest(TestCase):

    def setUp(self):

        self.manager = TestManager()

    def test_get_middleware(self):

        uri = '%s://' % (TestUnregisteredMiddleware.__protocol__)

        self.assertRaises(Middleware.Error, self.manager.get_middleware, uri)

        uri = '%s://' % (TestRegisteredMiddleware.__protocol__)

        middleware = self.manager.get_middleware(uri)

        self.assertTrue(type(middleware) is TestRegisteredMiddleware)

        middleware2 = self.manager.get_middleware(uri)

        self.assertTrue(middleware is middleware2)

        middleware3 = self.manager.get_middleware(uri, shared=False)

        self.assertFalse(middleware is middleware3)

        self.manager.shared = False

        middleware4 = self.manager.get_middleware(uri)

        self.assertFalse(middleware is middleware4)

        uri = '%s%s%s://' % (
            TestRegisteredWithDataTypeMiddleware.__protocol__,
            SCHEME_SEPARATOR,
            TestRegisteredWithDataTypeMiddleware.__datatype__)

        middleware_wd = self.manager.get_middleware(uri)

        self.assertTrue(
            type(middleware_wd) is TestRegisteredWithDataTypeMiddleware)

if __name__ == '__main__':
    main()
