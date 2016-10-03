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

"""Provides functions in order to weave/unweave/get advices from callable
objects."""

from re import compile as re_compile

from inspect import getmembers, isroutine, isclass

from opcode import opmap

try:
    from threading import Timer
except ImportError:
    from dummy_threading import Timer

from ..joinpoint.core import (
    _unapply_interception, is_intercepted, _get_function, Joinpoint, find_ctx,
    super_method, get_intercepted, base_ctx
)

from six import string_types, callable

__all__ = ['AdviceError', 'get_advices', 'weave', 'unweave', 'weave_on']

# consts for interception loading
LOAD_GLOBAL = opmap['LOAD_GLOBAL']
LOAD_CONST = opmap['LOAD_CONST']

WRAPPER_ASSIGNMENTS = ('__doc__', '__annotations__', '__dict__', '__module__')

_ADVICES = '_advices'  #: target advices attribute name


class AdviceError(Exception):
    """Handle Advice errors."""


class _Joinpoint(Joinpoint):
    """Manage target execution with Advices.

    Advices are callable objects which take in parameter a Joinpoint.

    Joinpoint provides to advices:
        - the target,
        - target call arguments as args and kwargs property,
        - a shared context during interception such as a dictionary.
    """

    def get_advices(self, target):

        result = get_advices(target, ctx=self.ctx)

        return result


def _add_advices(target, advices):
    """Add advices on input target.

    :param Callable target: target from where add advices.
    :param advices: advices to weave on input target.
    :type advices: routine or list.
    :param bool ordered: ensure advices to add will be done in input order
    """

    interception_fn = _get_function(target)

    target_advices = getattr(interception_fn, _ADVICES, [])
    target_advices += advices

    setattr(interception_fn, _ADVICES, target_advices)


def _remove_advices(target, advices, ctx):
    """Remove advices from input target.

    :param advices: advices to remove. If None, remove all advices.
    """
    # if ctx is not None
    if ctx is not None:  # check if intercepted ctx is ctx
        _, intercepted_ctx = get_intercepted(target)
        if intercepted_ctx is None or intercepted_ctx is not ctx:
            return

    interception_fn = _get_function(target)

    target_advices = getattr(interception_fn, _ADVICES, None)

    if target_advices is not None:

        if advices is None:
            target_advices = []
        else:
            target_advices = [
                advice for advice in target_advices if advice not in advices
            ]

        if target_advices:  # update target advices
            setattr(interception_fn, _ADVICES, target_advices)

        else:  # free target advices if necessary
            delattr(interception_fn, _ADVICES)
            _unapply_interception(target, ctx=ctx)


def get_advices(target, ctx=None, local=False):
    """Get element advices.

    :param target: target from where get advices.
    :param ctx: ctx from where get target.
    :param bool local: If ctx is not None or target is a method, if True
        (False by default) get only target advices without resolving super
        target advices in a super ctx.
    :return: list of advices.
    :rtype: list
    """

    result = []
    if is_intercepted(target):
        # find ctx if not given
        if ctx is None:
            ctx = find_ctx(target)
        # find advices among super ctx if same intercepted/interception
        if ctx is not None:
            # get target name
            target_name = target.__name__
            # resolve target and _target_ctx
            _target_ctx = ctx
            _target = getattr(_target_ctx, target_name, None)
            # check if _target is intercepted
            if is_intercepted(_target):
                # get intercepted_target in order to compare with super targets
                intercepted_target, intercepted_ctx = get_intercepted(_target)
                # if ctx is not a class
                if not isclass(_target_ctx):
                    if intercepted_ctx is _target_ctx:
                        interception_fn = _get_function(_target)
                        advices = getattr(interception_fn, _ADVICES, [])
                        result += advices
                        _target_ctx = _target_ctx.__class__
                        _target = getattr(_target_ctx, target_name, None)
                # climb back class hierarchy tree through all super targets
                while _target is not None and _target_ctx is not None:
                    # check if _target is intercepted
                    if is_intercepted(_target):
                        # get intercepted ctx
                        intercepted_fn, intercepted_ctx = get_intercepted(
                            _target
                        )
                        # if intercepted ctx is ctx
                        if intercepted_ctx is _target_ctx:
                            # get advices from _target interception
                            interception_fn = _get_function(_target)
                            advices = getattr(interception_fn, _ADVICES, [])
                            result += advices
                            # update _target
                            _target, _target_ctx = super_method(
                                name=target_name, ctx=_target_ctx
                            )
                        else:  # else _target_ctx is intercepted_ctx
                            _target_ctx = intercepted_ctx
                            # and update _target
                            _target = getattr(_target_ctx, target_name, None)
                    else:  # else, intercepted_fn is _target function
                        intercepted_fn = _get_function(_target)
                        _target, _target_ctx = super_method(
                            name=target_name, ctx=_target_ctx
                        )
                    # if intercepted are different, stop iteration
                    if intercepted_target is not intercepted_fn:
                        break

                    if local:  # break if local has been requested
                        break

        else:
            # get advices from interception function
            interception_function = _get_function(target)
            result = getattr(interception_function, _ADVICES, [])

    return result


