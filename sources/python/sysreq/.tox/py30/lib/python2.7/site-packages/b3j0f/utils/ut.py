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

"""Unit tests tools."""

from unittest import TestCase

from six import string_types, PY2

from .version import PY26

from re import match

__all__ = ['UTCase']


def _subset(subset, superset):
    """True if subset is a subset of superset.

    :param dict subset: subset to compare.
    :param dict superset: superset to compare.
    :return: True iif all pairs (key, value) of subset are in superset.
    :rtype: bool
    """

    result = True
    for k in subset:
        result = k in superset and subset[k] == superset[k]
        if not result:
            break
    return result


class UTCase(TestCase):
    """Class which enrichs TestCase with python version compatibilities."""

    def __init__(self, *args, **kwargs):

        super(UTCase, self).__init__(*args, **kwargs)

    if PY2:  # python 3 compatibility

        if PY26:  # python 2.7 compatibility

            def assertIs(self, first, second, msg=None):
                return self.assertTrue(first is second, msg=msg)

            def assertIsNot(self, first, second, msg=None):
                return self.assertTrue(first is not second, msg=msg)

            def assertIn(self, first, second, msg=None):
                return self.assertTrue(first in second, msg=msg)

            def assertNotIn(self, first, second, msg=None):
                return self.assertTrue(first not in second, msg=msg)

            def assertIsNone(self, expr, msg=None):
                return self.assertTrue(expr is None, msg=msg)

            def assertIsNotNone(self, expr, msg=None):
                return self.assertFalse(expr is None, msg=msg)

            def assertIsInstance(self, obj, cls, msg=None):
                return self.assertTrue(isinstance(obj, cls), msg=msg)

            def assertNotIsInstance(self, obj, cls, msg=None):
                return self.assertTrue(not isinstance(obj, cls), msg=msg)

            def assertGreater(self, first, second, msg=None):
                return self.assertTrue(first > second, msg=msg)

            def assertGreaterEqual(self, first, second, msg=None):
                return self.assertTrue(first >= second, msg=msg)

            def assertLess(self, first, second, msg=None):
                self.assertTrue(first < second, msg=msg)

            def assertLessEqual(self, first, second, msg=None):
                return self.assertTrue(first <= second, msg=msg)

            def assertRegexpMatches(self, text, regexp, msg=None):
                return self.assertTrue(
                    match(regexp, text) if isinstance(regexp, string_types)
                    else regexp.search(text),
                    msg=msg
                )

            def assertNotRegexpMatches(self, text, regexp, msg=None):
                return self.assertIsNone(
                    match(regexp, text) if isinstance(regexp, string_types)
                    else regexp.search(text),
                    msg=msg
                )

            def assertItemsEqual(self, actual, expected, msg=None):
                return self.assertEqual(
                    sorted(actual), sorted(expected), msg=msg
                )

            def assertDictContainsSubset(self, expected, actual, msg=None):
                return self.assertTrue(_subset(expected, actual), msg=msg)

            def assertCountEqual(self, first, second, msg=None):
                return self.assertEqual(len(first), len(second), msg=msg)

        def assertRegex(self, text, regexp, msg=None):
            return self.assertRegexpMatches(text, regexp, msg)

        def assertNotRegex(self, text, regexp, msg=None):
            return self.assertNotRegexpMatches(text, regexp, msg)

    else:  # python 2 compatibility
        def assertRegexpMatches(self, *args, **kwargs):
            return self.assertRegex(*args, **kwargs)

        def assertNotRegexpMatches(self, *args, **kwargs):
            return self.assertNotRegex(*args, **kwargs)

        def assertItemsEqual(self, actual, expected, msg=None):
            return self.assertEqual(sorted(actual), sorted(expected), msg=msg)

        def assertDictContainsSubset(self, expected, actual, msg=None):
            return self.assertTrue(_subset(expected, actual), msg=msg)
