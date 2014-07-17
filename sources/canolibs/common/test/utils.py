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

from canopsis.common.utils import resolve_element, path


def _test():
    pass


class UtilsTest(TestCase):

    def setUp(self):
        pass

    def test_resolve_element(self):

        # resolve builtin function
        _open = resolve_element('%s.open' % open.__module__)

        self.assertTrue(_open is open)

        # resolve function
        test = resolve_element('test.utils._test')

        self.assertTrue(_test is test)

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

        _path = 'b3j0f.common.utils.path'

        self.assertEqual(path(resolve_element(_path)), _path)

        self.assertEqual(resolve_element(path(resolve_element)), path)

if __name__ == '__main__':
    main()
