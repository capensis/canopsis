# -*- encoding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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

from grako.model import NodeWalker


class Walker(NodeWalker):
    def walk_forward_value(self, node):
        return self.walk(node.value)

    def walk_digits(self, node):
        return int(node.value)

    def walk_sign(self, node):
        """
        Convert sign to a positive or negative factor
        """
        if node.value == '-':
            return -1

        else:
            return 1

    def walk_characters(self, node):
        return node.value

    def walk_integer(self, node):
        sign = 1 if node.sign is None else self.walk(node.sign)
        value = self.walk(node.value)

        return sign * value

    def walk_floating(self, node):
        sign = 1 if node.sign is None else self.walk(node.sign)
        intpart = self.walk(node.intpart)
        floatpart = self.walk(node.floatpart)

        floating = sign * (intpart + 0.1 * floatpart)
        return floating

    def walk_true(self, node):
        return True

    def walk_false(self, node):
        return False

    def walk_none(self, node):
        return None

    def walk_condop(self, node):
        return node.value

    def walk_condition(self, node):
        left = self.walk(node.left)
        condop = self.walk(node.condop)
        right = self.walk(node.right)

        ret = 'left={} op={} right={}'.format(left, condop, right)

        return ret

    def walk_AST(self, node):
        return self.walk(node.request)
