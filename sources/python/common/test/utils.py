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

from canopsis.common.utils import \
    resolve_element, path, isiterable, isunicode, force_unicode, force_iterable

from sys import version as PYVER


def _test():
    pass


class UtilsTest(TestCase):

    def setUp(self):
        pass

    def test_resolve_element(self):

        # resolve builtin function
        _open = resolve_element('%s.open' % open.__module__)

        self.assertTrue(_open is open)

        # resolve resolve_element
        _resolve_element = resolve_element(
            'canopsis.common.utils.resolve_element')

        self.assertTrue(_resolve_element is resolve_element)

        # resolve package

        canopsis = resolve_element('canopsis')

        self.assertEqual(canopsis.__name__, 'canopsis')

        # resolve sub_module

        canopsis_common = resolve_element('canopsis.common')

        self.assertEqual(canopsis_common.__name__, 'canopsis.common')

    def test_path(self):

        # resolve built-in function
        open_path = path(open)

        self.assertEqual(open_path, '__builtin__.open')

        # resolve path
        self.assertEqual(path(path), 'canopsis.common.utils.path')

        # resolve package
        import canopsis
        self.assertEqual(path(canopsis), 'canopsis')

        # resolve sub-module
        import canopsis.common as canopsis_common
        self.assertEqual(path(canopsis_common), 'canopsis.common')

    def test_reciproc(self):

        _path = 'canopsis.common.utils.path'

        # Test if you can get the path _path using path() on the resolved element
        self.assertEqual(path(resolve_element(_path)), _path)

        # Test if you can retrieve the function by resolving the path got using path()
        self.assertEqual(resolve_element(path(path)), path)

    def test_isiterable(self):

        self.assertFalse(isiterable(2))

        self.assertTrue(isiterable([]))

        self.assertTrue(isiterable(""))

        self.assertFalse(isiterable('', is_str=False))

    def test_isunicode(self):

        if PYVER < '3':
            self.assertFalse(isunicode(str()))
            self.assertTrue(isunicode(unicode()))

    def test_forceunicode(self):

        if PYVER < '3':
            self.assertTrue(isinstance(force_unicode(str()), unicode))
            self.assertTrue(isinstance(force_unicode(unicode()), unicode))
            self.assertRaises(TypeError, force_unicode)

    def test_force_iterable(self):

        self.assertEqual(force_iterable(2), [2])
        self.assertEqual(force_iterable("2"), ["2"])
        self.assertEqual(force_iterable([2]), [2])
        self.assertEqual(force_iterable([2], iterable=set), {2})

if __name__ == '__main__':
    main()
