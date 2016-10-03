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

from six import PY3, PY2

from ..core import (
    Joinpoint,
    is_intercepted, get_intercepted,
    _apply_interception, _unapply_interception,
    _get_function, find_ctx, super_method
)

from types import MethodType, FunctionType

from inspect import isclass


class SuperMethodTest(UTCase):
    """
    Test super method function.
    """

    def _test_class(self, BaseTest):

        class Test(BaseTest):
            pass

        class FinalTest(Test):
            pass

        finaltest = FinalTest()

        super_elt, super_ctx = super_method(name='test', ctx=finaltest)

        self.assertIs(
            super_elt.__func__ if PY2 else super_elt,
            FinalTest.test.__func__ if PY2 else FinalTest.test
        )
        self.assertIs(super_ctx, FinalTest)

        super_elt, super_ctx = super_method(name='test', ctx=FinalTest)

        self.assertIs(
            super_elt.__func__ if PY2 else super_elt,
            Test.test.__func__ if PY2 else Test.test
        )
        self.assertIs(super_ctx, Test)

        super_elt, super_ctx = super_method(name='test', ctx=Test)

        self.assertIs(
            super_elt.__func__ if PY2 else super_elt,
            BaseTest.test.__func__ if PY2 else BaseTest.test
        )
        self.assertIs(super_ctx, BaseTest)

        super_elt, super_ctx = super_method(name='test', ctx=BaseTest)

        self.assertIs(super_elt, None)
        self.assertIs(super_ctx, None)

    def test_namespace(self):

        class BaseTest:
            def test(self):
                pass

        self._test_class(BaseTest)

    def test_class(self):

        class BaseTest(object):
            def test(self):
                pass

        self._test_class(BaseTest)


class FindCTXTest(UTCase):
    """
    Test find_ctx function.
    """

    def test_none(self):

        ctx = find_ctx(None)
        self.assertIsNone(ctx)

    def test_class_method(self):

        ctx = find_ctx(FindCTXTest.test_class_method)
        if PY2:
            self.assertIs(ctx, FindCTXTest)
        elif PY3:
            self.assertIsNone(ctx)
        else:
            raise NotImplementedError('Not implemented for python < 2 or > 3')

    def test_instance_method(self):

        ctx = find_ctx(self.test_instance_method)
        self.assertIs(ctx, self)

    def test_function(self):

        def test():
            pass

        ctx = find_ctx(test)
        self.assertIsNone(ctx)

    def test_namespace(self):

        class A:
            pass

        ctx = find_ctx(A)
        self.assertIs(ctx, A)

    def test_class(self):

        class A(object):
            pass

        ctx = find_ctx(A)
        self.assertIs(ctx, A)


