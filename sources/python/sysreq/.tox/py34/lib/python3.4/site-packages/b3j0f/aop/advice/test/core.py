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

from six import PY2

from b3j0f.utils.ut import UTCase

from ..core import weave, unweave, weave_on

from time import sleep


class WeaveTest(UTCase):

    def setUp(self):

        self.count = 0

    def joinpoint(self, joinpoint):
        """
        Default interceptor which increments self count
        """

        self.count += 1
        return joinpoint.proceed()

    def test_builtin(self):

        weave(target=min, advices=[self.joinpoint, self.joinpoint])
        weave(target=min, advices=self.joinpoint)

        min(5, 2)

        self.assertEqual(self.count, 3)

        unweave(min)

        min(5, 2)

        self.assertEqual(self.count, 3)

    def test_method(self):

        class A():
            def __init__(self):
                pass

            def a(self):
                pass

        weave(target=A.a, advices=[self.joinpoint, self.joinpoint], ctx=A)
        weave(target=A, advices=self.joinpoint, pointcut='__init__', ctx=A)
        weave(target=A.__init__, advices=self.joinpoint, ctx=A)

        a = A()
        a.a()

        self.assertEqual(self.count, 4)

        unweave(A.a, ctx=A)
        unweave(A, ctx=A)

        A()
        a.a()

        self.assertEqual(self.count, 4)

    def test_function(self):

        def f():
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f()

        self.assertEqual(self.count, 3)

        unweave(f)

        f()

        self.assertEqual(self.count, 3)

    def test_lambda(self):

        f = lambda: None

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f()

        self.assertEqual(self.count, 3)

        unweave(f)

        f()

        self.assertEqual(self.count, 3)

    def test_function_args(self):

        def f(a):
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f(1)

        self.assertEqual(self.count, 3)

        unweave(f)

        f(1)

        self.assertEqual(self.count, 3)

    def test_function_varargs(self):

        def f(*args):
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f()

        self.assertEqual(self.count, 3)

        unweave(f)

        f()

        self.assertEqual(self.count, 3)

    def test_function_args_varargs(self):

        def f(a, **args):
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)

        f(1)

        self.assertEqual(self.count, 3)

        unweave(f)

        f(1)

        self.assertEqual(self.count, 3)

    def test_function_kwargs(self):

        def f(**kwargs):
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f()

        self.assertEqual(self.count, 3)

        unweave(f)

        f()

        self.assertEqual(self.count, 3)

    def test_function_args_kwargs(self):

        def f(a, **args):
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f(1)

        self.assertEqual(self.count, 3)

        unweave(f)

        f(1)

        self.assertEqual(self.count, 3)

    def test_function_args_varargs_kwargs(self):

        def f(a, *args, **kwargs):
            pass

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)
        f(1)

        self.assertEqual(self.count, 3)

        unweave(f)

        f(1)

        self.assertEqual(self.count, 3)

    def _assert_class(self, cls):
        """
        Run assertion tests on input cls
        """

        weave(target=cls, advices=[self.joinpoint, self.joinpoint])
        weave(target=cls, advices=self.joinpoint, pointcut='__init__')
        weave(target=cls.B, advices=[self.joinpoint, self.joinpoint])
        weave(target=cls.B, advices=self.joinpoint, pointcut='__init__')
        weave(target=cls.C, advices=[self.joinpoint, self.joinpoint])
        weave(target=cls.C, advices=self.joinpoint, pointcut='__init__')

        cls()
        self.assertEqual(self.count, 3)
        cls.B()
        self.assertEqual(self.count, 6)
        cls.C()
        self.assertEqual(self.count, 9)

        unweave(cls)

        cls()

        self.assertEqual(self.count, 9)

        unweave(cls.B)

        cls.B()

        self.assertEqual(self.count, 9)

        unweave(cls.C)

        cls.C()

        self.assertEqual(self.count, 9)

    def test_class(self):

        class A(object):

            class B(object):
                def __init__(self):
                    pass

            class C(object):
                pass

            def __init__(self):
                pass

        self._assert_class(A)

    def test_namespace(self):

        class A:

            class B:
                def __init__(self):
                    pass

            class C:
                pass

            def __init__(self):
                pass

        self._assert_class(A)

    def test_multi(self):

        count = 5

        f = lambda: None

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)

        for i in range(count):
            f()

        self.assertEqual(self.count, 3 * count)

        unweave(f)

        for i in range(count):
            f()

        self.assertEqual(self.count, 3 * count)

        weave(target=f, advices=[self.joinpoint, self.joinpoint])
        weave(target=f, advices=self.joinpoint)

        for i in range(count):
            f()

        self.assertEqual(self.count, 6 * count)

        unweave(f)

        for i in range(count):
            f()

        self.assertEqual(self.count, 6 * count)

    def test_ttl(self):

        def test():
            pass

        weave(target=test, advices=self.joinpoint, ttl=0.1)

        test()

        self.assertEqual(self.count, 1)

        sleep(0.2)

        test()

        self.assertEqual(self.count, 1)

    def test_cancel_ttl(self):

        def test():
            pass

        _, timer = weave(target=test, advices=self.joinpoint, ttl=0.1)

        timer.cancel()

        sleep(0.2)

        test()

        self.assertEqual(self.count, 1)

    def test_inheritance(self):

        class BaseTest:
            def test(self):
                pass

        class Test(BaseTest):
            pass

        self.assertEqual(BaseTest.test, Test.test)

        weave(ctx=Test, target=Test.test, advices=lambda x: None)

        self.assertNotEqual(BaseTest.test, Test.test)

        unweave(ctx=Test, target=Test.test)

        self.assertEqual(BaseTest.test, Test.test)

    def test_inherited_method(self):

        self.count = 0

        class BaseTest:
            def __init__(self, testcase):
                self.testcase = testcase

            def test(self):
                self.testcase.count += 1

        class Test(BaseTest):
            pass

        basetest = BaseTest(self)

        test = Test(self)

        weave(ctx=Test, target=Test.test, advices=lambda x: None)

        basetest.test()

        self.assertEqual(self.count, 1)

        test.test()

        self.assertEqual(self.count, 1)

        unweave(ctx=Test, target=Test.test)

        basetest.test()

        self.assertEqual(self.count, 2)

        test.test()

        self.assertEqual(self.count, 3)

    def test_inherited_instance_method(self):

        class BaseTest(object):
            def test(self):
                pass

        self._test_inherited(BaseTest)

    def test_inherited_instance_method_with_container(self):

        class BaseTest:
            def test(self):
                pass

        self._test_inherited(BaseTest)

    def _test_inherited(self, BaseTest):

        self.count = 0

        class Test(BaseTest):
            pass

        def advice(jp):
            self.count += 1
            return jp.proceed()

        self.old_count = 0

        def assertCount(f, increment=0):
            """
            Assert incrementation of count in executing.
            """
            f()
            self.old_count += increment
            self.assertEqual(self.count, self.old_count)

        test = Test()
        test2 = Test()
        basetest = BaseTest()
        basetest2 = BaseTest()

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        weave(target=test.test, advices=advice, ctx=test)

        assertCount(test.test, 1)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        unweave(target=test.test, ctx=test)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        weave(target=basetest.test, advices=advice, ctx=basetest)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test)

        unweave(target=basetest.test, ctx=basetest)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        weave(target=BaseTest.test, advices=advice, ctx=BaseTest)

        assertCount(test.test, 1)
        assertCount(test2.test, 1)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)

        unweave(target=BaseTest.test, ctx=BaseTest)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        weave(target=Test.test, advices=advice, ctx=Test)

        assertCount(test.test, 1)
        assertCount(test2.test, 1)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        unweave(target=Test.test, ctx=Test)

        # weave all
        weave(target=BaseTest.test, advices=advice, ctx=BaseTest)
        weave(target=Test.test, advices=advice, ctx=Test)
        weave(target=test.test, advices=advice, ctx=test)

        assertCount(test.test, 3)
        assertCount(test2.test, 2)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove middle interceptor
        unweave(target=Test.test, ctx=Test)

        assertCount(test.test, 2)
        assertCount(test2.test, 1)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove the first
        unweave(target=BaseTest.test, ctx=BaseTest)

        assertCount(test.test, 1)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)
        # remove the last
        unweave(target=test.test, ctx=test)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        # weave all in opposite way
        weave(target=test.test, advices=advice, ctx=test)
        weave(target=Test.test, advices=advice, ctx=Test)
        weave(target=BaseTest.test, advices=advice, ctx=BaseTest)

        assertCount(test.test, 3)
        assertCount(test2.test, 2)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove middle interceptor
        unweave(target=Test.test, ctx=Test)

        assertCount(test.test, 2)
        assertCount(test2.test, 1)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove last
        unweave(target=BaseTest.test, ctx=BaseTest)

        assertCount(test.test, 1)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)
        # remove first
        unweave(target=test.test, ctx=test)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        # weave all in random way
        weave(target=Test.test, advices=advice, ctx=Test)
        weave(target=test.test, advices=advice, ctx=test)
        weave(target=BaseTest.test, advices=advice, ctx=BaseTest)

        assertCount(test.test, 3)
        assertCount(test2.test, 2)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove middle interceptor
        unweave(target=Test.test, ctx=Test)

        assertCount(test.test, 2)
        assertCount(test2.test, 1)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove last
        unweave(target=BaseTest.test, ctx=BaseTest)

        assertCount(test.test, 1)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)
        # remove first
        unweave(target=test.test, ctx=test)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        # weave all in random way
        weave(target=Test.test, advices=advice, ctx=Test)
        weave(target=test.test, advices=advice, ctx=test)
        weave(target=BaseTest.test, advices=advice, ctx=BaseTest)

        assertCount(test.test, 3)
        assertCount(test2.test, 2)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove first interceptor
        unweave(target=BaseTest.test, ctx=BaseTest)

        assertCount(test.test, 2)
        assertCount(test2.test, 1)
        assertCount(basetest.test)
        assertCount(basetest2.test)
        # remove second
        unweave(target=Test.test, ctx=Test)

        assertCount(test.test, 1)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)
        # remove last
        unweave(target=test.test, ctx=test)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

        # weave all in random way
        weave(target=Test.test, advices=advice, ctx=Test)
        weave(target=test.test, advices=advice, ctx=test)
        weave(target=BaseTest.test, advices=advice, ctx=BaseTest)

        assertCount(test.test, 3)
        assertCount(test2.test, 2)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove last interceptor
        unweave(target=test.test, ctx=test)

        assertCount(test.test, 2)
        assertCount(test2.test, 2)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove second
        unweave(target=Test.test, ctx=Test)

        assertCount(test.test, 1)
        assertCount(test2.test, 1)
        assertCount(basetest.test, 1)
        assertCount(basetest2.test, 1)
        # remove first
        unweave(target=BaseTest.test, ctx=BaseTest)

        assertCount(test.test)
        assertCount(test2.test)
        assertCount(basetest.test)
        assertCount(basetest2.test)

    def test_instance_method(self):

        class A:
            def __call__(self):
                return 1

        a = A()

        self.assertEqual(
            a.__call__.__func__,
            A.__call__.__func__ if PY2 else A.__call__
        )

        weave(target=a.__call__, advices=lambda ae: None, ctx=a)

        self.assertNotEqual(
            a.__call__.__func__,
            A.__call__.__func__ if PY2 else A.__call__
        )

        result = a.__call__()

        self.assertEqual(result, None)

        unweave(target=a.__call__, ctx=a)

        self.assertEqual(
            a.__call__.__func__,
            A.__call__.__func__ if PY2 else A.__call__
        )

        result = a()

        self.assertEqual(result, 1)

    def test_instance_method_with_pointcut(self):

        class A:
            def __call__(self):
                return 1

        a = A()

        weave(target=a, advices=lambda ae: None)

        result = a()

        self.assertEqual(result, None)

        unweave(target=a)

        result = a()

        self.assertEqual(result, 1)


