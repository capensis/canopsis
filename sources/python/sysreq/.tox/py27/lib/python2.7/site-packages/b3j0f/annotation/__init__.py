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

"""Annotation package."""

__all__ = [
    '__version__',  # version
    'Annotation',
    'Synchronized', 'SynchronizedClass',
    'Asynchronous', 'TimeOut', 'Wait', 'Observable',
    'Types', 'types', 'Curried', 'curried', 'Retries',
    'Condition', 'MaxCount', 'Target',
    'Interceptor',
    'PrivateInterceptor', 'CallInterceptor', 'PrivateCallInterceptor',
    'Transform', 'Mixin', 'Deprecated', 'Singleton', 'MethodMixin'
]

from .version import __version__

from .core import Annotation
from .async import (
    Synchronized, SynchronizedClass, Asynchronous, TimeOut, Wait, Observable
)
from .call import Types, types, Curried, curried, Retries
from .check import Condition, MaxCount, Target
from .interception import (
    Interceptor, PrivateInterceptor, CallInterceptor, PrivateCallInterceptor
)
from .oop import (
    Transform, Mixin, Deprecated, Singleton, MethodMixin
)