class JoinpointProceedingTest(UTCase):
    """
    Test to apply joinpoint poincut on any kinf of elements and parameters.
    """

    def setUp(self):
        """
        Create a default joinpoint with self joinpoint_proceeding such as the
        only one advice and an integer count attribute equals 0.
        """

        self.joinpoint = Joinpoint(advices=[self.joinpoint_proceeding])

        self.count = 0

    def joinpoint_proceeding(self, jp):
        """
        Base joinpoint proceeding which increments self count and return the
        proceeding result.

        :param Joinpoint jp: joinpoint to proceed.
        :return: joinpoint proceeding.
        """
        self.count += 1

        result = jp.proceed()

        self.count += 2

        return result

    def _test_joinpoint_proceeding(
        self,
        target, ctx=None, expected_result=None, elt=None,
        args=(), kwargs={}
    ):
        """
        Do a serie of tests on joinpoint proceeding related to input target.

        All targets may return result if not None, or themself.

        :param callable target: callable target to intercept.
        :param expected_result: target call result.
        :param elt: element to call if not None.
        :param args: target call var args.
        :param kwargs: target call keywords.
        """

        # check init context
        self.assertEqual(self.count, 0)
        # start to apply pointcut on joinpoint
        self.joinpoint.set_target(target, ctx=ctx)
        # update target with interception
        jp_interception = self.joinpoint._interception
        # call target
        result = jp_interception(*args, **kwargs)
        # compare result with target
        if expected_result is None and not isclass(target):
            expected_result = JoinpointProceedingTest
        self.assertEqual(result, expected_result)
        # compare count with before and after interception
        self.assertEqual(self.count, 3)
        # do the test a second time
        result = jp_interception(*args, **kwargs)
        # compare result with interception
        self.assertEqual(result, expected_result)
        # compare count with before and after interception
        self.assertEqual(self.count, 6)

    def test_namespace(self):
        """
        Test to intercept a namespace.
        """

        class Test:
            pass

        self._test_joinpoint_proceeding(
            target=Test, ctx=Test, args=[Test()]
        )

    def test_instance_method(self):
        """
        Test to intercept a method
        """

        class Test(object):
            def test(self):
                return JoinpointProceedingTest

        test = Test()

        self._test_joinpoint_proceeding(test.test)

    def test_builtin(self):
        """
        Test to intercept a builtin.
        """

        self._test_joinpoint_proceeding(max, expected_result=3, args=[2, 3])

    def test_function(self):
        """
        Test to intercept a function.
        """

        def test():
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test)

    def test_function_params(self):
        """
        Test to intercept a function with params.
        """

        def test(a, b):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test, args=[1, 2])

    def test_function_default_params(self):
        """
        Test to intercept a function with default params.
        """

        def test(a, b=None):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test, kwargs={'a': 1})

    def test_function_closure(self):
        """
        Test to intercept a function with closure.
        """

        closure = 1

        def test():
            return closure

        self._test_joinpoint_proceeding(test, expected_result=closure)

    def test_function_args(self):
        """
        Test to intercept a function with var args.
        """

        def test(*args):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test)

    def test_function_args_params(self):
        """
        Test to intercept a function with var args and params.
        """

        def test(a, b=1, *args):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test, args=[1])

    def test_function_args_params_closure(self):
        """
        Test to intercept a function with var args and params an closure.
        """

        closure = 0

        def test(a, b=1, *args):
            return closure

        self._test_joinpoint_proceeding(
            test, args=[2], expected_result=closure)

    def test_function_kwargs(self):
        """
        Test to intercept a function with kwargs.
        """

        def test(**kwargs):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test)

    def test_function_kwargs_params(self):
        """
        Test to intercept a function with kwargs and params.
        """

        def test(a, b=1, **kwargs):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test, kwargs={'a': 1})

    def test_function_kwargs_params_closure(self):
        """
        Test to intercept a function with kwargs and params.
        """

        closure = 0

        def test(a, b=1, **kwargs):
            return closure

        self._test_joinpoint_proceeding(test, expected_result=0, args=[1])

    def test_function_args_kwargs(self):
        """
        Test to intercept a function with var args and kwargs.
        """

        def test(*args, **kwargs):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test)

    def test_function_args_kwargs_params(self):
        """
        Test to intercept a function with var args, kwargs and params.
        """

        def test(a, b=1, *args, **kwargs):
            return JoinpointProceedingTest

        self._test_joinpoint_proceeding(test, args=[1])

    def test_function_args_kwargs_params_closure(self):
        """
        Test to intercept a function with var args, kwargs, params and closure.
        """

        closure = 0

        def test(a, b=1, *args, **kwargs):
            return closure

        self._test_joinpoint_proceeding(
            test, args=[1], expected_result=closure)


