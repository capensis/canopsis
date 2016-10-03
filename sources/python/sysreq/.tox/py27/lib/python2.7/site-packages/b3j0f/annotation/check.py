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

"""Interceptors dedicated to decorate decorations."""

from .core import Annotation
from .interception import PrivateInterceptor, PrivateCallInterceptor

from b3j0f.utils.iterable import ensureiterable

from types import FunctionType

from inspect import isclass, isroutine

from six import callable

__all__ = ['Condition', 'MaxCount', 'Target']


class Condition(PrivateInterceptor):
    """Apply a pre/post condition on an annotated element call."""

    class ConditionError(Exception):
        """Handle condition errors."""

    class PreConditionError(ConditionError):
        """Handle pre condition errors."""

    class PostConditionError(ConditionError):
        """Handle post condition errors."""

    #: pre condition attribute name
    PRE_COND = 'pre_cond'

    #: post condition attribute name
    POST_COND = 'post_cond'

    #: result attribute name
    RESULT = 'result'

    __slots__ = (PRE_COND, POST_COND) + PrivateInterceptor.__slots__

    def __init__(self, pre_cond=None, post_cond=None, *args, **kwargs):
        """
        :param pre_cond: function called before target call. Parameters
            are self annotation and AdvicesExecutor.
        :param post_cond: function called after target call. Parameters
            are self annotation, call result and AdvicesExecutor.
        """

        super(Condition, self).__init__(*args, **kwargs)

        self.pre_cond = pre_cond
        self.post_cond = post_cond

    def _interception(self, joinpoint):
        """Intercept call of joinpoint callee in doing pre/post conditions.
        """

        if self.pre_cond is not None:
            self.pre_cond(joinpoint)

        result = joinpoint.proceed()

        if self.post_cond is not None:
            joinpoint.exec_ctx[Condition.RESULT] = result
            self.post_cond(joinpoint)

        return result


class AnnotationChecker(PrivateInterceptor):
    """Annotation dedicated to intercept annotation target binding."""

    __slots = PrivateCallInterceptor.__slots__

    #: bind_target pointcut
    __BIND_TARGET__ = 'bind_target'

    def __init__(self, *args, **kwargs):

        super(AnnotationChecker, self).__init__(
            pointcut=AnnotationChecker.__BIND_TARGET__, *args, **kwargs
        )


class MaxCount(AnnotationChecker):
    """Set a maximum count of target instantiation by annotated element.
    """

    class Error(Exception):
        """Handle MaxCount errors."""

    __COUNT_KEY__ = 'count'

    DEFAULT_COUNT = 1

    __slots__ = (__COUNT_KEY__, ) + AnnotationChecker.__slots__

    def __init__(self, count=DEFAULT_COUNT, *args, **kwargs):
        """
        :param int count: maximal target instanciation by annotated element.
        """

        super(MaxCount, self).__init__(*args, **kwargs)

        self.count = MaxCount.DEFAULT_COUNT if count is None else count

    def _interception(self, joinpoint):

        target = joinpoint.kwargs['target']
        annotation = joinpoint.kwargs['self']

        annotation_class = annotation.__class__
        annotations = annotation_class.get_annotations(target)

        if len(annotations) >= self.count:
            raise MaxCount.Error(
                '{0} calls of {1} on {2}'.format(
                    self.count + 1, annotation, target)
            )

        result = joinpoint.proceed()

        return result

# apply MaxCount on itself
MaxCount()(MaxCount)


@MaxCount()
class Target(AnnotationChecker):
    """Check type of all decorated elements by this decorated Annotation in
    using a list of types.

    This list of types is checked related to logical rules such as OR and AND.

    - OR: true if at least one type is checked.
    - AND: true if all types are checked.
    """

    class Error(Exception):
        """Handle Target errors."""

    FUNC = FunctionType  #: function type.
    CALLABLE = callable  #: callable type.
    ROUTINE = 'routine'  #: routine type.

    OR = 'or'  #: or rule.
    AND = 'and'  #: and rule.

    DEFAULT_RULE = OR  #: default rule
    DEFAULT_INSTANCES = False  #: default instances condition

    TYPES = 'types'
    RULE = 'rule'
    INSTANCES = 'instances'

    __slots__ = (TYPES, RULE, INSTANCES) + PrivateCallInterceptor.__slots__

    def __init__(
            self, types=None, rule=DEFAULT_RULE, instances=False,
            *args, **kwargs
    ):
        """
        :param types: type(s) to check. The function ``callable`` can be used.
        :type types: type or list
        :param str rule: set condition on input types.
        :param bool instances: if True, check types such and instance types.
        """
        super(Target, self).__init__(*args, **kwargs)

        self.types = () if types is None else ensureiterable(types)
        self.rule = rule
        self.instances = instances

    def _interception(self, joinpoint):

        raiseexception = self.rule == Target.OR

        target = joinpoint.kwargs['target']
        annotation = joinpoint.kwargs['self']

        for _type in self.types:
            if (
                    (_type == Target.ROUTINE and isroutine(target))
                    or
                    (_type == Target.FUNC and isinstance(target, Target.FUNC))
                    or
                    (_type == type and isclass(target))
                    or
                    (_type is callable and callable(target))
                    or
                    (
                        _type is not callable
                        and (
                            (isclass(target) and issubclass(target, _type))
                            or (self.instances and isinstance(target, _type))
                        )
                    )
            ):
                if self.rule == Target.OR:
                    raiseexception = False
                    break

            elif self.rule == Target.AND:
                raiseexception = True
                break

        if raiseexception:
            interceptor_type = annotation.__class__

            raise Target.Error(
                "{0} is not allowed by {1}. Must be {2} {3}".format(
                    target,
                    interceptor_type,
                    'among' if self.rule == Target.OR else 'all',
                    self.types)
            )

        result = joinpoint.proceed()

        return result

# ensure AnnotationChecker is dedicated to Annotation
Target(Annotation)(AnnotationChecker)