class WeaveOnTest(UTCase):

    def setUp(self):

        self.count = 0

    def advice(self, joinpoint):
        """
        Default interceptor which increments self count
        """

        self.count += 1
        return joinpoint.proceed()

    def test_builtin(self):

        weave_on(advices=[self.advice, self.advice])(min)
        weave_on(advices=self.advice)(min)

        min(5, 2)

        self.assertEqual(self.count, 3)

        unweave(min)

        min(5, 2)

        self.assertEqual(self.count, 3)

    def test_method(self):

        @weave_on(advices=self.advice, pointcut='__init__')
        class A():

            @weave_on(advices=[self.advice, self.advice])
            def __init__(self):
                pass

            @weave_on(advices=[self.advice, self.advice, self.advice])
            def a(self):
                pass

        a = A()
        a.a()

        self.assertEqual(self.count, 6)

    def test_function(self):

        @weave_on(self.advice)
        @weave_on([self.advice, self.advice])
        def f():
            pass

        f()

        self.assertEqual(self.count, 3)

    def test_lambda(self):

        f = lambda: None

        weave_on(self.advice)(f)
        weave_on([self.advice, self.advice])(f)

        f()

        self.assertEqual(self.count, 3)

    def test_function_args(self):

        @weave_on(self.advice)
        @weave_on([self.advice, self.advice])
        def f(a):
            pass

        f(1)

        self.assertEqual(self.count, 3)

    def test_function_varargs(self):

        @weave_on(self.advice)
        @weave_on([self.advice, self.advice])
        def f(*args):
            pass

        f()

        self.assertEqual(self.count, 3)

    def test_function_args_varargs(self):

        @weave_on(self.advice)
        @weave_on([self.advice, self.advice])
        def f(a, **args):
            pass

        f(1)

        self.assertEqual(self.count, 3)

    def test_function_kwargs(self):

        @weave_on([self.advice, self.advice])
        @weave_on(self.advice)
        def f(**kwargs):
            pass

        f()

        self.assertEqual(self.count, 3)

    def test_function_args_kwargs(self):

        @weave_on(self.advice)
        @weave_on([self.advice, self.advice])
        def f(a, **args):
            pass

        f(1)

        self.assertEqual(self.count, 3)

    def test_function_args_varargs_kwargs(self):

        @weave_on(self.advice)
        @weave_on([self.advice, self.advice])
        def f(a, *args, **kwargs):
            pass

        f(1)

        self.assertEqual(self.count, 3)

    def _assert_class(self, cls):
        """
        Run assertion tests on input cls
        """

        weave_on(advices=[self.advice, self.advice])(cls)
        weave_on(advices=self.advice, pointcut='__init__')(cls)
        weave_on(advices=[self.advice, self.advice])(cls.B)
        weave_on(advices=self.advice, pointcut='__init__')(cls.B)
        weave_on(advices=[self.advice, self.advice])(cls.C)
        weave_on(advices=self.advice, pointcut='__init__')(cls.C)

        cls()
        cls.B()
        cls.C()

        self.assertEqual(self.count, 9)

    def test_class(self):

        class A(object):

            class B(object):
                def __init__(self):
                    pass

            class C(object):
                pass

            def __init__(self):
                pass

        self._assert_class(A)

    def test_namespace(self):

        class A:

            class B:
                def __init__(self):
                    pass

            class C:
                pass

            def __init__(self):
                pass

        self._assert_class(A)

    def test_multi(self):

        count = 5

        f = lambda: None

        weave_on(advices=[self.advice, self.advice])(f)
        weave_on(advices=self.advice)(f)

        for i in range(count):
            f()

        self.assertEqual(self.count, 3 * count)

        unweave(f)

        for i in range(count):
            f()

        self.assertEqual(self.count, 3 * count)

        weave_on(advices=[self.advice, self.advice])(f)
        weave_on(advices=self.advice)(f)

        for i in range(count):
            f()

        self.assertEqual(self.count, 6 * count)

        unweave(f)

        for i in range(count):
            f()

        self.assertEqual(self.count, 6 * count)

    def test_ttl(self):

        def test():
            pass

        weave_on(advices=self.advice, ttl=0.1)(test)

        test()

        sleep(0.2)

        test()

        self.assertEqual(self.count, 1)


if __name__ == '__main__':
    main()
