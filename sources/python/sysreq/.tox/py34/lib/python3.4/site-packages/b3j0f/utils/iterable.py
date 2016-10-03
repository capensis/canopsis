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

"""Provides tools to manage iterable types."""

from __future__ import absolute_import

__all__ = [
    'isiterable', 'ensureiterable', 'first', 'last', 'itemat', 'sliceit',
    'hashiter'
]

from collections import Iterable


def isiterable(element, exclude=None):
    """Check whatever or not if input element is an iterable.

    :param element: element to check among iterable types.
    :param type/tuple exclude: not allowed types in the test.

    :Example:

    >>> isiterable({})
    True
    >>> isiterable({}, exclude=dict)
    False
    >>> isiterable({}, exclude=(dict,))
    False
    """

    # check for allowed type
    allowed = exclude is None or not isinstance(element, exclude)
    result = allowed and isinstance(element, Iterable)

    return result


def ensureiterable(value, iterable=list, exclude=None):
    """Convert a value into an iterable if it is not.

    :param object value: object to convert
    :param type iterable: iterable type to apply (default: list)
    :param type/tuple exclude: types to not convert

    :Example:

    >>> ensureiterable([])
    []
    >>> ensureiterable([], iterable=tuple)
    ()
    >>> ensureiterable('test', exclude=str)
    ['test']
    >>> ensureiterable('test')
    ['t', 'e', 's', 't']
    """

    result = value

    if not isiterable(value, exclude=exclude):
        result = [value]
        result = iterable(result)

    else:
        result = iterable(value)

    return result


def first(iterable, default=None):
    """Try to get input iterable first item or default if iterable is empty.

    :param Iterable iterable: iterable to iterate on. Must provide the method
        __iter__.
    :param default: default value to get if input iterable is empty.
    :raises TypeError: if iterable is not an iterable value.

    :Example:

    >>> first('tests')
    't'
    >>> first('', default='test')
    'test'
    >>> first([])
    None
    """

    result = default

    # start to get the iterable iterator (raises TypeError if iter)
    iterator = iter(iterable)
    # get first element
    try:
        result = next(iterator)
    except StopIteration: # if no element exist, result equals default
        pass

    return result

def last(iterable, default=None):
    """Try to get the last iterable item by successive iteration on it.

    :param Iterable iterable: iterable to iterate on. Must provide the method
        __iter__.
    :param default: default value to get if input iterable is empty.
    :raises TypeError: if iterable is not an iterable value.

    :Example:

    >>> last('tests')
    's'
    >>> last('', default='test')
    'test'
    >>> last([])
    None"""

    result = default

    iterator = iter(iterable)

    while True:
        try:
            result = next(iterator)

        except StopIteration:
            break

    return result

def itemat(iterable, index):
    """Try to get the item at index position in iterable after iterate on
    iterable items.

    :param iterable: object which provides the method __getitem__ or __iter__.
    :param int index: item position to get.
    """

    result = None

    handleindex = True

    if isinstance(iterable, dict):
        handleindex = False

    else:
        try:
            result = iterable[index]
        except TypeError:
            handleindex = False

    if not handleindex:
        iterator = iter(iterable)

        if index < 0:  # ensure index is positive
            index += len(iterable)

        while index >= 0:
            try:
                value = next(iterator)

            except StopIteration:
                raise IndexError(
                    "{0} index {1} out of range".format(
                        iterable.__class__, index
                    )
                )

            else:
                if index == 0:
                    result = value
                    break
                index -= 1

    return result

def sliceit(iterable, lower=0, upper=None):
    """Apply a slice on input iterable.

    :param iterable: object which provides the method __getitem__ or __iter__.
    :param int lower: lower bound from where start to get items.
    :param int upper: upper bound from where finish to get items.
    :return: sliced object of the same type of iterable if not dict, or specific
        object. otherwise, simple list of sliced items.
    :rtype: Iterable
    """

    if upper is None:
        upper = len(iterable)

    try:
        result = iterable[lower: upper]

    except TypeError:  # if iterable does not implement the slice method
        result = []

        if lower < 0:  # ensure lower is positive
            lower += len(iterable)

        if upper < 0:  # ensure upper is positive
            upper += len(iterable)

        if upper > lower:
            iterator = iter(iterable)

            for index in range(upper):
                try:
                    value = next(iterator)

                except StopIteration:
                    break

                else:
                    if index >= lower:
                        result.append(value)

    iterablecls = iterable.__class__
    if not(isinstance(result, iterablecls) or issubclass(iterablecls, dict)):
        try:
            result = iterablecls(result)

        except TypeError:
            pass

    return result


def hashiter(iterable):
    """Try to hash input iterable in doing the sum of its content if not
    hashable.

    Hash method on not iterable depends on type:

    hash(iterable.__class__) + ...

        - dict: sum of (hash(key) + 1) * (hash(value) + 1).
        - Otherwise: sum of (pos + 1) * (hash(item) + 1)."""

    result = 0

    try:
        result = hash(iterable)

    except TypeError:

        result = hash(iterable.__class__)

        isdict = isinstance(iterable, dict)

        for index, entry in enumerate(list(iterable)):
            entryhash = hashiter(entry) + 1

            if isdict:
                entryhash *= hashiter(iterable[entry]) + 1

            else:
                entryhash *= index + 1

            result += entryhash

    return result
