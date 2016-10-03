#!/usr/bin/env python
# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2015 Jonathan Labéjof <jonathan.labejof@gmail.com>
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

from b3j0f.utils.ut import UTCase

from six.moves import range

from ..core import Annotation
from ..check import Condition, MaxCount, Target, AnnotationChecker

from random import randint


class ConditionTest(UTCase):
    """Test the Condition annotation.
    """
    def setUp(self):

        self.pre_count = 0
        self.post_count = 0
        self.condition = Condition(
            pre_cond=self.pre_cond, post_cond=self.post_cond
        )
        self.condition(self._test)

    def pre_cond(self, joinpoint):

        self.pre_count += 1

    def post_cond(self, joinpoint):

        self.post_count += 1

    def tearDown(self):

        self.condition.__del__()
        del self.condition

    def _test(self, **kwargs):
        """Test method which is intercepted by this condition.
        """

        return self

    def test_pre(self):
        """Test pre condition checking.
        """

        self.assertEqual(self.pre_count, 0)
        self.assertEqual(self.post_count, 0)
        result = self._test()
        self.assertEqual(self.pre_count, 1)
        self.assertEqual(self.post_count, 1)
        self.assertEqual(result, self)


class AnnotationCheckerTest(UTCase):
    """Test AnnotationChecker.
    """

    class TestAnnotationChecker(AnnotationChecker):

        def __init__(self, utcase, *args, **kwargs):

            super(AnnotationCheckerTest.TestAnnotationChecker, self).__init__(
                *args, **kwargs
            )
            self.utcase = utcase

        def _interception(self, joinpoint, *args, **kwargs):

            self.utcase.count += 1

            return joinpoint.proceed()

    def setUp(self):

        self.count = 0
        self.annotation = AnnotationCheckerTest.TestAnnotationChecker(self)

    def tearDown(self):

        self.annotation.__del__()
        del self.annotation

    def test_annotation_class(self):
        """Test to annotate an annotation class.
        """

        self.annotation(Annotation)
        annotation = Annotation()
        annotation(Annotation)
        annotation.__del__()
        self.assertEqual(self.count, 1)

    def test_not_Annotation(self):
        """Test to annotate an object which is not an annotation.
        """

        self.assertRaises(Target.Error, self.annotation, self)


class MaxCountTest(UTCase):
    """Test MaxCount annotation.
    """

    def _assertMaxCount(self, count=MaxCount.DEFAULT_COUNT):

        @MaxCount(count)
        class Test(Annotation):
            pass

        # weave count time test on None
        for i in range(count):
            Test()(None)

        # check if next weaving raise an Exception
        self.assertRaises(MaxCount.Error, Test(), None)

    def test_default(self):
        """Test default count.
        """

        self._assertMaxCount()

    def test_one(self):
        """Test one count.
        """

        self._assertMaxCount(1)

    def test_more_than_one(self):
        """Test > 1 count.
        """

        self._assertMaxCount(randint(2, 5))


class TargetTest(UTCase):
    """Test Target Annotation.
    """

    def test_class(self):
        """Test to use type such as types.
        """

        @Target(type)
        class Test(Annotation):
            pass

        # check to weave on a class
        @Test()
        class TestClassBis(object):
            pass

        # check to weave on a namespace
        @Test()
        class TestNSBis():
            pass

        # check to fail on a lambda expression
        self.assertRaises(Target.Error, Test(), lambda x: None)

    def test_callable(self):
        """Test the type callable.
        """

        @Target(callable)
        class Test(Annotation):
            pass

        @Test()
        def test():
            pass

        @Test()
        class TestClass:
            def __call__(self):
                pass

        testInstance = TestClass()

        Test()(testInstance)

        del TestClass.__call__
        testInstance = TestClass()

        self.assertRaises(Target.Error, Test(), testInstance)

    def test_function(self):
        """Test the type function.
        """

        @Target(Target.FUNC)
        class Test(Annotation):
            pass

        @Test()
        def test():
            pass

        class TestClass:
            pass

        self.assertRaises(Target.Error, Test(), TestClass)

    def test_multitypes(self):
        """Test multi types.
        """

        class TestA:
            pass

        class TestB:
            pass

        class TestAB(TestA, TestB):
            pass

        @Target([TestA, TestB])
        class Test(Annotation):
            pass

        # check types
        Test()(TestA)
        Test()(TestB)
        Test()(TestAB)
        # check instances is not checked
        self.assertRaises(Target.Error, Test(), TestA())
        self.assertRaises(Target.Error, Test(), TestB())
        self.assertRaises(Target.Error, Test(), TestAB())

        @Target([TestA, TestB], instances=True)
        class Test(Annotation):
            pass

        # check types
        Test()(TestA)
        Test()(TestB)
        Test()(TestAB)
        # check instances
        Test()(TestA())
        Test()(TestB())
        Test()(TestAB())

        @Target([TestA, TestB], rule=Target.AND)
        class Test(Annotation):
            pass

        # check types
        self.assertRaises(Target.Error, Test(), TestA)
        self.assertRaises(Target.Error, Test(), TestB)
        Test()(TestAB)
        # check instances is not checked
        self.assertRaises(Target.Error, Test(), TestA())
        self.assertRaises(Target.Error, Test(), TestB())
        self.assertRaises(Target.Error, Test(), TestAB())

        @Target([TestA, TestB], rule=Target.AND, instances=True)
        class Test(Annotation):
            pass

        # check types
        self.assertRaises(Target.Error, Test(), TestA)
        self.assertRaises(Target.Error, Test(), TestB)
        Test()(TestAB)
        # check instances is checked
        self.assertRaises(Target.Error, Test(), TestA())
        self.assertRaises(Target.Error, Test(), TestB())
        Test()(TestAB())


if __name__ == '__main__':
    main()
