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

"""Tools for managing path resolution of python objects."""

# ensure str are unicodes
from __future__ import unicode_literals, absolute_import

from inspect import ismodule, currentframe

from random import random

from .version import PY26
from .runtime import safe_eval

if PY26:
    import_module = __import__

else:
    from importlib import import_module

__all__ = ['clearcache', 'incache', 'lookup', 'getpath', 'alias']

#: lookup cache
__LOOKUP_CACHE = {}


def clearcache(path=None):
    """Clear cache memory for input path.

    :param str path: element path to remove from cache. If None clear all cache
    :param dict cache: cache to clear. Default is __LOOKUP_CACHE.

    :Example:

    >>> incache('b3j0f.utils')
    False
    >>> lookup('b3j0f.utils')
    >>> incache('b3j0f.utils')
    True
    >>> clearcache('b3j0f.path')
    >>> incache('b3j0f.utils')
    False
    >>> lookup('b3j0f.utils')
    >>> incache('b3j0f.utils')
    True
    >>> clearcache()
    >>> incache('b3j0f.utils')
    False
    """

    if path is None:
        __LOOKUP_CACHE.clear()

    else:
        __LOOKUP_CACHE.pop(path, None)


def incache(path):
    """Check if input path is in cache.

    :return: True if path is in cache
    :rtype: bool

    :Example:

    >>> incache('b3j0f.utils')
    False
    >>> lookup('b3j0f.utils')
    >>> incache('b3j0f.utils')
    True
    """

    return path in __LOOKUP_CACHE


def lookup(path, cache=True, scope=None, safe=False):
    """Get element reference from input element.

    The element can be a builtin/globals/scope object or is resolved from the
    current execution stack.

    :limitations: it does not resolve class methods or static values such as
        True, False, numbers, string and keywords.
    :param str path: full path to a python element.
    :param bool cache: if True (default), permits to reduce time complexity for
        lookup resolution in using cache memory to save resolved elements.
    :param dict scope: object scrope from where find path. For example, this
        scope can be locals(). Default is globals().
    :param bool safe: use lookup in a safe context. A safe context avoid to
        reach builtins function with I/O consequences.
    :return: python object which is accessible through input path
        or raise an exception if the path is wrong.
    :rtype: object
    :raises ImportError: if path is wrong
    """

    result = None

    found = path and cache and path in __LOOKUP_CACHE

    if found:
        result = __LOOKUP_CACHE[path]

    elif path:

        _eval = safe_eval if safe else eval

        try:  # search among scope
            result = _eval(path, scope)

        except NameError:

            # we generate a result in order to accept the result such as a None
            generated_result = random()
            result = generated_result

            components = path.split('.')
            index = 0
            components_len = len(components)

            module_name = components[0]

            # try to resolve an absolute path
            try:
                result = import_module(module_name)

            except ImportError:
                # resolve element globals or locals of the from previous frame
                previous_frame = currentframe().f_back

                if module_name in previous_frame.f_locals:
                    result = previous_frame.f_locals[module_name]

                elif module_name in previous_frame.f_globals:
                    result = previous_frame.f_globals[module_name]

            found = result is not generated_result

            if found:

                if components_len > 1:

                    index = 1

                    # try to import all sub-modules/packages
                    try:  # check if name is defined from an external module
                        # find the right module
                        while index < components_len:
                            module_name = '{0}.{1}'.format(
                                module_name, components[index]
                            )
                            result = import_module(module_name)
                            index += 1

                    except ImportError:
                        # path sub-module content
                        try:
                            if PY26:  # when __import__ is used
                                index = 1  # restart count of pathing
                            while index < components_len:
                                result = getattr(result, components[index])
                                index += 1

                        except AttributeError:
                            raise ImportError(
                                'Wrong path {0} at {1}'.format(
                                    path, components[:index]
                                )
                            )
                    else:  # in case of PY26
                        if PY26:
                            index = 1
                            while index < components_len:
                                result = getattr(result, components[index])
                                index += 1

        else:
            found = True

    if found:
        if cache:  # save in cache if found
            __LOOKUP_CACHE[path] = result

    else:
        raise ImportError('Wrong path {0}'.format(path))

    return result


def getpath(element):
    """Get full path of a given element such as the opposite of the
    resolve_path behaviour.

    :param element: must be directly defined into a module or a package and has
        the attribute '__name__'.

    :return: element absolute path.
    :rtype: str

    :raises AttributeError: if element has not the attribute __name__.

    :Example:

    >>> getpath(getpath)
    b3j0f.utils.path.getpath
    """

    if not hasattr(element, '__name__'):
        raise AttributeError(
            'element {0} must have the attribute __name__'.format(element)
        )

    result = element.__name__ if ismodule(element) else \
        '{0}.{1}'.format(element.__module__, element.__name__)

    return result


def alias(_id):
    """Decorator dedicated to make an alias of a decorated element in order to
    register it in the lookup cache.

    :param str _id: alias identifier.

    :Example:

    >>> alias('halfsonofzeus', 'hercules')
    'hercules'
    >>> lookup('halfsonofzeus')
    'hercules'
    >>> @alias('cube')
    >>> def cube(value): return value ** value ** value
    >>> lookpath('cube')(2)
    8
    """

    def _register_elt(elt):
        """Register a specific element in the lookup cache."""

        __LOOKUP_CACHE[_id] = elt

        return elt

    return _register_elt