def _namematcher(regex):
    """Checks if a target name matches with an input regular expression."""

    matcher = re_compile(regex)

    def match(target):
        target_name = getattr(target, '__name__', '')
        result = matcher.match(target_name)
        return result

    return match


def _publiccallable(target):
    """
    :return: True iif target is callable and name does not start with '_'
    """

    result = (
        callable(target)
        and not getattr(target, '__name__', '').startswith('_')
    )

    return result


def weave(
        target, advices, pointcut=None, ctx=None, depth=1, public=False,
        pointcut_application=None, ttl=None
):
    """Weave advices on target with input pointcut.

    :param callable target: target from where checking pointcut and
        weaving advices.
    :param advices: advices to weave on target.
    :param ctx: target ctx (class or instance).
    :param pointcut: condition for weaving advices on joinpointe.
        The condition depends on its type.
    :type pointcut:
        - NoneType: advices are weaved on target.
        - str: target name is compared to pointcut regex.
        - function: called with target in parameter, if True, advices will
            be weaved on target.
    :param int depth: class weaving depthing.
    :param bool public: (default True) weave only on public members.
    :param routine pointcut_application: routine which applies a pointcut when
        required. _Joinpoint().apply_pointcut by default. Such routine has
        to take in parameters a routine called target and its related
        function called function. Its result is the interception function.
    :param float ttl: time to leave for weaved advices.

    :return: the intercepted functions created from input target or a tuple
        with intercepted functions and ttl timer.
    :rtype: list

    :raises: AdviceError if pointcut is not None, not callable neither a str.
    """

    result = []

    # initialize advices
    if isroutine(advices):
        advices = [advices]

    if advices:
        # initialize pointcut
        # do nothing if pointcut is None or is callable
        if pointcut is None or callable(pointcut):
            pass
        # in case of str, use a name matcher
        elif isinstance(pointcut, string_types):
            pointcut = _namematcher(pointcut)
        else:
            error_msg = "Wrong pointcut to check weaving on {0}."
            error_msg = error_msg.format(target)
            advice_msg = "Must be None, or be a str or a function/method."
            right_msg = "Not {0}".format(type(pointcut))

            raise AdviceError(
                "{0} {1} {2}".format(error_msg, advice_msg, right_msg)
            )

        if ctx is None:
            ctx = find_ctx(elt=target)

        _weave(
            target=target, advices=advices, pointcut=pointcut, depth=depth,
            depth_predicate=_publiccallable if public else callable, ctx=ctx,
            intercepted=result, pointcut_application=pointcut_application
        )

        if ttl is not None:
            kwargs = {
                'target': target,
                'advices': advices,
                'pointcut': pointcut,
                'depth': depth,
                'public': public,
                'ctx': ctx
            }
            timer = Timer(ttl, unweave, kwargs=kwargs)
            timer.start()

            result = result, timer

    return result


