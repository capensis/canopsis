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

"""Decorators dedicated to class or functions calls."""

from __future__ import absolute_import

from .interception import PrivateInterceptor
from .check import Target

from b3j0f.utils.iterable import first
from b3j0f.utils.version import getcallargs

from six import get_function_code
from six.moves import range

from sys import stderr, maxsize

from time import sleep

from functools import wraps

__all__ = ['Types', 'types', 'Curried', 'curried', 'Retries', 'Memoize']


@Target(callable)
class Types(PrivateInterceptor):
    """
    Check routine parameters and return type.
    """

    class TypesError(Exception):
        """Handle Types error."""

    class SpecialCondition(object):
        """Handle SpecialCondition."""

        def __init__(self, _type):

            super(Types.SpecialCondition, self).__init__()

            self._type = _type

        def get_type(self):
            """Get special condition parameter type."""

            return self._type

    class NotNone(SpecialCondition):
        """Handle NotNone SpecialCondition."""

    class NotEmpty(SpecialCondition):
        """Handle NotEmpty SpecialCondition."""

    class NamedParameterType(object):
        """Handle Named Parameter Type."""

        def __init__(self, name, parameter_type):

            super(Types.NamedParameterType, self).__init__()

            self._name = name
            self._parameter_type = parameter_type

    class NamedParameterTypes(object):
        """Handle Named Parameter Types."""

        def __init__(self, target, named_parameter_types):

            super(Types.NamedParameterTypes, self).__init__()

            self._named_parameter_types = []

            target_code = get_function_code(target)

            for index in range(target_code.co_argcount):
                targetparamname = target_code.co_varnames[index]

                if targetparamname in named_parameter_types:
                    parameter_type = named_parameter_types[targetparamname]
                    named_parameter_type = Types.NamedParameterType(
                        targetparamname,
                        parameter_type
                    )
                    self._named_parameter_types.append(named_parameter_type)

                else:
                    self._named_parameter_types.append(None)

    #: return type attribute name
    RTYPE = 'rtype'

    #: parameter types attribute name
    PTYPES = 'ptypes'

    __slots__ = (RTYPE, PTYPES) + PrivateInterceptor.__slots__

    """
    Check parameter or result types of decorated class or function call.
    """
    def __init__(self, rtype=None, ptypes=None, *args, **kwargs):
        """
        :param rtype:
        """

        super(Types, self).__init__(*args, **kwargs)

        self.rtype = rtype
        self.ptypes = {} if ptypes is None else ptypes

    @staticmethod
    def check_value(value, expected_type):
        """Check Types parameters."""

        result = False

        if isinstance(expected_type, Types.NotNone):
            result = value is not None and Types.check_value(
                value,
                expected_type.get_type()
            )

        else:
            result = value is None

            if not result:

                value_type = type(value)

                if isinstance(expected_type, Types.NotEmpty):
                    try:
                        result = len(value) != 0
                        if result:
                            _type = expected_type.get_type()
                            result = Types.check_value(value, _type)
                    except TypeError:
                        result = False

                elif isinstance(expected_type, list):
                    result = issubclass(value_type, list)

                    if result:
                        if len(expected_type) == 0:
                            result = len(value) == 0
                        else:
                            _expected_type = expected_type[0]

                            for item in value:
                                result = Types.check_value(
                                    item,
                                    _expected_type
                                )

                                if not result:
                                    break

                elif isinstance(expected_type, set):
                    result = issubclass(value_type, set)

                    if result:
                        if len(expected_type) == 0:
                            result = len(value) == 0
                        else:

                            _expected_type = expected_type.copy().pop()
                            _value = value.copy()

                            value_length = len(_value)

                            for _ in range(value_length):
                                item = _value.pop()
                                result = Types.check_value(
                                    item,
                                    _expected_type)

                                if not result:
                                    break
                else:
                    result = issubclass(value_type, expected_type)

        return result

    def _interception(self, joinpoint):

        target = joinpoint.target
        args = joinpoint.args
        kwargs = joinpoint.kwargs

        if self.ptypes:
            callargs = getcallargs(target, *args, **kwargs)

            for arg in callargs:
                value = callargs[arg]
                expected_type = self.ptypes.get(arg)

                if (
                        expected_type is not None and
                        not Types.check_value(value, expected_type)
                ):
                    raise Types.TypesError(
                        "wrong typed parameter for arg {0} : {1} ({2}). \
                        Expected: {3}."
                        .format(
                            arg, value, type(value), expected_type
                        )
                    )

        result = joinpoint.proceed()

        target = joinpoint.target
        args = joinpoint.args
        kwargs = joinpoint.kwargs

        if self.rtype:
            if not Types.check_value(result, self.rtype):
                raise Types.TypesError(
                    "wrong result type for {0} with parameters {1}, {2}: {3} \
                    ({4}). Expected {5}."
                    .format(
                        target, args, kwargs, result, type(result),
                        self.rtype
                    )
                )

        return result


