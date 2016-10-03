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

"""Module which aims to manage python joinpoint interception.

A joinpoint is just a callable element.

functions allow to weave an interception function on any python callable
object.
"""

from __future__ import absolute_import

from inspect import (
    isbuiltin, ismethod, isclass, isfunction, getmodule, getmembers, getfile,
    getargspec, isroutine, getmro
)

from opcode import opmap

from six.moves import builtins

from types import MethodType, FunctionType

from functools import wraps

from time import time

from six import PY3, PY2

__all__ = [
    'Joinpoint', 'JoinpointError',
    'get_intercepted', 'is_intercepted',
    'super_method', 'find_ctx', 'base_ctx'
]

# consts for interception loading
LOAD_GLOBAL = opmap['LOAD_GLOBAL']
LOAD_CONST = opmap['LOAD_CONST']

#: attribute which binds the intercepted function from the interceptor function
_INTERCEPTED = '_intercepted'
#: ctx intercepted atttribute name
_INTERCEPTED_CTX = '_intercepted_ctx'

#: attribute which binds an interception function to its parent joinpoint
_INTERCEPTION = '_interception'

#: list of attributes to set/update after wrapping a function with a joinpoint
WRAPPER_ASSIGNMENTS = ['__doc__', '__module__', '__name__']
WRAPPER_UPDATES = ['__dict__']


def find_ctx(elt):
    """Find a Pointcut ctx which is a class/instance related to input
    function/method.

    :param elt: elt from where find a ctx.
    :return: elt ctx. None if no ctx available or if elt is a None method.
    """

    result = None

    if ismethod(elt):

        result = elt.__self__

        if result is None and PY2:
            result = elt.im_class

    elif isclass(elt):
        result = elt

    return result


def base_ctx(ctx):
    """Get base ctx.

    :param ctx: initial ctx.
    :return: base ctx.
    """

    result = None

    if isclass(ctx):
        result = getattr(ctx, '__base__', None)
        if result is None:
            mro = getmro(ctx)
            if len(mro) > 1:
                result = mro[1]
    else:
        result = ctx.__class__

    return result


def super_method(name, ctx):
    """Get super ctx method and ctx where name matches with input name.

    :param name: method name to find in super ctx.
    :param ctx: initial method ctx.
    :return: method in super ctx and super ctx.
    :rtype: tuple
    """

    result = None, None

    # get class ctx
    if isclass(ctx):
        _ctx = ctx
        first_mro = 1
    else:
        _ctx = ctx.__class__
        first_mro = 0
    # get class hierachy
    mro = getmro(_ctx)
    for cls in mro[first_mro:]:
        if hasattr(cls, name):
            result = getattr(cls, name), cls
            break

    return result


class JoinpointError(Exception):
    """Handle Joinpoint errors
    """