def _weave(
        target, advices, pointcut, ctx, depth, depth_predicate, intercepted,
        pointcut_application
):
    """Weave deeply advices in target.

    :param callable target: target from where checking pointcut and
        weaving advices.
    :param advices: advices to weave on target.
    :param ctx: target ctx (class or instance).
    :param pointcut: condition for weaving advices on joinpointe.
        The condition depends on its type.
    :type pointcut:
        - NoneType: advices are weaved on target.
        - str: target name is compared to pointcut regex.
        - function: called with target in parameter, if True, advices will
            be weaved on target.
    :param int depth: class weaving depthing.
    :param list intercepted: list of intercepted targets.
    :param routine pointcut_application: routine which applies a pointcut when
        required. _Joinpoint().apply_pointcut by default. Such routine has
        to take in parameters a routine called target and its related
        function called function. Its result is the interception function.
    """

    # if weaving has to be done
    if pointcut is None or pointcut(target):
        # get target interception function
        interception_fn = _get_function(target)
        # does not handle not python functions
        if interception_fn is not None:
            # flag which specifies if poincut has to by applied
            # True if target is not intercepted
            apply_poincut = not is_intercepted(target)
            # apply poincut if not intercepted
            if (not apply_poincut) and ctx is not None:
                # apply poincut if ctx is not intercepted_ctx
                intercepted_fn, intercepted_ctx = get_intercepted(target)
                # if previous weave was done directly on the function
                if intercepted_ctx is None:
                    # update intercepted_ctx on target
                    intercepted_ctx = interception_fn._intercepted_ctx = ctx
                # if old ctx and the new one are same
                if ctx is not intercepted_ctx:
                    # apply pointcut
                    apply_poincut = True
                    # and update interception_fn
                    interception_fn = intercepted_fn
            # if weave has to be done
            if apply_poincut:
                # instantiate a new joinpoint if pointcut_application is None
                if pointcut_application is None:
                    pointcut_application = _Joinpoint().apply_pointcut
                interception_fn = pointcut_application(
                    target=target, function=interception_fn, ctx=ctx
                )
            # add advices to the interception function
            _add_advices(
                target=interception_fn, advices=advices
            )
            # append interception function to the intercepted ones
            intercepted.append(interception_fn)

    # search inside the target
    elif depth > 0:  # for an object or a class, weave on methods
        # get the right ctx
        if ctx is None:
            ctx = target
        for _, member in getmembers(ctx, depth_predicate):
            _weave(
                target=member, advices=advices, pointcut=pointcut,
                depth_predicate=depth_predicate, intercepted=intercepted,
                pointcut_application=pointcut_application, depth=depth - 1,
                ctx=ctx
            )


def unweave(
    target, advices=None, pointcut=None, ctx=None, depth=1, public=False,
):
    """Unweave advices on target with input pointcut.

    :param callable target: target from where checking pointcut and
        weaving advices.

    :param pointcut: condition for weaving advices on joinpointe.
        The condition depends on its type.
    :type pointcut:
        - NoneType: advices are weaved on target.
        - str: target name is compared to pointcut regex.
        - function: called with target in parameter, if True, advices will
            be weaved on target.

    :param ctx: target ctx (class or instance).
    :param int depth: class weaving depthing.
    :param bool public: (default True) weave only on public members

    :return: the intercepted functions created from input target.
    """

    # ensure advices is a list if not None
    if advices is not None:

        if isroutine(advices):
            advices = [advices]

    # initialize pointcut

    # do nothing if pointcut is None or is callable
    if pointcut is None or callable(pointcut):
        pass

    # in case of str, use a name matcher
    elif isinstance(pointcut, string_types):
        pointcut = _namematcher(pointcut)

    else:
        error_msg = "Wrong pointcut to check weaving on {0}.".format(target)
        advice_msg = "Must be None, or be a str or a function/method."
        right_msg = "Not {0}".format(type(pointcut))

        raise AdviceError(
            "{0} {1} {2}".format(error_msg, advice_msg, right_msg)
        )

    # get the right ctx
    if ctx is None:
        ctx = find_ctx(target)

    _unweave(
        target=target, advices=advices, pointcut=pointcut,
        ctx=ctx,
        depth=depth, depth_predicate=_publiccallable if public else callable
    )


def _unweave(target, advices, pointcut, ctx, depth, depth_predicate):
    """Unweave deeply advices in target."""

    # if weaving has to be done
    if pointcut is None or pointcut(target):
        # do something only if target is intercepted
        if is_intercepted(target):
            _remove_advices(target=target, advices=advices, ctx=ctx)

    # search inside the target
    if depth > 0:  # for an object or a class, weave on methods
        # get base ctx
        _base_ctx = None
        if ctx is not None:
            _base_ctx = base_ctx(ctx)
        for _, member in getmembers(target, depth_predicate):
            _unweave(
                target=member, advices=advices, pointcut=pointcut,
                depth=depth - 1, depth_predicate=depth_predicate, ctx=_base_ctx
            )


def weave_on(advices, pointcut=None, ctx=None, depth=1, ttl=None):
    """Decorator for weaving advices on a callable target.

    :param pointcut: condition for weaving advices on joinpointe.
        The condition depends on its type.
    :param ctx: target ctx (instance or class).
    :type pointcut:
        - NoneType: advices are weaved on target.
        - str: target name is compared to pointcut regex.
        - function: called with target in parameter, if True, advices will
            be weaved on target.

    :param depth: class weaving depthing
    :type depth: int

    :param public: (default True) weave only on public members
    :type public: bool
    """

    def __weave(target):
        """Internal weave function."""
        weave(
            target=target, advices=advices, pointcut=pointcut,
            ctx=ctx, depth=depth, ttl=ttl
        )

        return target

    return __weave