class JoinpointTest(UTCase):

    def setUp(self):

        self.count = 0

        def a():
            self.count += 1
            return self.count

        self.joinpoint = Joinpoint(target=a, advices=[])

    def test_execution(self):

        result = self.joinpoint.start()

        self.assertEqual(result, 1)

    def test_execution_twice(self):

        result = self.joinpoint.start()

        self.assertEqual(result, 1)

        result = self.joinpoint.start()

        self.assertEqual(result, 2)

    def test_add_advices(self):

        def advice(joinpoint):

            proceed = joinpoint.proceed()

            return proceed, 3

        self.joinpoint._advices = [advice]

        result = self.joinpoint.start()

        self.assertEqual(result, (1, 3))

        result = self.joinpoint.start()

        self.assertEqual(result, (2, 3))


class GetFunctionTest(UTCase):

    def test_class(self):

        class A(object):
            pass

        func = _get_function(A)

        to_compare = A.__init__
        if hasattr(to_compare, '__func__'):
            to_compare = to_compare.__func__

        self.assertEqual(func, to_compare)

    def test_namespace(self):

        class A:
            pass

        func = _get_function(A)
        to_compare = A.__init__
        if PY2:
            to_compare = to_compare.__func__
        self.assertEqual(func, to_compare)

    def test_builtin(self):

        _max = _get_function(max)

        self.assertIs(_max, max)

    def test_method(self):

        class A:
            def method(self):
                pass

        func = _get_function(A.method)

        _func = A.method if PY3 else A.method.__func__

        self.assertIs(func, _func)

    def test_instancemethod(self):

        class A:
            def method(self):
                pass

        a = A()

        func = _get_function(a.method)

        self.assertIs(func, a.method.__func__)
        _func = A.method if PY3 else A.method.__func__
        self.assertIs(func, _func)

    def test_function(self):

        def function():
            pass

        func = _get_function(function)

        self.assertIs(func, function)

    def test_call(self):

        class A:
            def __call__(self):
                pass

        a = A()

        func = _get_function(a)

        self.assertEqual(func, a.__call__.__func__)


class ExecCtxTest(UTCase):
    """Test the execution context."""

    def _cmp(self, val, exec_ctx=None, start_exec_ctx=None):

        def advice(jp):

            jp.exec_ctx['test'] = jp.exec_ctx.setdefault('test', 0) + 1

        joinpoint = Joinpoint(
            advices=[advice], target=lambda: None, exec_ctx=exec_ctx
        )

        joinpoint.start(exec_ctx=start_exec_ctx)

        self.assertEqual(joinpoint.exec_ctx['test'], val)

        joinpoint.start(exec_ctx=start_exec_ctx)

        self.assertEqual(joinpoint.exec_ctx['test'], val)

        return joinpoint

    def test_no_exc_ctx(self):

        self._cmp(1)

    def test_exec_ctx(self):

        self._cmp(2, {'test': 1})

    def test_one_per_execution(self):

        joinpoint = self._cmp(2, {'test': 0, 'final': True}, {'test': 1})

        self.assertTrue(joinpoint.exec_ctx['final'])


