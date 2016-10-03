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

"""Module in charge of creating proxies like the design pattern ``proxy``.

A proxy is based on a callable element. It respects its signature but not the
implementation.
"""

from __future__ import absolute_import

__all__ = [
    'get_proxy', 'proxify_routine', 'proxify_elt', 'is_proxy', 'proxified_elt'
]

from time import time

from types import MethodType, FunctionType

from opcode import opmap

from random import randint

from sys import maxsize

from inspect import (
    getmembers, isroutine, ismethod, getargspec, getfile, isbuiltin, isclass
)

from six import (
    get_function_closure, get_function_code, get_function_defaults,
    get_function_globals, get_method_function, get_method_self, exec_, PY2, PY3,
    string_types, wraps
)

from .path import lookup
from .runtime import getcodeobj

# consts for interception loading
LOAD_GLOBAL = opmap['LOAD_GLOBAL']
LOAD_CONST = opmap['LOAD_CONST']

#: list of attributes to set after proxifying a function
WRAPPER_ASSIGNMENTS = ['__doc__', '__name__']
#: list of attributes to update after proxifying a function
WRAPPER_UPDATES = ['__dict__']
#: lambda function name
__LAMBDA_NAME__ = (lambda: None).__name__
#: proxy class name
__PROXY_CLASS__ = 'Proxy'
#: attribute name for proxified element
__PROXIFIED__ = '__proxified__'
#: instance method name for delegating proxy generation to the elt to proxify
__GETPROXY__ = '__getproxy__'


def proxify_elt(elt, bases=None, _dict=None, public=False):
    """Proxify input elt.

    :param elt: elt to proxify.
    :param bases: elt class base classes. If None, use elt type.
    :param dict _dict: specific elt class content to use.
    :param bool public: if True (default False), proxify only public members
        (where name starts with the character '_').
    :return: proxified element.
    :raises: TypeError if elt does not implement all routines of bases and
        _dict.
    """

    # ensure _dict is a dictionary
    proxy_dict = {} if _dict is None else _dict.copy()
    # set of proxified attribute names which are proxified during bases parsing
    # and avoid to proxify them twice during _dict parsing
    proxified_attribute_names = set()
    # ensure bases is a tuple of types
    if bases is None:
        bases = (elt if isclass(elt) else elt.__class__,)

    if isinstance(bases, string_types):
        bases = (lookup(bases),)

    elif isclass(bases):
        bases = (bases,)

    else:
        bases = tuple(bases)

    # fill proxy_dict with routines of bases
    for base in bases:
        # exclude object
        if base is object:
            continue
        for name, member in getmembers(base, isroutine):
            # check if name is public
            if public and not name.startswith('_'):
                continue
            eltmember = getattr(elt, name, None)
            if eltmember is None:
                raise TypeError(
                    'Wrong elt {0}. Must implement {1} ({2}) of {3}.'.
                    format(elt, name, member, base)
                )
            # proxify member if member is not a constructor
            if name not in ['__new__', '__init__']:
                # get routine from proxy_dict or eltmember
                routine = proxy_dict.get(name, eltmember)
                # exclude object methods
                if getattr(routine, '__objclass__', None) is not object:
                    # get routine proxy
                    routine_proxy = proxify_routine(routine)
                    if ismethod(routine_proxy):
                        routine_proxy = get_method_function(routine_proxy)
                    # update proxy_dict
                    proxy_dict[name] = routine_proxy
                    # and save the proxified attribute flag
                    proxified_attribute_names.add(name)

    # proxify proxy_dict
    for name in proxy_dict:
        value = proxy_dict[name]
        if not hasattr(elt, name):
            raise TypeError(
                'Wrong elt {0}. Must implement {1} ({2}).'.format(
                    elt, name, value
                )
            )
        if isroutine(value):
            # if member has not already been proxified
            if name not in proxified_attribute_names:
                # proxify it
                value = proxify_routine(value)
            proxy_dict[name] = value

    # set default constructors if not present in proxy_dict
    if '__new__' not in proxy_dict:
        proxy_dict['__new__'] = object.__new__

    if '__init__' not in proxy_dict:
        proxy_dict['__init__'] = object.__init__

    # generate a new proxy class
    cls = type('Proxy', bases, proxy_dict)
    # instantiate proxy cls
    result = cls if isclass(elt) else cls()
    # bind elt to proxy
    setattr(result, __PROXIFIED__, elt)

    return result


