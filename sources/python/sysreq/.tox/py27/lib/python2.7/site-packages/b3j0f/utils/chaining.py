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

"""This module aims to provide tools to chaining of calls.

It is inspired from method chaining pattern in embedding objects to chain
methods calls in a dedicated Chaining object. Such method calls return the
Chaining object itself, allowing multiple calls to object methods to be invoked
in a concise statement.
"""

from __future__ import unicode_literals, absolute_import

__all__ = ['Chaining', 'ListChaining']


class Chaining(object):
    """Class which permits to process chaining of routines/attr such as call
    chaining in javascript.

    In order to chain calls, start to embed an object, then call embedded
    object methods which returns all the Chaining object. Finally, if you want
    to get back method results, use Chaining.__getitem__ method where the index
    corresponds to the order of calls.

    :Example:

    >>> chaining = Chaining('example').upper().capitalize()
    >>> chaining[0]
    'EXAMPLE'
    >>> chaining[1]
    'Example'
    >>> chaining[-1]
    'example.'
    >>> chaining[:]
    ['EXAMPLE', 'Example', 'example.']
    >>> chaining._
    'example'
    >>> chaining.__iadd__('.').__iadd__('!')
    >>> chaining._
    'example.!'
    >>> chaining[:]
    ['EXAMPLE', 'Example', 'example.', None, None]
    """

    CONTENT = '_'  #: content attribute name
    RESULTS = '___'  #: chained method results attribute name

    __slots__ = (CONTENT, RESULTS)

    def __init__(self, content):

        super(Chaining, self).__init__()

        self._ = content
        self.___ = []

    def __getitem__(self, key):

        return self.___[key]

    def __getattribute__(self, key):

        if key in Chaining.__slots__:
            # check if key is in self slots
            result = super(Chaining, self).__getattribute__(key)

        else:  # else try to get key from self.content
            attr = getattr(self._, key)
            # embed routine in self._processing_name method
            result = _process_function(self, attr)

        return result


def _process_function(chaining, routine):
    """Chain function which returns a function.

    :param routine: routine to process.
    :return: routine embedding execution function.
    """

    def processing(*args, **kwargs):
        """Execute routine with input args and kwargs and add reuslt in
        chaining.___.

        :param tuple args: routine varargs.
        :param dict kwargs: routine kwargs.
        :return: chaining chaining.
        :rtype: Chaining
        """

        result = routine(*args, **kwargs)
        chaining.___.append(result)

        return chaining

    return processing


class ListChaining(Chaining):
    """Apply chaining on a list of objects.

    According to content length, chaining results are saved in a list where
        values are call result or exception if an exeception occured.

    :Example:

    >>> chaining = ListChaining('example', 'test').upper().capitalize()
    >>> chaining[0]
    ['EXAMPLE', 'TEST']
    >>> chaining[1]
    ['Example', 'Test']
    >>> chaining[-1]
    ['Example', 'Test']
    >>> chaining._
    ['example', 'test']
    >>> chaining += '.'
    >>> chaining._
    ['example.', 'test.']
    >>> chaining[2]
    [None, None]
    >>> chaining[:]
    [['EXAMPLE', 'TEST'], ['Example', 'Test'], [None, None]]
    """

    __slots__ = Chaining.__slots__

    def __init__(self, *content):

        super(ListChaining, self).__init__(content=content)

    def __getattribute__(self, key):

        if key in ListChaining.__slots__:
            result = super(ListChaining, self).__getattribute__(key)

        else:
            # list of routines to execute
            self_content = self._
            routines = [None] * len(self_content)
            # get routines from self.content and input key
            for index, content in enumerate(self_content):
                routine = None
                try:
                    routine = getattr(content, key)
                except AttributeError as excp:
                    # in case of exception, routine is the exception
                    routine = excp
                # in all cases, put routine in routines
                routines[index] = routine
            result = _process_function_list(self, routines)

        return result


def _process_function_list(self, routines):
    """Chain function which returns a function.

    :param routines: routines to process.
    """

    def processing(*args, **kwargs):
        """Execute routines with input args and kwargs and add result in
        chaining.___.

        :param tuple args: routines varargs.
        :param dict kwargs: routines kwargs.
        :return: chaining chaining.
        :rtype: Chaining
        """

        results = [None] * len(routines)
        for index, routine in enumerate(routines):
            if isinstance(routine, Exception):
                result = routine
            else:
                try:
                    result = routine(*args, **kwargs)
                except AttributeError as excp:
                    result = excp
            results[index] = result
        self.___.append(results)

        return self

    return processing
