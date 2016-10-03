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

"""Annotations dedicated to object oriented programming."""

from six import PY2, get_method_function, get_function_code

from .core import Annotation
from .interception import PrivateInterceptor

from types import MethodType

from inspect import getmembers, isfunction, ismethod, isclass

from warnings import warn_explicit

__all__ = ['Transform', 'Mixin', 'Deprecated', 'Singleton', 'MethodMixin']


class Transform(Annotation):
    """Transform a class into an annotation or something else if parameters are
    different.
    """

    NAME = 'name'
    BASES = 'bases'
    DICT = 'dict'
    UPDATE = 'update'

    __slots__ = (NAME, BASES, DICT, UPDATE) + Annotation.__slots__

    def __init__(
            self, name=None, bases=None, _dict=None, update=True,
            *args, **kwargs
    ):
        """According to type constructor parameters, you can specifiy name,
        bases and dict.

        :param str name: new class name.
        :param tuple bases: new class bases.
        :param dict dict: class __dict__.
        :param bool update: if True, update new properties on the new class.
            Otherwise replace old ones.
        """

        super(Transform, self).__init__(*args, **kwargs)

        self.name = name
        self.bases = (Annotation,) if bases is None else bases
        self.dict = {} if _dict is None else _dict
        self.update = update

    def _bind_target(self, target, ctx=None):

        result = target
        # do something only for class objects
        if isclass(target):

            name = target.__name__ if self.name is None else self.name
            # if update is True, update target properties
            bases = tuple(set((target,) + self.bases))

            if self.update:
                _dict = target.__dict__.copy()
                _dict.update(self.dict)

            else:  # else use new properties
                _dict = self.dict

            # get the new type related to name, bases and dict
            result = type(name, bases, _dict)

        else:
            raise TypeError(
                'Wrong annotated element type {0}, class expected.'
                .format(target)
            )

        return result


