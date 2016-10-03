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

from __future__ import absolute_import

from unittest import main

from ..ut import UTCase
from ..iterable import (
    first, ensureiterable, isiterable, last, itemat, sliceit, hashiter
)

from random import random

from six import string_types

from ..version import OrderedDict


class EnsureIterableTest(UTCase):
    """Test ensure iterable function."""

    def test_list(self):
        """test list."""

        value = []
        itererable = ensureiterable(value)
        self.assertEqual(itererable, value)

    def test_dict(self):
        """test dict."""

        value = []
        iterable = ensureiterable(value, iterable=dict)
        self.assertTrue(isinstance(iterable, dict))
        self.assertFalse(iterable)

    def test_exclude(self):
        """test exclude."""

        value = ""
        iterable = ensureiterable(value, exclude=string_types)
        self.assertTrue(iterable)


class IsIterable(UTCase):
    """Test isiterable."""

    def test_iterable(self):
        """Test an iterable value."""

        self.assertTrue(isiterable([]))

    def test_exclude(self):
        """Test iterable and not allowed types."""

        self.assertFalse(isiterable([], exclude=list))

    def test_excludes(self):
        """Test iterable with a tuple of exclude types."""

        self.assertFalse(isiterable([], exclude=(list,) + string_types))

    def test_not_iterable(self):
        """Test not iterable element."""

        self.assertFalse(isiterable(None))

def _randlist():
    """Generate a random tuple of float."""

    return list((random(), random()) for _ in range(5))


class _Set(object):
    """Base test class for first, last, itemat and slice."""

    def _testfunctionandparams(self):
        """Get the function to test with kwargs."""

        raise NotImplementedError()

    def _assertvalue(self, _type):
        """Assert input value."""

        raise NotImplementedError()

    def test_dict(self):
        """Test dict."""

        self._assertvalue(dict)

    def test_str(self):
        """Test str."""

        self._assertvalue(str)

    def test_list(self):
        """Test list."""

        self._assertvalue(list)

    def test_tuple(self):
        """Test tuple."""

        self._assertvalue(tuple)

    def test_set(self):
        """Test set."""

        self._assertvalue(set)

    def test_ordereddict(self):
        """Test ordered dict."""

        self._assertvalue(OrderedDict)

    def test_object(self):
        """Test object."""

        class Test(object):
            """Test Object."""

            def __init__(self, value=None):

                self.value = [] if value is None else value

            def __iter__(self):

                return iter(self.value)

            def __len__(self):

                return len(self.value)

            def __getslice__(self, lower, upper):

                return self.value[lower: upper]

            def __eq__(self, other):

                valuetocmp = other.value if isinstance(other, Test) else other

                return self.value == valuetocmp

        self._assertvalue(Test)

    def test_notiterable(self):
        """Test not iterable."""

        testfunction, params = self._testfunctionandparams()

        self.assertRaises(TypeError, testfunction, *params)


class First(UTCase, _Set):
    """Test the function first."""

    def _testfunctionandparams(self):

        return first, (None, )

    def _assertvalue(self, _type):

        default = 'test'

        # test empty iterable
        empty = _type()
        val = first(empty, default=default)
        self.assertEqual(val, default)

        # test with not empty iterable
        randlist = _randlist()
        iterable = _type(randlist)
        val = first(iterable, default=default)
        value = next(iter(iterable))

        self.assertEqual(value, val)



class Last(UTCase, _Set):
    """Test the function last."""

    def _testfunctionandparams(self):

        return last, (None, )

    def _assertvalue(self, _type):
        """Assert input value."""

        default = 'test'

        # test empty iterable
        empty = _type()
        val = last(empty, default=default)
        self.assertEqual(val, default)

        # test with not empty iterable
        randlist = _randlist()
        iterable = _type(randlist)
        val = last(iterable, default=default)
        iterator = iter(iterable)
        while True:
            try:
                value = next(iterator)
            except StopIteration:
                break

        self.assertEqual(value, val)


