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

"""Advice utilities."""

from uuid import uuid4 as uuid

from .core import weave, unweave, get_advices

__all__ = ['Advice']


class Advice(object):
    """Advice class which aims to embed an advice function with disabling
    property.
    """

    __slots__ = ('_impl', '_enable', '_uid')

    def __init__(self, impl, uid=None, enable=True):

        self._impl = impl
        self._enable = enable
        self._uid = uuid() if uid is None else uid

    @property
    def uid(self):
        """Get advice uid."""

        return self._uid

    @property
    def enable(self):
        """Get self enable state. Change state if input enable is a boolean.

        TODO: change of method execution instead of saving a state.
        """

        return self._enable

    @enable.setter
    def enable(self, value):
        """Change of enable status."""

        self._enable = value

    def apply(self, joinpoint):
        """Apply this advice on input joinpoint.

        TODO: improve with internal methods instead of conditional test.
        """

        if self._enable:
            result = self._impl(joinpoint)

        else:
            result = joinpoint.proceed()

        return result

    @staticmethod
    def set_enable(target, enable=True, advice_ids=None):
        """Enable or disable all target Advices designated by input advice_ids.

        If advice_ids is None, apply (dis|en)able state to all advices.
        """

        advices = get_advices(target)

        for advice in advices:
            try:
                if isinstance(Advice) \
                        and (advice_ids is None or advice.uuid in advice_ids):
                    advice.enable = enable
            except ValueError:
                pass

    @staticmethod
    def weave(target, advices, pointcut=None, depth=1, public=False):
        """Weave advices such as Advice objects."""

        advices = (
            advice if isinstance(advice, Advice) else Advice(advice)
            for advice in advices
        )

        weave(
            target=target, advices=advices, pointcut=pointcut,
            depth=depth, public=public
        )

    @staticmethod
    def unweave(target, *advices):
        """Unweave advices from input target."""

        advices = (
            advice if isinstance(advice, Advice) else Advice(advice)
            for advice in advices
        )

        unweave(target=target, *advices)

    def __call__(self, joinpoint):
        """Shortcut for self apply."""

        return self.apply(joinpoint)

    def __hash__(self):
        """Return self uid hash."""

        result = hash(self._uid)

        return result

    def __eq__(self, other):
        """Compare with self uid."""

        result = isinstance(other, Advice) and other.uid == self._uid

        return result