class Joinpoint(object):
    """Manage joinpoint execution with Advices.

    Advices are callable objects which take in parameter a Joinpoint.

    Joinpoint provides to advices:
        - the joinpoint,
        - joinpoint call arguments as args and kwargs property,
        - a shared context during interception such as a dictionary.
    """

    #: lambda function name
    __LAMBDA_NAME__ = (lambda: None).__name__

    #: lambda function interception name
    __INTERCEPTION__ = 'interception'

    #: context execution attribute name
    EXEC_CTX = 'exec_ctx'

    #: interception attribute name
    _INTERCEPTION = '_interception'

    #: interception args attribute name
    ARGS = 'args'
    #: interception kwargs attribute name
    KWARGS = 'kwargs'

    #: target element attribute name
    TARGET = 'target'
    #: target element ctx attribute name
    CTX = 'ctx'

    #: private attribute name for internal iterator for advices execution
    _ADVICES_ITERATOR = '_advices_iterator'

    #: private attribute name for advices
    _ADVICES = '_advices'

    __slots__ = (
        EXEC_CTX, ARGS, KWARGS, TARGET, CTX,  # public attributes
        _ADVICES_ITERATOR, _ADVICES, _INTERCEPTION  # private attributes
    )

    def __init__(
            self,
            target=None, args=None, kwargs=None, advices=None, ctx=None,
            exec_ctx=None
    ):
        """Initialize a new Joinpoint with optional parameters such as a
        target, its calling arguments (args and kwargs) and a list of
            advices (callable which take self in parameter).

        If target, args and kwargs are not None, self Joinpoint use them in a
            static context. Otherwise, they will be resolved at proceeding
            time.

        :param callable target: target which is intercepted by advices.
        :param tuple args: target call varargs argument.
        :param dict kwargs: target call keywords argument.
        :param Iterable advices: iterable of advices which take in parameters
            this joinpoint. If None, they will be dynamically loaded during
            self proceeding time related to target.
        :param dict exec_ctx: execution context. Empty dict by default.
        :param ctx: target ctx if target is an class/instance attribute.
        """

        super(Joinpoint, self).__init__()

        self._advices_iterator = None

        # init critical parameters
        self._interception = None
        self.ctx = None
        self.target = None

        # set target arguments
        self.args = () if args is None else args
        self.kwargs = {} if kwargs is None else kwargs

        # set advices
        self._advices = advices

        # set context
        self.exec_ctx = {} if exec_ctx is None else exec_ctx

        # set target
        self.set_target(target=target, ctx=ctx)

    def __repr__(self):

        self_type = type(self)
        result = "{0}(".format(self_type.__name__)

        for slot in self_type.__slots__:
            # do not display advices iterator
            if slot != Joinpoint._ADVICES_ITERATOR:
                result += "{0}:{1},".format(slot, getattr(self, slot))
        else:
            result = "{0})".format(result[:-2])

        return result

    def set_target(self, target, ctx=None):
        """Set target.

        :param target: new target to use.
        :param target ctx: target ctx if target is an class/instance attribute.
        """

        if target is not None:
            # check if target is already intercepted
            if is_intercepted(target):
                # set self interception last target reference
                self._interception = target
                # and targets, ctx
                self.target, self.ctx = get_intercepted(target)
            else:
                # if not, update target reference with new interception
                self.apply_pointcut(target, ctx=ctx)

    def start(
            self, target=None, args=None, kwargs=None, advices=None,
            exec_ctx=None, ctx=None
    ):
        """ Start to proceed this Joinpoint in initializing target, its
        arguments and advices. Call self.proceed at the end.

        :param callable target: new target to use in proceeding. self.target by
            default.
        :param tuple args: new target args to use in proceeding. self.args by
            default.
        :param dict kwargs: new target kwargs to use in proceeding. self.kwargs
            by default.
        :param list advices: advices to use in proceeding. self advices by
            default.
        :param dict exec_ctx: execution context.
        :param target_ctx: target ctx to use in proceeding.
        :return: self.proceed()
        """

        # init target and _interception if not None as set_target method do
        if target is not None:
            self.set_target(target=target, ctx=ctx)

        # init args if not None
        if args is not None:
            self.args = args

        # init kwargs if not None
        if kwargs is not None:
            self.kwargs = kwargs

        # get advices to process
        if advices is None:
            if self._advices is not None:
                advices = self._advices
            else:
                advices = self.get_advices(self._interception)

        # initialize self._advices_iterator
        self._advices_iterator = iter(advices)

        # initialize execution context
        self.exec_ctx = self.exec_ctx if exec_ctx is None else exec_ctx

        result = self.proceed()

        return result

    def proceed(self):
        """Proceed this Joinpoint in calling all advices with this joinpoint
        as the only one parameter, and call at the end the target.
        """

        try:
            # get next advice
            advice = next(self._advices_iterator)

        except StopIteration:  # if no advice can be applied
            # call target
            result = self.target(*self.args, **self.kwargs)

        else:
            # if has next, apply advice on self
            result = advice(self)

        return result

    def apply_pointcut(self, target, function=None, ctx=None):
        """Apply pointcut on input target and returns final interception.

        The poincut respects all meta properties such as:
        - function signature,
        - module path,
        - file path,
        - __dict__ reference.
        """

        try:
            __file__ = getfile(target)
        except TypeError:
            __file__ = '<string>'

        if function is None:
            function = _get_function(target)

        # flag which indicates that the function is not a pure python function
        # and has to be wrapped
        wrap_function = not hasattr(function, '__code__')

        try:
            # get params from target
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
            def wrapper(*args, **kwargs):
                """Default wrapper."""

            function = wrapper

            # get params from target wrapper
            args, varargs, kwargs, _ = getargspec(function)

        # get params from target
        name = function.__name__

        # if target has not name, use 'function'
        if name == Joinpoint.__LAMBDA_NAME__:
            name = Joinpoint.__INTERCEPTION__

        # get join method for reducing concatenation time execution
        join = "".join

        # default indentation
        indent = '    '

        newcodestr = "def {0}(".format(name)
        if args:
            newcodestr = join((newcodestr, "{0}".format(args[0])))
        for arg in args[1:]:
            newcodestr = join((newcodestr, ", {0}".format(arg)))

        if varargs is not None:
            if args:
                newcodestr = join((newcodestr, ", "))
            newcodestr = join((newcodestr, "*{0}".format(varargs)))

        if kwargs is not None:
            if args or varargs is not None:
                newcodestr = join((newcodestr, ", "))
            newcodestr = join((newcodestr, "**{0}".format(kwargs)))

        newcodestr = join((newcodestr, "):\n"))

        # unique id which will be used for advicesexecutor and kwargs
        generated_id = repr(time()).replace('.', '_')

        # if kwargs is None
        if kwargs is None and args:
            kwargs = "kwargs_{0}".format(generated_id)  # generate a name
            # initialize a new dict with args
            newcodestr = join(
                (newcodestr, "{0}{1} = {{\n".format(indent, kwargs)))
            for arg in args:
                newcodestr = join(
                    (newcodestr, "{0}{0}'{1}': {1},\n".format(indent, arg))
                )
            newcodestr = join((newcodestr, "{0}}}\n".format(indent)))

        else:
            # fill args in kwargs
            for arg in args:
                newcodestr = join(
                    (newcodestr, "{0}{1}['{2}'] = {2}\n".format(
                        indent, kwargs, arg))
                )

        # advicesexecutor name
        joinpoint = "joinpoint_{0}".format(generated_id)

        if varargs:
            newcodestr = join(
                (newcodestr, "{0}{1}.args = {2}\n".format(
                    indent, joinpoint, varargs))
            )

        # set kwargs in advicesexecutor
        if kwargs is not None:
            newcodestr = join(
                (newcodestr, "{0}{1}.kwargs = {2}\n".format(
                    indent, joinpoint, kwargs))
            )

        # return advicesexecutor proceed result
        start = "start_{0}".format(generated_id)
        newcodestr = join(
            (newcodestr, "{0}return {1}()\n".format(indent, start))
        )
        # compile newcodestr
        code = compile(newcodestr, __file__, 'single')

        _globals = {}

        # define the code with the new function
        exec(code, _globals)

        # get new code
        newco = _globals[name].__code__

        # get new consts list
        newconsts = list(newco.co_consts)

        if PY3:
            newcode = list(newco.co_code)
        else:
            newcode = map(ord, newco.co_code)

        consts_values = {joinpoint: self, start: self.start}

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
                    if name == start:
                        break  # stop when start is encountered
            index += 1

        # get code string
        codestr = bytes(newcode) if PY3 else join(map(chr, newcode))

        # get vargs
        vargs = [
            newco.co_argcount, newco.co_nlocals, newco.co_stacksize,
            newco.co_flags, codestr, tuple(newconsts), newco.co_names,
            newco.co_varnames, newco.co_filename, newco.co_name,
            newco.co_firstlineno, newco.co_lnotab,
            getattr(function.__code__, 'co_freevars', ()),
            newco.co_cellvars
        ]
        if PY3:
            vargs.insert(1, newco.co_kwonlyargcount)

        # instanciate a new code object
        codeobj = type(newco)(*vargs)
        # instanciate a new function
        if function is None or isbuiltin(function):
            interception_fn = FunctionType(codeobj, {})

        else:
            interception_fn = type(function)(
                codeobj,
                {} if function.__globals__ is None else function.__globals__,
                function.__name__,
                function.__defaults__,
                function.__closure__
            )

        # set wrapping assignments
        for wrapper_assignment in WRAPPER_ASSIGNMENTS:
            try:
                value = getattr(function, wrapper_assignment)
            except AttributeError:
                pass
            else:
                setattr(interception_fn, wrapper_assignment, value)
        # update wrapping updating
        for wrapper_update in WRAPPER_UPDATES:
            try:
                value = getattr(function, wrapper_update)
            except AttributeError:
                pass
            else:
                getattr(interception_fn, wrapper_update).update(value)

        # set interception, target function and ctx
        self._interception, self.target, self.ctx = _apply_interception(
            target=target, interception_fn=interception_fn, ctx=ctx
        )

        return self._interception

    def get_advices(self, target):
        """Get target advices.

        :param target: target from where getting advices.
        """

        raise NotImplementedError()


