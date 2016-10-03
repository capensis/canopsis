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

from ..interception import Interceptor
from ..call import Types, Curried, Retries, Memoize


class CallTests(UTCase):

    def _assertCall(self, f, *args, **kwargs):

        self.assertRaises(Exception, f, *args, **kwargs)

    def testTypes(self):

        @Types(rtype=int)
        def a(p=None):
            return p

        a(1)
        a()
        self._assertCall(a, '')

        @Types(rtype=Types.NotNone(int))
        def b(p=None):
            return p

        b(2)
        self._assertCall(b, '')
        self._assertCall(b)

        @Types(rtype=[int])
        def c(p=None):
            return [2, 3] \
                if p == 1 else \
                None if p == 2 else \
                [] if p == 3 else \
                [2, None] if p == 4 else [2, '']

        c(1)
        c(2)
        c(3)
        c(4)
        self._assertCall(c)

        @Types([Types.NotNone(int)])
        def d(p=None):
            return [2, 3] \
                if p == 1 else \
                None if p == 2 else \
                [] if p == 3 else \
                [2, None] if p == 4 else [2, '']

        d(1)
        d(2)
        d(3)
        self._assertCall(d, 4)
        self._assertCall(d)

        @Types(Types.NotEmpty([int]))
        def e(p=None):
            return [2, 3] \
                if p == 1 else \
                None if p == 2 else \
                [] if p == 3 else \
                [2, None] if p == 4 else [2, '']

        e(1)
        e(2)
        self._assertCall(e, 3)
        e(4)
        self._assertCall(e)

        @Types(Types.NotEmpty([Types.NotNone(int)]))
        def f(p=None):
            return [2, 3] \
                if p == 1 else \
                None if p == 2 else \
                [] if p == 3 else \
                [2, None] if p == 4 else [2, '']

        f(1)
        f(2)
        self._assertCall(f, 3)
        self._assertCall(f, 4)
        self._assertCall(f)

        @Types(Types.NotNone(Types.NotEmpty([int])))
        def g(p=None):
            return [2, 3] \
                if p == 1 else \
                None if p == 2 else \
                [] if p == 3 else \
                [2, None] if p == 4 else [2, '']

        g(1)
        self._assertCall(g, 2)
        self._assertCall(g, 3)
        g(4)
        self._assertCall(g)

        @Types(ptypes={'p': int})
        def f(p=None):
            pass

        a(1)
        a()
        self._assertCall(a, '')

        @Types(ptypes={'p': Types.NotNone(int)})
        def g(p=None):
            pass

        g(1)
        self._assertCall(b)
        self._assertCall(b, '')

        @Types(ptypes={'p': [int]})
        def h(p=None):
            pass

        h([1])
        h([2, 3])
        h([None])
        h([])
        h()
        self._assertCall(h, [2, ''])

        @Types(ptypes={'p': [Types.NotNone(int)]})
        def i(p=None):
            pass

        i([1])
        i([2, 3])
        self._assertCall(i, [None])
        i([])
        i()
        self._assertCall(i, [2, ''])

        @Types(ptypes={'p': Types.NotEmpty([int])})
        def d(p=None):
            pass

        d()
        d([1])
        d([2, 3])
        self._assertCall(d, 3)
        self._assertCall(d, [])
        self._assertCall(d, [2, ''])

    def testInterceptor(self):

        def interception(target, args, kwargs):
            raise Exception()

        @Interceptor(interception)
        def a(a, b=2):
            pass

        self._assertCall(a, a=None)
        self._assertCall(a)

    def testCurried(self):

        @Curried()
        def a(b, c=None):
            return 2

        result = a()

        self.assertTrue(isinstance(result, Curried.CurriedResult))

        self.assertEqual(a(None), 2)

        result.curried.args = []

        self.assertEqual(a(b=None), 2)

        self.assertEqual(a(b=None, c=None), 2)

        self.assertTrue(isinstance(a(None), Curried.CurriedResult))

        pass

    def testRetries(self):

        global count
        count = 10

        @Retries(10, delay=0, backoff=0)
        def a():
            global count
            count -= 1
            if count > 0:
                raise Exception()
            else:
                return ""

        result = a()

        self.assertTrue(count == 0)
        self.assertTrue(result == "")


class MemoizeTest(UTCase):
    """Test the memoize annotation."""

    def setUp(self):

        self.memoize = Memoize()

        self.n = 0

        def func(*args, **kwargs):

            self.n += 1

            return self.n

        self.func = self.memoize(func)

    def test_empty(self):

        result = self.func()
        self.assertEqual(result, 1)

        params = self.memoize.getparams(1)
        self.assertEqual(params, ((), {}))

        result = self.func()
        self.assertEqual(result, 1)

    def test_maxsize(self):

        self.memoize.max_size = 2

        result = self.func()
        self.assertEqual(result, 1)

        params = self.memoize.getparams(1)
        self.assertEqual(params, ((), {}))

        result = self.func()
        self.assertEqual(result, 1)

        params = self.memoize.getparams(1)
        self.assertEqual(params, ((), {}))

        result = self.func(1, 2, a=3)
        self.assertEqual(result, 2)

        params = self.memoize.getparams(2)
        self.assertEqual(params, ((1, 2), {'a': 3}))

        result = self.func(3, 4, b=5)
        self.assertEqual(result, 3)

        self.assertRaises(ValueError, self.memoize.getparams, 3)


if __name__ == '__main__':
    main()
