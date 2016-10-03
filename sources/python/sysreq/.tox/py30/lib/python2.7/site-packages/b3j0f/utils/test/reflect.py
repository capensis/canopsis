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
from ..reflect import base_elts, find_embedding, isoldstyle
from ..version import PY3

from inspect import getmodule


class IsOldStyle(UTCase):
    """Test the function isoldstyle."""

    def test_true_simple(self):
        """Test true."""

        class Test:
            """Test class."""
            pass

        self.assertTrue(PY3 or isoldstyle(Test))

    def test_true_multi(self):
        """Test true multi."""

        class Test:
            pass

        class SubTest(Test):
            pass

        self.assertTrue(PY3 or isoldstyle(SubTest))

    def test_false_simple(self):
        """Test false."""
        class Test(object):
            pass

        self.assertFalse(isoldstyle(Test))

    def test_false_multi(self):
        """Test false multi."""

        class Test(object):
            pass

        class SubTest(Test):
            pass

        self.assertFalse(isoldstyle(Test))


class BaseEltsTest(UTCase):
    """
    Test base_elts function
    """

    def test_not_inherited(self):
        """
        Test with a not inherited element.
        """

        bases = base_elts(None)
        self.assertFalse(bases)

    def test_function(self):
        """
        Test function
        """

        bases = base_elts(lambda: None)
        self.assertFalse(bases)

    def test_class(self):
        """
        Test class
        """

        class A:
            pass

        class B(A, dict):
            pass

        bases = base_elts(B)
        self.assertEqual(bases, list(B.__bases__) + [object])

    def test_method(self):
        """
        Test method
        """

        class A:
            def a(self):
                pass

        class B(A):
            pass

        bases = base_elts(B.a, cls=A)
        self.assertEqual(len(bases), 1)
        base = bases.pop()
        self.assertEqual(base, A.a)

    def test_not_method(self):
        """
        Test when method has been overriden
        """

        class A:
            def a(self):
                pass

        class B(A):
            def a(self):
                pass

        bases = base_elts(B.a, cls=A)
        self.assertFalse(bases)

    def test_boundmethod(self):
        """
        Test bound method
        """

        class Test:
            def test(self):
                pass

        test = Test()

        bases = base_elts(test.test)
        self.assertEqual(len(bases), 1)
        self.assertEqual(bases.pop(), Test.test)

    def test_not_boundmethod(self):
        """
        Test with a bound method which is only defined in the instance
        """

        class Test:
            def test(self):
                pass

        test = Test()
        test.test = lambda self: None

        bases = base_elts(test.test)
        self.assertFalse(bases)


class FindEmbeddingTest(UTCase):

    def test_none(self):

        embedding = find_embedding(None)

        self.assertFalse(embedding)

    def test_wrong_embedding(self):

        embedding = find_embedding(None, embedding=FindEmbeddingTest)

        self.assertFalse(embedding)

    def test_module(self):

        FindEmbeddingTestModule = getmodule(FindEmbeddingTest)
        embedding = find_embedding(FindEmbeddingTestModule)

        self.assertEqual(len(embedding), 1)
        self.assertIs(FindEmbeddingTestModule, embedding[0])

    def test_function(self):

        embedding = find_embedding(find_embedding)

        self.assertEqual(len(embedding), 2)
        self.assertIs(embedding[0], getmodule(find_embedding))
        self.assertIs(embedding[1], find_embedding)

    def test_class(self):

        embedding = find_embedding(FindEmbeddingTest)

        self.assertEqual(len(embedding), 2)
        self.assertIs(getmodule(FindEmbeddingTest), embedding[0])
        self.assertIs(FindEmbeddingTest, embedding[1])

    def test_method(self):

        embedding = find_embedding(FindEmbeddingTest.test_method)

        self.assertEqual(len(embedding), 3)
        self.assertIs(getmodule(FindEmbeddingTest), embedding[0])
        self.assertIs(FindEmbeddingTest, embedding[1])
        self.assertEqual(FindEmbeddingTest.test_method, embedding[2])

if __name__ == '__main__':
    main()