def _apply_interception(
        target, interception_fn, ctx=None, _globals=None
):
    """Apply interception on input target and return the final target.

    :param Callable target: target on applying the interception_fn.
    :param function interception_fn: interception function to apply on
        target
    :param ctx: target ctx (instance or class) if not None.

    :return: both interception and intercepted
        - if target is a builtin function,
            the result is a (wrapper function, builtin).
        - if target is a function, interception is target where
            code is intercepted code, and interception is a new function where
            code is target code.
    :rtype: tuple(callable, function, ctx)
    :raises: TypeError if target is not a routine.
    """

    if not callable(target):
        raise TypeError('target {0} is not callable.'.format(target))

    intercepted = target
    interception = interception_fn

    # try to get the right ctx
    if ctx is None:
        ctx = find_ctx(elt=target)

    # if target is a builtin
    if isbuiltin(target) or getmodule(target) is builtins:
        # update builtin function reference in module with wrapper
        module = getmodule(target)
        found = False  # check for found function

        if module is not None:
            # update all references by value
            for name, _ in getmembers(
                    module, lambda member: member is target):
                setattr(module, name, interception_fn)
                found = True

            if not found:  # raise Exception if not found
                raise JoinpointError(
                    "Impossible to weave on not modifiable function {0}. \
                    Must be contained in module {1}".format(target, module)
                )

    elif ctx is None:
        # update code with interception code
        target_fn = _get_function(target)
        # switch interception and intercepted
        interception, intercepted = target, interception_fn
        # switch of code between target_fn and
        # interception_fn
        target_fn.__code__, interception_fn.__code__ = \
            interception_fn.__code__, target_fn.__code__

    else:
        # get target name
        if isclass(target):  # if target is a class, get constructor name
            target_name = _get_function(target).__name__
        else:  # else get target name
            target_name = target.__name__
        # get the right intercepted
        intercepted = getattr(ctx, target_name)
        # in case of method
        if ismethod(intercepted):  # in creating eventually a new method
            args = [interception, ctx]
            if PY2:  # if py2, specify the ctx class
                # and unbound method type
                if intercepted.__self__ is None:
                    args = [interception, None, ctx]
                else:
                    args.append(ctx.__class__)
            # instantiate a new method
            interception = MethodType(*args)
        # get the right intercepted function
        if is_intercepted(intercepted):
            intercepted, _ = get_intercepted(intercepted)
        else:
            intercepted = _get_function(intercepted)
        # set in ctx the new method
        setattr(ctx, target_name, interception)

    # add intercepted into interception_fn globals and attributes
    interception_fn = _get_function(interception)
    # set intercepted
    setattr(interception_fn, _INTERCEPTED, intercepted)
    # set intercepted ctx
    if ctx is not None:
        setattr(interception_fn, _INTERCEPTED_CTX, ctx)

    interception_fn.__globals__[_INTERCEPTED] = intercepted
    interception_fn.__globals__[_INTERCEPTION] = interception

    if _globals is not None:
        interception_fn.__globals__.update(_globals)

    return interception, intercepted, ctx