def proxify_routine(routine, impl=None):
    """Proxify a routine with input impl.

    :param routine: routine to proxify.
    :param impl: new impl to use. If None, use routine.
    """

    # init impl
    impl = routine if impl is None else impl

    is_method = ismethod(routine)
    if is_method:
        function = get_method_function(routine)
    else:
        function = routine

    # flag which indicates that the function is not a pure python function
    # and has to be wrapped
    wrap_function = not hasattr(function, '__code__')

    try:
        # get params from routine
        args, varargs, kwargs, _ = getargspec(function)
    except TypeError:
        # in case of error, wrap the function
        wrap_function = True

    if wrap_function:
        # if function is not pure python, create a generic one
        # with assignments
        assigned = []
        for wrapper_assignment in WRAPPER_ASSIGNMENTS:
            if hasattr(function, wrapper_assignment):
                assigned.append(wrapper_assignment)
        # and updates
        updated = []
        for wrapper_update in WRAPPER_UPDATES:
            if hasattr(function, wrapper_update):
                updated.append(wrapper_update)

        @wraps(function, assigned=assigned, updated=updated)
        def wrappedfunction(*args, **kwargs):
            """Default wrap function."""

        function = wrappedfunction
        # get params from function
        args, varargs, kwargs, _ = getargspec(function)

    name = function.__name__

    result = _compilecode(
        function=function, name=name, impl=impl,
        args=args, varargs=varargs, kwargs=kwargs
    )

    # set wrapping assignments
    for wrapper_assignment in WRAPPER_ASSIGNMENTS:
        try:
            value = getattr(function, wrapper_assignment)
        except AttributeError:
            pass
        else:
            setattr(result, wrapper_assignment, value)

    # set proxy module
    result.__module__ = proxify_routine.__module__

    # update wrapping updating
    for wrapper_update in WRAPPER_UPDATES:
        try:
            value = getattr(function, wrapper_update)
        except AttributeError:
            pass
        else:
            getattr(result, wrapper_update).update(value)

    # set proxyfied element on proxy
    setattr(result, __PROXIFIED__, routine)

    if is_method:  # create a new method
        args = [result, get_method_self(routine)]
        if PY2:
            args.append(routine.im_class)
        result = MethodType(*args)

    return result


def _compilecode(function, name, impl, args, varargs, kwargs):
    """Get generated code.

    :return: function proxy generated code.
    :rtype: str
    """

    newcodestr, generatedname, impl_name = _generatecode(
        function=function, name=name, impl=impl,
        args=args, varargs=varargs, kwargs=kwargs
    )

    try:
        __file__ = getfile(function)
    except TypeError:
        __file__ = '<string>'

    # compile newcodestr
    code = compile(newcodestr, __file__, 'single')

    # define the code with the new function
    _globals = {}
    exec_(code, _globals)

    # get new code
    _var = _globals[generatedname]
    newco = get_function_code(_var)

    # get new consts list
    newconsts = list(newco.co_consts)

    if PY3:
        newcode = list(newco.co_code)
    else:
        newcode = [ord(co) for co in newco.co_code]

    consts_values = {impl_name: impl}

    # change LOAD_GLOBAL to LOAD_CONST
    index = 0
    newcodelen = len(newcode)
    while index < newcodelen:
        if newcode[index] == LOAD_GLOBAL:
            oparg = newcode[index + 1] + (newcode[index + 2] << 8)
            name = newco.co_names[oparg]
            if name in consts_values:
                const_value = consts_values[name]
                if const_value in newconsts:
                    pos = newconsts.index(const_value)
                else:
                    pos = len(newconsts)
                    newconsts.append(consts_values[name])
                newcode[index] = LOAD_CONST
                newcode[index + 1] = pos & 0xFF
                newcode[index + 2] = pos >> 8
        index += 1

    codeobj = getcodeobj(newconsts, newcode, newco, get_function_code(function))
    # instanciate a new function
    if function is None or isbuiltin(function):
        result = FunctionType(codeobj, {})

    else:
        result = type(function)(
            codeobj,
            get_function_globals(function),
            function.__name__,
            get_function_defaults(function),
            get_function_closure(function)
        )

    return result


