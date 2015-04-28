#!/usr/bin/env python
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

from canopsis.common.utils import (
    lookup, path, isiterable, isunicode, ensure_unicode, ensure_iterable,
    forceUTF8
)

from sys import version as PYVER


def _test():
    pass


class UtilsTest(TestCase):

    def setUp(self):
        pass

    def test_lookup(self):

        # resolve builtin function
        _open = lookup('{0}.open'.format(open.__module__))

        self.assertTrue(_open is open)

        # resolve lookup
        _lookup = lookup(
            'canopsis.common.utils.lookup')

        self.assertTrue(_lookup is lookup)

        # resolve package

        canopsis = lookup('canopsis')

        self.assertEqual(canopsis.__name__, 'canopsis')

        # resolve sub_module

        canopsis_common = lookup('canopsis.common')

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
        self.assertEqual(path(lookup(_path)), _path)

        # Test if you can retrieve the function by resolving the path got using path()
        self.assertEqual(lookup(path(path)), path)

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
            self.assertTrue(isinstance(ensure_unicode(str()), unicode))
            self.assertTrue(isinstance(ensure_unicode(unicode()), unicode))
            self.assertRaises(TypeError, ensure_unicode)

    def test_ensure_iterable(self):

        self.assertEqual(ensure_iterable(2), [2])
        self.assertEqual(ensure_iterable("2"), ["2"])
        self.assertEqual(ensure_iterable([2]), [2])
        self.assertEqual(ensure_iterable([2], iterable=set), {2})

    def test_forceUTF8(self):

        notUTF8 = "é"
        utf8 = unicode(notUTF8, "utf-8") if PYVER < "3" else notUTF8

        data_to_check = notUTF8
        result = forceUTF8(data_to_check)
        self.assertEqual(result, utf8)

        data_to_check = {notUTF8: notUTF8, utf8: utf8, 1: 1}
        result = forceUTF8(data_to_check)
        data_to_compere = data_to_check if PYVER < "3" else {utf8: utf8, 1: 1}
        self.assertEqual(str(result), str(data_to_compere))

        data_to_check = [notUTF8, utf8, 1]
        result = forceUTF8(data_to_check)
        self.assertEqual(result, [utf8, utf8, 1])

        data_to_check = (notUTF8, utf8, 1)
        result = forceUTF8(data_to_check)
        self.assertEqual(result, (utf8, utf8, 1))

        data_to_check = {notUTF8, utf8, 1}
        result = forceUTF8(data_to_check)
        self.assertEqual(result, {utf8, 1})

if __name__ == '__main__':
    main()
