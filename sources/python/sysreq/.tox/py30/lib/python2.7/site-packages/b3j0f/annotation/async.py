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

"""Decorators dedicated to asynchronous programming."""

from __future__ import absolute_import

try:
    from threading import Thread, RLock
except ImportError:
    from dummythreading import Thread, RLock

from time import sleep

from signal import signal, SIGALRM, alarm

from six import callable
from six.moves.queue import Queue

from .core import Annotation
from .interception import PrivateInterceptor
from .oop import Mixin

__all__ = [
    'Synchronized', 'SynchronizedClass',
    'Asynchronous', 'TimeOut', 'Wait', 'Observable'
]


class Synchronized(PrivateInterceptor):
    """Transform a target into a thread safe target."""

    #: lock attribute name
    _LOCK = '_lock'

    __slots__ = (_LOCK,) + PrivateInterceptor.__slots__

    def __init__(self, lock=None, *args, **kwargs):

        super(Synchronized, self).__init__(*args, **kwargs)

        self._lock = RLock() if lock is None else lock

    def _interception(self, joinpoint):

        self._lock.acquire()

        result = joinpoint.proceed()

        self._lock.release()

        return result


class SynchronizedClass(Synchronized):
    """Transform a class into a thread safe class."""

    def on_bind_target(self, target, ctx=None):

        for attribute in target.__dict__:
            if callable(attribute):
                Synchronized(attribute, self._lock)


class Asynchronous(Annotation):
    """Transform a target into an asynchronous callable target."""

    def __init__(self, *args, **kwargs):

        super(Asynchronous, self).__init__(*args, **kwargs)

        self.queue = None

    def _threaded(self, *args, **kwargs):
        """Call the target and put the result in the Queue."""

        for target in self.targets:
            result = target(*args, **kwargs)
            self.queue.put(result)

    def on_bind_target(self, target, ctx=None):

        # add start function to wrapper
        super(Asynchronous, self).on_bind_target(target, ctx=ctx)

        setattr(target, 'start', self.start)

    def start(self, *args, **kwargs):
        """Start execution of the function."""

        self.queue = Queue()
        thread = Thread(target=self._threaded, args=args, kwargs=kwargs)
        thread.start()

        return Asynchronous.Result(self.queue, thread)

    class NotYetDoneException(Exception):
        """Handle when a result is not yet available."""

    class Result(object):
        """In charge of receive asynchronous function result."""

        __slots__ = ('queue', 'thread', 'result')

        def __init__(self, queue, thread):

            super(Asynchronous.Result, self).__init__()

            self.result = None
            self.queue = queue
            self.thread = thread

        def is_done(self):
            """True if result is available."""

            return not self.thread.is_alive()

        def get_result(self, wait=-1):
            """Get result value.

            Wait for it if necessary.

            :param int wait: maximum wait time.
            :return: result value.
            """

            if not self.is_done():

                if wait >= 0:
                    self.thread.join(wait)

                else:
                    raise Asynchronous.NotYetDoneException(
                        'the call has not yet completed its task'
                    )

            if self.result is None:
                self.result = self.queue.get()

            return self.result


class TimeOut(PrivateInterceptor):
    """Raise an Exception if the target call has not finished in time."""

    class TimeOutError(Exception):
        """Exception thrown if time elapsed before the end of the target call.
        """

        #: Default time out error message.
        DEFAULT_MESSAGE = \
            'Call of {0} with parameters {1} and {2} is timed out in frame {3}'

        def __init__(self, timeout_interceptor, frame):

            super(TimeOut.TimeOutError, self).__init__(
                timeout_interceptor.message.format(
                    timeout_interceptor.target,
                    timeout_interceptor.args,
                    timeout_interceptor.kwargs,
                    frame
                )
            )

    SECONDS = 'seconds'
    ERROR_MESSAGE = 'error_message'

    __slots__ = (SECONDS, ERROR_MESSAGE) + PrivateInterceptor.__slots__

    def __init__(
            self,
            seconds, error_message=TimeOutError.DEFAULT_MESSAGE,
            *args, **kwargs
    ):

        super(TimeOut, self).__init__(*args, **kwargs)

        self.seconds = seconds
        self.error_message = error_message

    def _handle_timeout(self, frame=None, **_):
        """Sig ALARM timeout function."""

        raise TimeOut.TimeOutError(self, frame)

    def _interception(self, joinpoint):

        signal(SIGALRM, self._handle_timeout)
        alarm(self.seconds)

        try:
            result = joinpoint.proceed()

        finally:
            alarm(0)

        return result


class Wait(PrivateInterceptor):
    """Define a time to wait before and after a target call."""

    DEFAULT_BEFORE = 1  #: default seconds to wait before the target call.
    DEFAULT_AFTER = 1  #: default seconds to wait after the target call.

    BEFORE = 'before'  #: before attribute name.

    AFTER = 'after'  #: after attribute name.

    __slots__ = (BEFORE, AFTER) + PrivateInterceptor.__slots__

    def __init__(
            self, before=DEFAULT_BEFORE, after=DEFAULT_AFTER, *args, **kwargs
    ):

        super(Wait, self).__init__(*args, **kwargs)

        self.before = before
        self.after = after

    def _interception(self, joinpoint):

        sleep(self.before)

        result = joinpoint.proceed()

        sleep(self.after)

        return result


class Observable(PrivateInterceptor):
    """Imlementation of the observer design pattern.

    It transforms a target into an observable object in adding method
    register_observer, unregister_observer and notify_observers.
    Observers listen to pre/post target interception.
    """

    def __init__(self, *args, **kwargs):

        super(Observable, self).__init__(*args, **kwargs)

        self.observers = set()

    def register_observer(self, observer):
        """Register an observer."""

        self.observers.add(observer)

    def unregister_observer(self, observer):
        """Unregister an observer."""

        self.observers.remove(observer)

    def notify_observers(self, joinpoint, post=False):
        """Notify observers with parameter calls and information about
        pre/post call.
        """

        _observers = tuple(self.observers)

        for observer in _observers:
            observer.notify(joinpoint=joinpoint, post=post)

    def on_bind_target(self, target, ctx=None):

        Mixin.set_mixin(target, self.register_observer)
        Mixin.set_mixin(target, self.unregister_observer)
        Mixin.set_mixin(target, self.notify_observers)

    def _interception(self, joinpoint):

        self.notify_observers(joinpoint=joinpoint)

        result = joinpoint.proceed()

        self.notify_observers(joinpoint=joinpoint, post=True)

        return result
