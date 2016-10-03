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

from ..runtime import (
    SAFE_BUILTINS, safe_eval, safe_exec,
    make_constants, bind_all
)

import random


class SafeTestCase(UTCase):
    """Test the function about safe coding."""

    def test_safe_builtins_max(self):
        """Test max in SAFE_BUILTINS."""

        self.assertIn('max', SAFE_BUILTINS['__builtins__'])

    def test_safe_builtins_open(self):
        """Test open not in SAFE_BUILTINS."""

        self.assertNotIn('open', SAFE_BUILTINS['__builtins__'])

    def test_eval_error(self):
        """Test safe eval function with open (error)."""

        self.assertRaises(NameError, safe_eval, 'open')

    def test_eval(self):
        """Test safe eval function with max."""

        res = safe_eval('max')
        self.assertIs(res, max)

    def test_eval_globals(self):
        """Test safe eval function with global open."""

        res = safe_eval('open', {'open': 2})
        self.assertIs(res, 2)

    def test_eval_locals(self):
        """Test safe eval function."""

        res = safe_eval('open', None, {'open': 2})
        self.assertIs(res, 2)

    def test_safe_exec_error(self):
        """Test safe exec function with open (error)."""

        self.assertRaises(NameError, safe_exec, 'res = open')
        self.assertNotIn('res', globals())
        self.assertNotIn('res', locals())

    def test_safe_exec(self):
        """Test safe exec function with max."""

        safe_exec('res = max')
        self.assertNotIn('res', globals())
        self.assertNotIn('res', locals())

    def test_safe_exec_globals(self):
        """Test safe exec function with open in globals."""

        properties = {'open': 2}
        safe_exec('res = open', properties)
        self.assertNotIn('res', globals())
        self.assertNotIn('res', locals())
        self.assertIn('res', properties)
        self.assertIs(properties['res'], 2)

    def test_safe_exec_locals(self):
        """Test safe exec function with open in locals."""

        properties = {'open': 2}
        safe_exec('res = open', None, properties)
        self.assertNotIn('res', globals())
        self.assertNotIn('res', locals())
        self.assertIn('res', properties)

    def test_safe_exec_empty_globals(self):
        """Test safe exec function with empty globals."""

        properties = {}
        safe_exec('res = max', properties)
        self.assertNotIn('res', globals())
        self.assertNotIn('res', locals())
        self.assertIn('res', properties)
        self.assertIs(properties['res'], max)

    def test_safe_exec_empty_locals(self):
        """Test safe exec function with empty locals."""

        properties = {}
        safe_exec('res = max', None, properties)
        self.assertNotIn('res', globals())
        self.assertNotIn('res', locals())
        self.assertIn('res', properties)
        self.assertIs(properties['res'], max)


class MakeConstants(UTCase):

    def setUp(self):
        self.output = []
        self.stoplist = ['range']

    def verbose(self, message):
        """
        Verbose function to apply when using make_constants
        """

        self.output.append(message)

    def sample(self):

        def sample(population, k):
            "Choose k unique random elements from a population sequence."
            if not isinstance(population, (list, tuple, str)):
                raise TypeError('Cannot handle type', type(population))
            n = len(population)
            if not 0 <= k <= n:
                raise ValueError("sample larger than population")
            result = [None] * k
            pool = list(population)
            for i in range(k):         # invariant:  non-selected at [0,n-i)
                j = int(random.random() * (n - i))
                result[i] = pool[j]
                pool[j] = pool[n - i - 1]
            return result

        return sample

    def _test_verbose(self):

        verbose_message = [
            "isinstance --> {0}".format(isinstance),
            "list --> {0}".format(list), "tuple --> {0}".format(tuple),
            "str --> {0}".format(str), "TypeError --> {0}".format(TypeError),
            "type --> {0}".format(type), "len --> {0}".format(len),
            "ValueError --> {0}".format(ValueError),
            "list --> {0}".format(list),
            #"range --> {0}".format(range),
            "int --> {0}".format(int),
            "random --> {0}".format(random),
            "new folded constant:{0}".format((list, tuple, str)),
            "new folded constant:{0}".format(random.random)
        ]

        self.assertEqual(self.output, verbose_message)

    def test_function(self):

        make_constants(
            verbose=self.verbose, stoplist=self.stoplist)(self.sample())

    def test_class(self):

        class A(object):

            pass

        A.a = self.sample()

        bind_all(A, verbose=self.verbose, stoplist=self.stoplist)

        self._test_verbose()


if __name__ == '__main__':
    main()