class Mixin(Annotation):
    """Annotation which enrichs a target with Mixin.

    For every defined mixin, a private couple of (name, array of mixed items)
    is created into the target in order to go back in a no mixin state.
    """

    class MixInError(Exception):
        """Raised for any Mixin error."""

    #: key to organize a dictionary of mixed elements by name in targets.
    __MIXEDIN_KEY__ = '__MIXEDIN__'

    #: Key to identify if a mixin add or replace a content.
    __NEW_CONTENT_KEY__ = '__NEW_CONTENT__'

    def __init__(self, classes=(), *attributes, **named_attributes):

        super(Mixin, self).__init__()

        self.classes = classes if isinstance(classes, tuple) else (classes,)
        self.attributes = attributes
        self.named_attributes = named_attributes

    def on_bind_target(self, target, ctx=None):

        super(Mixin, self).on_bind_target(target)

        for cls in self.classes:
            Mixin.mixin(target, cls)

        for attribute in self.attributes:
            Mixin.mixin(target, attribute)

        for name in self.named_attributes:
            attribute = self.named_attributes[name]
            Mixin.mixin(target, attribute, name)

    @staticmethod
    def mixin_class(target, cls):
        """Mix cls content in target."""

        for name, field in getmembers(cls):
            Mixin.mixin(target, field, name)

    @staticmethod
    def mixin_function_or_method(target, routine, name=None, isbound=False):
        """Mixin a routine into the target.

        :param routine: routine to mix in target.
        :param str name: mixin name. Routine name by default.
        :param bool isbound: If True (False by default), the mixin result is a
            bound method to target.
        """

        function = None

        if isfunction(routine):
            function = routine

        elif ismethod(routine):
            function = get_method_function(routine)

        else:
            raise Mixin.MixInError(
                "{0} must be a function or a method.".format(routine))

        if name is None:
            name = routine.__name__

        if not isclass(target) or isbound:
            _type = type(target)
            method_args = [function, target]

            if PY2:
                method_args += _type

            result = MethodType(*method_args)

        else:

            if PY2:
                result = MethodType(function, None, target)

            else:
                result = function

        Mixin.set_mixin(target, result, name)

        return result

    @staticmethod
    def mixin(target, resource, name=None):
        """Do the correct mixin depending on the type of input resource.

        - Method or Function: mixin_function_or_method.
        - class: mixin_class.
        - other: set_mixin.

        And returns the result of the choosen method (one or a list of mixins).
        """

        result = None

        if ismethod(resource) or isfunction(resource):
            result = Mixin.mixin_function_or_method(target, resource, name)

        elif isclass(resource):
            result = list()

            for name, content in getmembers(resource):
                if isclass(content):
                    mixin = Mixin.set_mixin(target, content, name)

                else:
                    mixin = Mixin.mixin(target, content, name)

                result.append(mixin)

        else:
            result = Mixin.set_mixin(target, resource, name)

        return result

    @staticmethod
    def get_mixedins_by_name(target):
        """Get a set of couple (name, field) of target mixedin."""

        result = getattr(target, Mixin.__MIXEDIN_KEY__, None)

        if result is None:
            result = dict()
            setattr(target, Mixin.__MIXEDIN_KEY__, result)

        return result

    @staticmethod
    def set_mixin(target, resource, name=None, override=True):
        """Set a resource and returns the mixed one in target content or
        Mixin.__NEW_CONTENT_KEY__ if resource name didn't exist in target.

        :param str name: target content item to mix with the resource.
        :param bool override: If True (default) permits to replace an old
            resource by the new one.
        """

        if name is None and not hasattr(resource, '__name__'):
            raise Mixin.MixInError(
                "name must be given or resource {0} can't be anonymous"
                .format(resource))

        result = None

        name = resource.__name__ if name is None else name

        mixedins_by_name = Mixin.get_mixedins_by_name(target)

        result = getattr(target, name, Mixin.__NEW_CONTENT_KEY__)

        if override or result == Mixin.__NEW_CONTENT_KEY__:

            if name not in mixedins_by_name:
                mixedins_by_name[name] = (result,)
            else:
                mixedins_by_name[name] += (result,)

            try:
                setattr(target, name, resource)

            except (AttributeError, TypeError):

                if len(mixedins_by_name[name]) == 1:
                    del mixedins_by_name[name]
                    if len(mixedins_by_name) == 0:
                        delattr(target, Mixin.__MIXEDIN_KEY__)

                else:
                    mixedins_by_name[name] = mixedins_by_name[name][:-2]

                result = None
        else:
            result = None

        return result

    @staticmethod
    def remove_mixin(target, name, mixedin=None, replace=True):
        """Remove a mixin with name (and reference) from targetand returns the
        replaced one or None.

        :param mixedin: a mixedin value or the last defined mixedin if is None
            (by default).
        :param bool replace: If True (default), the removed mixedin replaces
            the current mixin.
        """

        try:
            result = getattr(target, name)

        except AttributeError:
            raise Mixin.MixInError(
                "No mixin {0} exists in {1}".format(name, target)
            )

        mixedins_by_name = Mixin.get_mixedins_by_name(target)

        mixedins = mixedins_by_name.get(name)

        if mixedins:

            if mixedin is None:
                mixedin = mixedins[-1]
                mixedins = mixedins[:-2]

            else:

                try:
                    index = mixedins.index(mixedin)

                except ValueError:
                    raise Mixin.MixInError(
                        "Mixedin {0} with name {1} does not exist \
                        in target {2}"
                        .format(mixedin, name, target)
                    )

                mixedins = mixedins[0:index] + mixedins[index + 1:]

            if len(mixedins) == 0:
                # force to replace/delete the mixin even if replace is False
                # in order to stay in a consistent state
                if mixedin != Mixin.__NEW_CONTENT_KEY__:
                    setattr(target, name, mixedin)

                else:
                    delattr(target, name)

                del mixedins_by_name[name]

            else:
                if replace:
                    setattr(target, name, mixedin)

                mixedins_by_name[name] = mixedins

        else:
            # shouldn't be raised except if removing has been done
            # manually
            raise Mixin.MixInError(
                "No mixin {0} exists in {1}".format(name, target))

        # clean mixedins if no one exists
        if len(mixedins_by_name) == 0:
            delattr(target, Mixin.__MIXEDIN_KEY__)

        return result

    @staticmethod
    def remove_mixins(target):
        """Tries to get back target in a no mixin consistent state.
        """

        mixedins_by_name = Mixin.get_mixedins_by_name(target).copy()

        for _name in mixedins_by_name:  # for all named mixins
            while True:  # remove all mixins named _name
                try:
                    Mixin.remove_mixin(target, _name)
                except Mixin.MixInError:
                    break  # leave the loop 'while'


class MethodMixin(Annotation):
    """Apply a mixin on a method."""

    def __init__(self, function, *args, **kwargs):

        super(MethodMixin, self).__init__(*args, **kwargs)

        self.function = function

    def _bind_target(self, target, ctx, *args, **kwargs):

        result = super(MethodMixin, self)._bind_target(
            target=target, ctx=ctx, *args, **kwargs
        )

        if ismethod(result):
            cls = target.im_class
            name = target.__name__
            Mixin.mixin_function_or_method(cls, self.function, name)
            result = target

        else:

            if ctx is not None:
                cls = ctx
                name = target.__name__
                Mixin.mixin_function_or_method(cls, self.function, name)

            result = self.function

        return result


class Deprecated(PrivateInterceptor):
    """Decorator which can be used to mark functions as deprecated. It will
    result in a warning being emitted when the function is used.
    """

    def _interception(self, joinpoint):

        target = joinpoint.target
        target_code = get_function_code(target)
        warn_explicit(
            "Call to deprecated function {0}.".format(target.__name__),
            category=DeprecationWarning,
            filename=target_code.co_filename,
            lineno=target_code.co_firstlineno + 1
        )
        result = joinpoint.proceed()
        return result


class Singleton(Annotation):
    """Transforms cls into a singleton.

    Reference to cls, or to any instance is the same reference.
    """

    def __init__(self, *args, **kwargs):
        super(Singleton, self).__init__()
        self.args = args
        self.kwargs = kwargs

    def _bind_target(self, target, ctx, *args, **kwargs):

        target = super(Singleton, self)._bind_target(
            target, ctx=ctx, *args, **kwargs
        )

        instance = target(*self.args, **self.kwargs)
        instance.__call__ = lambda x=None: instance
        target.__call__ = lambda x=None: instance

        return instance