class ApplyInterceptionTest(UTCase):

    def test_function(self):

        def function():
            pass

        __code__ = function.__code__

        self.assertFalse(is_intercepted(function))
        self.assertEqual(get_intercepted(function), (None, None))

        interception, intercepted, ctx = _apply_interception(
            function, lambda x: None)

        self.assertTrue(isinstance(interception, FunctionType))

        self.assertIs(interception, function)
        self.assertEqual(get_intercepted(interception), (intercepted, ctx))
        self.assertTrue(is_intercepted(interception))
        self.assertIsNot(interception.__code__, __code__)
        self.assertIs(intercepted.__code__, __code__)

        _unapply_interception(function)

        self.assertFalse(is_intercepted(function))
        self.assertIs(interception, function)
        self.assertIs(function.__code__, __code__)
        self.assertEqual(get_intercepted(function), (None, None))

    def test_method(self):
        class A(object):
            def method(self):
                pass

        self.assertEqual(get_intercepted(A.method), (None, None))
        self.assertFalse(is_intercepted(A.method))

        interception_fn = lambda: None

        interception, intercepted, ctx = _apply_interception(
            target=A.method,
            interception_fn=interception_fn,
            ctx=A)

        joinpoint_type = FunctionType if PY3 else MethodType

        self.assertTrue(isinstance(interception, joinpoint_type))
        self.assertTrue(is_intercepted(A.method))
        self.assertEqual(interception, A.method)
        self.assertEqual((intercepted, ctx), get_intercepted(A.method))

        _unapply_interception(target=A.method, ctx=A)

        self.assertFalse(is_intercepted(A.method))
        self.assertEqual(get_intercepted(A.method), (None, None))

    def test_class_container(self):
        class A(object):
            def method(self):
                pass

        class B(A):
            pass

        self.assertEqual(A.method, B.method)

        _apply_interception(
            target=B.method,
            interception_fn=lambda: None,
            ctx=B)

        self.assertNotEqual(A.method, B.method)

        _unapply_interception(
            target=B.method, ctx=B)

        self.assertEqual(A.method, B.method)

    def test_instance(self):
        class A(object):
            def method(self):
                pass

        class B(A):
            pass

        a = A()
        b = B()

        self.assertEqual(a.__dict__, b.__dict__)

        _apply_interception(
            target=b.method,
            interception_fn=lambda: None)

        self.assertNotEqual(a.__dict__, b.__dict__)

        _unapply_interception(target=b.method)

        self.assertEqual(a.__dict__, b.__dict__)

    def test_builtin(self):

        function = min

        self.assertEqual(get_intercepted(min), (None, None))
        self.assertFalse(is_intercepted(min))

        interception, intercepted, ctx = _apply_interception(min, lambda: None)

        self.assertTrue(isinstance(interception, FunctionType))

        self.assertTrue(is_intercepted(min))
        self.assertIs(interception, min)
        self.assertIsNot(min, function)
        self.assertEqual(get_intercepted(min), (intercepted, None))

        _unapply_interception(interception)

        self.assertFalse(is_intercepted(min))
        self.assertIsNot(interception, min)
        self.assertIs(min, function)
        self.assertEqual(get_intercepted(min), (None, None))

    def test_inheritance(self):
        """
        Test interception in an inheritance context.
        """

        def new_test():
            def test():
                pass
            return test

        def test():
            pass

        class Test0(object):
            def test(self):
                pass

        class Test1(Test0):
            pass

        class Test2(Test1):
            def test(self):
                pass

        test0 = Test0()
        test1 = Test1()
        test1.test = test
        test2 = Test2()

        # check for inherited method
        _apply_interception(Test1.test, ctx=Test1, interception_fn=new_test())
        self.assertFalse(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertFalse(is_intercepted(Test2.test))
        self.assertFalse(is_intercepted(test0.test))
        self.assertFalse(is_intercepted(test1.test))
        self.assertFalse(is_intercepted(test2.test))

        # check for base method
        _apply_interception(Test0.test, ctx=Test0, interception_fn=new_test())
        self.assertTrue(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(Test1.test))
        self.assertIsNot(Test0.test.__dict__, Test1.test.__dict__)
        self.assertFalse(is_intercepted(Test2.test))
        self.assertTrue(is_intercepted(test0.test))
        self.assertIs(_get_function(Test0.test), _get_function(test0.test))
        self.assertFalse(is_intercepted(test1.test))
        self.assertFalse(is_intercepted(test2.test))

        # check for overriden method
        _apply_interception(Test2.test, ctx=Test2, interception_fn=new_test())
        self.assertTrue(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(Test1.test))
        self.assertIsNot(Test0.test.__dict__, Test1.test.__dict__)
        self.assertTrue(is_intercepted(Test2.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(Test2.test))
        self.assertIsNot(Test1.test.__dict__, Test2.test.__dict__)
        self.assertTrue(is_intercepted(test0.test))
        self.assertIs(_get_function(Test0.test), _get_function(test0.test))
        self.assertFalse(is_intercepted(test1.test))
        self.assertTrue(is_intercepted(test2.test))
        self.assertIs(_get_function(Test2.test), _get_function(test2.test))

        # check for inherited instance method
        _apply_interception(test0.test, interception_fn=new_test())
        self.assertTrue(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(Test1.test))
        self.assertIsNot(Test0.test.__dict__, Test1.test.__dict__)
        self.assertTrue(is_intercepted(Test2.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(Test2.test))
        self.assertIsNot(Test1.test.__dict__, Test2.test.__dict__)
        self.assertTrue(is_intercepted(test0.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(test0.test))
        self.assertIsNot(test0.test.__dict__, Test0.test.__dict__)
        self.assertFalse(is_intercepted(test1.test))
        self.assertTrue(is_intercepted(test2.test))
        self.assertIs(_get_function(Test2.test), _get_function(test2.test))

        # check for overriden instance method
        _apply_interception(test1.test, interception_fn=new_test())
        self.assertTrue(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(Test1.test))
        self.assertIsNot(Test0.test.__dict__, Test1.test.__dict__)
        self.assertTrue(is_intercepted(Test2.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(Test2.test))
        self.assertIsNot(Test1.test.__dict__, Test2.test.__dict__)
        self.assertTrue(is_intercepted(test0.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(test0.test))
        self.assertIsNot(test0.test.__dict__, Test0.test.__dict__)
        self.assertTrue(is_intercepted(test1.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(test1.test))
        self.assertIsNot(test1.test.__dict__, Test1.test.__dict__)
        self.assertTrue(is_intercepted(test2.test))
        self.assertIs(_get_function(Test2.test), _get_function(test2.test))

        # check to unapply base class method
        _unapply_interception(Test0.test, ctx=Test0)
        self.assertFalse(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(Test1.test))
        self.assertIsNot(Test0.test.__dict__, Test1.test.__dict__)
        self.assertTrue(is_intercepted(Test2.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(Test2.test))
        self.assertIsNot(Test1.test.__dict__, Test2.test.__dict__)
        self.assertTrue(is_intercepted(test0.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(test0.test))
        self.assertIsNot(test0.test.__dict__, Test0.test.__dict__)
        self.assertTrue(is_intercepted(test1.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(test1.test))
        self.assertIsNot(test1.test.__dict__, Test1.test.__dict__)
        self.assertTrue(is_intercepted(test2.test))
        self.assertIs(_get_function(Test2.test), _get_function(test2.test))

        # check to unapply overriden method
        _unapply_interception(Test2.test, ctx=Test2)
        self.assertFalse(is_intercepted(Test0.test))
        self.assertTrue(is_intercepted(Test1.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(Test1.test))
        self.assertIsNot(Test0.test.__dict__, Test1.test.__dict__)
        self.assertFalse(is_intercepted(Test2.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(Test2.test))
        self.assertIsNot(Test1.test.__dict__, Test2.test.__dict__)
        self.assertTrue(is_intercepted(test0.test))
        self.assertIsNot(_get_function(Test0.test), _get_function(test0.test))
        self.assertIsNot(test0.test.__dict__, Test0.test.__dict__)
        self.assertTrue(is_intercepted(test1.test))
        self.assertIsNot(_get_function(Test1.test), _get_function(test1.test))
        self.assertIsNot(test1.test.__dict__, Test1.test.__dict__)
        self.assertFalse(is_intercepted(test2.test))

        _unapply_interception(test0.test)
        self.assertFalse(is_intercepted(test0.test))

        # check to unapply a method class which is overriden by an instance
        _unapply_interception(Test1.test, ctx=Test1)
        self.assertFalse(is_intercepted(Test1.test))
        self.assertTrue(is_intercepted(test1.test))

        # check to unapply an instance method
        _unapply_interception(test1.test)
        self.assertFalse(is_intercepted(test1.test))


if __name__ == '__main__':
    main()
