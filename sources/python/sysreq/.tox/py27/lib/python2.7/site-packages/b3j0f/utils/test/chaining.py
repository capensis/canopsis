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

from ..ut import UTCase
from ..chaining import Chaining, ListChaining


class ChainingTest(UTCase):
    """
    Test Chaining object.
    """

    def setUp(self):
        """
        Initialize object to embed into a Chaining.
        """

        self.content = None
        self.chaining = Chaining(self.content)

    def test_empty(self):
        """
        Test Chaining result without calls
        """

        results = self.chaining[:]

        self.assertFalse(results)

    def test_one(self):
        """
        Test one chaining of calls.
        """

        self.chaining.__hash__()
        results = self.chaining[:]

        self.assertEqual(len(results), 1)
        self.assertEqual(results[0], hash(self.content))

        result = self.chaining[0]

        self.assertEqual(result, hash(self.content))

    def test_many(self):
        """
        Test several chaining of calls.
        """

        self.chaining.__hash__().__repr__()

        results = self.chaining[:]

        self.assertEqual(len(results), 2)
        self.assertEqual(results[0], hash(self.content))
        self.assertEqual(results[1], repr(self.content))
        self.assertEqual(results[-1], repr(self.content))

        first_result = self.chaining[0]
        self.assertEqual(first_result, hash(self.content))

        second_result = self.chaining[1]
        self.assertEqual(second_result, repr(self.content))

        last_result = self.chaining[-1]
        self.assertEqual(last_result, repr(self.content))

    def test_exception(self):
        """
        Test when an exception occured.
        """

        self.assertRaises(
            AttributeError,
            lambda: self.chaining.raiseexception
        )


class ListChainingTest(UTCase):
    """
    Test ListChaining object.
    """

    def setUp(self):
        """
        Initialize object to embed into a Chaining.
        """

        self.content = [None, '']
        self.chaining = ListChaining(*self.content)

    def test_empty(self):
        """
        Test Chaining result without calls
        """

        results = self.chaining[:]

        self.assertFalse(results)

    def test_one(self):
        """
        Test one chaining of calls.
        """

        self.chaining.__hash__()
        results = self.chaining[:]

        self.assertEqual(len(results), 1)
        self.assertEqual(
            results[0],
            [hash(self.content[0]), hash(self.content[1])]
        )

        result = self.chaining[0]

        self.assertEqual(result[0], hash(self.content[0]))
        self.assertEqual(result[1], hash(self.content[1]))

    def test_many(self):
        """
        Test several chaining of calls.
        """

        self.chaining.__hash__().__repr__()

        results = self.chaining[:]

        self.assertEqual(len(results), 2)
        self.assertEqual(
            results[0],
            [hash(self.content[0]), hash(self.content[1])]
        )
        self.assertEqual(
            results[1],
            [repr(self.content[0]), repr(self.content[1])]
        )
        self.assertEqual(
            results[-1],
            [repr(self.content[0]), repr(self.content[1])]
        )

        first_result = self.chaining[0]
        self.assertEqual(
            first_result,
            [hash(self.content[0]), hash(self.content[1])]
        )

        second_result = self.chaining[1]
        self.assertEqual(
            second_result,
            [repr(self.content[0]), repr(self.content[1])]
        )

        last_result = self.chaining[-1]
        self.assertEqual(
            last_result,
            [repr(self.content[0]), repr(self.content[1])]
        )

    def test_exception(self):
        """
        Test when an exception occured.
        """
        self.chaining.upper()

        results = self.chaining[:]

        self.assertEqual(len(results), 1)
        self.assertIsInstance(results[0][0], Exception)
        self.assertEqual(results[0][1], self.content[1].upper())


if __name__ == '__main__':
    main()