def types(*args, **kwargs):
    """Quick alias for the Types Annotation with only args and kwargs
    parameters.

    :param tuple args: may contain rtype.
    :param dict kwargs: may contain ptypes.
    """

    rtype = first(args)

    return Types(rtype=rtype, ptypes=kwargs)


class Curried(PrivateInterceptor):
    """Annotation that returns a function that keeps returning functions
    until all arguments are supplied; then the original function is
    evaluated.

    Inspirated from Jeff Laughlin Consulting LLC projects.
    """

    ARGS = 'args'  #: args attribute name
    KWARGS = 'kwargs'  #: kwargs attribute name
    DEFAULT_ARGS = 'default_args'  #: default args attribute name
    DEFAULT_KWARGS = 'default_kwargs'  #: default kwargs attribute name

    __slots__ = (
        ARGS, KWARGS, DEFAULT_ARGS, DEFAULT_KWARGS
    ) + PrivateInterceptor.__slots__

    class CurriedResult(object):
        """Curried result in case of missing arguments."""

        __slots__ = ('curried', 'exception')

        def __init__(self, curried, exception):

            super(Curried.CurriedResult, self).__init__()

            self.curried = curried
            self.exception = exception

    def __init__(self, varargs=None, keywords=None, *args, **kwargs):
        """
        :param tuple varargs: function call varargs.
        :param dict keywords: function call keywords.
        """

        super(Curried, self).__init__(*args, **kwargs)

        # initialize arguments
        if varargs is None:
            varargs = ()
        if keywords is None:
            keywords = {}
        # set attributes
        self.args = self.default_args = varargs
        self.kwargs = self.default_kwargs = keywords

    def _bind_target(self, target, *args, **kwargs):

        @wraps(target)
        def wrapper(*args, **kwargs):
            """Target wrapper."""

            return target(*args, **kwargs)

        result = super(Curried, self)._bind_target(
            target=wrapper, *args, **kwargs
        )

        return result

    def _interception(self, joinpoint, *args, **kwargs):

        result = None

        target = joinpoint.target

        args = joinpoint.args
        kwargs = joinpoint.kwargs

        self.kwargs.update(kwargs)
        self.args += args

        try:
            # check if all arguments are given
            getcallargs(target, *self.args, **self.kwargs)
            joinpoint.args = self.args
            joinpoint.kwargs = self.kwargs
            result = joinpoint.proceed()
        except TypeError as ex:
            # in case of problem, returns curried decorater and exception
            result = Curried.CurriedResult(self, ex)

        return result


def curried(*args, **kwargs):
    """Curried annotation with varargs and kwargs.
    """

    return Curried(varargs=args, keywords=kwargs)


def example_exc_handler(tries_remaining, exception, delay):
    """Example exception handler; prints a warning to stderr.

    tries_remaining: The number of tries remaining.
    exception: The exception instance which was raised.
    """

    print >> stderr, "Caught '{0}', {1} tries remaining, \
    sleeping for {2} seconds".format(exception, tries_remaining, delay)