def _generatecode(function, name, impl, args, varargs, kwargs):

    code = ''

    # flag for lambda function
    islambda = __LAMBDA_NAME__ == name
    if islambda:
        generatedname = '_{0}'.format(int(time()))

    else:
        generatedname = name

    # get join method for reducing concatenation time execution
    join = "".join

    # default indentation
    indent = '    '

    if islambda:
        code = '{0} = lambda '.format(generatedname)

    else:
        code = 'def {0}('.format(generatedname)

    if args:
        code = join((code, '{0}'.format(args[0])))

    for arg in args[1:]:
        code = join((code, ', {0}'.format(arg)))

    if varargs is not None:
        if args:
            code = join((code, ', '))
        code = join((code, '*{0}'.format(varargs)))

    if kwargs is not None:
        if args or varargs is not None:
            code = join((code, ', '))

        code = join((code, '**{0}'.format(kwargs)))

    impl_name = '_{0}'.format(randint(0, maxsize))

    # insert impl call
    if islambda:
        code = join((code, ': {0}('.format(impl_name)))

    else:
        code = join(
            (
                code,
                '):\n{0}return {1}('.format(indent, impl_name)
            )
        )

    impl_args = args[1:] if ismethod(impl) else args
    if impl_args:
        code = join((code, '{0}'.format(impl_args[0])))
    for arg in impl_args[1:]:
        code = join((code, ', {0}'.format(arg)))

    if varargs is not None:
        if args:
            code = join((code, ', '))
        code = join((code, '*{0}'.format(varargs)))

    if kwargs is not None:
        if args or varargs is not None:
            code = join((code, ', '))
        code = join((code, '**{0}'.format(kwargs)))

    code = join((code, ')\n'))

    result = code, generatedname, impl_name

    return result


def get_proxy(elt, bases=None, _dict=None):
    """Get proxy from an elt.

    If elt implements the proxy generator method (named ``__getproxy__``), use
    it instead of using this module functions.

    :param elt: elt to proxify.
    :type elt: object or function/method
    :param bases: base types to enrich in the result cls if not None.
    :param _dict: class members to proxify if not None.
    """

    # try to find an instance proxy generator
    proxygenerator = getattr(elt, __GETPROXY__, None)

    # if a proxy generator is not found, use this module
    if proxygenerator is None:
        if isroutine(elt):
            result = proxify_routine(elt)

        else:  # in case of object, result is a Proxy
            result = proxify_elt(elt, bases=bases, _dict=_dict)

    else:  # otherwise, use the specific proxy generator
        result = proxygenerator()

    return result


def proxified_elt(proxy):
    """Get proxified element.

    :param proxy: proxy element from where get proxified element.
    :return: proxified element. None if proxy is not proxified.
    """

    if ismethod(proxy):
        proxy = get_method_function(proxy)
    result = getattr(proxy, __PROXIFIED__, None)

    return result


def is_proxy(elt):
    """Return True if elt is a proxy.

    :param elt: elt to check such as a proxy.
    :return: True iif elt is a proxy.
    :rtype: bool
    """

    if ismethod(elt):
        elt = get_method_function(elt)

    result = hasattr(elt, __PROXIFIED__)

    return result
