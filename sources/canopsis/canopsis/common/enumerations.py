# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2017 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from __future__ import unicode_literals

from enum import Enum


class DefaultEnum(Enum):
    """
    Default simple Enum class

    Extends:
        Enum
    """

    def __str__(self):
        return str(self.value)


class FastEnumMeta(type):
    """
    https://github.com/Leryan/leryan.types/blob/v0.0.18/leryan/types/fastenum.py
    """

    def __new__(metacls, name, bases, attrs):
        slots = []
        _members = []
        _values = []
        _items = {}
        for attr in attrs:
            if not attr.startswith('_'):
                slots.append(attr)
                _members.append(attr)
                _values.append(attrs[attr])
                _items[attr] = attrs[attr]

        attrs['__slots__'] = tuple(slots)
        for member in _members:
            del attrs[member]

        klass = type.__new__(metacls, name, bases, attrs)

        super(metacls, klass).__setattr__('_members', frozenset(_members))
        super(metacls, klass).__setattr__('_values', tuple(_values))
        super(metacls, klass).__setattr__('_items', _items)

        return klass

    def __init__(cls, name, bases, attrs):
        type.__init__(cls, name, bases, attrs)

        for member, value in cls:
            super(FastEnumMeta, cls).__setattr__(member, value)

    def __setattr__(cls, attr, value):
        raise AttributeError("can't set attribute {}".format(attr))

    def __contains__(cls, attr):
        return hasattr(cls, attr)

    def __iter__(cls):
        for item in cls._items:
            yield (item, cls._items[item])

    def items(cls):
        return dict(cls.__iter__())

    @property
    def members(cls):
        return cls._members

    @property
    def values(cls):
        return cls._values


class FastEnum:
    """
    Fast and simple Enum implementation.

    .. code-block:: python
        class MyEnum(FastEnum):

            MEMBER = 'value'
            OTHER_MEMBER = 0

        MyEnum.MEMBER
        myenum_members = MyEnum.members
    """
    __metaclass__ = FastEnumMeta
