#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from unittest import TestCase, main

from time import time

from dateutil import rrule

from datetime import datetime

from canopsis.rule import register_task
from canopsis.rule.condition import any, all, during


@register_task
def condition(event, ctx, **kwargs):
    pass


@register_task
def condition_raise(event, ctx, **kwargs):
    raise Exception()


@register_task
def condition_true(event, ctx, **kwargs):
    return True


@register_task
def condition_false(event, ctx, **kwargs):
    return False


class DuringTest(TestCase):

    def test_inside(self):
        duration = {"seconds": 5}
        now = time()
        now = datetime.fromtimestamp(now)
        rrule_p = {"freq": rrule.DAILY, "dtstart": now}
        result = during(event=None, ctx=None, rrule=rrule_p, duration=duration)
        self.assertTrue(result)

    def test_outside(self):
        now = time()
        now = datetime.fromtimestamp(now + 5)
        rrule_p = {"freq": rrule.DAILY, "dtstart": now}
        result = during(event=None, ctx=None, rrule=rrule_p)
        self.assertFalse(result)


class AnyTest(TestCase):

    def test_empty(self):
        result = any(event=None, ctx=None)
        self.assertFalse(result)

    def test_unary(self):
        conditions = ["condition_true"]
        result = any(event=None, ctx=None, conditions=conditions)
        self.assertTrue(result)

    def test_false(self):
        conditions = ["condition_false", "condition_false"]
        result = any(event=None, ctx=None, conditions=conditions)
        self.assertFalse(result)

    def test_true(self):
        conditions = ["condition_false", "condition_true", "condition_true"]
        result = any(event=None, ctx=None, conditions=conditions)
        self.assertTrue(result)

    def test_all_true(self):
        conditions = ["condition_true", "condition_true", "condition_true"]
        result = any(event=None, ctx=None, conditions=conditions)
        self.assertTrue(result)

    def test_raise(self):
        conditions = ["condition_raise"]
        self.assertRaises(
            Exception,
            any,
            event=None, ctx=None, conditions=conditions)


class AllTest(TestCase):

    def test_empty(self):
        result = all(event=None, ctx=None)
        self.assertFalse(result)

    def test_unary(self):
        conditions = ["condition_true"]
        result = all(event=None, ctx=None, conditions=conditions)
        self.assertTrue(result)

    def test_false(self):
        conditions = ["condition_false", "condition_true", "condition_false"]
        result = all(event=None, ctx=None, conditions=conditions)
        self.assertFalse(result)

    def test_true(self):
        conditions = ["condition_true", "condition_true", "condition_true"]
        result = all(event=None, ctx=None, conditions=conditions)
        self.assertTrue(result)

    def test_raise(self):
        conditions = ["condition_raise"]
        self.assertRaises(
            Exception,
            all,
            event=None, ctx=None, conditions=conditions)


if __name__ == '__main__':
    main()
