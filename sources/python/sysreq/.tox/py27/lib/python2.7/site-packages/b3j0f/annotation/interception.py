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

"""Definition of annotation dedicated to intercept annotated element calls."""

from .core import Annotation

from b3j0f.aop import weave, unweave

__all__ = [
    'Interceptor',
    'PrivateInterceptor', 'CallInterceptor', 'PrivateCallInterceptor'
]


class Interceptor(Annotation):
    """Annotation able to intercept annotated elements.

    This interception can be disabled at any time and specialized with a
        pointcut.
    """

    #: interception attribute name
    INTERCEPTION = 'interception'

    #: pointcut attribute name
    POINTCUT = 'pointcut'

    #: attribute name for enable the interception
    ENABLE = 'enable'

    #: interceptor attribute name
    INTERCEPTOR = '__interceptor__'

    #: private attribute name for pointcut
    _POINTCUT = '_pointcut'

    __slots__ = (
        INTERCEPTION, ENABLE,  # public attributes
        _POINTCUT  # private attribute
    ) + Annotation.__slots__

    class InterceptorError(Exception):
        """Handle Interceptor errors."""

    def __init__(
            self, interception=None, pointcut=None, enable=True,
            *args, **kwargs
    ):
        """Default constructor with interception function and enable property.

        :param callable interception: called if self is enabled and if any
            annotated element is called. Parameters of call are self annotation
            and an AdvicesExecutor.
        :param pointcut: pointcut to use in order to weave interception.
        :param bool enable:
        """

        super(Interceptor, self).__init__(*args, **kwargs)

        self.interception = interception
        self._pointcut = pointcut
        self.enable = enable

    @property
    def pointcut(self):
        """Get pointcut."""

        return self._pointcut

    @pointcut.setter
    def pointcut(self, value):
        """Change of pointcut.
        """

        pointcut = getattr(self, Interceptor.POINTCUT)

        # for all targets
        for target in self.targets:
            # unweave old advices
            unweave(target, pointcut=pointcut, advices=self.intercepts)
            # weave new advices with new pointcut
            weave(target, pointcut=value, advices=self.intercepts)

        # and save new pointcut
        setattr(self, Interceptor._POINTCUT, value)

    def _bind_target(self, target, ctx=None, *args, **kwargs):
        """Weave self.intercepts among target advices with pointcut."""

        result = super(Interceptor, self)._bind_target(
            target=target, ctx=ctx, *args, **kwargs
        )

        pointcut = getattr(self, Interceptor.POINTCUT)

        weave(result, pointcut=pointcut, advices=self.intercepts, ctx=ctx)

        return result

    def remove_from(self, target, ctx=None, *args, **kwargs):

        super(Interceptor, self).remove_from(target, *args, **kwargs)

        unweave(
            target, pointcut=self.pointcut, advices=self.intercepts, ctx=ctx
        )

    def intercepts(self, joinpoint):
        """Self target interception if self is enabled

        :param joinpoint: advices executor
        """

        result = None

        if self.enable:

            interception = getattr(self, Interceptor.INTERCEPTION)

            joinpoint.exec_ctx[Interceptor.INTERCEPTION] = self

            result = interception(joinpoint)

        else:
            result = joinpoint.proceed()

        return result

    @classmethod
    def set_enable(cls, target, enable=True):
        """(Dis|En)able annotated interceptors."""

        interceptors = cls.get_annotations(target)

        for interceptor in interceptors:
            setattr(interceptor, Interceptor.ENABLE, enable)


class PrivateInterceptor(Interceptor):
    """Interceptor with a private interception resource.
    """

    __slots__ = Interceptor.__slots__

    def __init__(self, *args, **kwargs):
        """
        Do nothing except set self._interception such as its interception
        """
        super(PrivateInterceptor, self).__init__(
            interception=self._interception, *args, **kwargs
        )

    def _interception(self, joinpoint):
        """
        Default interception which raises a NotImplementedError.

        :param Annotation annotation: intercepted annotation
        """

        raise NotImplementedError()


class CallInterceptor(Interceptor):
    """Interceptor dedicated to intercept call of annotated element.

    Instead of Interceptor, the pointcut equals '__call__', therefore, the
    target may be a class instead of an instance beceause python does not
    handle system methods (such as __call__, __repr__, etc.) if they are
    overriden on instance.

    It could be used to intercepts annotation binding annotated by self.
    """

    __CALL__ = '__call__'

    __slots__ = Interceptor.__slots__

    def __init__(self, *args, **kwargs):

        super(CallInterceptor, self).__init__(
            pointcut=CallInterceptor.__CALL__, *args, **kwargs
        )


class PrivateCallInterceptor(CallInterceptor):
    """Interceptor dedicated to apply a private interception on target calls.

    Instead of Interceptor, the pointcut equals '__call__', therefore, the
    target may be a class instead of an instance beceause python does not
    handle system methods (such as __call__, __repr__, etc.) if they are
    overriden on instance.
    """

    __slots__ = CallInterceptor.__slots__ + PrivateInterceptor.__slots__

    def __init__(self, *args, **kwargs):

        super(PrivateCallInterceptor, self).__init__(
            interception=self._interception, *args, **kwargs
        )

    def _interception(self, joinpoint):
        """
        Default interception which raise a NotImplementedError.
        """

        raise NotImplementedError()
