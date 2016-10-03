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

from unittest import main

from b3j0f.utils.ut import UTCase

from ..core import Annotation
from ..oop import Transform, Mixin, MethodMixin, Deprecated, Singleton


class TransformTest(UTCase):
    """Test Transform Annotation.
    """

    def test_default(self):

        @Transform()
        class A:
            pass

        self.assertTrue(issubclass(A, Annotation))


class OOPTests(UTCase):

    def setUp(self):
        pass

    def testMixIn(self):

        class A(object):
            def plop(self):
                pass

        plop = A.plop

        mixedins_by_name = Mixin.get_mixedins_by_name(A)

        self.assertEqual(len(mixedins_by_name), 0)

        Mixin.set_mixin(A, 2, 'a')
        Mixin.set_mixin(A, lambda: None, 'b')
        Mixin.set_mixin(A, None, 'plop')
        Mixin.set_mixin(A, None, 'plop')
        Mixin.set_mixin(A, None, 'plop')

        mixedins_by_name = Mixin.get_mixedins_by_name(A)

        self.assertEqual(len(mixedins_by_name), 3)
        self.assertEqual(A.a, 2)
        self.assertIsNone(A.plop)

        Mixin.remove_mixin(A, 'plop')

        mixedins_by_name = Mixin.get_mixedins_by_name(A)

        self.assertEqual(len(mixedins_by_name), 3)

        Mixin.remove_mixins(A)

        mixedins_by_name = Mixin.get_mixedins_by_name(A)

        self.assertEqual(len(mixedins_by_name), 0)
        self.assertFalse(hasattr(A, 'a'))
        self.assertEqual(plop, A.plop)

    def testClassMixIn(self):

        class ClassForMixIn(object):
            def get_1(self):
                return 1

            def get_2(self):
                return 2

        def get_3(self):
            return self.a

        a = None

        @Mixin(ClassForMixIn, lambda: None, get_3=get_3, a=a)
        class MixedInClass(object):
            def get_1(self):
                return '1'

            def get_3(self):
                return 3

        self.assertTrue(hasattr(MixedInClass, 'get_2'))
        self.assertTrue(hasattr(MixedInClass, 'a'))
        self.assertTrue(hasattr(MixedInClass, '<lambda>'))

        mixedInstance = MixedInClass()

        self.assertEqual(mixedInstance.a, None)
        self.assertEqual(mixedInstance.get_3(), a)
        self.assertIs(mixedInstance.get_1(), 1)

        Mixin.remove_mixins(MixedInClass)
        self.assertEqual(mixedInstance.get_3(), 3)
        self.assertIs(mixedInstance.get_1(), '1')

        Mixin(ClassForMixIn, lambda: None, get_3=get_3, a=a)(MixedInClass)

        self.assertTrue(hasattr(MixedInClass, 'get_2'))
        self.assertTrue(hasattr(MixedInClass, 'a'))
        self.assertTrue(hasattr(MixedInClass, '<lambda>'))

        self.assertIsNone(mixedInstance.a)
        self.assertEqual(mixedInstance.get_3(), a)
        self.assertIs(mixedInstance.get_1(), 1)

    def testMethodMixIn(self):

        def test(self):
            return None

        class MixedInClass(object):
            @MethodMixin(test)
            def get_1(self):
                return 1

        self.assertTrue(hasattr(MixedInClass, 'get_1'))

        mixedInstance = MixedInClass()

        self.assertIsNone(mixedInstance.get_1())

        class MixedInClass(object):
            def get_1(self):
                return 1

        MethodMixin(test)(MixedInClass.get_1, ctx=MixedInClass)

        self.assertTrue(hasattr(MixedInClass, 'get_1'))

        mixedInstance = MixedInClass()

        self.assertIsNone(mixedInstance.get_1())

        Mixin.remove_mixins(MixedInClass)

        self.assertEqual(mixedInstance.get_1(), 1)

    def testDeprecated(self):

        @Deprecated()
        def b():
            pass

        b()

    def testSingleton(self):

        @Singleton(a=1)
        class A(object):
            def __init__(self, a):
                self.a = a

            def get_a(self):
                return self.a

        self.assertEqual(A, A())

        self.assertEqual(A(), A())

        self.assertEqual(A.get_a(), A().get_a())


if __name__ == '__main__':
    main()