class Retries(PrivateInterceptor):
    """Function decorator implementing retrying logic.

    condition: retry condition, among execution success or failure or both.
    delay: Sleep this many seconds * backoff * try number after failure
    backoff: Multiply delay by this factor after each failure
    exceptions: A tuple of exception classes; default (Exception,)
    hook: A function with the signature myhook(data, condition, tries_remaining
    , mydelay) where data is result function or raised Exception, condition is
    ON_ERROR or ON_SUCCESS depending on error or success execution function,
    tries_remaining is tries remaining, and finally, mydelay is waiting seconds
    between calls; default None.

    The decorator will call the function up to max_tries times if it raises
    an exception or if it simply execute the function, depending on state
    condition.

    By default it catches instances of the Exception class and subclasses.
    This will recover after all but the most fatal errors. You may specify a
    custom tuple of exception classes with the 'exceptions' argument; the
    function will only be retried if it raises one of the specified
    exceptions.

    Additionally you may specify a hook function which will be called prior
    to retrying with the number of remaining tries and the exception instance;
    see given example. This is primarily intended to give the opportunity to
    log the failure. Hook is not called after failure if no retries remain.
    """

    MAX_TRIES = 'max_tries'  #: max_tries attribute name.
    DELAY = 'delay'  #: delay attribute name.
    BACKOFF = 'backoff'  #: backoff attribute name.
    EXCEPTIONS = 'exceptions'  #: exceptions attribute name.
    HOOK = 'hook'  #: hook attribute name.
    CONDITION = 'condition'  #: condition attribute name.

    DEFAULT_DELAY = 1
    DEFAULT_BACKOFF = 2
    DEFAULT_EXCEPTIONS = (Exception,)

    ON_ERROR = 1  #: on error retries condition.
    ON_SUCCESS = 2  #: on success retries condition.
    ALL = ON_ERROR | ON_SUCCESS  #: all retries condition.

    __slots__ = (
        MAX_TRIES, DELAY, BACKOFF, EXCEPTIONS, HOOK, CONDITION
    ) + PrivateInterceptor.__slots__

    def __init__(
            self,
            max_tries,
            delay=DEFAULT_DELAY,
            backoff=DEFAULT_BACKOFF,
            exceptions=DEFAULT_EXCEPTIONS,
            hook=None,
            condition=ALL,
            *args, **kwargs
    ):

        super(Retries, self).__init__(*args, **kwargs)

        self.max_tries = max_tries
        self.delay = delay
        self.backoff = backoff
        self.exceptions = exceptions
        self.hook = hook
        self.condition = condition

    def _interception(self, joinpoint):

        result = None

        mydelay = self.delay

        for tries_remaining in range(self.max_tries - 1, -1, -1):

            try:
                result = joinpoint.proceed()

            except self.exceptions as ex:

                mydelay = self._checkretry(
                    mydelay=mydelay, condition=Retries.ON_ERROR, data=ex,
                    tries_remaining=tries_remaining
                )

            else:

                mydelay = self._checkretry(
                    mydelay=mydelay, condition=Retries.ON_SUCCESS, data=result,
                    tries_remaining=tries_remaining
                )

                if mydelay is None:
                    break  # stop execution if mydelay is None

        return result

    def _checkretry(self, mydelay, condition, tries_remaining, data):
        """Check if input parameters allow to retries function execution.

        :param float mydelay: waiting delay between two execution.
        :param int condition: condition to check with this condition.
        :param int tries_remaining: tries remaining.
        :param data: data to hook.
        """

        result = mydelay

        if self.condition & condition and tries_remaining > 0:

            # hook data with tries_remaining and mydelay
            if self.hook is not None:
                self.hook(data, condition, tries_remaining, mydelay)
            # wait mydelay seconds
            sleep(mydelay)
            result *= self.backoff  # increment mydelay with this backoff

        elif condition is Retries.ON_ERROR:
            raise data  # raise data if no retries and on_error

        else:  # else Nonify mydelay to prevent callee function to stop
            result = None

        return result


class Memoize(PrivateInterceptor):
    """Save funtion results related to called parameters.

    Parameters must be hashable."""

    MAX_SIZE = 'max_size'  #: max size result.
    _CACHE = '_cache'  #: cache object which stores results and params.

    DEFAULT_MAX_SIZE = maxsize  #: default max size value.

    __slots__ = (MAX_SIZE, _CACHE) + PrivateInterceptor.__slots__

    def __init__(self, max_size=DEFAULT_MAX_SIZE, *args, **kwargs):

        super(Memoize, self).__init__(*args, **kwargs)

        self.max_size = max_size
        self._cache = {}

    def _getkey(self, args, kwargs):
        """Get hash key from args and kwargs.

        args and kwargs must be hashable.

        :param tuple args: called vargs.
        :param dict kwargs: called keywords.
        :return: hash(tuple(args) + tuple((key, val) for key in sorted(kwargs)).
        :rtype: int."""

        values = list(args)

        keys = sorted(list(kwargs))

        for key in keys:
            values.append((key, kwargs[key]))

        result = hash(tuple(values))

        return result

    def _interception(self, joinpoint):

        result = None

        args = joinpoint.args
        kwargs = joinpoint.kwargs
        key = self._getkey(args, kwargs)

        _cache = self._cache

        if key in _cache:
            _, _, result = _cache[key]

        else:
            result = joinpoint.proceed()

            if len(self._cache) < self.max_size:
                _cache[key] = (args, kwargs, result)

        return result

    def getparams(self, result):
        """Get result parameters.

        :param result: cached result.
        :raises: ValueError if result is not cached.
        :return: args and kwargs registered with input result.
        :rtype: tuple"""

        for key in self._cache:
            if self._cache[key][2] == result:
                args, kwargs, _ = self._cache[key]
                return args, kwargs

        else:
            raise ValueError('Result is not cached')

    def clearcache(self):
        """Clear cache."""

        self._cache.clear()