class ItemAt(UTCase, _Set):
    """Test the function itemat."""

    def _testfunctionandparams(self):

        return itemat, (None, 0)

    def _assertvalue(self, _type):
        """Assert input value."""

        # test empty iterable
        empty = _type()
        self.assertRaises(IndexError, itemat, empty, 0)

        # test with not empty iterable
        randlist = _randlist()
        iterable = _type(randlist)

        for index in range(len(iterable)):  # check positive indexes
            value = itemat(iterable, index)
            iterator = iter(iterable)
            for _ in range(index + 1):
                val = next(iterator)
            self.assertEqual(value, val)

        for index in range(-1, -len(iterable), -1):  # check negative indexes
            value = itemat(iterable, index)
            iterator = iter(iterable)
            for _ in range(index + len(iterable) + 1):
                val = next(iterator)
            self.assertEqual(value, val)

        # assert raise IndexError
        self.assertRaises(IndexError, itemat, empty, len(iterable))


class SliceIt(UTCase, _Set):
    """Test the function sliceit."""

    def _testfunctionandparams(self):

        return sliceit, (None, )

    def _assertvalue(self, _type):
        """Assert input value."""

        isdict = issubclass(_type, dict)

        # test empty iterable
        empty = _type()
        value = sliceit(iterable=empty)
        if isdict:
            empty = []
        self.assertEqual(empty, value)

        # test with not empty iterable
        randlist = _randlist()
        iterable = _type(randlist)

        # check upper >= lower
        value = sliceit(iterable, 10, 0)
        self.assertEqual(len(value), 0)

        value = sliceit(iterable, 1, 1)
        self.assertEqual(len(value), 0)

        len_iterable = len(iterable)

        # check for all lower and upper
        for lower in range(len_iterable):

            for upper in range(lower, len_iterable):

                if upper <= lower:
                    continue

                value = sliceit(iterable, lower, upper)

                if isinstance(iterable, string_types):
                    val = iterable[lower:upper]

                else:

                    val = []

                    index = lower

                    for index in range(lower, upper):
                        item = itemat(iterable, index)
                        val.append(item)

                    if not isdict:
                        val = _type(val)

                self.assertEqual(val, value)

        # check for all negatives lower and upper
        for lower in range(-1, -len_iterable, -1):

            for upper in range(-1, -len_iterable, -1):

                if upper <= lower:
                    continue

                value = sliceit(iterable, lower, upper)

                if isinstance(iterable, string_types):
                    val = iterable[lower:upper]

                else:

                    val = []

                    index = lower

                    for index in range(lower, upper):
                        item = itemat(iterable, index)
                        val.append(item)

                    if not isdict:
                        val = _type(val)

                self.assertEqual(val, value)


class HashIterTest(UTCase):
    """Test the hashiter function."""

    def test_hashable(self):
        """Test to hash an hashable object."""

        test = 'test'

        result = hashiter(test)

        self.assertEqual(result, hash(test))

    def test_list(self):
        """Test to hash a list."""

        test = ['test', 1, list()]

        result = hashiter(test)

        self.assertEqual(
            result,
            hash(list) +
            (hash('test') + 1) * 1 +
            (hash(1) + 1) * 2 + (hashiter([]) + 1) * 3
        )

    def test_set(self):
        """Test to hash a set."""

        test = set([1, 2, 3])

        result = hashiter(test)

        self.assertEqual(
            result,
            hash(set) +
            (hash(1) + 1) * 1 + (hash(2) + 1) * 2 + (hash(3) + 1) * 3
        )

    def test_dict(self):
        """Test to hash a dict."""

        test = {'test0': 0, 'test1': 1}

        result = hashiter(test)

        self.assertEqual(
            result,
            hash(dict) +
            (hash('test0') + 1) * (hash(0) + 1) +
            (hash('test1') + 1) * (hash(1) + 1)
        )


if __name__ == '__main__':
    main()
