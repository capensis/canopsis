# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
from imp import load_source

from inspect import ismodule

from collections import Iterable

from sys import version as PYVER

from os.path import expanduser, splitext
from os.path import join as joinpath
from os import listdir

from re import search as regsearch

from .init import basestring

__RESOLVED_ELEMENTS = {}  #: dictionary of resolved elements by name


def dynmodloads(_path='.', subdef=False, pattern='.*', logger=None):
    loaded = {}
    _path = expanduser(_path)

    for mfile in listdir(_path):
        name, ext = splitext(mfile)

        # Ignore "." and "__init__.py" and everything not matched by "*.py"
        if name in ['.', '__init__'] or ext != '.py':
            continue

        logger.info("Load '{0}' ...".format(name))

        try:
            module = load_source(name, joinpath(_path, mfile))

        except ImportError as err:
            logger.error('Impossible to import {0}: {1}'.format(name, err))

        else:
            loaded[name] = module

            if subdef:
                alldefs = dir(module)
                builtindefs =  [
                    '__builtins__',
                    '__doc__',
                    '__file__',
                    '__name__',
                    '__package__'
                ]

                for mydef in alldefs:
                    if mydef not in builtindefs and regsearch(pattern, mydef):
                        logger.debug('from {0} import {1}'.format(
                            name, mydef
                        ))

                        loaded[mydef] = getattr(module, mydef)

    return loaded


def setdefaultattr(obj, attr, value):
    """
    Set attribute in object if not present.

    :param obj: Object to set attribute to
    :type obj: anything

    :param attr: Attribute's name to set
    :type attr: str

    :param value: Value to set
    :type value: anything

    :returns: current value if attribute exists, or new value otherwise
    """

    if hasattr(obj, attr):
        return getattr(obj, attr)

    else:
        setattr(obj, attr, value)
        return value


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

    path = expanduser(path)
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
        except ImportError as e:
            print(
                'Error while importing module {} : {}'.format(module_name, e)
            )

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

                except ImportError as ie:
                    # path sub-module content
                    try:

                        while index < components_len:
                            result = getattr(result, components[index])
                            index += 1

                    except AttributeError as ae:
                        raise ImportError(
                            ('Wrong path {0} at {1}, ' +
                                'errors when importing module {2} ' +
                                ': {3}, {4}').format(
                                path,
                                components[:index],
                                module_name,
                                ie,
                                ae
                            ))

            if result is not None and cached:
                __RESOLVED_ELEMENTS[path] = result

        else:  # get relative object from current module
            raise ImportError(
                'Does not handle relative path: {0}'.format(path)
            )

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
    """Convert string to unicode.

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


def forceUTF8(data, _memory=None):
    """Return a copy of data where all embedded strings are UTF8.

    :param data: data from where convert str to UTF8 format.
    :return: data copy where all str are in UTF8 format.
    """

    # by default, result is data
    result = data
    # do something only if python version is 2
    if PYVER < '3':
        # initialize memory
        if _memory is None:
            _memory = {}
        # if data has already been processed
        data_id = id(data)
        if data_id in _memory:
            # result is the previous result
            result = _memory[data_id]
        else:  # else process data
            # if data is a basestring, decode it to an utf8
            if isinstance(data, basestring):
                if not isinstance(data, unicode):
                    result = data.decode('utf-8', 'ignore')
            # if data is a dict
            elif isinstance(data, dict):
                data_id = id(data)
                if data_id in _memory:
                    result = _memory[data_id]
                else:
                    # copy data
                    result = data.copy()
                    for param in data:
                        value = data[param]
                        # convert param and value
                        param = forceUTF8(param)
                        value = forceUTF8(value)
                        result[param] = value
            # if data is an iterable
            elif isinstance(data, Iterable):
                result = []
                # convert all values of data
                for d in data:
                    value = forceUTF8(d)
                    result.append(value)
                # and convert result to data type
                result = type(data)(result)
            # save the result in memory
            _memory[data_id] = result

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


class dictproperty(object):
    """
        Property decorator for dict-like attributes.
    """

    class _proxy(object):
        """
            Proxy interface to dict-like objects.
        """

        def __init__(self, obj, fget, fset, fdel, *args, **kwargs):
            super(dictproperty._proxy, self).__init__(*args, **kwargs)

            self._obj = obj
            self._fget = fget
            self._fset = fset
            self._fdel = fdel

        def __getitem__(self, key):
            if self._fget is None:
                raise TypeError("Impossible to get key: {0}".format(key))

            return self._fget(self._obj, key)

        def __setitem__(self, key, value):
            if self._fset is None:
                raise TypeError("Impossible to set key: {0} = {1}".format(key, value))

            self._fset(self._obj, key, value)

        def __delitem__(self, key):
            if self._fdel is None:
                raise TypeError("Impossible to delete key: {0}".format(key))

            self._fdel(self._obj, key)

    def __init__(self, fget=None, fset=None, fdel=None, doc=None, *args, **kwargs):
        super(dictproperty, self).__init__(*args, **kwargs)

        self._fget = fget
        self._fset = fset
        self._fdel = fdel
        self.__doc__ = doc

    def __get__(self, obj, objtype=None):
        if obj is None:
            return self

        return self._proxy(obj, self._fget, self._fset, self._fdel)