def _unapply_interception(target, ctx=None):
    """Unapply interception on input target in cleaning it.

    :param routine target: target from where removing an interception
        function. is_joinpoint(target) must be True.
    :param ctx: target ctx.
    """

    # try to get the right ctx
    if ctx is None:
        ctx = find_ctx(elt=target)

    # get previous target
    intercepted, old_ctx = get_intercepted(target)

    # if ctx is None and old_ctx is not None, update ctx with old_ctx
    if ctx is None and old_ctx is not None:
        ctx = old_ctx

    if intercepted is None:
        raise JoinpointError('{0} must be intercepted'.format(target))

    # flag to deleting of joinpoint_function
    del_joinpoint_function = False

    # if old target is a not modifiable resource
    if isbuiltin(intercepted):
        module = getmodule(intercepted)
        found = False

        # update references to target to not modifiable element in module
        for name, member in getmembers(module):
            if member is target:
                setattr(module, name, intercepted)
                found = True

        # if no reference found, raise an Exception
        if not found:
            raise JoinpointError(
                "Impossible to unapply interception on not modifiable element \
                {0}. Must be contained in module {1}".format(target, module)
            )

    elif ctx is None:
        # get joinpoint function
        joinpoint_function = _get_function(target)
        # update old code on target
        joinpoint_function.__code__ = intercepted.__code__
        # ensure to delete joinpoint_function
        del_joinpoint_function = True

    else:
        # flag for joinpoint recovering
        recover = False
        # get interception name in order to update/delete interception from ctx
        intercepted_name = intercepted.__name__
        # should we change of target or is it inherited ?
        if isclass(ctx):
            base_interception, _ = super_method(name=intercepted_name, ctx=ctx)
        else:
            base_interception = getattr(ctx.__class__, intercepted_name, None)
        # if base interception does not exist
        if base_interception is None:  # recover intercepted
            recover = True

        else:
            # get joinpoint_function
            joinpoint_function = _get_function(target)
            # get base function
            if is_intercepted(base_interception):
                base_intercepted, _ = get_intercepted(base_interception)
            else:
                base_intercepted = _get_function(base_interception)
            # is interception inherited ?
            if base_intercepted is joinpoint_function:
                pass  # do nothing
            # is intercepted inherited
            elif base_intercepted is intercepted:
                # del interception
                delattr(ctx, intercepted_name)
                del_joinpoint_function = True
            else:  # base function is something else
                recover = True

        if recover:  # if recover is required
            # new content to put in ctx
            new_content = intercepted
            if ismethod(target):  # in creating eventually a new method
                args = [new_content, ctx]
                if PY2:  # if py2, specify the ctx class
                    # and unbound method type
                    if target.__self__ is None:
                        args = [new_content, None, ctx]
                    else:  # or instance method
                        args.append(ctx.__class__)
                # instantiate a new method
                new_content = MethodType(*args)
            # update ctx with intercepted
            setattr(ctx, intercepted_name, new_content)
            joinpoint_function = _get_function(target)
            del_joinpoint_function = True

    if del_joinpoint_function:
        # delete _INTERCEPTED and _INTERCEPTED_CTX from joinpoint_function
        if hasattr(joinpoint_function, _INTERCEPTED):
            delattr(joinpoint_function, _INTERCEPTED)
            if hasattr(joinpoint_function, _INTERCEPTED_CTX):
                delattr(joinpoint_function, _INTERCEPTED_CTX)
        del joinpoint_function


