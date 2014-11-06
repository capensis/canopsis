#!/usr/bin/env python
# -*- coding: utf-8 -*-
#--------------------------------
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

from canopsis.rule import register_task, ActionError
from canopsis.rule.action import actions


@register_task
def action(event, ctx, **kwargs):
    pass


@register_task
def action_raise(event, ctx, **kwargs):
    raise Exception()


@register_task
def action_one(event, ctx, **kwargs):
    return 1


class ActionsTest(TestCase):

    def test_empty(self):
        result = actions(event=None, ctx=None)
        self.assertEqual(len(result), 0)

    def test_unary(self):
        actions_conf = ["action"]
        result = actions(event=None, ctx=None, actions=actions_conf)
        self.assertEqual(result, [None])

    def test_exception(self):
        actions_conf = ["action_raise"]
        result = actions(event=None, ctx=None, actions=actions_conf)
        self.assertIs(type(result[0]), ActionError)

    def test_raiseError(self):
        actions_conf = ["action_raise"]
        self.assertRaises(
            ActionError,
            actions,
            event=None, ctx=None, actions=actions_conf, raiseError=True)

    def test_many(self):
        actions_conf = ["action", "action_raise", "action_one"]
        result = actions(event=None, ctx=None, actions=actions_conf)
        self.assertIs(result[0], None)
        self.assertIs(type(result[1]), ActionError)
        self.assertIs(result[2], 1)


if __name__ == '__main__':
    main()
