# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2016 Jonathan Labéjof <jonathan.labejof@gmail.com>
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

"""Read module."""

__all__ = ['Read', 'Cursor', 'Join']

from .base import CRUDElement

from enum import IntEnum, unique

from collections import Iterable

from six import string_types


@unique
class Join(IntEnum):

    INNER = 0  #: inner join.
    LEFT = 1  #: left join.
    LEFTEX = 2  #: left exclusive join.
    RIGHT = 3  #: right exclusive join.
    RIGHTEX = 4  #: right exclusive join.
    FULL = 5  #: full join.
    FULLEX = 6  #: full exclusive join.
    CROSS = 7  #: cross join.
    SELF = 8  #: self join.
    NATURAL = 9  #: natural join.
    UNION = 10  #: union join.


class Read(CRUDElement):
    """In charge of parameterize a reading request.

    Execution is done in calling it or in using the getslice method.
    Result is a Cursor."""

    __slots__ = [
        '_select', '_offset', '_limit', '_orderby', '_groupby', '_join'
    ] + CRUDElement.__slots__

    def __init__(
            self,
            select=None, offset=None, limit=None, orderby=None, groupby=None,
            join=None, callback=None, *args, **kwargs
    ):
        """
        :param list select: data to select.
        :param int offset: data to avoid.
        :param int limit: max number of data to retrieve.
        :param list orderby: data sorting.
        :param list groupby: data field group.
        :param join: join type (INNER, LEFT, etc.).
        :type join: str or Join
        """

        super(Read, self).__init__(*args, **kwargs)

        # initialize protected properties
        self._select = None
        self._offset = None
        self._limit = None
        self._orderby = None
        self._groupby = None
        self._join = None

        # set parameters
        if select is not None:
            self.select(*select)

        if offset is not None:
            self.offset(offset)

        if limit is not None:
            self.limit(limit)

        if orderby is not None:
            self.orderby(*orderby)

        if groupby is not None:
            self.groupby(*groupby)

        if join is not None:
            self.join(join)

    def offset(self, value=None):
        """Get or set offset if value is not None.

        :param int value: value to set. Default is None.
        :return: depending on value. If None, return this offset, otherwise
            this.
        :rtype: int or Read
        """

        if value is None:
            result = self._offset

        else:
            if not isinstance(value, int):
                raise TypeError(
                    'Wrong value {0}. {1} expected'.format(value, int)
                )

            result = self
            self._offset = value

        return result

    def limit(self, value=None):
        """Get or set limit if value is not None.

        :param int value: value to set. Default is None.
        :return: depending on value. If None, return this offset, otherwise
            this.
        :rtype: int or Read
        """

        if value is None:
            result = self._limit

        else:
            if not isinstance(value, int):
                raise TypeError(
                    'Wrong value {0}. {1} expected'.format(value, int)
                )

            result = self
            self._limit = value

        return result

    def orderby(self, *values):
        """Get or set orderby if value is not None.

        :param tuple value: value to set. Default is None.
        :return: depending on value. If None, return this offset, otherwise
            this.
        :rtype: tuple or Read
        """

        if values:
            self._orderby = values
            result = self

        else:
            result = self._orderby

        return result

    def groupby(self, *values):
        """Get or set groupby if value is not None.

        :param tuple value: value to set. Default is None.
        :return: depending on value. If None, return this offset, otherwise
            this.
        :rtype: tuple or Read
        """

        if values:
            self._groupby = values
            result = self

        else:
            result = self._groupby

        return result

    def select(self, *values):
        """Get or set select if value is not None.

        :param tuple value: value to set. Default is None.
        :return: depending on value. If None, return this offset, otherwise
            this.
        :rtype: tuple or Read
        """

        if values:
            self._select = values
            result = self

        else:
            result = self._select

        return result

    def join(self, value=None):
        """Get or set join if value is not None.

        :param value: value to set. Default is None.
        :type value: str or Join
        :return: depending on value. If None, return this offset, otherwise
            this.
        :rtype: str or Join or Read
        """

        if value is None:
            result = self._join

        else:
            if not isinstance(value, string_types + (Join,)):
                raise TypeError(
                    'Wrong value {0}. {1} expected'.format(
                        value, string_types + (Join,)
                    )
                )

            self._join = value.name if isinstance(value, Join) else value
            result = self

        return result

    def __getslice__(self, i, j):
        """Set offset and limit and execute the selection.

        :param int i: offset property.
        :param int j: limit property.
        :return: selection execution result.
        :rtype: Cursor"""

        if i is not None:
            self._offset = i

        if j is not None:
            self._limit = j

        return self()


class Cursor(Iterable):
    """Read request result."""

    def __init__(self, cursor, *args, **kwargs):

        super(Cursor, self).__init__(*args, **kwargs)

        self._cursor = cursor

    def __len__(self):

        return len(self._cursor)

    def __iter__(self):

        return iter(self._cursor)

    def __getitem__(self, key):

        return self._cursor[key]

    def __getslice__(self, i, j):

        return self._cursor[i:j]
