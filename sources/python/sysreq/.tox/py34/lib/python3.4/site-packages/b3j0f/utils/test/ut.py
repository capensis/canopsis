#!/usr/bin/env python
# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2014 Jonathan Labéjof <jonathan.labejof@gmail.com>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
# --------------------------------------------------------------------

from unittest import main

from b3j0f.utils.ut import UTCase

from re import compile as re_compile


class UTCaseTest(UTCase):

    def test_assertIs(self):
        self.assertIs(1, 1)

    def test_assertIsNot(self):
        self.assertIsNot(1, 2)

    def test_assertIn(self):
        self.assertIn(1, {1: 1})

    def test_assertNotIn(self):
        self.assertNotIn(1, {})

    def test_assertIsNone(self):
        self.assertIsNone(None)

    def test_assertIsNotNone(self):
        self.assertIsNotNone(1)

    def test_assertIsInstance(self):
        self.assertIsInstance(1, int)

    def test_assertNotIsInstance(self):
        self.assertNotIsInstance(None, type)

    def test_assertGreater(self):
        self.assertGreater(2, 1)

    def test_assertGreaterEqual(self):
        self.assertGreaterEqual(1, 1)
        self.assertGreaterEqual(2, 1)

    def test_assertLess(self):
        self.assertLess(1, 2)

    def test_assertLessEqual(self):
        self.assertLessEqual(1, 2)
        self.assertLessEqual(1, 1)

    def test_assertRegexpMatches(self):
        self.assertRegexpMatches('a', 'a')
        self.assertRegexpMatches('a', re_compile('a'))

    def test_assertRegex(self):
        self.assertRegex('a', 'a')
        self.assertRegex('a', re_compile('a'))

    def test_assertNotRegexpMatches(self):
        self.assertNotRegexpMatches('b', 'a')
        self.assertNotRegexpMatches('a', re_compile('b'))

    def test_assertNotRegex(self):
        self.assertNotRegex('a', 'b')
        self.assertNotRegex('a', re_compile('b'))

    def test_assertItemsEqual(self):
        self.assertItemsEqual([1, 2], [2, 1])

    def test_assertDictContainsSubset(self):
        self.assertDictContainsSubset({1: 2}, {1: 2, 3: 4})


if __name__ == '__main__':
    main()
