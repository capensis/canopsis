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

from ..interception import (
    Interceptor, PrivateInterceptor, PrivateCallInterceptor, CallInterceptor
)


class InterceptorTest(UTCase):
    """
    Base class for other tests
    """

    def get_interceptor(self):

        raise NotImplementedError()

    def get_target(self):

        raise NotImplementedError()

    def setUp(self):

        self.count = 0

        self.interceptor = self.get_interceptor()
        self.target = self.get_target()
        self.interceptor(self.target)

    def interception(self, joinpoint):

        self.count += 1

        return joinpoint.proceed()

    def tearDown(self):
        """
        Delete self.interceptor
        """

        self.interceptor.__del__()
        del self.interceptor


class InterceptionTest(InterceptorTest):
    """
    Test interception
    """

    def get_interceptor(self):

        return Interceptor(interception=self.interception)

    def get_target(self):

        return lambda: None

    def call_target(self):
        """Call target.
        """

        self.target()

    def test_one_one(self):
        """
        Test one annotation and one call
        """

        self.call_target()

        self.assertEqual(self.count, 1)

    def test_one_two(self):
        """
        Test one annotation and two calls
        """
        self.call_target()
        self.call_target()

        self.assertEqual(self.count, 2)

    def test_two_one(self):
        """
        Test two annotation and one call
        """

        self.interceptor(self.target)

        self.call_target()

        self.assertEqual(self.count, 2)

    def test_two_two(self):
        """
        Test two annotations and two calls
        """

        self.interceptor(self.target)

        self.call_target()
        self.call_target()

        self.assertEqual(self.count, 4)

    def test_enable(self):

        self.interceptor.enable = True

        self.call_target()

        self.assertEqual(self.count, 1)

    def test_disable(self):

        self.interceptor.enable = False

        self.call_target()

        self.assertEqual(self.count, 0)

    def test_enables(self):

        Interceptor.set_enable(self.target, enable=True)

        self.call_target()

        self.assertEqual(self.count, 1)

    def test_disables(self):

        Interceptor.set_enable(self.target, enable=False)

        self.call_target()

        self.assertEqual(self.count, 0)


class PrivateInterceptorTest(InterceptorTest):
    """
    Test interception
    """

    class TestPrivateInterceptor(PrivateInterceptor):

        def __init__(self, utcase):

            super(
                PrivateInterceptorTest.TestPrivateInterceptor, self
            ).__init__()

            self.utcase = utcase

        def _interception(self, joinpoint, *args, **kwargs):

            return self.utcase.interception(joinpoint)

    def get_interceptor(self):

        return PrivateInterceptorTest.TestPrivateInterceptor(self)

    def get_target(self):

        return lambda: None


class CallInterceptorTest(InterceptorTest):
    """
    Test interception
    """
    class Test(object):

        def __call__(self):

            pass

    def get_interceptor(self):

        return CallInterceptor(interception=self.interception)

    def get_target(self):

        return CallInterceptorTest.Test

    def call_target(self):

        self.target()()


class PrivateCallInterceptorTest(InterceptorTest):
    """
    Test interception
    """

    class TestCallInterceptor(PrivateCallInterceptor):

        def __init__(self, utcase):

            super(
                PrivateCallInterceptorTest.TestCallInterceptor,
                self
            ).__init__()
            self.utcase = utcase

        def _interception(self, joinpoint):
            return self.utcase.interception(joinpoint)

    class Test(object):

        def __call__(self):
            pass

    def get_interceptor(self):

        return PrivateCallInterceptorTest.TestCallInterceptor(self)

    def get_target(self):

        return CallInterceptorTest.Test

    def call_target(self):

        self.target()()


if __name__ == '__main__':
    main()
