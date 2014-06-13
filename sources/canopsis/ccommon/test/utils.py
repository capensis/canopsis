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

from ccommon.utils import resolve_element, path


def _test():
    pass


class UtilsTest(TestCase):

    def setUp(self):
        pass

    def test_resolve_element(self):

        # resolve builtin function
        _open = resolve_element('__builtin__.open')

        self.assertTrue(_open is open)

        # resolve class
        utilsTest = resolve_element('test.utils.UtilsTest')

        self.assertTrue(utilsTest is UtilsTest)

        # do not resolve method
        setUp = resolve_element('test.utils.UtilsTest.setUp')

        self.assertFalse(setUp is UtilsTest.setUp)

        # resolve function
        test = resolve_element('test.utils._test')

        self.assertTrue(_test is test)

        # resolve resolve element
        _resolve_element = resolve_element('ccommon.utils.resolve_element')

        self.assertTrue(_resolve_element is resolve_element)

    def test_path(self):

        open_path = path(open)

        self.assertEqual(open_path, '__builtin__.open')

        self.assertEqual(path(UtilsTest), 'test.utils.UtilsTest')

        self.assertEqual(resolve_element(path(open)), open)

if __name__ == '__main__':
    main()
