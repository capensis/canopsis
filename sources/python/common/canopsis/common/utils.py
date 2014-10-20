# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

"""
Python utility library.
"""

from importlib import import_module

from inspect import ismodule

from collections import Iterable

from sys import version as PYVER

__RESOLVED_ELEMENTS = {}  #: dictionary of resolved elements by name


def free_cache(path=None):
    """
    Remove an element from cache memory

    :param str path: path to element to remove from cache. If None, remove all
        elements from cache.
    """

    if path is None:
        __RESOLVED_ELEMENTS = {}
    else:
        __RESOLVED_ELEMENTS[path].pop(path, '')


def lookup(path, cached=True):
    """
    Get element reference from input full path element.

    :limitations: does not resolve class method.

    :param str path: full path to a python element.
        Examples:
            - __builtin__.open
            - canopsis.common.utils.lookup

    :para bool cached: if True (by default), use __RESOLVED_ELEMENTS cache
        memory to quickly load elements

    :return: python object which is accessible thourgh input path.
    :rtype: object
    """

    element_in_cache = cached and path in __RESOLVED_ELEMENTS

    result = __RESOLVED_ELEMENTS[path] if element_in_cache else None

    if result is None and path:

        components = path.split('.')
        index = 0
        components_len = len(components)

        module_name = components[0]

        # try to import the first component name
        try:
            result = import_module(module_name)
        except ImportError:
            pass

        if result is not None:

            if components_len > 1:

                index = 1

                # try to import all sub-modules/packages
                try:  # check if name is defined from an external module
                    # find the right module

                    while index < components_len:
                        module_name = '{0}.{1}'.format(
                            module_name, components[index])
                        result = import_module(module_name)
                        index += 1

                except ImportError:
                    # path sub-module content
                    try:

                        while index < components_len:
                            result = getattr(result, components[index])
                            index += 1

                    except AttributeError:
                        raise ImportError(
                            'Wrong path {0} at {1}'.format(
                                path, components[:index]))

            if result is not None and cached:
                __RESOLVED_ELEMENTS[path] = result

        else:  # get relative object from current module
            raise ImportError('Does not handle relative path')

    return result


def path(element):
    """
    Get full path of a given element.

    Do the inverse of lookup

    :param element: must be directly defined into a module or a package
    :type element: object
    """

    if not hasattr(element, '__name__'):
        raise AttributeError(
            'element {0} must have the attribute __name__'.format(element))

    result = element.__name__ if ismodule(element) else \
        '{0}.{1}'.format(element.__module__, element.__name__)

    return result


def isiterable(element, is_str=True):
    """
    Check whatever or not if input element is an iterable.

    :param is_str: check if element is also a str
    :type is_str: bool
    """
    result = isinstance(element, Iterable) \
        and (is_str or not isinstance(element, basestring))

    return result


def isunicode(s):
    """
    Check if string is unicode.

    :param s: string to check
    :type s: basestring

    :return: True if unicode or Python3, False otherwise
    """

    if PYVER < '3':
        return isinstance(s, unicode)

    else:
        return True


def ensure_unicode(s):
    """
    Convert string to unicode.

    :param s: string to convert
    :type s: basestring

    :return: unicode (or the same string if Python3)
    """

    result = s

    if PYVER < '3':
        if isinstance(s, basestring):
            if not isinstance(s, unicode):
                result = s.decode()
        else:
            raise TypeError('Expecting a string as argument')

    return result


def ensure_iterable(value, iterable=list):
    """
    Convert a value into an iterable if it is not.

    :param value: value to convert
    :type value: object

    :param iterable: iterable type to apply (default: list)
    :type iterable: type
    """

    result = value

    if not isiterable(value, is_str=False) or isinstance(value, dict):
        result = [value]
        result = iterable(result)

    else:
        result = iterable(value)

    return result


def get_first(iterable, default=None):
    """
    Try to get input iterable first item or default if iterable is empty.
    """

    result = iterable[0] if iterable else default

    return result


def prototype(typed_args=None, typed_kwargs=None, typed_return=None):
    """
    Decorate a function to check its parameters type.

    :param typed_args: Types for *args
    :type typed_args: tuple

    :param typed_kwargs: Types for **kwargs
    :type typed_kwargs: dict

    :param typed_return: Types for return
    :type typed_return: tuple of type, or type

    :raises: TypeError
    """

    if typed_args is None:
        typed_args = ()

    if typed_kwargs is None:
        typed_kwargs = {}

    if typed_return is None:
        typed_return = type(None)

    if isinstance(typed_return, list):
        typed_return = tuple(typed_return)

    def decorator(func):
        def wrapper(*args, **kwargs):
            i = 0
            l = len(args)

            for i in range(l):
                types = typed_args[i]

                if isinstance(types, list):
                    types = tuple(types)

                if not isinstance(args[i], types):
                    raise TypeError(
                        'Invalid type for arg#{}, got {} instead of {}'.format(
                            type(args[i]),
                            types[i]
                        )
                    )

            for key in typed_kwargs:
                types = typed_kwargs[key]

                if isinstance(types, list):
                    types = tuple(types)

                arg = kwargs.get(key, None)

                if not isinstance(arg, types):
                    raise TypeError(
                        'Invalid type for {}, got {} instead of {}'.format(
                            key,
                            type(arg),
                            types
                        )
                    )

            ret = func(*args, **kwargs)

            if not isinstance(ret, typed_return):
                raise TypeError(
                    'Invalid type for return, got {0} instead of {}'.format(
                        type(ret),
                        typed_return
                    )
                )

            return ret

        return wrapper

    return decorator