def is_intercepted(target):
    """True iif input target is intercepted.

    :param target: target to check such as an intercepted target.

    :return: True iif input target is intercepted.
    :rtype: bool
    """

    result = False

    # get interception function from input target
    function = _get_function(target)

    result = hasattr(function, _INTERCEPTED)

    return result


def get_intercepted(target):
    """Get intercepted function and ctx from input target.

    :param target: target from where getting the intercepted function and ctx.

    :return: target intercepted function and ctx.
        (None, None) if no intercepted function exist.
        (fn, None) if not ctx exists.
    :rtype: tuple
    """

    function = _get_function(target)

    intercepted = getattr(function, _INTERCEPTED, None)
    ctx = getattr(function, _INTERCEPTED_CTX, None)

    return intercepted, ctx


def _get_function(target):
    """Get target function.

    :param callable target: target from where get function.

    :return: depending on target type::

        - class: constructor.
        - method: method function.
        - function: function.
        - else: __call__ method.

    :raises: TypeError if target is not callable or is a class without a
        constructor.
    """

    result = None

    # raise TypeError if target is not callable
    if not callable(target):
        raise TypeError('target {0} must be callable'.format(target))

    # in case of class, final target is its constructor
    if isclass(target):
        constructor = getattr(
            target, '__init__',  # try to find __init__
            getattr(target, '__new__', None)
        )  # try to find __new__ | target
        # if no constructor exists
        if constructor is None:
            # create one
            def __init__(self):
                pass
            if PY2:
                target.__init__ = MethodType(__init__, None, target)
            else:
                target.__init__ = __init__
            constructor = target.__init__
        # if constructor is a method, return function method
        if ismethod(constructor):
            result = constructor.__func__
        # else return constructor
        else:
            result = constructor

    elif ismethod(target):  # if target is a method, return function method
        result = target.__func__

    # return target if target is function or builtin
    elif isfunction(target) or isbuiltin(target) or isroutine(target):
        result = target

    else:  # otherwise, return __call__ method
        __call__ = getattr(target, '__call__')

        if ismethod(__call__):  # if __call__ is a method, return its function
            result = __call__.__func__
        else:  # otherwise return __call__
            result = __call__

    return result
